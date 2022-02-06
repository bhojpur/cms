$(function() {
    "use strict";
  
    let $body = $("body"),
      Slideout,
      BottomSheets,
      CLASS_IS_SELECTED = "is-selected",
      isSlideoutOpened = function() {
        return $body.hasClass("bhojpur-slideout-open");
      },
      isBottomsheetsOpened = function() {
        return $body.hasClass("bhojpur-bottomsheets-open");
      };
  
    $body.bhojpurBottomSheets();
    $body.bhojpurSlideout();
  
    Slideout = $body.data("bhojpur.slideout");
    BottomSheets = $body.data("bhojpur.bottomsheets");
  
    function toggleSelectedCss(ele) {
      $("[data-url]").removeClass(CLASS_IS_SELECTED);
      ele && ele.length && ele.addClass(CLASS_IS_SELECTED);
    }
  
    function collectSelectID() {
      let $checked = $(".bhojpur-js-table tbody").find(
          ".mdl-checkbox__input:checked"
        ),
        IDs = [];
  
      if (!$checked.length) {
        return false;
      }
  
      $checked.each(function() {
        IDs.push(
          $(this)
            .closest("tr")
            .data("primary-key")
        );
      });
  
      return IDs;
    }
  
    $(document).on("click.bhojpur.openUrl", "[data-url]", function(e) {
      let $this = $(this),
        $target = $(e.target),
        isNewButton = $this.hasClass("bhojpur-button--new"),
        isEditButton = $this.hasClass("bhojpur-button--edit"),
        isInTable =
          ($this.is(".bhojpur-table tr[data-url]") ||
            $this.closest(".bhojpur-js-table").length) &&
          !$this.closest(".bhojpur-slideout").length, // if table is in slideout, will open bottom sheet
        openData = $this.data(),
        actionData,
        openType = openData.openType,
        hasSlideoutTheme = $this.parents(".bhojpur-theme-slideout").length,
        isInSlideout = $this.closest(".bhojpur-slideout").length,
        isActionButton =
          $this.hasClass("bhojpur-action-button") ||
          $this.hasClass("bhojpur-action--button");
  
      e.stopPropagation();
  
      // if clicking item's menu actions
      if (
        $this.data("ajax-form") ||
        $target.closest(".bhojpur-table--bulking").length ||
        $target.closest(".bhojpur-button--actions").length ||
        (!$target.data("url") && $target.is("a")) ||
        (isInTable && isBottomsheetsOpened())
      ) {
        return;
      }
  
      if (openType == "window") {
        window.location.href = openData.url;
        return;
      }
  
      if (openType == "new_window") {
        window.open(openData.url, "_blank");
        return;
      }
  
      if (isActionButton) {
        actionData = collectSelectID();
        if (actionData) {
          openData = $.extend({}, openData, {
            actionData: actionData
          });
        }
      }
  
      openData.$target = $target;
  
      if (!openData.method || openData.method.toUpperCase() == "GET") {
        // Open in BottmSheet: is action button, open type is bottom-sheet
        // is action button  but opentype == slideout, should open in slideout\
        // open type is No.1 priority
  
        if (
          (openType == "bottomsheet" || isActionButton) &&
          openType != "slideout"
        ) {
          // if is bulk action and no item selected
          if (
            isActionButton &&
            !actionData &&
            $this.closest('[data-toggle="bhojpur.action.bulk"]').length &&
            !isInSlideout
          ) {
            window.bhojpur.bhojpurConfirm(openData.errorNoItem);
            return false;
          }
  
          BottomSheets.open(openData);
          return false;
        }
  
        // Slideout or New Page: table items, new button, edit button
        if (
          openType == "slideout" ||
          isInTable ||
          (isNewButton && !isBottomsheetsOpened()) ||
          isEditButton
        ) {
          if (openType == "slideout" || hasSlideoutTheme) {
            if ($this.hasClass(CLASS_IS_SELECTED)) {
              Slideout.hide();
              toggleSelectedCss();
              return false;
            } else {
              Slideout.open(openData);
              toggleSelectedCss($this);
              return false;
            }
          } else {
            window.location.href = openData.url;
            return false;
          }
        }
  
        // Open in BottmSheet: slideout is opened or openType is Bottom Sheet
        if (isSlideoutOpened() || (isNewButton && isBottomsheetsOpened())) {
          BottomSheets.open(openData);
          return false;
        }
  
        // Other clicks
        if (hasSlideoutTheme) {
          Slideout.open(openData);
          return false;
        } else {
          BottomSheets.open(openData);
          return false;
        }
      }
    });
  });