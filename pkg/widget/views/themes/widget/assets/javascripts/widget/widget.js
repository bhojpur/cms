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

    let $body = $('body'),
        NAMESPACE = 'bhojpur.widget',
        EVENT_ENABLE = 'enable.' + NAMESPACE,
        EVENT_DISABLE = 'disable.' + NAMESPACE,
        EVENT_CHANGE = 'change.' + NAMESPACE,
        EVENT_CLICK = 'click.' + NAMESPACE,
        TARGET_WIDGET = 'select[name="BhojpurResource.Widgets"]',
        TARGET_WIDGET_KIND = '[name="BhojpurResource.Kind"]',
        SELECT_FILTER = '[name="BhojpurResource.Widgets"],[name="BhojpurResource.Template"]',
        CLASS_IS_NEW = 'bhojpur-layout__widget-new',
        CLASS_FORM_SECTION = '.bhojpur-form-section',
        CLASS_LABEL = '.bhojpur-field__label',
        CLASS_LABEL_VISIBLE = '.bhojpur-field__label:visible',
        CLASS_FORM_SETTING = '.bhojpur-layout__widget-setting';

    function BhojpurWidget(element, options) {
        this.$element = $(element);
        this.options = $.extend({}, BhojpurWidget.DEFAULTS, $.isPlainObject(options) && options);
        this.init();
    }

    BhojpurWidget.prototype = {
        constructor: BhojpurWidget,

        init: function() {
            this.bind();
            this.isNewForm = this.$element.hasClass(CLASS_IS_NEW);
            if (!this.isNewForm && this.$element.find(TARGET_WIDGET).length) {
                this.addWidgetSlideout();
            }
            this.initSelect();
        },

        bind: function() {
            this.$element
                .on(EVENT_CHANGE, 'select', this.change.bind(this))
                .on(EVENT_CLICK, '.bhojpur-widget__new', this.getFormHtml.bind(this))
                .on(EVENT_CLICK, '.bhojpur-widget__cancel', this.cancelForm.bind(this));
        },

        unbind: function() {
            this.$element.off(EVENT_CHANGE, 'select', this.change.bind(this)).off(EVENT_CLICK, '.bhojpur-widget__new', this.getFormHtml.bind(this));
        },

        initSelect: function() {
            let $element = this.$element,
                $select = $element.find('select').filter(SELECT_FILTER),
                $kind = $(TARGET_WIDGET_KIND),
                hint = $element.data('hint'),
                NO_SETTINGS = `<h2 class="bhojpur-page__tips">${hint}</h2>`;

            $select.closest(CLASS_FORM_SECTION).hide();
            $select.each(function() {
                if (
                    $(this)
                        .find('option')
                        .filter('[value!=""]').length >= 2
                ) {
                    $(this)
                        .closest(CLASS_FORM_SECTION)
                        .show();
                }
            });

            if (!this.isNewForm && !$element.find('.bhojpur-bannereditor').length) {
                let $kindNext = $kind.parent().next('.bhojpur-form-section-rows'),
                    $kindParent = $kind.closest('.bhojpur-fieldset'),
                    $kindSection = $kind.closest('.bhojpur-form-section');

                // if settings dont have any field. will show NO SETTING HINT in settings container
                if (!$kindNext.children().length || !$kindParent.find(CLASS_LABEL_VISIBLE).length) {
                    $kindNext.append(NO_SETTINGS);

                    // if no other fields, just has empty settings fields, hide all elements and show NO SETTING HINT.
                    if (
                        !$element
                            .find(CLASS_LABEL)
                            .not($kindSection.find(CLASS_LABEL))
                            .is(':visible')
                    ) {
                        $kindSection.hide();
                        $element
                            .append(NO_SETTINGS)
                            .parent()
                            .find('.bhojpur-form__actions')
                            .remove();
                    }
                }
            }
        },

        addWidgetSlideout: function() {
            var $select = this.$element.find(TARGET_WIDGET),
                tabScopeActive = $body.data('tabScopeActive'),
                isInSlideout = $('.bhojpur-slideout').is(':visible'),
                $form = $select.closest('form'),
                actionUrl = $form.data('action-url') || $form.prop('action'),
                separator = actionUrl && actionUrl.indexOf('?') !== -1 ? '&' : '?',
                url,
                clickTmpl;

            $select.find('option').each(function() {
                var $this = $(this),
                    val = $this.val();

                if (val) {
                    url = `${actionUrl}${separator}widget_type=${val}`;

                    if (tabScopeActive) {
                        url = `${url}&widget_scope=${tabScopeActive}`;
                    }

                    if (isInSlideout) {
                        clickTmpl = `<a href=${url} style="display: none;" class="bhojpur-widget-${val}" data-open-type="slideout" data-url="${url}">${val}</a>`;
                    } else {
                        clickTmpl = `<a href=${url} style="display: none;" class="bhojpur-widget-${val}">${val}</a>`;
                    }

                    $select.after(clickTmpl);
                }
            });
        },

        change: function(e) {
            let $target = $(e.target),
                widgetValue = $target.val(),
                isInSlideout = $('.bhojpur-slideout').is(':visible'),
                clickClass = '.bhojpur-widget-' + widgetValue,
                $link = $(clickClass),
                url = $link.prop('href');

            if (!$target.is(TARGET_WIDGET)) {
                return;
            }

            $.fn.bhojpurSlideoutBeforeHide = null;
            window.onbeforeunload = null;

            if (!this.isNewForm) {
                if (isInSlideout) {
                    $link.trigger('click');
                } else {
                    location.href = url;
                }
            }

            return false;
        },

        getFormHtml: function(e) {
            let $target = $(e.target).closest('a'),
                widgetType = $target.data('widget-type'),
                $element = this.$element,
                url = $target.attr('href'),
                $title = $element.find('.bhojpur-layout__widget-selector'),
                $selector = $title.find('select'),
                $setting = $(CLASS_FORM_SETTING),
                $sectionSetting = $element.find('[data-section-title="Settings"]'),
                $loading = $(BhojpurWidget.TEMPLATE_LOADING);

            if ($sectionSetting.length) {
                $setting = $sectionSetting;
            }

            $loading.appendTo($setting).trigger('enable');

            $element.find('.bhojpur-slideout__lists-item a').hide();
            $element.find('.bhojpur-slideout__lists-groupname').hide();
            $element.find('.bhojpur-layout__widget-actions').show();

            $.get(url, function(html) {
                $title.find('.bhojpur-layout__widget-name').html($target.data('widget-name'));
                $title.show();
                $selector
                    .val(widgetType)
                    .closest('.bhojpur-form-section')
                    .hide();
                $setting.html(html).trigger('enable');
            }).fail(function() {
                window.alert('server error, please try again!');
            });

            return false;
        },

        cancelForm: function() {
            let $element = this.$element;

            if ($element.closest('.bhojpur-bottomsheets').length) {
                $element.closest('.bhojpur-bottomsheets').removeClass('bhojpur-bottomsheets__fullscreen');
            }
            $element.find('.bhojpur-slideout__lists-item a').show();
            $element.find('.bhojpur-slideout__lists-groupname').show();
            $element.find('.bhojpur-layout__widget-actions, .bhojpur-layout__widget-selector').hide();
            $element.find(CLASS_FORM_SETTING).html('');
        },

        destroy: function() {
            this.unbind();
            this.$element.removeData(NAMESPACE);
        }
    };

    BhojpurWidget.DEFAULTS = {};

    BhojpurWidget.TEMPLATE_LOADING =
        '<div style="text-align: center; margin-top: 30px;"><div class="mdl-spinner mdl-js-spinner is-active bhojpur-layout__bottomsheet-spinner"></div></div>';

    BhojpurWidget.plugin = function(options) {
        return this.each(function() {
            let $this = $(this),
                data = $this.data(NAMESPACE),
                fn;

            if (!data) {
                if (/destroy/.test(options)) {
                    return;
                }

                $this.data(NAMESPACE, (data = new BhojpurWidget(this, options)));
            }

            if (typeof options === 'string' && $.isFunction((fn = data[options]))) {
                fn.apply(data);
            }
        });
    };

    $(function() {
        let selector = '[data-toggle="bhojpur.widget"]';

        $(document)
            .on(EVENT_DISABLE, function(e) {
                BhojpurWidget.plugin.call($(selector, e.target), 'destroy');
            })
            .on(EVENT_ENABLE, function(e) {
                BhojpurWidget.plugin.call($(selector, e.target));
            })
            .triggerHandler(EVENT_ENABLE);

        if ($('.bhojpur-page__header .bhojpur-page-subnav__header').length) {
            $('.mdl-layout__content').addClass('has-subnav');
        }
    });

    return BhojpurWidget;
});