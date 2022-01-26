"use strict";
var _typeof="function"==typeof Symbol && "symbol"==typeof Symbol.iterator
? function(e) {
    return typeof e
}
: function(e) {
    return e&&
    "function"==typeof Symbol&&
    e.constructor===Symbol&&
    e!==Symbol.prototype
    ? "symbol"
    : typeof e
};
!function(e) {
    "function"==typeof define&&
    define.amd
    ? define(["jquery"],e)
    : e("object"===("undefined"==typeof exports
    ? "undefined"
    : _typeof(exports))
    ? require("jquery")
    : jQuery)
}
(function(e) {
    function t(i,o) {
        this.$element=e(i),
        this.options=e.extend({},
            t.DEFAULTS,
            e.isPlainObject(o)&&
            o),
        this.init()
    }
    var i = e("body"), o = 'select[name="BhojpurResource.Widgets"]';
    return t.prototype = {
        constructor : t,
        init : function() {
            this.bind(),
            this.isNewForm=this.$element.hasClass(
                "bhojpur-layout__widget-new"
            ),
            !this.isNewForm&&
            this.$element.find(o).length&&
            this.addWidgetSlideout(),
            this.initSelect()
        },
        bind : function() {
            this.$element.on(
                "change.bhojpur.widget",
                "select",
            this.change.bind(this)).on(
                "click.bhojpur.widget",
                ".bhojpur-widget__new",
            this.getFormHtml.bind(this)).on(
                "click.bhojpur.widget",
                ".bhojpur-widget__cancel",
            this.cancelForm.bind(this))
        },
        unbind : function(){
            this.$element.off(
                "change.bhojpur.widget",
                "select",
            this.change.bind(this)).off(
                "click.bhojpur.widget",
                ".bhojpur-widget__new",
            this.getFormHtml.bind(this))
        },
        initSelect : function(){
            var t = this.$element,
            i = t.find("select").filter(
                    '[name="BhojpurResource.Widgets"],[name="BhojpurResource.Template"]'
                ),
            o = e('[name="BhojpurResource.Kind"]'),
            n = t.data("hint"),
            r = '<h2 class="bhojpur-page__tips">'+
            n+
            "</h2>";
            if (i.closest(".bhojpur-form-section").hide(),
                i.each(function(){
                    e(this).find("option").filter(
                        '[value!=""]'
                ).length>=2&&
                e(this).closest(
                    ".bhojpur-form-section"
                ).show()
            }),
            !this.isNewForm&&!t.find(".bhojpur-bannereditor").length){
                var s=o.parent().next(".bhojpur-form-section-rows"),
                l=o.closest(".bhojpur-fieldset"),
                a=o.closest(".bhojpur-form-section");
                s.children().length&&
                l.find(".bhojpur-field__label:visible").length||
                (s.append(r),
                t.find(
                    ".bhojpur-field__label").not(a.find(
                    ".bhojpur-field__label")
                    ).is(":visible")||
                (a.hide(),
                t.append(r).parent().find(".bhojpur-form__actions").remove()))
            }
        },
        addWidgetSlideout : function(){
            var t, n, r=this.$element.find(o),
            s=i.data("tabScopeActive"),
            l=e(".bhojpur-slideout").is(":visible"),
            a=r.closest("form"),
            d=a.data("action-url")||
            a.prop("action"),
            c=d&&-1!==d.indexOf("?")
            ? "&"
            : "?";
            r.find("option").each(function(){
                var i=e(this), o=i.val();
                o && (t="" + d + c + "widget_type=" + o, s &&(
                    t = t + "&widget_scope=" + s),
                    n = l
                    ? "<a href=" + t + ' style="display: none;" class="bhojpur-widget-'
                        + o + '" data-open-type="slideout" data-url="' + t + '">'
                        + o + "</a>"
                    : "<a href=" + t + ' style="display: none;" class="bhojpur-widget-'
                        + o + '">'
                        + o + "</a>",
                        r.after(n))
            })
        },
        change : function(t) {
            var i=e(t.target),
            n=i.val(),
            r=e(".bhojpur-slideout").is(":visible"),
            s=".bhojpur-widget-" + n,
            l = e(s),
            a = l.prop("href");
            if(i.is(o))
                return e.fn.bhojpurSlideoutBeforeHide=null,
                window.onbeforeunload=null,
                this.isNewForm ||
                (r ? l.trigger("click")
                    : location.href=a),
                    !1
        },
        getFormHtml : function(i) {
            var o = e(i.target).closest("a"),
            n = o.data("widget-type"),
            r = this.$element,
            s = o.attr("href"),
            l = r.find(".bhojpur-layout__widget-selector"),
            a = l.find("select"),
            d = e(".bhojpur-layout__widget-setting"),
            c = r.find('[data-section-title="Settings"]'),
            f = e(t.TEMPLATE_LOADING);
            return c.length && (d=c),
            f.appendTo(d).trigger("enable"),
            r.find(".bhojpur-slideout__lists-item a").hide(),
            r.find(".bhojpur-slideout__lists-groupname").hide(),
            r.find(".bhojpur-layout__widget-actions").show(),
            e.get(s,function(e){
                l.find(".bhojpur-layout__widget-name").html(
                    o.data("widget-name")),
                    l.show(),
                    a.val(n).closest(".bhojpur-form-section").hide(),
                    d.html(e).trigger("enable")}).fail(function(){
                        window.alert("server error, please try again!")
                    }),!1
        },
        cancelForm : function() {
            var e = this.$element;
            e.closest(".bhojpur-bottomsheets").length&&
            e.closest(".bhojpur-bottomsheets").removeClass("bhojpur-bottomsheets__fullscreen"),
            e.find(".bhojpur-slideout__lists-item a").show(),
            e.find(".bhojpur-slideout__lists-groupname").show(),
            e.find(".bhojpur-layout__widget-actions, .bhojpur-layout__widget-selector").hide(),
            e.find(".bhojpur-layout__widget-setting").html("")},
        destroy : function(){
            this.unbind(),
            this.$element.removeData("bhojpur.widget")
        }
    },
    t.DEFAULTS={},
    t.TEMPLATE_LOADING='<div style="text-align: center; margin-top: 30px;"><div class="mdl-spinner mdl-js-spinner is-active bhojpur-layout__bottomsheet-spinner"></div></div>',
    t.plugin=function(i){
        return this.each(function(){
            var o=e(this),
            n=o.data("bhojpur.widget"),
            r=void 0;
            if(!n){
                if(/destroy/.test(i))
                return;
            o.data("bhojpur.widget",
            n=new t(this,i))
            }
            "string"==typeof i&&
            e.isFunction(r=n[i])&&
            r.apply(n)
        })
    },
    e(function(){
        var i='[data-toggle="bhojpur.widget"]';
        e(document).on("disable.bhojpur.widget",
        function(o){
            t.plugin.call(e(i,o.target),"destroy")
        }).on("enable.bhojpur.widget",
        function(o){
            t.plugin.call(e(i,o.target))
        }).triggerHandler("enable.bhojpur.widget"),
        e(".bhojpur-page__header .bhojpur-page-subnav__header").length&&
        e(".mdl-layout__content").addClass("has-subnav")
    }),t
});