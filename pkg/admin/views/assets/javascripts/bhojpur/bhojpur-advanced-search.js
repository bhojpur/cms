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
  
    let location = window.location,
      BHOJPUR = window.BHOJPUR,
      NAMESPACE = "bhojpur.advancedsearch",
      EVENT_ENABLE = "enable." + NAMESPACE,
      EVENT_DISABLE = "disable." + NAMESPACE,
      EVENT_CLICK = "click." + NAMESPACE,
      EVENT_SHOWN = "shown.bhojpur.modal",
      EVENT_SUBMIT = "submit." + NAMESPACE;
  
    function getExtraPairs(names) {
      let pairs = decodeURIComponent(location.search.substr(1)).split("&"),
        pairsObj = {},
        pair,
        i;
  
      if (pairs.length == 1 && pairs[0] == "") {
        return false;
      }
  
      for (i in pairs) {
        if (pairs[i] === "") continue;
  
        pair = pairs[i].split("=");
        pairsObj[pair[0]] = pair[1];
      }
  
      names.forEach(function(item) {
        delete pairsObj[item];
      });
  
      return pairsObj;
    }
  
    function BhojpurAdvancedSearch(element, options) {
      this.$element = $(element);
      this.options = $.extend(
        {},
        BhojpurAdvancedSearch.DEFAULTS,
        $.isPlainObject(options) && options
      );
      this.init();
    }
  
    BhojpurAdvancedSearch.prototype = {
      constructor: BhojpurAdvancedSearch,
  
      init: function() {
        this.$form = this.$element.find("form");
        this.$modal = $(BhojpurAdvancedSearch.MODAL).appendTo("body");
        this.bind();
      },
  
      bind: function() {
        this.$element
          .on(EVENT_SUBMIT, "form", this.submit.bind(this))
          .on(
            EVENT_CLICK,
            ".bhojpur-advanced-filter__save",
            this.showSaveFilter.bind(this)
          )
          .on(
            EVENT_CLICK,
            ".bhojpur-advanced-filter__toggle",
            this.toggleFilterContent
          )
          .on(EVENT_CLICK, ".bhojpur-advanced-filter__close", this.closeFilter)
          .on(
            EVENT_CLICK,
            ".bhojpur-advanced-filter__delete",
            this.deleteSavedFilter
          );
  
        this.$modal.on(EVENT_SHOWN, this.start.bind(this));
      },
  
      closeFilter: function() {
        $(".bhojpur-advanced-filter__dropdown").hide();
      },
  
      toggleFilterContent: function(e) {
        $(e.target)
          .closest(".bhojpur-advanced-filter__toggle")
          .parent()
          .find(">[advanced-search-toggle]")
          .toggle();
      },
  
      showSaveFilter: function() {
        this.$modal.bhojpurModal("show");
      },
  
      deleteSavedFilter: function(e) {
        let $target = $(e.target).closest(".bhojpur-advanced-filter__delete"),
          $savedFilter = $target.closest(".bhojpur-advanced-filter__savedfilter"),
          name = $target.data("filter-name"),
          url = location.pathname,
          message = {
            confirm: "Are you sure you want to delete this saved filter?"
          };
  
        BHOJPUR.bhojpurConfirm(message, function(confirm) {
          if (confirm) {
            $.get(url, $.param({ delete_saved_filter: name }))
              .done(function() {
                $target.closest("li").remove();
                if ($savedFilter.find("li").length === 0) {
                  $savedFilter.remove();
                }
              })
              .fail(function() {
                BHOJPUR.bhojpurConfirm("Server error, please try again!");
              });
          }
        });
        return false;
      },
  
      start: function() {
        this.$modal
          .trigger("enable.bhojpur.material")
          .on(
            EVENT_CLICK,
            ".bhojpur-advanced-filter__savefilter",
            this.saveFilter.bind(this)
          );
      },
  
      saveFilter: function() {
        let name = this.$modal.find("#bhojpur-advanced-filter__savename").val();
  
        if (!name) {
          return;
        }
  
        this.$form
          .prepend(
            `<input type="hidden" name="filter_saving_name" value="${name}"  />`
          )
          .submit();
      },
  
      submit: function() {
        let $form = this.$form,
          formArr = $form.find("input[name],select[name]"),
          names = [],
          extraPairs,
          $bottomsheet = $form.closest(".bhojpur-bottomsheets"),
          params = $form.serialize();
  
        formArr.each(function() {
          names.push($(this).attr("name"));
        });
  
        extraPairs = getExtraPairs(names);
  
        if (!$.isEmptyObject(extraPairs)) {
          for (let key in extraPairs) {
            if (extraPairs.hasOwn(key)) {
              $form.prepend(
                `<input type="hidden" name=${key} value=${extraPairs[key]}  />`
              );
            }
          }
        }
  
        this.$element.find(".bhojpur-advanced-filter__dropdown").hide();
  
        this.removeEmptyPairs($form);
  
        if ($bottomsheet.length) {
          if ($bottomsheet.data().url) {
            let reloadUrl = `${$bottomsheet.data().url}?${params}`;
            $bottomsheet.trigger("reloadFromUrl.bhojpur.bottomsheets", [reloadUrl]);
            return false;
          } else {
            console.log("dont have base URL! advancedsearch reload failed");
          }
        }
      },
  
      removeEmptyPairs: function($form) {
        $form.find("advanced-filter-group").each(function() {
          let $this = $(this),
            $input = $this.find("[filter-required]");
          if ($input.val() == "") {
            $this.remove();
          }
        });
      },
  
      destroy: function() {
        this.$element.removeData(NAMESPACE);
      }
    };
  
    BhojpurAdvancedSearch.DEFAULTS = {};
  
    BhojpurAdvancedSearch.MODAL = `<div class="bhojpur-modal fade" tabindex="-1" role="dialog" aria-hidden="true">
              <div class="mdl-card mdl-shadow--2dp" role="document">
                  <div class="mdl-card__title">
                      <h2 class="mdl-card__title-text">Save advanced filter</h2>
                  </div>
                  <div class="mdl-card__supporting-text">
                          
                      <div class="mdl-textfield mdl-textfield--full-width mdl-js-textfield">
                          <input class="mdl-textfield__input" type="text" id="bhojpur-advanced-filter__savename">
                          <label class="mdl-textfield__label" for="bhojpur-advanced-filter__savename">Please enter name for this filter</label>
                      </div>
  
                  </div>
                  <div class="mdl-card__actions">
                      <a class="mdl-button mdl-button--colored mdl-button--raised bhojpur-advanced-filter__savefilter">Save This Filter</a>
                      <a class="mdl-button mdl-button--colored" data-dismiss="modal">Cancel</a>
                  </div>
                  <div class="mdl-card__menu">
                      <button class="mdl-button mdl-button--icon" data-dismiss="modal" aria-label="close">
                          <i class="material-icons">close</i>
                      </button>
                  </div>
              </div>
          </div>`;
  
    BhojpurAdvancedSearch.plugin = function(options) {
      return this.each(function() {
        let $this = $(this),
          data = $this.data(NAMESPACE),
          fn;
  
        if (!data) {
          if (/destroy/.test(options)) {
            return;
          }
  
          $this.data(NAMESPACE, (data = new BhojpurAdvancedSearch(this, options)));
        }
  
        if (typeof options === "string" && $.isFunction((fn = data[options]))) {
          fn.apply(data);
        }
      });
    };
  
    $(function() {
      let selector = '[data-toggle="bhojpur.advancedsearch"]',
        options;
  
      $(document)
        .on(EVENT_DISABLE, function(e) {
          BhojpurAdvancedSearch.plugin.call($(selector, e.target), "destroy");
        })
        .on(EVENT_ENABLE, function(e) {
          BhojpurAdvancedSearch.plugin.call($(selector, e.target), options);
        })
        .triggerHandler(EVENT_ENABLE);
    });
  
    return BhojpurAdvancedSearch;
  });