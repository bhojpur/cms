(function(factory) {
    if (typeof define === "function" && define.amd) {
      // AMD. Register as anonymous module.
      define(["jquery"], factory);
    } else if (typeof exports === "object") {
      // Node / CommonJS
      factory(require("jquery"));
    } else {
      // Browser globals.
      factory(jQuery);
    }
  })(function($) {
    "use strict";
  
    let $document = $(document),
      FormData = window.FormData,
      BHOJPUR_Translations = window.BHOJPUR_Translations,
      _ = window._,
      BHOJPUR = window.BHOJPUR,
      NAMESPACE = "bhojpur.slideout",
      EVENT_KEYUP = "keyup." + NAMESPACE,
      EVENT_CLICK = "click." + NAMESPACE,
      EVENT_SUBMIT = "submit." + NAMESPACE,
      EVENT_SHOW = "show." + NAMESPACE,
      EVENT_SLIDEOUT_SUBMIT_COMPLEMENT = "slideoutSubmitComplete." + NAMESPACE,
      EVENT_SLIDEOUT_CLOSED = "slideoutClosed." + NAMESPACE,
      EVENT_SLIDEOUT_LOADED = "slideoutLoaded." + NAMESPACE,
      EVENT_SLIDEOUT_BEFORESEND = "slideoutBeforeSend." + NAMESPACE,
      EVENT_SHOWN = "shown." + NAMESPACE,
      EVENT_HIDE = "hide." + NAMESPACE,
      EVENT_HIDDEN = "hidden." + NAMESPACE,
      EVENT_TRANSITIONEND = "transitionend",
      CLASS_OPEN = "bhojpur-slideout-open",
      CLASS_MINI = "bhojpur-slideout-mini",
      CLASS_IS_SHOWN = "is-shown",
      CLASS_IS_SLIDED = "is-slided",
      CLASS_IS_SELECTED = "is-selected",
      CLASS_MAIN_CONTENT = ".mdl-layout__content.bhojpur-page",
      CLASS_HEADER_LOCALE = ".bhojpur-actions__locale",
      CLASS_BODY_LOADING = ".bhojpur-body__loading";
  
    function replaceHtml(el, html) {
      let oldEl = typeof el === "string" ? document.getElementById(el) : el,
        newEl = oldEl.cloneNode(false);
      newEl.innerHTML = html;
      oldEl.parentNode.replaceChild(newEl, oldEl);
      return newEl;
    }
  
    function pushArrary($ele, isScript) {
      let array = [],
        prop = "href";
  
      isScript && (prop = "src");
      $ele.each(function() {
        array.push($(this).attr(prop));
      });
      return _.uniq(array);
    }
  
    function execSlideoutEvents(url, response) {
      // exec bhojpurSliderAfterShow after script loaded
      var bhojpurSliderAfterShow = $.fn.bhojpurSliderAfterShow;
      for (var name in bhojpurSliderAfterShow) {
        if (
          bhojpurSliderAfterShow.hasOwn(name) &&
          !bhojpurSliderAfterShow[name]["isLoaded"]
        ) {
          bhojpurSliderAfterShow[name]["isLoaded"] = true;
          bhojpurSliderAfterShow[name].call(this, url, response);
        }
      }
    }
  
    function loadScripts(srcs, data, callback) {
      let scriptsLoaded = 0;
  
      for (let i = 0, len = srcs.length; i < len; i++) {
        let script = document.createElement("script");
  
        script.onload = function() {
          scriptsLoaded++;
  
          if (scriptsLoaded === srcs.length) {
            if ($.isFunction(callback)) {
              callback();
            }
          }
  
          if (data && data.url && data.response) {
            execSlideoutEvents(data.url, data.response);
          }
        };
  
        script.src = srcs[i];
        document.body.appendChild(script);
      }
    }
  
    function loadStyles(srcs) {
      let ss = document.createElement("link"),
        src = srcs.shift();
  
      ss.type = "text/css";
      ss.rel = "stylesheet";
      ss.onload = function() {
        if (srcs.length) {
          loadStyles(srcs);
        }
      };
      ss.href = src;
      document.getElementsByTagName("head")[0].appendChild(ss);
    }
  
    function compareScripts($scripts) {
      let $currentPageScripts = $("script"),
        slideoutScripts = pushArrary($scripts, true),
        currentPageScripts = pushArrary($currentPageScripts, true),
        scriptDiff = _.difference(slideoutScripts, currentPageScripts);
      return scriptDiff;
    }
  
    function compareLinks($links) {
      let $currentStyles = $("link"),
        slideoutStyles = pushArrary($links),
        currentStyles = pushArrary($currentStyles),
        styleDiff = _.difference(slideoutStyles, currentStyles);
  
      return styleDiff;
    }
  
    function BhojpurSlideout(element, options) {
      this.$element = $(element);
      this.options = $.extend(
        {},
        BhojpurSlideout.DEFAULTS,
        $.isPlainObject(options) && options
      );
      this.slided = false;
      this.disabled = false;
      this.slideoutType = false;
      this.init();
    }
  
    BhojpurSlideout.prototype = {
      constructor: BhojpurSlideout,
  
      init: function() {
        this.build();
        this.bind();
      },
  
      build: function() {
        var $slideout;
  
        this.$slideout = $slideout = $(BhojpurSlideout.TEMPLATE).appendTo("body");
        this.$slideoutTemplate = $slideout.html();
      },
  
      unbuild: function() {
        this.$slideout.remove();
      },
  
      bind: function() {
        this.$slideout
          .on(EVENT_SUBMIT, "form", this.submit.bind(this))
          .on(
            EVENT_CLICK,
            ".bhojpur-slideout__fullscreen",
            this.toggleSlideoutMode.bind(this)
          )
          .on(EVENT_CLICK, '[data-dismiss="slideout"]', this.hide.bind(this));
  
        $document.on(EVENT_KEYUP, $.proxy(this.keyup, this));
      },
  
      unbind: function() {
        this.$slideout.off(EVENT_SUBMIT, this.submit).off(EVENT_CLICK);
  
        $document.off(EVENT_KEYUP, this.keyup);
      },
  
      keyup: function(e) {
        if (e.which === 27) {
          if (
            $(".bhojpur-bottomsheets").is(":visible") ||
            $(".bhojpur-modal").is(":visible") ||
            $("#redactor-modal-box").length ||
            $("#dialog").is(":visible")
          ) {
            return;
          }
  
          this.hide();
          this.removeSelectedClass();
        }
      },
  
      loadExtraResource: function(data) {
        let styleDiff = compareLinks(data.$links),
          scriptDiff = compareScripts(data.$scripts);
  
        if (styleDiff.length) {
          loadStyles(styleDiff);
        }
  
        if (scriptDiff.length) {
          loadScripts(scriptDiff, data);
        }
      },
  
      removeSelectedClass: function() {
        this.$element.find("[data-url]").removeClass(CLASS_IS_SELECTED);
      },
  
      addLoading: function() {
        $(CLASS_BODY_LOADING).remove();
        var $loading = $(BhojpurSlideout.TEMPLATE_LOADING);
        $loading.appendTo($("body")).trigger("enable.bhojpur.material");
      },
  
      toggleSlideoutMode: function() {
        this.$slideout
          .toggleClass("bhojpur-slideout__fullscreen")
          .find(".bhojpur-slideout__fullscreen i")
          .toggle();
      },
  
      checkRichedutorHTMLTags: function(source){    
      var DOMHolderArray = new Array();
      var tagsArray = new Array();
      var lines = source.value.split('\n');
      for (var x = 0; x < lines.length; x++) {
          tagsArray = lines[x].match(/<(\/{1})?\w+((\s+\w+(\s*=\s*(?:".*?"|'.*?'|[^'">\s]+))?)+\s*|\s*)>/g);
          if (tagsArray) {
              for (var i = 0; i < tagsArray.length; i++) {
                  if (tagsArray[i].indexOf('</') >= 0) {
                      let elementToPop = tagsArray[i].substr(2, tagsArray[i].length - 3);
                      elementToPop = elementToPop.replace(/ /g, '');
                      for (var j = DOMHolderArray.length - 1; j >= 0; j--) {
                          if (DOMHolderArray[j].element == elementToPop) {
                              DOMHolderArray.splice(j, 1);
                              if (elementToPop != 'html') {
                                  break;
                              }
                          }
                      }
                  } else {
                      var tag = new Object();
                      tag.full = tagsArray[i];
                      tag.line = x + 1;
                      if (tag.full.indexOf(' ') > 0) {
                          tag.element = tag.full.substr(1, tag.full.indexOf(' ') - 1);
                      } else {
                          tag.element = tag.full.substr(1, tag.full.length - 2);
                      }
                      var selfClosingTags = new Array('area', 'base', 'br', 'col', 'command', 'embed', 'hr', 'img', 'input', 'keygen', 'link', 'meta', 'param', 'source', 'track', 'wbr');
                      var isSelfClosing = false;
                      for (var y = 0; y < selfClosingTags.length; y++) {
                          if (selfClosingTags[y].localeCompare(tag.element) == 0) {
                              isSelfClosing = true;
                          }
                      }
                      if (isSelfClosing == false) {
                          DOMHolderArray.push(tag);
                      }
                  }
  
              }
          }
        }
  
        return DOMHolderArray.length;
      
      },
  
      submit: function(e) {
        let $slideout = this.$slideout,
          form = e.target,
          $form = $(form),
          _this = this,
          $loading = $(BHOJPUR.$formLoading),
          $submit = $form.find(":submit"),
          hasNotClosedTags = false;
  
        if ($form.data("normal-submit")) {
          return;
        }
  
        $slideout.trigger(EVENT_SLIDEOUT_BEFORESEND);
  
        if (!FormData) {
          return;
        }
        e.preventDefault();
  
        document.querySelectorAll('.bhojpur-redactor-box .redactor-source').forEach(function(item) {
          if(_this.checkRichedutorHTMLTags(item)){
            hasNotClosedTags=true;
          }
        });
  
        if(hasNotClosedTags){
          BHOJPUR.bhojpurConfirm(BHOJPUR_Translations.slideoutCheckHTMLTagsError);
          return false;
        }
  
        this.submitXHR = $.ajax($form.prop("action"), {
          method: $form.prop("method"),
          data: new FormData(form),
          dataType: "html",
          processData: false,
          contentType: false,
          beforeSend: function() {
            $(".bhojpur-submit-loading").remove();
            $loading
              .appendTo(
                $submit.prop("disabled", true).closest(".bhojpur-form__actions")
              )
              .trigger("enable.bhojpur.material");
            $.fn.bhojpurSlideoutBeforeHide = null;
          },
          success: function() {
            let returnUrl = $form.data("returnUrl"),
              refreshUrl = $form.data("refreshUrl");
  
            $slideout.trigger(EVENT_SLIDEOUT_SUBMIT_COMPLEMENT);
  
            if (refreshUrl) {
              window.location.href = refreshUrl;
              return;
            }
  
            if (returnUrl == "refresh") {
              _this.refresh();
              return;
            }
  
            if (returnUrl && returnUrl != "refresh") {
              _this.load(returnUrl);
            } else {
              _this.refresh();
            }
          },
          error: function(err) {
            BHOJPUR.handleAjaxError(err);
          },
          complete: function() {
            $submit.prop("disabled", false);
          }
        });
      },
  
      load: function(url, data) {
        var options = this.options;
        var method;
        var dataType;
        var load;
        var $slideout = this.$slideout;
        var $title;
  
        if (!url) {
          return;
        }
  
        data = $.isPlainObject(data) ? data : {};
  
        method = data.method ? data.method : "GET";
        dataType = data.datatype ? data.datatype : "html";
  
        load = $.proxy(function() {
          $.ajax(url, {
            method: method,
            dataType: dataType,
            cache: true,
            ifModified: true,
            success: $.proxy(function(response) {
              let $response,
                $content,
                $bhojpurFormContainer,
                $scripts,
                $links,
                bodyClass;
  
              $(CLASS_BODY_LOADING).remove();
  
              if (method === "GET") {
                $response = $(response);
                $content = $response.find(CLASS_MAIN_CONTENT);
                $bhojpurFormContainer = $content.find(".bhojpur-form-container");
                this.slideoutType =
                  $bhojpurFormContainer.length &&
                  $bhojpurFormContainer.data().slideoutType;
  
                if (!$content.length) {
                  return;
                }
  
                let bodyHtml = response.match(
                  /<\s*body.*>[\s\S]*<\s*\/body\s*>/gi
                );
                if (bodyHtml) {
                  bodyHtml = bodyHtml
                    .join("")
                    .replace(/<\s*body/gi, "<div")
                    .replace(/<\s*\/body/gi, "</div");
                  bodyClass = $(bodyHtml).prop("class");
                  $("body").addClass(bodyClass);
  
                  let data = {
                    $scripts: $response.filter("script"),
                    $links: $response.filter("link"),
                    url: url,
                    response: response
                  };
  
                  this.loadExtraResource(data);
                }
  
                $content
                  .find(".bhojpur-button--cancel")
                  .attr("data-dismiss", "slideout")
                  .removeAttr("href");
  
                $scripts = compareScripts($content.find("script[src]"));
                $links = compareLinks($content.find("link[href]"));
  
                if ($scripts.length) {
                  let data = {
                    url: url,
                    response: response
                  };
  
                  loadScripts($scripts, data, function() {});
                }
  
                if ($links.length) {
                  loadStyles($links);
                }
  
                $content.find("script[src],link[href]").remove();
  
                // reset slideout header and body
                $slideout.html(this.$slideoutTemplate);
                $title = $slideout.find(".bhojpur-slideout__title");
                this.$body = $slideout.find(".bhojpur-slideout__body");
  
                $title.html($response.find(options.title).html());
                replaceHtml(
                  $slideout.find(".bhojpur-slideout__body")[0],
                  $content.html()
                );
                this.$body.find(CLASS_HEADER_LOCALE).remove();
  
                $slideout
                  .one(EVENT_SHOWN, function() {
                    $(this).trigger("enable");
                  })
                  .one(EVENT_HIDDEN, function() {
                    $(this).trigger("disable");
                  });
  
                $slideout.find(".bhojpur-slideout__opennew").attr("href", url);
                this.show();
  
                // callback for after slider loaded HTML
                // this callback is deprecated, use slideoutLoaded.bhojpur.slideout event.
                var bhojpurSliderAfterShow = $.fn.bhojpurSliderAfterShow;
                if (bhojpurSliderAfterShow) {
                  for (var name in bhojpurSliderAfterShow) {
                    if (
                      bhojpurSliderAfterShow.hasOwn(name) &&
                      $.isFunction(bhojpurSliderAfterShow[name])
                    ) {
                      bhojpurSliderAfterShow[name]["isLoaded"] = true;
                      bhojpurSliderAfterShow[name].call(this, url, response);
                    }
                  }
                }
  
                // will trigger slideoutLoaded.bhojpur.slideout event after slideout loaded
                $slideout.trigger(EVENT_SLIDEOUT_LOADED, [url, response]);
              } else {
                if (data.returnUrl) {
                  this.load(data.returnUrl);
                } else {
                  this.refresh();
                }
              }
            }, this),
  
            error: $.proxy(function() {
              var errors;
              $(CLASS_BODY_LOADING).remove();
              if ($(".bhojpur-error span").length > 0) {
                errors = $(".bhojpur-error span")
                  .map(function() {
                    return $(this).text();
                  })
                  .get()
                  .join(", ");
              } else {
                errors = BHOJPUR_Translations.serverError;
              }
              window.alert(errors);
            }, this)
          });
        }, this);
  
        if (this.slided) {
          this.hide(true);
          this.$slideout.one(EVENT_HIDDEN, load);
        } else {
          load();
        }
      },
  
      open: function(options) {
        this.addLoading();
        this.load(options.url, options.data);
      },
  
      reload: function(url) {
        this.hide();
        this.load(url);
      },
  
      show: function() {
        var $slideout = this.$slideout;
        var showEvent;
  
        if (this.slided) {
          return;
        }
  
        showEvent = $.Event(EVENT_SHOW);
        $slideout.trigger(showEvent);
  
        if (showEvent.isDefaultPrevented()) {
          return;
        }
  
        $slideout.removeClass(CLASS_MINI);
        this.slideoutType == "mini" && $slideout.addClass(CLASS_MINI);
  
        $slideout.addClass(CLASS_IS_SHOWN).get(0).offsetWidth;
        $slideout
          .one(EVENT_TRANSITIONEND, $.proxy(this.shown, this))
          .addClass(CLASS_IS_SLIDED)
          .scrollTop(0);
      },
  
      shown: function() {
        this.slided = true;
        // Disable to scroll body element
        $("body").addClass(CLASS_OPEN);
        this.$slideout
          .trigger("beforeEnable.bhojpur.slideout")
          .trigger(EVENT_SHOWN)
          .trigger("afterEnable.bhojpur.slideout");
      },
  
      hide: function() {
        let message = {
          confirm: BHOJPUR_Translations.slideoutCloseWarning
        };
  
        if ($.fn.bhojpurSlideoutBeforeHide) {
          BHOJPUR.bhojpurConfirm(
            message,
            function(confirm) {
              if (confirm) {
                this.hideSlideout();
              }
            }.bind(this)
          );
        } else {
          this.hideSlideout();
        }
  
        this.removeSelectedClass();
      },
  
      hideSlideout: function() {
        var $slideout = this.$slideout;
        var hideEvent;
        var $datePicker = $(".bhojpur-datepicker").not(".hidden");
  
        // remove onbeforeunload event
        window.onbeforeunload = null;
        $.fn.bhojpurSlideoutBeforeHide = null;
  
        this.submitXHR && this.submitXHR.abort();
  
        if ($datePicker.length) {
          $datePicker.addClass("hidden");
        }
  
        if (!this.slided) {
          return;
        }
  
        hideEvent = $.Event(EVENT_HIDE);
        $slideout.trigger(hideEvent);
  
        if (hideEvent.isDefaultPrevented()) {
          return;
        }
  
        $slideout
          .one(EVENT_TRANSITIONEND, $.proxy(this.hidden, this))
          .removeClass(`${CLASS_IS_SLIDED} bhojpur-slideout__fullscreen`);
        $slideout.trigger(EVENT_SLIDEOUT_CLOSED);
      },
  
      hidden: function() {
        this.slided = false;
  
        // Enable to scroll body element
        $("body").removeClass(CLASS_OPEN);
  
        this.$slideout.removeClass(CLASS_IS_SHOWN).trigger(EVENT_HIDDEN);
      },
  
      refresh: function() {
        this.hide();
  
        setTimeout(function() {
          window.location.reload();
        }, 350);
      },
  
      destroy: function() {
        this.unbind();
        this.unbuild();
        this.$element.removeData(NAMESPACE);
      }
    };
  
    BhojpurSlideout.DEFAULTS = {
      title: ".bhojpur-form-title, .mdl-layout-title",
      content: false
    };
  
    BhojpurSlideout.TEMPLATE = `<div class="bhojpur-slideout">
              <div class="bhojpur-slideout__header">
                  <div class="bhojpur-slideout__header-link">
                      <a href="#" target="_blank" class="mdl-button mdl-button--icon mdl-js-button mdl-js-repple-effect bhojpur-slideout__opennew"><i class="material-icons">open_in_new</i></a>
                      <a href="#" class="mdl-button mdl-button--icon mdl-js-button mdl-js-repple-effect bhojpur-slideout__fullscreen">
                          <i class="material-icons">fullscreen</i>
                          <i class="material-icons" style="display: none;">fullscreen_exit</i>
                      </a>
                  </div>
                  <button type="button" class="mdl-button mdl-button--icon mdl-js-button mdl-js-repple-effect bhojpur-slideout__close" data-dismiss="slideout">
                      <span class="material-icons">close</span>
                  </button>
                  <h3 class="bhojpur-slideout__title"></h3>
              </div>
              <div class="bhojpur-slideout__body"></div>
          </div>`;
  
    BhojpurSlideout.TEMPLATE_LOADING = `<div class="bhojpur-body__loading">
              <div><div class="mdl-spinner mdl-js-spinner is-active bhojpur-layout__bottomsheet-spinner"></div></div>
          </div>`;
  
    BhojpurSlideout.plugin = function(options) {
      return this.each(function() {
        var $this = $(this);
        var data = $this.data(NAMESPACE);
        var fn;
  
        if (!data) {
          if (/destroy/.test(options)) {
            return;
          }
  
          $this.data(NAMESPACE, (data = new BhojpurSlideout(this, options)));
        }
  
        if (typeof options === "string" && $.isFunction((fn = data[options]))) {
          fn.apply(data);
        }
      });
    };
  
    $.fn.bhojpurSlideout = BhojpurSlideout.plugin;
  
    return BhojpurSlideout;
  });