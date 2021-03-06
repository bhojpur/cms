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
  
    var NAMESPACE = 'bhojpur.autoheight';
    var EVENT_ENABLE = 'enable.' + NAMESPACE;
    var EVENT_DISABLE = 'disable.' + NAMESPACE;
    var EVENT_INPUT = 'input';
  
    function BhojpurAutoheight(element, options) {
      this.$element = $(element);
      this.options = $.extend({}, BhojpurAutoheight.DEFAULTS, $.isPlainObject(options) && options);
      this.init();
    }
  
    BhojpurAutoheight.prototype = {
      constructor: BhojpurAutoheight,
  
      init: function () {
        var $this = this.$element;
  
        this.paddingTop = parseInt($this.css('padding-top'), 10);
        this.paddingBottom = parseInt($this.css('padding-bottom'), 10);
        this.resize();
        this.bind();
      },
  
      bind: function () {
        this.$element.on(EVENT_INPUT, $.proxy(this.resize, this));
      },
  
      unbind: function () {
        this.$element.off(EVENT_INPUT, this.resize);
      },
  
      resize: function () {
        var $this = this.$element;
        var scrollHeight = $this.prop('scrollHeight');
  
        if(scrollHeight){
          $this.height('auto').height(scrollHeight - this.paddingTop - this.paddingBottom);
        } else {
          $this.height('40px');
        }
      },
  
      destroy: function () {
        this.unbind();
        this.$element.removeData(NAMESPACE);
      }
    };
  
    BhojpurAutoheight.DEFAULTS = {};
  
    BhojpurAutoheight.plugin = function (options) {
      return this.each(function () {
        var $this = $(this);
        var data = $this.data(NAMESPACE);
        var fn;
  
        if (!data) {
          if (/destroy/.test(options)) {
            return;
          }
  
          $this.data(NAMESPACE, (data = new BhojpurAutoheight(this, options)));
        }
  
        if (typeof options === 'string' && $.isFunction(fn = data[options])) {
          fn.apply(data);
        }
      });
    };
  
    $(function () {
      var selector = 'textarea.bhojpur-js-autoheight';
  
      $(document).
        on(EVENT_DISABLE, function (e) {
          BhojpurAutoheight.plugin.call($(selector, e.target), 'destroy');
        }).
        on(EVENT_ENABLE, function (e) {
          BhojpurAutoheight.plugin.call($(selector, e.target));
        }).
        triggerHandler(EVENT_ENABLE);
    });
  
    return BhojpurAutoheight;
  
  });