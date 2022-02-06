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

    let $document = $(document),
        NAMESPACE = 'bhojpur.modal',
        EVENT_ENABLE = 'enable.' + NAMESPACE,
        EVENT_DISABLE = 'disable.' + NAMESPACE,
        EVENT_CLICK = 'click.' + NAMESPACE,
        EVENT_KEYUP = 'keyup.' + NAMESPACE,
        EVENT_SHOW = 'show.' + NAMESPACE,
        EVENT_SHOWN = 'shown.' + NAMESPACE,
        EVENT_HIDE = 'hide.' + NAMESPACE,
        EVENT_HIDDEN = 'hidden.' + NAMESPACE,
        EVENT_TRANSITION_END = 'transitionend',
        CLASS_OPEN = 'bhojpur-modal-open',
        CLASS_SHOWN = 'shown',
        CLASS_FADE = 'fade',
        CLASS_IN = 'in',
        ARIA_HIDDEN = 'aria-hidden';

    function BhojpurModal(element, options) {
        this.$element = $(element);
        this.options = $.extend({}, BhojpurModal.DEFAULTS, $.isPlainObject(options) && options);
        this.transitioning = false;
        this.fadable = false;
        this.init();
    }

    BhojpurModal.prototype = {
        constructor: BhojpurModal,

        init: function() {
            this.fadable = this.$element.hasClass(CLASS_FADE);

            if (this.options.show) {
                this.show();
            } else {
                this.toggle();
            }
        },

        bind: function() {
            this.$element.on(EVENT_CLICK, $.proxy(this.click, this));

            if (this.options.keyboard) {
                $document.on(EVENT_KEYUP, $.proxy(this.keyup, this));
            }
        },

        unbind: function() {
            this.$element.off(EVENT_CLICK, this.click);

            if (this.options.keyboard) {
                $document.off(EVENT_KEYUP, this.keyup);
            }
        },

        click: function(e) {
            var element = this.$element[0];
            var target = e.target;

            if (target === element && this.options.backdrop) {
                this.hide();
                return;
            }

            while (target !== element) {
                if ($(target).data('dismiss') === 'modal') {
                    this.hide();
                    break;
                }

                target = target.parentNode;
            }
        },

        keyup: function(e) {
            if (e.which === 27) {
                this.hide();
            }
        },

        show: function(noTransition) {
            var $this = this.$element,
                showEvent;

            if (this.transitioning || $this.hasClass(CLASS_IN)) {
                return;
            }

            showEvent = $.Event(EVENT_SHOW);
            $this.trigger(showEvent);

            if (showEvent.isDefaultPrevented()) {
                return;
            }

            $document.find('body').addClass(CLASS_OPEN);

            /*jshint expr:true */
            $this
                .addClass(CLASS_SHOWN)
                .scrollTop(0)
                .get(0).offsetHeight; // reflow for transition
            this.transitioning = true;

            if (noTransition || !this.fadable) {
                $this.addClass(CLASS_IN);
                this.shown();
                return;
            }

            $this.one(EVENT_TRANSITION_END, $.proxy(this.shown, this));
            $this.addClass(CLASS_IN);
        },

        shown: function() {
            this.transitioning = false;
            this.bind();
            this.$element
                .attr(ARIA_HIDDEN, false)
                .trigger(EVENT_SHOWN)
                .focus();
        },

        hide: function(noTransition) {
            var $this = this.$element,
                hideEvent;

            if (this.transitioning || !$this.hasClass(CLASS_IN)) {
                return;
            }

            hideEvent = $.Event(EVENT_HIDE);
            $this.trigger(hideEvent);

            if (hideEvent.isDefaultPrevented()) {
                return;
            }

            $document.find('body').removeClass(CLASS_OPEN);
            this.transitioning = true;

            if (noTransition || !this.fadable) {
                $this.removeClass(CLASS_IN);
                this.hidden();
                return;
            }

            $this.one(EVENT_TRANSITION_END, $.proxy(this.hidden, this));
            $this.removeClass(CLASS_IN);
        },

        hidden: function() {
            this.transitioning = false;
            this.unbind();
            this.$element
                .removeClass(CLASS_SHOWN)
                .attr(ARIA_HIDDEN, true)
                .trigger(EVENT_HIDDEN);
        },

        toggle: function() {
            if (this.$element.hasClass(CLASS_IN)) {
                this.hide();
            } else {
                this.show();
            }
        },

        destroy: function() {
            this.$element.removeData(NAMESPACE);
        }
    };

    BhojpurModal.DEFAULTS = {
        backdrop: false,
        keyboard: true,
        show: true
    };

    BhojpurModal.plugin = function(options) {
        return this.each(function() {
            var $this = $(this);
            var data = $this.data(NAMESPACE);
            var fn;

            if (!data) {
                if (/destroy/.test(options)) {
                    return;
                }

                $this.data(NAMESPACE, (data = new BhojpurModal(this, options)));
            }

            if (typeof options === 'string' && $.isFunction((fn = data[options]))) {
                fn.apply(data);
            }
        });
    };

    $.fn.bhojpurModal = BhojpurModal.plugin;

    $(function() {
        var selector = '.bhojpur-modal';

        $(document)
            .on(EVENT_CLICK, '[data-toggle="bhojpur.modal"]', function() {
                var $this = $(this);
                var data = $this.data();
                var $target = $(data.target || $this.attr('href'));

                BhojpurModal.plugin.call($target, $target.data(NAMESPACE) ? 'toggle' : data);
            })
            .on(EVENT_DISABLE, function(e) {
                BhojpurModal.plugin.call($(selector, e.target), 'destroy');
            })
            .on(EVENT_ENABLE, function(e) {
                BhojpurModal.plugin.call($(selector, e.target));
            });
    });

    return BhojpurModal;
});