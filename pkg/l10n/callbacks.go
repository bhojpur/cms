package l10n

// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

import (
	"fmt"
	"reflect"

	"github.com/bhojpur/application/pkg/utils"
	orm "github.com/bhojpur/orm/pkg/engine"
)

func beforeQuery(scope *orm.Scope) {
	if IsLocalizable(scope) {
		quotedTableName := scope.QuotedTableName()
		quotedPrimaryKey := scope.Quote(scope.PrimaryKey())
		_, hasDeletedAtColumn := scope.FieldByName("deleted_at")

		locale, isLocale := getQueryLocale(scope)
		switch mode, _ := scope.DB().Get("l10n:mode"); mode {
		case "unscoped":
		case "global":
			scope.Search.Where(fmt.Sprintf("%v.language_code = ?", quotedTableName), Global)
		case "locale":
			scope.Search.Where(fmt.Sprintf("%v.language_code = ?", quotedTableName), locale)
		case "reverse":
			if !scope.Search.Unscoped && hasDeletedAtColumn {
				scope.Search.Where(fmt.Sprintf(
					"(%v.%v NOT IN (SELECT DISTINCT(%v) FROM %v t2 WHERE t2.language_code = ? AND t2.deleted_at IS NULL) AND %v.language_code = ?)", quotedTableName, quotedPrimaryKey, quotedPrimaryKey, quotedTableName, quotedTableName), locale, Global)
			} else {
				scope.Search.Where(fmt.Sprintf("(%v.%v NOT IN (SELECT DISTINCT(%v) FROM %v t2 WHERE t2.language_code = ?) AND %v.language_code = ?)", quotedTableName, quotedPrimaryKey, quotedPrimaryKey, quotedTableName, quotedTableName), locale, Global)
			}
		case "fallback":
			fallthrough
		default:
			if isLocale {
				if !scope.Search.Unscoped && hasDeletedAtColumn {
					scope.Search.Where(fmt.Sprintf("((%v.%v NOT IN (SELECT DISTINCT(%v) FROM %v t2 WHERE t2.language_code = ? AND t2.deleted_at IS NULL) AND %v.language_code = ?) OR %v.language_code = ?) AND %v.deleted_at IS NULL", quotedTableName, quotedPrimaryKey, quotedPrimaryKey, quotedTableName, quotedTableName, quotedTableName, quotedTableName), locale, Global, locale)
				} else {
					scope.Search.Where(fmt.Sprintf("(%v.%v NOT IN (SELECT DISTINCT(%v) FROM %v t2 WHERE t2.language_code = ?) AND %v.language_code = ?) OR (%v.language_code = ?)", quotedTableName, quotedPrimaryKey, quotedPrimaryKey, quotedTableName, quotedTableName, quotedTableName), locale, Global, locale)
				}
				scope.Search.Order(orm.Expr(fmt.Sprintf("%v.language_code = ? DESC", quotedTableName), locale))
			} else {
				scope.Search.Where(fmt.Sprintf("%v.language_code = ?", quotedTableName), Global)
			}
		}
	}
}

func beforeCreate(scope *orm.Scope) {
	if IsLocalizable(scope) {
		if locale, ok := getLocale(scope); ok { // is locale
			if isLocaleCreatable(scope) || !scope.PrimaryKeyZero() {
				setLocale(scope, locale)
			} else {
				err := fmt.Errorf("the resource %v cannot be created in %v", scope.GetModelStruct().ModelType.Name(), locale)
				scope.Err(err)
			}
		} else {
			setLocale(scope, Global)
		}
	}
}

func beforeUpdate(scope *orm.Scope) {
	if IsLocalizable(scope) {
		locale, isLocale := getLocale(scope)

		switch mode, _ := scope.DB().Get("l10n:mode"); mode {
		case "unscoped":
		default:
			scope.Search.Where(fmt.Sprintf("%v.language_code = ?", scope.QuotedTableName()), locale)
			setLocale(scope, locale)
		}

		if isLocale {
			scope.Search.Omit(syncColumns(scope)...)
		}
	}
}

func afterUpdate(scope *orm.Scope) {
	if !scope.HasError() {
		if IsLocalizable(scope) {
			if locale, ok := getLocale(scope); ok {
				if scope.DB().RowsAffected == 0 && !scope.PrimaryKeyZero() { //is locale and nothing updated
					var count int
					var query = fmt.Sprintf("%v.language_code = ? AND %v.%v = ?", scope.QuotedTableName(), scope.QuotedTableName(), scope.PrimaryKey())

					// if enabled soft delete, delete soft deleted records
					if scope.HasColumn("DeletedAt") {
						scope.NewDB().Unscoped().Where("deleted_at is not null").Where(query, locale, scope.PrimaryKeyValue()).Delete(scope.Value)
					}

					// if no localized records exist, localize it
					if scope.NewDB().Table(scope.TableName()).Where(query, locale, scope.PrimaryKeyValue()).Count(&count); count == 0 {
						scope.DB().RowsAffected = scope.DB().Create(scope.Value).RowsAffected
					}
				}
			} else if syncColumns := syncColumns(scope); len(syncColumns) > 0 { // is global
				if mode, _ := scope.DB().Get("l10n:mode"); mode != "unscoped" {
					if scope.DB().RowsAffected > 0 {
						var primaryField = scope.PrimaryField()
						var syncAttrs = map[string]interface{}{}

						if updateAttrs, ok := scope.InstanceGet("orm:update_attrs"); ok {
							for key, value := range updateAttrs.(map[string]interface{}) {
								for _, syncColumn := range syncColumns {
									if syncColumn == key {
										syncAttrs[syncColumn] = value
										break
									}
								}
							}
						} else {
							for _, syncColumn := range syncColumns {
								if field, ok := scope.FieldByName(syncColumn); ok && field.IsNormal {
									syncAttrs[syncColumn] = field.Field.Interface()
								}
							}
						}

						if len(syncAttrs) > 0 {
							db := scope.DB().Model(reflect.New(utils.ModelType(scope.Value)).Interface()).Set("l10n:mode", "unscoped").Where("language_code <> ?", Global)
							if !primaryField.IsBlank {
								db = db.Where(fmt.Sprintf("%v = ?", primaryField.DBName), primaryField.Field.Interface())
							}
							scope.Err(db.UpdateColumns(syncAttrs).Error)
						}
					}
				}
			}
		}
	}
}

func beforeDelete(scope *orm.Scope) {
	if IsLocalizable(scope) {
		if locale, ok := getQueryLocale(scope); ok { // is locale
			scope.Search.Where(fmt.Sprintf("%v.language_code = ?", scope.QuotedTableName()), locale)
		}
	}
}

// RegisterCallbacks register callback into Bhojpur ORM DB
func RegisterCallbacks(db *orm.DB) {
	callback := db.Callback()

	if callback.Create().Get("l10n:before_create") == nil {
		callback.Create().Before("orm:before_create").Register("l10n:before_create", beforeCreate)
	}

	if callback.Update().Get("l10n:before_update") == nil {
		callback.Update().Before("orm:before_update").Register("l10n:before_update", beforeUpdate)
	}
	if callback.Update().Get("l10n:after_update") == nil {
		callback.Update().After("orm:after_update").Register("l10n:after_update", afterUpdate)
	}

	if callback.Delete().Get("l10n:before_delete") == nil {
		callback.Delete().Before("orm:before_delete").Register("l10n:before_delete", beforeDelete)
	}

	if callback.RowQuery().Get("l10n:before_query") == nil {
		callback.RowQuery().Before("orm:row_query").Register("l10n:before_query", beforeQuery)
	}
	if callback.Query().Get("l10n:before_query") == nil {
		callback.Query().Before("orm:query").Register("l10n:before_query", beforeQuery)
	}
}
