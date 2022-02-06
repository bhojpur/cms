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

    let FormData = window.FormData,
        BHOJPUR = window.BHOJPUR,
        NAMESPACE = 'bhojpur.selectcore',
        EVENT_SELECTCORE_BEFORESEND = 'selectcoreBeforeSend.' + NAMESPACE,
        EVENT_ONSELECT = 'afterSelected.' + NAMESPACE,
        EVENT_ONSUBMIT = 'afterSubmitted.' + NAMESPACE,
        EVENT_CLICK = 'click.' + NAMESPACE,
        EVENT_SUBMIT = 'submit.' + NAMESPACE,
        CLASS_TABLE = 'table.bhojpur-js-table tr',
        CLASS_FORM = 'form';

    function BhojpurSelectCore(element, options) {
        this.$element = $(element);
        this.options = $.extend({}, BhojpurSelectCore.DEFAULTS, $.isPlainObject(options) && options);
        this.init();
    }

    BhojpurSelectCore.prototype = {
        constructor: BhojpurSelectCore,

        init: function() {
            this.bind();
        },

        bind: function() {
            this.$element.on(EVENT_CLICK, CLASS_TABLE, this.processingData.bind(this)).on(EVENT_SUBMIT, CLASS_FORM, this.submit.bind(this));
        },

        unbind: function() {
            this.$element.off(EVENT_CLICK, CLASS_TABLE).off(EVENT_SUBMIT, CLASS_FORM);
        },

        processingData: function(e) {
            let $this = $(e.target).closest('tr'),
                $bottomsheets = $this.closest('.bhojpur-bottomsheets'),
                data = {},
                url,
                options = this.options,
                onSelect = options.onSelect,
                loading = options.loading;

            data = $.extend({}, data, $this.data());
            data.$clickElement = $this;

            url = data.mediaLibraryUrl || data.url;

            if (loading && $.isFunction(loading)) {
                loading($bottomsheets);
            }

            if (url) {

                $.getJSON(url, function(json) {
                    json.MediaOption && (json.MediaOption = JSON.parse(json.MediaOption));
                    data = $.extend({}, json, data);
                    if (onSelect && $.isFunction(onSelect)) {
                        onSelect(data, e);
                        $(document).trigger(EVENT_ONSELECT);
                    }
                }).always(function() {
                    $bottomsheets.find('.bhojpur-media-loading').remove();
                  });

            } else {
                if (onSelect && $.isFunction(onSelect)) {
                    onSelect(data, e);
                    $(document).trigger(EVENT_ONSELECT);
                }
            }
            return false;
        },

        submit: function(e) {
            let form = e.target,
                $form = $(form),
                _this = this,
                $submit = $form.find(':submit'),
                data,
                $loading = $(BHOJPUR.$formLoading),
                onSubmit = this.options.onSubmit;

            $(document).trigger(EVENT_SELECTCORE_BEFORESEND);

            $form.find('.bhojpur-fieldset--new').remove();

            if (FormData) {
                e.preventDefault();

                $.ajax($form.prop('action'), {
                    method: $form.prop('method'),
                    data: new FormData(form),
                    dataType: 'json',
                    processData: false,
                    contentType: false,
                    beforeSend: function() {
                        $('.bhojpur-submit-loading').remove();
                        $loading.appendTo($submit.prop('disabled', true).closest('.bhojpur-form__actions')).trigger('enable.bhojpur.material');
                    },
                    success: function(json) {
                        json.MediaOption && (json.MediaOption = JSON.parse(json.MediaOption));
                        data = json;
                        data.primaryKey = data.ID;

                        $('.bhojpur-error').remove();

                        if (onSubmit && $.isFunction(onSubmit)) {
                            onSubmit(data, e);
                            $(document).trigger(EVENT_ONSUBMIT);
                        } else {
                            _this.refresh();
                        }
                    },
                    error: function(err) {
                        BHOJPUR.handleAjaxError(err);
                    },
                    complete: function() {
                        $submit.prop('disabled', false);
                    }
                });
            }
        },

        refresh: function() {
            setTimeout(function() {
                window.location.reload();
            }, 350);
        },

        destroy: function() {
            this.unbind();
        }
    };

    BhojpurSelectCore.plugin = function(options) {
        return this.each(function() {
            let $this = $(this),
                data = $this.data(NAMESPACE),
                fn;

            if (!data) {
                if (/destroy/.test(options)) {
                    return;
                }
                $this.data(NAMESPACE, (data = new BhojpurSelectCore(this, options)));
            }

            if (typeof options === 'string' && $.isFunction((fn = data[options]))) {
                fn.apply(data);
            }
        });
    };

    $.fn.bhojpurSelectCore = BhojpurSelectCore.plugin;

    return BhojpurSelectCore;
});