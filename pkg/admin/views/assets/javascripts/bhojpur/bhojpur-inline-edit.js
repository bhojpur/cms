(function(factory) {
    if (typeof define === 'function' && define.amd) {
        // AMD. Register as anonymous module.
        define(['jquery'], factory);
    } else if (typeof exports === 'object') {
        // Node / CommonJS
        factory(require('jquery'));
    } else {
        // Browser globals.
        factory(jQuery);
    }
})(function($) {
    'use strict';

    const NAMESPACE = 'bhojpur.inlineEdit',
        EVENT_ENABLE = 'enable.' + NAMESPACE,
        EVENT_DISABLE = 'disable.' + NAMESPACE,
        EVENT_CLICK = 'click.' + NAMESPACE,
        EVENT_MOUSEENTER = 'mouseenter.' + NAMESPACE,
        EVENT_MOUSELEAVE = 'mouseleave.' + NAMESPACE,
        CLASS_FIELD = '.bhojpur-field',
        CLASS_FIELD_SHOW = '.bhojpur-field__show',
        CLASS_FIELD_SHOW_INNER = '.bhojpur-field__show-inner',
        CLASS_EDIT = '.bhojpur-inlineedit__edit',
        CLASS_SAVE = '.bhojpur-inlineedit__save',
        CLASS_BUTTONS = '.bhojpur-inlineedit__buttons',
        CLASS_CANCEL = '.bhojpur-inlineedit__cancel',
        CLASS_CONTAINER = 'bhojpur-inlineedit__field';

    function BhojpurInlineEdit(element, options) {
        this.$element = $(element);
        this.options = $.extend({}, BhojpurInlineEdit.DEFAULTS, $.isPlainObject(options) && options);
        this.init();
    }

    function getJsonData(names, data) {
        let key,
            value = data[names[0].slice(1)];

        if (names.length > 1) {
            for (let i = 1; i < names.length; i++) {
                key = names[i].slice(1);
                value = $.isArray(value) ? value[0][key] : value[key];
            }
        }

        return value;
    }

    BhojpurInlineEdit.prototype = {
        constructor: BhojpurInlineEdit,

        init: function() {
            let $element = this.$element,
                saveButton = $element.data('button-save'),
                cancelButton = $element.data('button-cancel');

            this.TEMPLATE_SAVE = `<div class="bhojpur-inlineedit__buttons">
                                        <button class="mdl-button mdl-button--colored mdl-js-button bhojpur-button--small bhojpur-inlineedit__cancel" type="button">${cancelButton}</button>
                                        <button class="mdl-button mdl-button--colored mdl-js-button bhojpur-button--small bhojpur-inlineedit__save" type="button">${saveButton}</button>
                                      </div>`;
            this.bind();
        },

        bind: function() {
            this.$element
                .on(EVENT_MOUSEENTER, CLASS_FIELD_SHOW, this.showEditButton)
                .on(EVENT_MOUSELEAVE, CLASS_FIELD_SHOW, this.hideEditButton)
                .on(EVENT_CLICK, CLASS_CANCEL, this.hideEdit)
                .on(EVENT_CLICK, CLASS_SAVE, this.saveEdit)
                .on(EVENT_CLICK, CLASS_EDIT, this.showEdit.bind(this));
        },

        unbind: function() {
            this.$element
                .off(EVENT_MOUSEENTER)
                .off(EVENT_MOUSELEAVE)
                .off(EVENT_CLICK);
        },

        showEditButton: function(e) {
            let $edit = $(BhojpurInlineEdit.TEMPLATE_EDIT),
                $field = $(e.target).closest(CLASS_FIELD);

            if ($field.find('input:disabled, textarea:disabled,select:disabled').length) {
                return false;
            }

            $edit.appendTo($(this));
        },

        hideEditButton: function() {
            $('.bhojpur-inlineedit__edit').remove();
        },

        showEdit: function(e) {
            let $parent = $(e.target)
                    .closest(CLASS_EDIT)
                    .hide()
                    .closest(CLASS_FIELD)
                    .addClass(CLASS_CONTAINER),
                $save = $(this.TEMPLATE_SAVE);

            $save.appendTo($parent);
        },

        hideEdit: function() {
            let $parent = $(this)
                .closest(CLASS_FIELD)
                .removeClass(CLASS_CONTAINER);
            $parent.find(CLASS_BUTTONS).remove();
        },

        saveEdit: function() {
            let $btn = $(this),
                $parent = $btn.closest(CLASS_FIELD),
                $form = $btn.closest('form'),
                $hiddenInput = $parent.closest('.bhojpur-fieldset').find('input.bhojpur-hidden__primary_key[type="hidden"]'),
                $input = $parent.find('input[name*="BhojpurResource"],textarea[name*="BhojpurResource"],select[name*="BhojpurResource"]'),
                names = $input.length && $input.prop('name').match(/\.\w+/g),
                inputData = $input.serialize();

            if ($hiddenInput.length) {
                inputData = `${inputData}&${$hiddenInput.serialize()}`;
            }

            if (names.length) {
                $.ajax($form.prop('action'), {
                    method: $form.prop('method'),
                    data: inputData,
                    dataType: 'json',
                    beforeSend: function() {
                        $btn.prop('disabled', true);
                    },
                    success: function(data) {
                        let newValue = getJsonData(names, data),
                            $show = $parent.removeClass(CLASS_CONTAINER).find(CLASS_FIELD_SHOW),
                            $inner = $show.find(CLASS_FIELD_SHOW_INNER);

                        if (typeof newValue === 'string' || newValue instanceof String){
                            newValue = newValue.escapeSymbol();
                        }
                        
                        if ($inner.length) {
                            $inner.html(newValue);
                        } else {
                            $show.html(newValue);
                        }

                        $parent.find(CLASS_BUTTONS).remove();
                        $btn.prop('disabled', false);
                    },
                    error: function(err) {
                        window.BHOJPUR.handleAjaxError(err);
                        $btn.prop('disabled', false);
                    }
                });
            }
        },

        destroy: function() {
            this.unbind();
            this.$element.removeData(NAMESPACE);
        }
    };

    BhojpurInlineEdit.DEFAULTS = {};

    BhojpurInlineEdit.TEMPLATE_EDIT = `<button class="mdl-button mdl-js-button mdl-button--icon mdl-button--colored bhojpur-inlineedit__edit" type="button"><i class="material-icons">mode_edit</i></button>`;

    BhojpurInlineEdit.plugin = function(options) {
        return this.each(function() {
            var $this = $(this);
            var data = $this.data(NAMESPACE);
            var fn;

            if (!data) {
                $this.data(NAMESPACE, (data = new BhojpurInlineEdit(this, options)));
            }

            if (typeof options === 'string' && $.isFunction((fn = data[options]))) {
                fn.call(data);
            }
        });
    };

    $(function() {
        let selector = '[data-toggle="bhojpur.inlineEdit"]',
            options = {};

        $(document)
            .on(EVENT_DISABLE, function(e) {
                BhojpurInlineEdit.plugin.call($(selector, e.target), 'destroy');
            })
            .on(EVENT_ENABLE, function(e) {
                BhojpurInlineEdit.plugin.call($(selector, e.target), options);
            })
            .triggerHandler(EVENT_ENABLE);
    });

    return BhojpurInlineEdit;
});