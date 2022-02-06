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
  
    var _ = window._;
    var $body = $('body');
    var NAMESPACE = 'bhojpur.tabbar';
    var EVENT_ENABLE = 'enable.' + NAMESPACE;
    var EVENT_DISABLE = 'disable.' + NAMESPACE;
    var EVENT_CLICK = 'click.' + NAMESPACE;
    var CLASS_TAB = '.bhojpur-layout__tab-button';
    var CLASS_TAB_CONTENT = '.bhojpur-layout__tab-content';
    var CLASS_TAB_BAR = '.mdl-layout__tab-bar-container';
    var CLASS_TAB_BAR_RIGHT = '.bhojpur-layout__tab-right';
    var CLASS_TAB_BAR_LEFT = '.bhojpur-layout__tab-left';
    var CLASS_ACTIVE = 'is-active';
  
    function BhojpurTab(element, options) {
      this.$element = $(element);
      this.options = $.extend({}, BhojpurTab.DEFAULTS, $.isPlainObject(options) && options);
      this.init();
    }
  
    BhojpurTab.prototype = {
      constructor: BhojpurTab,
  
      init: function () {
        this.initTab();
        this.bind();
      },
  
      bind: function () {
        this.$element.on(EVENT_CLICK, CLASS_TAB, this.switchTab.bind(this));
        this.$element.on(EVENT_CLICK, CLASS_TAB_BAR_RIGHT, this.scrollTabRight.bind(this));
        this.$element.on(EVENT_CLICK, CLASS_TAB_BAR_LEFT, this.scrollTabLeft.bind(this));
      },
  
      unbind: function () {
        this.$element.off(EVENT_CLICK, CLASS_TAB, this.switchTab);
        this.$element.off(EVENT_CLICK, CLASS_TAB_BAR_RIGHT, this.scrollTabRight);
        this.$element.off(EVENT_CLICK, CLASS_TAB_BAR_LEFT, this.scrollTabLeft);
      },
  
      initTab: function () {
        var data = this.$element.data();
  
        if (!data.scopeActive) {
          $(CLASS_TAB).first().addClass(CLASS_ACTIVE);
          $body.data('tabScopeActive',$(CLASS_TAB).first().data('name'));
        } else {
          $body.data('tabScopeActive',data.scopeActive);
        }
  
        this.tabWidth = 0;
        this.slideoutWidth = $(CLASS_TAB_CONTENT).outerWidth();
  
        _.each($(CLASS_TAB), function(ele) {
          this.tabWidth = this.tabWidth + $(ele).outerWidth();
        }.bind(this));
  
        if (this.tabWidth > this.slideoutWidth) {
          this.$element.find(CLASS_TAB_BAR).append(BhojpurTab.ARROW_RIGHT);
        }
  
      },
  
      scrollTabLeft: function (e) {
        e.stopPropagation();
  
        var $scrollBar = $(CLASS_TAB_BAR),
            scrollLeft = $scrollBar.scrollLeft(),
            jumpDistance = scrollLeft - this.slideoutWidth;
  
        if (scrollLeft > 0){
          $scrollBar.animate({scrollLeft:jumpDistance}, 400, function () {
  
            $(CLASS_TAB_BAR_RIGHT).show();
            if ($scrollBar.scrollLeft() == 0) {
              $(CLASS_TAB_BAR_LEFT).hide();
            }
  
          });
        }
      },
  
      scrollTabRight: function (e) {
        e.stopPropagation();
  
        var $scrollBar = $(CLASS_TAB_BAR),
            scrollLeft = $scrollBar.scrollLeft(),
            tabWidth = this.tabWidth,
            slideoutWidth = this.slideoutWidth,
            jumpDistance = scrollLeft + slideoutWidth;
  
        if (jumpDistance < tabWidth){
          $scrollBar.animate({scrollLeft:jumpDistance}, 400, function () {
  
            $(CLASS_TAB_BAR_LEFT).show();
            if ($scrollBar.scrollLeft() + slideoutWidth >= tabWidth) {
              $(CLASS_TAB_BAR_RIGHT).hide();
            }
  
          });
  
          !$(CLASS_TAB_BAR_LEFT).length && this.$element.find(CLASS_TAB_BAR).prepend(BhojpurTab.ARROW_LEFT);
        }
      },
  
      switchTab: function (e) {
        var $target = $(e.target),
            $element = this.$element,
            data = $target.data(),
            tabScopeActive = $body.data().tabScopeActive,
            isInSlideout = $('.bhojpur-slideout').is(':visible');
  
        if (!isInSlideout) {
          return;
        }
  
        if ($target.hasClass(CLASS_ACTIVE)){
          return false;
        }
  
        $element.find(CLASS_TAB).removeClass(CLASS_ACTIVE);
        $target.addClass(CLASS_ACTIVE);
  
        $.ajax(data.tabUrl, {
            method: 'GET',
            dataType: 'html',
            processData: false,
            contentType: false,
            beforeSend: function () {
              $('.bhojpur-layout__tab-spinner').remove();
              var $spinner = '<div class="mdl-spinner mdl-js-spinner is-active bhojpur-layout__tab-spinner"></div>';
              $(CLASS_TAB_CONTENT).hide().before($spinner);
              window.componentHandler.upgradeElement($('.bhojpur-layout__tab-spinner')[0]);
            },
            success: function (html) {
              $('.bhojpur-layout__tab-spinner').remove();
              $body.data('tabScopeActive',$target.data('name'));
              var $content = $(html).find(CLASS_TAB_CONTENT).html();
              $(CLASS_TAB_CONTENT).show().html($content).trigger('enable');
  
            },
            error: function () {
              $('.bhojpur-layout__tab-spinner').remove();
              $body.data('tabScopeActive',tabScopeActive);
            }
          });
        return false;
      },
  
      destroy: function () {
        this.unbind();
        $body.removeData('tabScopeActive');
      }
    };
  
    BhojpurTab.ARROW_RIGHT = '<a href="javascript://" class="bhojpur-layout__tab-right"></a>';
    BhojpurTab.ARROW_LEFT = '<a href="javascript://" class="bhojpur-layout__tab-left"></a>';
  
    BhojpurTab.DEFAULTS = {};
  
    BhojpurTab.plugin = function (options) {
      return this.each(function () {
        var $this = $(this);
        var data = $this.data(NAMESPACE);
        var fn;
  
        if (!data) {
          if (/destroy/.test(options)) {
            return;
          }
  
          $this.data(NAMESPACE, (data = new BhojpurTab(this, options)));
        }
  
        if (typeof options === 'string' && $.isFunction(fn = data[options])) {
          fn.apply(data);
        }
      });
    };
  
    $(function () {
      var selector = '[data-toggle="bhojpur.tab"]';
  
      $(document)
        .on(EVENT_DISABLE, function (e) {
          BhojpurTab.plugin.call($(selector, e.target), 'destroy');
        })
        .on(EVENT_ENABLE, function (e) {
          BhojpurTab.plugin.call($(selector, e.target));
        })
        .triggerHandler(EVENT_ENABLE);
    });
  
    return BhojpurTab;
  
  });