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
        $document = $(document),
        Mustache = window.Mustache,
        NAMESPACE = 'bhojpur.selectone',
        PARENT_NAMESPACE = 'bhojpur.bottomsheets',
        EVENT_CLICK = 'click.' + NAMESPACE,
        EVENT_ENABLE = 'enable.' + NAMESPACE,
        EVENT_DISABLE = 'disable.' + NAMESPACE,
        EVENT_RELOAD = 'reload.' + PARENT_NAMESPACE,
        CLASS_CLEAR_SELECT = '.bhojpur-selected__remove',
        CLASS_CHANGE_SELECT = '.bhojpur-selected__change',
        CLASS_SELECT_FIELD = '.bhojpur-field__selected',
        CLASS_SELECT_INPUT = '.bhojpur-field__selectone-input',
        CLASS_SELECT_TRIGGER = '.bhojpur-field__selectone-trigger',
        CLASS_PARENT = '.bhojpur-field__selectone',
        CLASS_SELECTED = 'is_selected',
        CLASS_ONE = 'bhojpur-bottomsheets__select-one';

    function BhojpurSelectOne(element, options) {
        this.$element = $(element);
        this.options = $.extend({}, BhojpurSelectOne.DEFAULTS, $.isPlainObject(options) && options);
        this.init();
    }

    BhojpurSelectOne.prototype = {
        constructor: BhojpurSelectOne,

        init: function() {
            this.bind();
        },

        bind: function() {
            $document.on(EVENT_RELOAD, `.${CLASS_ONE}`, this.reloadData.bind(this));
            this.$element
                .on(EVENT_CLICK, CLASS_CLEAR_SELECT, this.clearSelect.bind(this))
                .on(EVENT_CLICK, '[data-selectone-url]', this.openBottomSheets.bind(this))
                .on(EVENT_CLICK, CLASS_CHANGE_SELECT, this.changeSelect);
        },

        unbind: function() {
            $document.off(EVENT_CLICK, '[data-selectone-url]').off(EVENT_RELOAD, `.${CLASS_ONE}`);
            this.$element.off(EVENT_CLICK, CLASS_CLEAR_SELECT).off(EVENT_CLICK, CLASS_CHANGE_SELECT);
        },

        clearSelect: function(e) {
            var $target = $(e.target),
                $parent = $target.closest(CLASS_PARENT);

            $parent.find(CLASS_SELECT_FIELD).remove();
            $parent.find(CLASS_SELECT_INPUT).html('');
            $parent.find(CLASS_SELECT_INPUT)[0].value = '';
            $parent.find(CLASS_SELECT_TRIGGER).show();

            $parent.trigger('bhojpur.selectone.unselected');
            return false;
        },

        changeSelect: function() {
            var $target = $(this),
                $parent = $target.closest(CLASS_PARENT);

            $parent.find(CLASS_SELECT_TRIGGER).trigger('click');
        },

        openBottomSheets: function(e) {
            var $this = $(e.target),
                data = $this.data();

            this.BottomSheets = $body.data('bhojpur.bottomsheets');
            this.$parent = $this.closest(CLASS_PARENT);

            data.url = data.selectoneUrl;

            this.SELECT_ONE_SELECTED_ICON = $('[name="select-one-selected-icon"]').html();
            this.BottomSheets.open(data, this.handleSelectOne.bind(this));
        },

        initItem: function() {
            var $selectFeild = this.$parent.find(CLASS_SELECT_FIELD),
                selectedID;

            if (!$selectFeild.length) {
                return;
            }

            selectedID = $selectFeild.data().primaryKey;

            if (selectedID) {
                this.$bottomsheets
                    .find('tr[data-primary-key="' + selectedID + '"]')
                    .addClass(CLASS_SELECTED)
                    .find('td:first')
                    .append(this.SELECT_ONE_SELECTED_ICON);
            }
        },

        reloadData: function() {
            this.initItem();
        },

        renderSelectOne: function(data) {
            return Mustache.render($('[name="select-one-selected-template"]').html(), data);
        },

        handleSelectOne: function($bottomsheets) {
            var options = {
                onSelect: this.onSelectResults.bind(this), //render selected item after click item lists
                onSubmit: this.onSubmitResults.bind(this) //render new items after new item form submitted
            };

            $bottomsheets.bhojpurSelectCore(options).addClass(CLASS_ONE);
            this.$bottomsheets = $bottomsheets;
            this.initItem();
        },

        onSelectResults: function(data) {
            this.handleResults(data);
        },

        onSubmitResults: function(data) {
            this.handleResults(data, true);
        },

        handleResults: function(data) {
            var template,
                $parent = this.$parent,
                $select = $parent.find('select'),
                $selectFeild = $parent.find(CLASS_SELECT_FIELD);

            data.displayName = data.Text || data.Name || data.Title || data.Code || data[Object.keys(data)[0]];
            data.selectoneValue = data.primaryKey || data.ID;

            data.displayName = (data.displayName).escapeSymbol();

            if (!$select.length) {
                return;
            }

            template = this.renderSelectOne(data);

            if ($selectFeild.length) {
                $selectFeild.remove();
            }

            $parent.prepend(template);
            $parent.find(CLASS_SELECT_TRIGGER).hide();

            $select.html(Mustache.render(BhojpurSelectOne.SELECT_ONE_OPTION_TEMPLATE, data));
            $select[0].value = data.primaryKey || data.ID;

            $parent.trigger('bhojpur.selectone.selected', [data]);

            this.$bottomsheets.bhojpurSelectCore('destroy').remove();
            if (!$('.bhojpur-bottomsheets').is(':visible')) {
                $('body').removeClass('bhojpur-bottomsheets-open');
            }
        },

        destroy: function() {
            this.unbind();
            this.$element.removeData(NAMESPACE);
        }
    };

    BhojpurSelectOne.SELECT_ONE_OPTION_TEMPLATE = '<option value="[[ selectoneValue ]]" selected>[[ displayName ]]</option>';

    BhojpurSelectOne.plugin = function(options) {
        return this.each(function() {
            var $this = $(this);
            var data = $this.data(NAMESPACE);
            var fn;

            if (!data) {
                if (/destroy/.test(options)) {
                    return;
                }

                $this.data(NAMESPACE, (data = new BhojpurSelectOne(this, options)));
            }

            if (typeof options === 'string' && $.isFunction((fn = data[options]))) {
                fn.apply(data);
            }
        });
    };

    $(function() {
        var selector = '[data-toggle="bhojpur.selectone"]';
        $(document)
            .on(EVENT_DISABLE, function(e) {
                BhojpurSelectOne.plugin.call($(selector, e.target), 'destroy');
            })
            .on(EVENT_ENABLE, function(e) {
                BhojpurSelectOne.plugin.call($(selector, e.target));
            })
            .triggerHandler(EVENT_ENABLE);
    });

    return BhojpurSelectOne;
});