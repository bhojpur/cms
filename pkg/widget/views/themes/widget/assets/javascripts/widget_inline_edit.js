"use strict";
var _typeof="function"==typeof Symbol&&
    "symbol"==typeof Symbol.iterator
    ? function(t){
        return typeof t
    }
    : function(t){
        return t&&
        "function"==typeof Symbol&&
        t.constructor===Symbol&&
        t!==Symbol.prototype
        ? "symbol"
        : typeof t
    };
    !function(t){
        "function"==typeof define&&
        define.amd
        ? define(["jquery"],t)
        : t("object"===("undefined"==typeof exports
        ? "undefined"
        : _typeof(exports))
        ? require("jquery")
        : jQuery)
    }
    (function(t){
        function e(i,n){
            this.$element=t(i),
            this.options=t.extend({
            },
            e.DEFAULTS, t.isPlainObject(n)&&n
            ),
            this.init()
        }
        var i, n=t("body"), o="bhojpur.widget.inlineEdit";
        return e.prototype =
        {
            constructor:e,
            init: function(){
                this.bind(),
                this.initStatus()
            },
            bind:function(){
                this.$element.on(
                    "click.bhojpur.widget.inlineEdit",
                    ".bhojpur-widget-button",
                    this.click),
                    t(document).on("keyup",
                    this.keyup)
            },
            initStatus:function(){
                var e=document.createElement("iframe");
                e.src=i,
                e.id="bhojpur-widget-iframe",
                e.attachEvent
                ? e.attachEvent("onload",
                    function(){
                        t(".bhojpur-widget-button").show()
                    })
                : e.onload=function(){
                    t(".bhojpur-widget-button").show()
                },
                document.body.appendChild(e)
            },
            keyup:function(t){
                var e=document.getElementById("bhojpur-widget-iframe");
                27===t.keyCode&&
                e&&
                e.contentDocument.querySelector(".bhojpur-slideout__close").click()
            },
            click:function(){
                var e=t(this),
                i=document.getElementById("bhojpur-widget-iframe"),
                o=i.contentWindow.$,
                d=i.contentDocument.querySelector(".js-widget-edit-link");
                if(d)
                    return i.classList.add("show"),
                        o
                        ? o(".js-widget-edit-link").data("url", e.data("url")).click()
                        : (d.setAttribute("data-url", e.data("url")), d.click()),
                n.addClass("open-widget-editor"), !1
            }
        },
        e.plugin=function(i){
            return this.each(function(){
                var n,
                d=t(this),
                r=d.data(o);
                if(!r){
                    if(/destroy/.test(i))
                        return;
                    d.data(o, r=new e(this,i))
                }
                "string"==typeof i&&
                t.isFunction(n=r[i])&&
                n.apply(r)
            })
        },
        t(function(){
            n.attr("data-toggle","bhojpur.widgets"),
            t(".bhojpur-widget").each(function(){
                var e=t(this),
                n=e.children().eq(0);
                i=e.data(
                    "widget-inline-edit-url"),
                    "static"===n.css("position")&&
                        n.css("position", "relative"),
                        n.addClass("bhojpur-widget").unwrap(),
                n.append(
                    '<div class="bhojpur-widget-embed-wrapper"><button style="display: none;" data-url="'+
                e.data("url")+
                '" class="bhojpur-widget-button">Edit</button></div>')
            });
            var o='[data-toggle="bhojpur.widgets"]';
        t(document).on("disable.bhojpur.widget.inlineEdit",
            function(i){
                e.plugin.call(t(o,i.target),
                "destroy")
            }).on("enable.bhojpur.widget.inlineEdit",
            function(i){
                e.plugin.call(t(o,i.target))
        }).triggerHandler("enable.bhojpur.widget.inlineEdit")
        }),
        e
});