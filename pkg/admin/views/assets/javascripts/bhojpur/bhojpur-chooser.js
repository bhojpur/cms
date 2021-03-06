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

    var NAMESPACE = 'bhojpur.chooser';
    var EVENT_ENABLE = 'enable.' + NAMESPACE;
    var EVENT_DISABLE = 'disable.' + NAMESPACE;

    function BhojpurChooser(element, options) {
        this.$element = $(element);
        this.options = $.extend({}, BhojpurChooser.DEFAULTS, $.isPlainObject(options) && options);
        this.init();
    }

    BhojpurChooser.prototype = {
        constructor: BhojpurChooser,

        init: function() {
            let $this = this.$element,
                select2Data = $this.data(),
                resetSelect2Width,
                option = {
                    minimumResultsForSearch: 8,
                    dropdownParent: $this.parent()
                };
            let getSelect2AjaxDynamicURL = window.getSelect2AjaxDynamicURL;
            let remoteImage = select2Data.remoteImage;

            if (select2Data.remoteData) {
                option.ajax = $.fn.select2.ajaxCommonOptions(select2Data);
                if (getSelect2AjaxDynamicURL && $.isFunction(getSelect2AjaxDynamicURL)) {
                    option.ajax.url = function() {
                        return getSelect2AjaxDynamicURL(select2Data);
                    };
                } else {
                    option.ajax.url = select2Data.remoteUrl;
                }

                option.templateResult = function(data) {
                    let tmpl = $this.parents('.bhojpur-field').find('[name="select2-result-template"]');
                    return $.fn.select2.ajaxFormatResult(data, tmpl, remoteImage);
                };

                option.templateSelection = function(data) {
                    if (data.loading) return data.text;
                    let tmpl = $this.parents('.bhojpur-field').find('[name="select2-selection-template"]');
                    return $.fn.select2.ajaxFormatResult(data, tmpl, remoteImage);
                };
            }

            $this
                .on('select2:select', function(evt) {
                    $(evt.target).attr('chooser-selected', 'true');
                })
                .on('select2:unselect', function(evt) {
                    $(evt.target).attr('chooser-selected', '');
                });

            $this.select2(option);

            // reset select2 container width
            this.resetSelect2Width();
            resetSelect2Width = window._.debounce(this.resetSelect2Width.bind(this), 300);
            $(window).resize(resetSelect2Width);

            if ($this.val()) {
                $this.attr('chooser-selected', 'true');
            }
        },

        resetSelect2Width: function() {
            var $container,
                select2 = this.$element.data().select2;
            if (select2 && select2.$container) {
                $container = select2.$container;
                $container.width($container.parent().width());
            }
        },

        destroy: function() {
            this.$element.select2('destroy').removeData(NAMESPACE);
        }
    };

    BhojpurChooser.DEFAULTS = {};

    BhojpurChooser.plugin = function(options) {
        return this.each(function() {
            var $this = $(this);
            var data = $this.data(NAMESPACE);
            var fn;

            if (!data) {
                if (/destroy/.test(options)) {
                    return;
                }

                $this.data(NAMESPACE, (data = new BhojpurChooser(this, options)));
            }

            if (typeof options === 'string' && $.isFunction((fn = data[options]))) {
                fn.apply(data);
            }
        });
    };

    $(function() {
        var selector = 'select[data-toggle="bhojpur.chooser"]';

        $(document)
            .on(EVENT_DISABLE, function(e) {
                BhojpurChooser.plugin.call($(selector, e.target), 'destroy');
            })
            .on(EVENT_ENABLE, function(e) {
                BhojpurChooser.plugin.call($(selector, e.target));
            })
            .triggerHandler(EVENT_ENABLE);
    });

    return BhojpurChooser;
});