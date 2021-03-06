"use strict";
$R.add("plugin", "medialibrary", {
    init : function(e) {
        this.app=e,
        this.opts=e.opts,
        this.lang=e.lang,
        this.inline=e.inline,
        this.toolbar=e.toolbar,
        this.func=new this.funcInit(this)
    },
    start : function() {
        var e=this.toolbar.addButton("medialibrary", {
            title:"MediaLibrary",
            api:"plugin.medialibrary.addMedialibrary"
        });
        e.setIcon('<i class="material-icons">photo_library</i>'),
        this.buttonElement=e.nodes[0],
        this.$currentTag=!1,
        $(document).on(
            "reload.bhojpur.bottomsheets",
            ".bhojpur-bottomsheets__mediabox",
        this.func.initItem)
    },
    addMedialibrary : function() {
        this.$currentTag=this.app.selection.getCurrent(),
        this.func.addMedialibrary()
    },
    funcInit : function(e) {
        var f=e, a=this;
        this.addMedialibrary=function() {
            var e, t={
                selectModal:"mediabox",
                maxItem:"1"
            },
            i=$(f.app.rootElement).data().redactorSettings.medialibraryUrl;
            a.BottomSheets=e=$("body").data("bhojpur.bottomsheets"),
            t.url=i,
            e.open(t, function(e) {
                a.handleMediaLibrary(e)
            })
        },
        this.handleMediaLibrary=function(e) {
            var t={
                onSelect:a.handleResults,
                onSubmit:a.handleResults
            };
            (a.$bottomsheets=e).bhojpurSelectCore(t).addClass("bhojpur-bottomsheets__mediabox"),
            a.initItem()
        },
        this.initItem=function(){
            var e,t;
            $(".bhojpur-bottomsheets").find("tbody tr").each(function() {
                e=$(this),
                t=e.find(".bhojpur-table--ml-slideout p img").first(),
                e.find(".bhojpur-table__actions").remove(),
                t.length&&(e.find(".bhojpur-table--medialibrary-item").css("background-image",
                "url("+t.prop("src")+")"),
                t.parent().remove())
            })
        },
        this.handleResults=function(e) {
            var t=e.MediaOption;
            "video_link"==e.SelectedType || t.Video || t.URL.match(/\.mp4$|\.m4p$|\.m4v$|\.m4v$|\.mov$|\.mpeg$|\.webm$|\.avi$|\.ogg$|\.ogv$/)
            ? a.insertVideo(e)
            : a.insertImage(e),
            a.$bottomsheets.remove(),
            $(".bhojpur-bottomsheets").is(":visible") || $("body").removeClass("bhojpur-bottomsheets-open")
        },
        this.insertVideo=function(e){
            var t=$(f.app.editor.$editor.nodes[0]),
            i=$(f.app.rootElement);
            f.opts.mediaContainerClass=void 0===f.opts.mediaContainerClass
            ? "bhojpur-video-container"
            : f.opts.mediaContainerClass;
            var a,o,r,n,d,s,c,l={},
            m=f.opts.mediaContainerClass,
            p=f.opts.regex.youtube,
            u=f.opts.regex.vimeo,
            h=/(\/id_)(\w+)/,
            b="bhojpur-video-"+(Math.random()+1).toString(36).substring(7),
            g=e.MediaOption,
            v=g.Description;
            r='<figure class="'.concat(m,
                ' video-scale"><iframe title="').concat(v,
                    '" data-media-id="').concat(e.ID || e.primaryKey,
                        '" width="100%" height="380px" src="'),
                        n='" frameborder="0" allowfullscreen="true"></iframe></figure>',
                        "video_link"==e.SelectedType
                        ? a=(o=g.Video).match(p)
                        ? (d="youtube",
                            o.replace(p,
                            r+"//www.youtube.com/embed/$1"+n))
                        : o.match(u)?(d="vimeo",
                            o.replace(u,
                            r+"//player.vimeo.com/video/$2"+n))
                        : o.match(/http?:\/\/(www\.)|(v\.)youku.com/) && h.test(o)
                        ? (d="youku",
                            c=o.match(h)[2],
                            '<div class="video-scale"><iframe width=100% height=400 data-media-id="'.concat(e.ID||
                                e.primaryKey,'" src="http://player.youku.com/embed/').concat(c,
                                '" frameborder=0 allowfullscreen="true"></iframe></div>'))
                        : (d="others",
                            '<div class="video-scale"><iframe data-media-id="'.concat(e.ID||
                                e.primaryKey,'" width=100% height=400 src="').concat(o,
                                    '" frameborder=0 allowfullscreen="true"></iframe></div>'))
                        : g.URL.match(/\.mp4$|\.m4p$|\.m4v$|\.m4v$|\.mov$|\.mpeg$|\.webm$|\.avi$|\.ogg$|\.ogv$/) &&
                        (d="uploadedVideo",
                        a='<figure class="'
                        +m+
                        '"><div class="video-scale" role="application"><video width="100%" title="'
                        +v+
                        '" aria-label="'
                        +v+
                        '" height="380px" controls="controls" aria-describedby="'
                        +b+
                        '" tabindex="0"><source src="'
                        +g.URL+
                        '"></video></div></figure>'),
                        a &&
                        (s=$(a).addClass(b),
                        f.$currentTag
                        ? $(f.$currentTag).after(s)
                        : t.prepend(s),
                        l.type=d,
                        l.videoLink=o || g.URL,
                        l.videoIdentification=b,
                        l.description=v,
                        l.$editor=t,
                        (l.$element=i).trigger("insertedVideo.redactor",
                        [l]))
                    },
                    this.insertImage=function(e){
                        var t,
                        i=$(f.app.editor.$editor.nodes[0]),
                        a=$(f.app.rootElement),
                        o=$("<img>"),
                        r=$("<figure>"),
                        n=e.MediaOption,
                        d={},
                        s=e.File&&JSON.parse(e.File);
                        t=n.URL.replace(/image\..+\./,"image."),
                        o.attr({
                            src:t,
                            alt:n.
                            Description||
                            s.Description
                        }),
                        r.append(o),
                        f.$currentTag
                        ? $(f.$currentTag).after(r)
                        : i.prepend(r),
                        r.addClass("redactor-component"),
                        r.attr({"data-redactor-type"
                        : "image",
                        tabindex : "-1",
                        contenteditable:!1
                    }),
                    d.description=n.Description,
                    d.$img=r,
                    d.$editor=i,
                    (d.$element=a).trigger("insertedImage.redactor",
                    [d])
                }
            }
        });