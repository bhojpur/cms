$(function() {
    let _ = window._,
      BHOJPUR = window.BHOJPUR,
      BHOJPUR_Translations = window.BHOJPUR_Translations,
      html = `<div id="dialog" style="display: none;">
                    <div class="mdl-dialog-bg"></div>
                    <div class="mdl-dialog">
                        <div class="mdl-dialog__content">
                          <p><i class="material-icons">warning</i></p>
                          <p class="mdl-dialog__message dialog-message">
                          </p>
                        </div>
                        <div class="mdl-dialog__actions">
                          <button type="button" class="mdl-button mdl-button--raised mdl-button--colored dialog-ok dialog-button" data-type="confirm">
                            ${BHOJPUR_Translations.okButton}
                          </button>
                          <button type="button" class="mdl-button dialog-cancel dialog-button" data-type="">
                            ${BHOJPUR_Translations.cancelButton}
                          </button>
                        </div>
                      </div>
                  </div>`,
      $dialog = $(html).appendTo("body");
  
    // ************************************ Refactor window.confirm ************************************
    $(document)
      .on("keyup.bhojpur.confirm", function(e) {
        if (!$dialog.is(":visible")) {
          return;
        }
        if (e.which === 27) {
          setTimeout(function() {
            $dialog.hide();
            BHOJPUR.bhojpurConfirmCallback = undefined;
          }, 100);
        }
        if (e.which === 13) {
          setTimeout(function() {
            $('.dialog-button[data-type="confirm"]').click();
          }, 100);
        }
      })
      .on("click.bhojpur.confirm", ".dialog-button", function() {
        let value = $(this).data("type"),
          callback = BHOJPUR.bhojpurConfirmCallback;
  
        $.isFunction(callback) && callback(value);
        $dialog.hide();
        BHOJPUR.bhojpurConfirmCallback = undefined;
        return false;
      });
  
    BHOJPUR.bhojpurConfirm = function(data, callback) {
      let okBtn = $dialog.find(".dialog-ok"),
        cancelBtn = $dialog.find(".dialog-cancel");
  
      if (_.isString(data)) {
        $dialog.find(".dialog-message").text(data);
        okBtn.text(BHOJPUR_Translations.okButton);
        cancelBtn.text(BHOJPUR_Translations.cancelButton);
      } else if (_.isObject(data)) {
        if (data.confirmOk && data.confirmCancel) {
          okBtn.text(data.confirmOk);
          cancelBtn.text(data.confirmCancel);
        } else {
          okBtn.text(BHOJPUR_Translations.okButton);
          cancelBtn.text(BHOJPUR_Translations.cancelButton);
        }
  
        if(data.icon){
          $dialog.find('i.material-icons').addClass(data.icon).html(data.icon);
        }
  
        $dialog.find(".dialog-message").text(data.confirm);
      }
  
      $dialog.show();
      BHOJPUR.bhojpurConfirmCallback = callback;
      return false;
    };
  
    // *******************************************************************************
  
    // ****************Handle download file from AJAX POST****************************
    let objectToFormData = function(obj, form) {
      let formdata = form || new FormData(),
        key;
  
      for (var variable in obj) {
        if (obj.hasOwn(variable) && obj[variable]) {
          key = variable;
        }
  
        if (obj[variable] instanceof Date) {
          formdata.append(key, obj[variable].toISOString());
        } else if (
          typeof obj[variable] === "object" &&
          !(obj[variable] instanceof File)
        ) {
          objectToFormData(obj[variable], formdata);
        } else {
          formdata.append(key, obj[variable]);
        }
      }
  
      return formdata;
    };
  
    BHOJPUR.bhojpurAjaxHandleFile = function(url, contentType, fileName, data) {
      let request = new XMLHttpRequest();
  
      request.responseType = "arraybuffer";
      request.open("POST", url, true);
      request.onload = function() {
        if (this.status === 200) {
          let blob = new Blob([this.response], {
              type: contentType
            }),
            url = window.URL.createObjectURL(blob),
            a = document.createElement("a");
  
          document.body.appendChild(a);
          a.href = url;
          a.download = fileName || "download-" + $.now();
          a.click();
        } else {
          window.alert(BHOJPUR_Translations.serverError);
        }
      };
  
      if (_.isObject(data)) {
        if (Object.prototype.toString.call(data) != "[object FormData]") {
          data = objectToFormData(data);
        }
  
        request.send(data);
      }
    };
  
    // ********************************convert video link********************
    // linkyoutube: /https?:\/\/(?:[0-9A-Z-]+\.)?(?:youtu\.be\/|youtube\.com\S*[^\w\-\s])([\w\-]{11})(?=[^\w\-]|$)(?![?=&+%\w.\-]*(?:['"][^<>]*>|<\/a>))[?=&+%\w.-]*/ig,
    // linkvimeo: /https?:\/\/(www\.)?vimeo.com\/(\d+)($|\/)/,
  
    let converVideoLinks = function() {
      let $ele = $(".bhojpur-linkify-object"),
        linkyoutube = /https?:\/\/(?:[0-9A-Z-]+\.)?(?:youtu\.be\/|youtube\.com\S*[^\w\-\s])([\w-]{11})(?=[^\w-]|$)(?![?=&+%\w.-]*(?:['"][^<>]*>|<\/a>))[?=&+%\w.-]*/gi;
  
      if (!$ele.length) {
        return;
      }
  
      $ele.each(function() {
        let url = $(this).data("video-link");
        if (url.match(linkyoutube)) {
          $(this).html(
            `<iframe width="100%" height="100%" src="//www.youtube.com/embed/${url.replace(
              linkyoutube,
              "$1"
            )}" frameborder="0" allowfullscreen></iframe>`
          );
        }
      });
    };
  
    $.fn.bhojpurSliderAfterShow.converVideoLinks = converVideoLinks;
    converVideoLinks();
  
    // ********************************Bhojpur CMS Handle AJAX error********************
    BHOJPUR.handleAjaxError = function(err) {
      let $body = $("body"),
        rJSON = err.responseJSON,
        rText = err.responseText,
        $error = $(`<ul class="bhojpur-alert bhojpur-error" data-dismissible="true"><button type="button" class="mdl-button mdl-button--icon" data-dismiss="alert">
                              <i class="material-icons">close</i>
                          </button></ul>`);
  
      $body.find(".bhojpur-alert").remove();
  
      if (err.status === 422) {
        if (rJSON) {
          let errors = rJSON.errors,
            $errorContent = "";
  
          if ($.isArray(errors)) {
            for (let i = 0; i < errors.length; i++) {
              $errorContent += `<li>
                                            <i class="material-icons">error</i>
                                            <span>${errors[i]}</span>
                                        </li>`;
            }
          } else {
              $errorContent = `<li>
                                <i class="material-icons">error</i>
                                <span>${errors}</span>
                            </li>`;
          }
          $error.append($errorContent);
        } else {
          $error = $(rText).find(".bhojpur-error");
        }
      } else {
        $error.append(`<li>
                              <i class="material-icons">error</i>
                              <span>${err.statusText}</span>
                          </li>`);
      }
  
      $error.prependTo($body);
      setTimeout(function() {
        $error.addClass("bhojpur-alert__active");
      }, 50);
  
      setTimeout(function() {
        $('.bhojpur-alert[data-dismissible="true"]').removeClass("bhojpur-alert__active");
        $("#bhojpur-submit-loading").remove();
      }, 6000);
    };
  });