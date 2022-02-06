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

    var location = window.location;

    var NAMESPACE = 'bhojpur.i18n';
    var EVENT_CLICK = 'click.' + NAMESPACE;
    var EVENT_CHANGE = 'change.' + NAMESPACE;

    // For Bhojpur Autoheight plugin
    var EVENT_INPUT = 'input';

    function I18N(element, options) {
        this.$element = $(element);
        this.options = $.extend({}, I18N.DEFAULTS, $.isPlainObject(options) && options);
        this.multiple = false;
        this.init();
    }

    function encodeSearch(data) {
        var params = [];

        if ($.isPlainObject(data)) {
            $.each(data, function(name, value) {
                params.push([name, value].join('='));
            });
        }

        return params.join('&');
    }

    function decodeSearch(search) {
        var data = {};

        if (search) {
            search = search.replace('?', '').split('&');

            $.each(search, function(i, param) {
                param = param.split('=');
                i = param[0];
                data[i] = param[1];
            });
        }

        return data;
    }

    I18N.prototype = {
        contructor: I18N,

        init: function() {
            var $this = this.$element;

            this.$languages = $this.find('.bhojpur-js-language');
            this.$items = $this.find('.i18n-list > li');
            this.bind();
        },

        bind: function() {
            this.$element.on(EVENT_CLICK, $.proxy(this.click, this)).on(EVENT_CHANGE, $.proxy(this.change, this));
            this.$languages.on(EVENT_CHANGE, $.proxy(this.reload, this));
        },

        unbind: function() {
            this.$element.off(EVENT_CLICK).off(EVENT_CHANGE);
            this.$languages.off(EVENT_CHANGE);
        },

        click: function(e) {
            var $target = $(e.target);
            var $items = this.$items;
            var $item;
            var $btn;

            if (!$target.is('button')) {
                $btn = $target.closest('button');
            }
            // event target is a button
            if ($btn && $btn.length === 1) {
                $target = $btn;
            } else {
                // event target is a item
                $item = $target.closest('.i18n-list-item');
                if (!$item.hasClass('active highlight')) {
                    $target = $item;
                }
            }

            if (($target.data() && $target.data().url) || !$target.length) {
                return;
            }

            switch (String($target.data('toggle')).replace('.' + NAMESPACE, '')) {
                case 'bulk':
                    this.multiple = true;
                    $target
                        .addClass('hidden')
                        .siblings('button')
                        .removeClass('hidden');
                    $items
                        .addClass('active highlight')
                        .find('.bhojpur-js-translator')
                        .trigger(EVENT_INPUT);
                    break;

                case 'exit':
                    this.multiple = false;
                    $target.addClass('hidden');
                    $target
                        .siblings('button')
                        .addClass('hidden')
                        .filter('.bhojpur-js-bulk')
                        .removeClass('hidden');
                    $items.removeClass('active highlight');
                    break;

                case 'edit':
                    $items.removeClass('active highlight');
                    $target
                        .closest('li')
                        .addClass('active highlight')
                        .find('.bhojpur-js-translator')
                        .trigger(EVENT_INPUT);
                    break;

                case 'save':
                    $item = $target.closest('li');

                    this.submit($item.find('form'), function() {
                        $item.removeClass('active highlight');
                    });
                    break;

                case 'cancel':
                    $target.closest('li').removeClass('active highlight');
                    break;

                case 'copy':
                    $item = $target.closest('li');
                    $item
                        .find('.bhojpur-js-translator')
                        .val($item.find('.bhojpur-js-translation-source').text())
                        .trigger(EVENT_INPUT);
                    break;

                case 'copyall':
                    $items.find('.bhojpur-js-copy').click();
                    break;
            }
        },

        change: function(e) {
            var $target = $(e.target);

            if ($target.is('.bhojpur-js-translator')) {
                if (this.multiple) {
                    this.submit($target.closest('form'), function($form) {
                        var $help = $form.find('.bhojpur-js-help');

                        $help.addClass('in');

                        setTimeout(function() {
                            $help.removeClass('in');
                        }, 3000);
                    });
                }

                // Resize textarea height
                $target.trigger(EVENT_INPUT);
            }
        },

        reload: function(e) {
            var $target = $(e.target);
            var search = decodeSearch(location.search);

            search[$target.attr('name')] = $target.val();
            location.search = encodeSearch(search);
        },

        submit: function($form, callback) {
            window.onbeforeunload = null;
            $.fn.bhojpurSlideoutBeforeHide = null;

            if ($form.is('form')) {
                $.ajax(location.pathname, {
                    method: 'POST',
                    data: $form.serialize(),
                    success: function() {
                        $form.siblings('.bhojpur-js-translation-target').text($form.find('.bhojpur-js-translator').val());

                        if ($.isFunction(callback)) {
                            callback($form);
                        }
                    }
                });
            }
        },

        destroy: function() {
            this.unbind();
            this.$element.removeData(NAMESPACE);
        }
    };

    I18N.DEFAULTS = {};

    I18N.plugin = function(options) {
        return this.each(function() {
            var $this = $(this);
            var data = $this.data(NAMESPACE);
            var fn;

            if (!data) {
                $this.data(NAMESPACE, (data = new I18N(this, options)));
            }

            if (typeof options === 'string' && $.isFunction((fn = data[options]))) {
                fn.apply(data);
            }
        });
    };

    $(function() {
        I18N.plugin.call($('.bhojpur-i18n'));
    });
});