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

    let NAMESPACE = 'bhojpur.datepicker',
        EVENT_ENABLE = 'enable.' + NAMESPACE,
        EVENT_DISABLE = 'disable.' + NAMESPACE,
        EVENT_CHANGE = 'pick.' + NAMESPACE,
        EVENT_CLICK = 'click.' + NAMESPACE,
        CLASS_EMBEDDED = '.bhojpur-datepicker__embedded',
        CLASS_SAVE = '.bhojpur-datepicker__save',
        CLASS_PARENT = '[data-picker-type]';

    function replaceText(str, data) {
        if (typeof str === 'string') {
            if (typeof data === 'object') {
                $.each(data, function(key, val) {
                    str = str.replace('$[' + String(key).toLowerCase() + ']', val);
                });
            }
        }

        return str;
    }

    function BhojpurDatepicker(element, options) {
        this.$element = $(element);
        this.options = $.extend(true, {}, BhojpurDatepicker.DEFAULTS, $.isPlainObject(options) && options);
        this.date = null;
        this.formatDate = null;
        this.built = false;
        this.pickerData = this.$element.data();
        this.$parent = this.$element.closest(CLASS_PARENT);
        this.isDateTimePicker = this.$parent.data('picker-type') == 'datetime';
        this.$targetInput = this.$parent.find(this.pickerData.targetInput || (this.isDateTimePicker ? '.bhojpur-datetimepicker__input' : '.bhojpur-datepicker__input'));
        this.init();
    }

    BhojpurDatepicker.prototype = {
        init: function() {
            if (this.$targetInput.is(':disabled')) {
                this.$element.remove();
                return;
            }
            this.bind();
        },

        bind: function() {
            this.$element.on(EVENT_CLICK, $.proxy(this.show, this));
        },

        unbind: function() {
            this.$element.off(EVENT_CLICK, this.show);
        },

        build: function() {
            let $modal,
                $ele = this.$element,
                $targetInput = this.$targetInput,
                defaultDate = $targetInput.val(),
                datepickerOptions = {
                    date: new Date(),
                    inline: true
                };

            if (this.built) {
                return;
            }

            if ($ele.is(':input') && Date.parse($ele.val())) {
                datepickerOptions.date = new Date($ele.val());
            } else if (defaultDate && Date.parse(defaultDate)) {
                datepickerOptions.date = new Date(defaultDate);
            }

            this.$modal = $modal = $(replaceText(BhojpurDatepicker.TEMPLATE, this.options.text)).appendTo('body');

            if ($targetInput.data('start-date')) {
                datepickerOptions.startDate = new Date($targetInput.data('start-date'));
            }

            if ($targetInput.data('end-date')) {
                datepickerOptions.endDate = new Date($targetInput.data('end-date'));
            }


            $modal
                .find(CLASS_EMBEDDED)
                .on(EVENT_CHANGE, $.proxy(this.change, this))
                .bhojpurDatepicker(datepickerOptions)
                .triggerHandler(EVENT_CHANGE);

            $modal.find(CLASS_SAVE).on(EVENT_CLICK, $.proxy(this.pick, this));

            this.built = true;
        },

        unbuild: function() {
            if (!this.built) {
                return;
            }

            this.$modal
                .find(CLASS_EMBEDDED)
                .off(EVENT_CHANGE, this.change)
                .bhojpurDatepicker('destroy')
                .end()
                .find(CLASS_SAVE)
                .off(EVENT_CLICK, this.pick)
                .end()
                .remove();
        },

        change: function(e) {
            var $modal = this.$modal;
            var $target = $(e.target);
            var date;

            this.date = date = $target.bhojpurDatepicker('getDate');
            this.formatDate = $target.bhojpurDatepicker('getDate', true);

            $modal.find('.bhojpur-datepicker__picked-year').text(date.getFullYear());
            $modal
                .find('.bhojpur-datepicker__picked-date')
                .text(
                    [$target.bhojpurDatepicker('getDayName', date.getDay(), true) + ',', String($target.bhojpurDatepicker('getMonthName', date.getMonth(), true)), date.getDate()].join(' ')
                );
        },

        show: function() {
            if (!this.built) {
                this.build();
            }

            this.$modal.bhojpurModal('show');
        },

        pick: function() {
            let $targetInput = this.$targetInput,
                newValue = this.formatDate;

            if (this.isDateTimePicker) {
                var regDate = /^\d{4}-\d{1,2}-\d{1,2}/;
                var oldValue = $targetInput.val();
                var hasDate = regDate.test(oldValue);

                if (hasDate) {
                    newValue = oldValue.replace(regDate, newValue);
                } else {
                    newValue = newValue + ' 00:00';
                }
            }

            $targetInput.val(newValue).trigger('change');
            this.$modal.bhojpurModal('hide');
        },

        destroy: function() {
            this.unbind();
            this.unbuild();
            this.$element.removeData(NAMESPACE);
        }
    };

    BhojpurDatepicker.DEFAULTS = {
        text: {
            title: 'Pick a date',
            ok: 'OK',
            cancel: 'Cancel'
        }
    };

    BhojpurDatepicker.TEMPLATE = `<div class="bhojpur-modal fade bhojpur-datepicker" tabindex="-1" role="dialog" aria-hidden="true">
            <div class="mdl-card mdl-shadow--2dp" role="document">
                <div class="mdl-card__title">
                    <h2 class="mdl-card__title-text">$[title]</h2>
                </div>
                <div class="mdl-card__supporting-text">
                    <div class="bhojpur-datepicker__picked">
                        <div class="bhojpur-datepicker__picked-year"></div>
                        <div class="bhojpur-datepicker__picked-date"></div>
                    </div>
                    <div class="bhojpur-datepicker__embedded"></div>
                </div>
                <div class="mdl-card__actions">
                    <a class="mdl-button mdl-button--colored  mdl-button--raised bhojpur-datepicker__save">$[ok]</a>
                    <a class="mdl-button mdl-button--colored " data-dismiss="modal">$[cancel]</a>
                </div>
            </div>
        </div>`;

    BhojpurDatepicker.plugin = function(option) {
        return this.each(function() {
            var $this = $(this);
            var data = $this.data(NAMESPACE);
            var options;
            var fn;

            if (!data) {
                if (!$.fn.bhojpurDatepicker) {
                    return;
                }

                if (/destroy/.test(option)) {
                    return;
                }

                options = $.extend(true, {}, $this.data(), typeof option === 'object' && option);
                $this.data(NAMESPACE, (data = new BhojpurDatepicker(this, options)));
            }

            if (typeof option === 'string' && $.isFunction((fn = data[option]))) {
                fn.apply(data);
            }
        });
    };

    $(function() {
        var selector = '[data-toggle="bhojpur.datepicker"]';

        $(document)
            .on(EVENT_DISABLE, function(e) {
                BhojpurDatepicker.plugin.call($(selector, e.target), 'destroy');
            })
            .on(EVENT_ENABLE, function(e) {
                BhojpurDatepicker.plugin.call($(selector, e.target));
            })
            .triggerHandler(EVENT_ENABLE);
    });

    return BhojpurDatepicker;
});