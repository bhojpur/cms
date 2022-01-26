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

    var $body = $("body");
    var NAMESPACE = 'bhojpur.widget.inlineEdit';
    var EVENT_ENABLE = 'enable.' + NAMESPACE;
    var EVENT_DISABLE = 'disable.' + NAMESPACE;
    var EVENT_CLICK = 'click.' + NAMESPACE;
    var EDIT_WIDGET_BUTTON = '.bhojpur-widget-button';
    var ID_WIDGET = 'bhojpur-widget-iframe';
    var INLINE_EDIT_URL;

    function BhojpurWidgetInlineEdit(element, options) {
        this.$element = $(element);
        this.options = $.extend({}, BhojpurWidgetInlineEdit.DEFAULTS, $.isPlainObject(options) && options);
        this.init();
    }

    BhojpurWidgetInlineEdit.prototype = {
        constructor: BhojpurWidgetInlineEdit,

        init: function() {
            this.bind();
            this.initStatus();
        },

        bind: function() {
            this.$element.on(EVENT_CLICK, EDIT_WIDGET_BUTTON, this.click);
            $(document).on('keyup', this.keyup);
        },

        initStatus: function() {
            var iframe = document.createElement("iframe");

            iframe.src = INLINE_EDIT_URL;
            iframe.id = ID_WIDGET;

            // show edit button after iframe totally loaded.
            if (iframe.attachEvent) {
                iframe.attachEvent("onload", function() {
                    $('.bhojpur-widget-button').show();
                });
            } else {
                iframe.onload = function() {
                    $('.bhojpur-widget-button').show();
                };
            }

            document.body.appendChild(iframe);
        },

        keyup: function(e) {
            var iframe = document.getElementById('bhojpur-widget-iframe');
            if (e.keyCode === 27) {
                iframe && iframe.contentDocument.querySelector('.bhojpur-slideout__close').click();
            }
        },

        click: function() {
            var $this = $(this);
            var iframe = document.getElementById('bhojpur-widget-iframe');
            var $iframe = iframe.contentWindow.$;
            var editLink = iframe.contentDocument.querySelector('.js-widget-edit-link');

            if (!editLink) {
                return;
            }

            iframe.classList.add('show');

            if ($iframe) {
                $iframe(".js-widget-edit-link").data("url", $this.data("url")).click();
            } else {
                editLink.setAttribute("data-url", $this.data("url"));
                editLink.click();
            }

            $body.addClass("open-widget-editor");

            return false;
        }
    };

    BhojpurWidgetInlineEdit.plugin = function(options) {
        return this.each(function() {
            var $this = $(this);
            var data = $this.data(NAMESPACE);
            var fn;

            if (!data) {
                if (/destroy/.test(options)) {
                    return;
                }
                $this.data(NAMESPACE, (data = new BhojpurWidgetInlineEdit(this, options)));
            }

            if (typeof options === 'string' && $.isFunction(fn = data[options])) {
                fn.apply(data);
            }
        });
    };


    $(function() {
        $body.attr("data-toggle", "bhojpur.widgets");

        $(".bhojpur-widget").each(function() {
            var $this = $(this);
            var $wrap = $this.children().eq(0);

            INLINE_EDIT_URL = $this.data("widget-inline-edit-url");

            if ($wrap.css("position") === "static") {
                $wrap.css("position", "relative");
            }

            $wrap.addClass("bhojpur-widget").unwrap();

            $wrap.append('<div class="bhojpur-widget-embed-wrapper"><button style="display: none;" data-url=\"' + $this.data("url") + '\" class="bhojpur-widget-button">Edit</button></div>');
        });

        var selector = '[data-toggle="bhojpur.widgets"]';

        $(document).
        on(EVENT_DISABLE, function(e) {
            BhojpurWidgetInlineEdit.plugin.call($(selector, e.target), 'destroy');
        }).
        on(EVENT_ENABLE, function(e) {
            BhojpurWidgetInlineEdit.plugin.call($(selector, e.target));
        }).
        triggerHandler(EVENT_ENABLE);
    });


    return BhojpurWidgetInlineEdit;
});