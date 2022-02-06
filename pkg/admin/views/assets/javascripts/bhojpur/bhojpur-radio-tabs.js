(function (factory) {
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
  })(function ($) {
  
    'use strict';
  
    var NAMESPACE = 'bhojpur.tabbar.radio';
    var EVENT_ENABLE = 'enable.' + NAMESPACE;
    var EVENT_DISABLE = 'disable.' + NAMESPACE;
    var EVENT_CLICK = 'click.' + NAMESPACE;
    var EVENT_SWITCHED = 'switched.' + NAMESPACE;
    var CLASS_TAB = '[data-tab-target]';
    var CLASS_TAB_SOURCE = '[data-tab-source]';
    var CLASS_ACTIVE = 'is-active';
  
    function BhojpurTabRadio(element, options) {
      this.$element = $(element);
      this.options = $.extend({}, BhojpurTabRadio.DEFAULTS, $.isPlainObject(options) && options);
      this.init();
    }
  
    BhojpurTabRadio.prototype = {
      constructor: BhojpurTabRadio,
  
      init: function () {
        this.bind();
      },
  
      bind: function () {
        this.$element.on(EVENT_CLICK, CLASS_TAB, this.switchTab.bind(this));
      },
  
      unbind: function () {
        this.$element.off(EVENT_CLICK, CLASS_TAB, this.switchTab);
      },
  
  
      switchTab: function (e) {
        var $target = $(e.target),
            $element = this.$element,
            $tabs = $element.find(CLASS_TAB),
            $tabSources = $element.find(CLASS_TAB_SOURCE),
            data = $target.data(),
            tabTarget = data.tabTarget;
  
        if ($target.hasClass(CLASS_ACTIVE)){
          return;
        }
  
        $tabs.removeClass(CLASS_ACTIVE);
        $target.addClass(CLASS_ACTIVE);
  
        $tabSources.hide().filter('[data-tab-source="' + tabTarget + '"]').show();
        $element.trigger(EVENT_SWITCHED, [$element, tabTarget]);
  
      },
  
      destroy: function () {
        this.unbind();
      }
    };
  
    BhojpurTabRadio.DEFAULTS = {};
  
    BhojpurTabRadio.plugin = function (options) {
      return this.each(function () {
        var $this = $(this);
        var data = $this.data(NAMESPACE);
        var fn;
  
        if (!data) {
          if (/destroy/.test(options)) {
            return;
          }
  
          $this.data(NAMESPACE, (data = new BhojpurTabRadio(this, options)));
        }
  
        if (typeof options === 'string' && $.isFunction(fn = data[options])) {
          fn.apply(data);
        }
      });
    };
  
    $(function () {
      var selector = '[data-toggle="bhojpur.tab.radio"]';
  
      $(document)
        .on(EVENT_DISABLE, function (e) {
          BhojpurTabRadio.plugin.call($(selector, e.target), 'destroy');
        })
        .on(EVENT_ENABLE, function (e) {
          BhojpurTabRadio.plugin.call($(selector, e.target));
        })
        .triggerHandler(EVENT_ENABLE);
    });
  
    return BhojpurTabRadio;
  
  });