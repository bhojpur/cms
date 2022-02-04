package publish

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

	orm "github.com/bhojpur/orm/pkg/engine"
)

func isProductionModeAndNewScope(scope *orm.Scope) (isProduction bool, clone *orm.Scope) {
	if !IsDraftMode(scope.DB()) {
		if _, ok := scope.InstanceGet("publish:supported_model"); ok {
			table := OriginalTableName(scope.TableName())
			clone := scope.New(scope.Value)
			clone.Search.Table(table)
			return true, clone
		}
	}
	return false, nil
}

func setTableAndPublishStatus(ensureDraftMode bool) func(*orm.Scope) {
	return func(scope *orm.Scope) {
		if scope.Value == nil {
			return
		}

		if IsPublishableModel(scope.Value) {
			scope.InstanceSet("publish:supported_model", true)

			if ensureDraftMode {
				scope.Set("publish:force_draft_table", true)
				scope.Search.Table(DraftTableName(scope.TableName()))

				// Only set publish status when updating data from draft tables
				if IsDraftMode(scope.DB()) {
					if _, ok := scope.DB().Get(publishEvent); ok {
						scope.InstanceSet("publish:creating_publish_event", true)
					} else {
						if attrs, ok := scope.InstanceGet("orm:update_attrs"); ok {
							updateAttrs := attrs.(map[string]interface{})
							updateAttrs["publish_status"] = DIRTY
							scope.InstanceSet("orm:update_attrs", updateAttrs)
						} else {
							scope.SetColumn("PublishStatus", DIRTY)
						}
					}
				}
			}
		}
	}
}

func syncCreateFromProductionToDraft(scope *orm.Scope) {
	if !scope.HasError() {
		if ok, clone := isProductionModeAndNewScope(scope); ok {
			scope.DB().Callback().Create().Get("orm:create")(clone)
		}
	}
}

func syncUpdateFromProductionToDraft(scope *orm.Scope) {
	if !scope.HasError() {
		if ok, clone := isProductionModeAndNewScope(scope); ok {
			if updateAttrs, ok := scope.InstanceGet("orm:update_attrs"); ok {
				table := OriginalTableName(scope.TableName())
				clone.Search = scope.Search
				clone.Search.Table(table)
				clone.InstanceSet("orm:update_attrs", updateAttrs)
			}
			scope.DB().Callback().Update().Get("orm:update")(clone)
		}
	}
}

func syncDeleteFromProductionToDraft(scope *orm.Scope) {
	if !scope.HasError() {
		if ok, clone := isProductionModeAndNewScope(scope); ok {
			scope.DB().Callback().Delete().Get("orm:delete")(clone)
		}
	}
}

func deleteScope(scope *orm.Scope) {
	if !scope.HasError() {
		_, supportedModel := scope.InstanceGet("publish:supported_model")

		if !scope.Search.Unscoped && supportedModel && IsDraftMode(scope.DB()) {
			scope.Raw(
				fmt.Sprintf("UPDATE %v SET deleted_at=%v, publish_status=%v %v",
					scope.QuotedTableName(),
					scope.AddToVars(orm.NowFunc()),
					scope.AddToVars(DIRTY),
					scope.CombinedConditionSql(),
				))
			scope.Exec()
		} else {
			scope.DB().Callback().Delete().Get("orm:delete")(scope)
		}
	}
}

func createPublishEvent(scope *orm.Scope) {
	if _, ok := scope.InstanceGet("publish:creating_publish_event"); ok {
		if event, ok := scope.Get(publishEvent); ok {
			if event, ok := event.(*PublishEvent); ok {
				event.PublishStatus = DIRTY
				scope.Err(scope.NewDB().Save(&event).Error)
			}
		}
	}
}
