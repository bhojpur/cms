"use strict";

let gulp = require("gulp"),
  babel = require("gulp-babel"),
  eslint = require("gulp-eslint"),
  plumber = require("gulp-plumber"),
  cleanCSS = require("gulp-clean-css"),
  concat = require("gulp-concat"),
  sass = require("gulp-sass"),
  uglify = require("gulp-uglify"),
  autoprefixer = require("gulp-autoprefixer"),
  fs = require("fs"),
  path = require("path"),
  es = require("event-stream"),
  rename = require("gulp-rename");

let moduleName = (function() {
  let args = process.argv,
    length = args.length,
    i = 0,
    name,
    subName,
    useSubName;

  while (i++ < length) {
    if (/^--+(\w+)/i.test(args[i])) {
      name = args[i].split("--")[1];
      subName = args[i].split("--")[2];
      useSubName = args[i].split("--")[3];
      break;
    }
  }
  return {
    name: name,
    subName: subName,
    useSubName: useSubName,
  };
})();

// Bhojpur CMS - System Adminitrator module
// Command: gulp [task]
// Admin is the default task
// Watch Admin module: gulp
// -----------------------------------------------------------------------------

function adminTasks() {
  let pathto = function(file) {
    return "pkg/admin/views/assets/" + file;
  };
  let scripts = {
    src: pathto("javascripts/app/*.js"),
    dest: pathto("javascripts"),
    bhojpur: pathto("javascripts/bhojpur/*.js"),
    bhojpurInit: pathto("javascripts/bhojpur/bhojpur-config.js"),
    bhojpurCommon: pathto("javascripts/bhojpur/bhojpur-common.js"),
    bhojpurAdmin: [pathto("javascripts/bhojpur.js"), pathto("javascripts/app.js")],
    all: ["gulpfile.js", pathto("javascripts/bhojpur/*.js")],
  };
  let styles = {
    src: pathto("stylesheets/scss/{app,bhojpur}.scss"),
    dest: pathto("stylesheets"),
    vendors: pathto("stylesheets/vendors"),
    main: pathto("stylesheets/{bhojpur,app}.css"),
    bhojpurAdmin: [
      pathto("stylesheets/vendors.css"),
      pathto("stylesheets/bhojpur.css"),
      pathto("stylesheets/app.css"),
    ],
    scss: pathto("stylesheets/scss/**/*.scss"),
  };

  gulp.task("bhojpur", function() {
    return gulp
      .src([scripts.bhojpurInit, scripts.bhojpurCommon, scripts.bhojpur])
      .pipe(plumber())
      .pipe(concat("bhojpur.js"))
      .pipe(uglify())
      .pipe(gulp.dest(scripts.dest));
  });

  gulp.task("js",
    gulp.series("bhojpur", function() {
      return gulp
        .src(scripts.src)
        .pipe(plumber())
        .pipe(
          eslint({
            configFile: ".eslintrc",
          })
        )
        .pipe(concat("app.js"))
        .pipe(uglify())
        .pipe(gulp.dest(scripts.dest));
    })
  );

  gulp.task("bhojpur+", function() {
    return gulp
      .src([scripts.bhojpurInit, scripts.bhojpurCommon, scripts.bhojpur])
      .pipe(plumber())
      .pipe(
        eslint({
          configFile: ".eslintrc",
        })
      )
      .pipe(
        babel({
          presets: ["@babel/env"],
        })
      )
      .pipe(eslint.format())
      .pipe(concat("bhojpur.js"))
      .pipe(uglify())
      .pipe(gulp.dest(scripts.dest));
  });

  gulp.task("js+", function() {
    return gulp
      .src(scripts.src)
      .pipe(plumber())
      .pipe(
        babel({
          presets: ["@babel/env"],
        })
      )
      .pipe(eslint.format())
      .pipe(concat("app.js"))
      .pipe(uglify())
      .pipe(gulp.dest(scripts.dest));
  });

  gulp.task("sass", function() {
    return gulp
      .src(styles.src)
      .pipe(plumber())
      .pipe(sass().on("error", sass.logError))
      .pipe(gulp.dest(styles.dest));
  });

  gulp.task("css",
    gulp.series("sass", function() {
      return gulp
        .src(styles.main)
        .pipe(plumber())
        .pipe(autoprefixer())
        .pipe(cleanCSS())
        .pipe(gulp.dest(styles.dest));
    })
  );

  gulp.task("release_js", function() {
    return gulp
      .src(scripts.bhojpurAdmin)
      .pipe(concat("bhojpur_admin_default.js"))
      .pipe(gulp.dest(scripts.dest));
  });

  gulp.task("release_css", function() {
    return gulp
      .src(styles.bhojpurAdmin)
      .pipe(concat("bhojpur_admin_default.css"))
      .pipe(gulp.dest(styles.dest));
  });

  gulp.task("release",
    gulp.series("bhojpur+", "js+", "css", "release_js", "release_css")
  );

  let watcher = gulp.task("watch", function() {
    let watch_bhojpur = gulp.watch(scripts.bhojpur, gulp.series("bhojpur+")),
      watch_js = gulp.watch(scripts.src, gulp.series("js+")),
      watch_css = gulp.watch(styles.scss, gulp.series("css"));

    gulp.watch(styles.bhojpurAdmin, gulp.series("release_css"));
    gulp.watch(scripts.bhojpurAdmin, gulp.series("release_js"));

    watch_bhojpur.on("change", function(path, stats) {
      console.log(":==> File " + path + " was changed, running tasks...");
    });
    watch_js.on("change", function(path, stats) {
      console.log(":==> File " + path + " was changed, running tasks...");
    });
    watch_css.on("change", function(path, stats) {
      console.log(":==> File " + path + " was changed, running tasks...");
    });
  });

  gulp.task("default", gulp.series("watch"));
}

// -----------------------------------------------------------------------------
// Other the Bhojpur CMS modules
// Command: gulp [task] --moduleName--subModuleName
//
// For example:
// Watch Worker module: gulp --worker
//
// if the Bhojpur CMS module's assets just as normal path:
// moduleName/views/themes/moduleName/assets/javascripts(stylesheets)
// just use gulp --worker
//
// if the Bhojpur CMS module's assets in enterprise as normal path:
// moduleName/views/themes/moduleName/assets/javascripts(stylesheets)
// just use gulp --microsite--enterprise
//
// if the Bhojpur CMS module's assets path as Administrator module:
// moduleName/views/assets/javascripts(stylesheets)
// you need set subModuleName as admin
// gulp --worker--admin
//
// if you need run task for subModule in modules
// example: worker module inline_edit subModule:
// gulp --worker--inline_edit
//
// gulp --media--media_library--true
//
// -----------------------------------------------------------------------------

function moduleTasks(moduleNames) {
  let moduleName = moduleNames.name,
    subModuleName = moduleNames.subName,
    useSubName = moduleNames.useSubName;

  let pathto = function(file) {
    if (moduleName && subModuleName) {
      if (subModuleName == "admin") {
        return "pkg/" + moduleName + "/views/assets/" + file;
      } else if (subModuleName == "enterprise") {
        return (
          "app.bhojpur.net/" +
          moduleName +
          "/views/themes/" +
          moduleName +
          "/assets/" +
          file
        );
      } else if (useSubName) {
        if (useSubName == "admin") {
          return (
            "pkg/" + moduleName + "/" + subModuleName + "/views/assets/" + file
          );
        } else {
          return (
            "pkg/" +
            moduleName +
            "/" +
            subModuleName +
            "/views/themes/" +
            subModuleName +
            "/assets/" +
            file
          );
        }
      } else {
        return (
          "pkg/" +
          moduleName +
          "/" +
          subModuleName +
          "/views/themes/" +
          moduleName +
          "/assets/" +
          file
        );
      }
    }
    return (
      "pkg/" + moduleName + "/views/themes/" + moduleName + "/assets/" + file
    );
  };

  let scripts = {
    src: pathto("javascripts/"),
    watch: pathto("javascripts/**/*.js"),
  };
  let styles = {
    src: pathto("stylesheets/"),
    watch: pathto("stylesheets/**/*.scss"),
  };

  function getFolders(dir) {
    return fs.readdirSync(dir).filter(function(file) {
      return fs.statSync(path.join(dir, file)).isDirectory();
    });
  }

  gulp.task("js", function() {
    let scriptPath = scripts.src,
      folders = getFolders(scriptPath);

    let task = folders.map(function(folder) {
      return gulp
        .src(path.join(scriptPath, folder, "/*.js"))
        .pipe(plumber())
        .pipe(
          eslint({
            configFile: ".eslintrc",
          })
        )
        .pipe(
          babel({
            presets: ["@babel/env"],
          })
        )
        .pipe(eslint.format())
        .pipe(concat(folder + ".js"))
        .pipe(uglify())
        .pipe(gulp.dest(scriptPath));
    });

    return es.concat.apply(null, task);
  });

  gulp.task("css", function() {
    let stylePath = styles.src,
      folders = getFolders(stylePath);

    let task = folders.map(function(folder) {
      return gulp
        .src(path.join(stylePath, folder, "/*.scss"))
        .pipe(plumber())
        .pipe(
          plugins
            .sass({
              outputStyle: "compressed",
            })
            .on("error", sass.logError)
        )
        .pipe(cleanCSS())
        .pipe(rename(folder + ".css"))
        .pipe(gulp.dest(stylePath));
    });

    return es.concat.apply(null, task);
  });

  gulp.task("watch", function() {
    let moduleScript = gulp.watch(
      scripts.watch,
      { debounceDelay: 2000 },
      gulp.series("js")
    );
    gulp.watch(styles.watch, gulp.series("css"));

    moduleScript.on("change", function(event) {
      console.log(
        ":==> File " + event.path + " was " + event.type + ", running tasks..."
      );
    });
  });

  gulp.task("default", gulp.series("watch"));
  gulp.task("release", gulp.series("js", "css"));
}

// Init
// -----------------------------------------------------------------------------

if (moduleName.name) {
  let taskPath =
      moduleName.name + "/views/themes/" + moduleName.name + "/assets/",
    runModuleName =
      'Running "' + moduleName.name + '" module task in "' + taskPath + '"...';

  if (moduleName.subName) {
    if (moduleName.subName == "admin") {
      taskPath = moduleName.name + "/views/assets/";
      runModuleName =
        'Running "' +
        moduleName.name +
        '" module task in "' +
        taskPath +
        '"...';
    } else if (moduleName.subName == "enterprise") {
      taskPath =
        "app.bhojpur.net/" +
        moduleName.name +
        "/views/themes/" +
        moduleName.name +
        "/assets/";
      runModuleName =
        'Running "' +
        moduleName.name +
        '" module task in "' +
        taskPath +
        '"...';
    } else if (moduleName.useSubName) {
      if (moduleName.useSubName == "admin") {
        taskPath =
          moduleName.name + "/" + moduleName.subName + "/views/assets/";
      } else {
        taskPath =
          moduleName.name +
          "/" +
          moduleName.subName +
          "/views/themes/" +
          moduleName.subName +
          "/assets/";
      }

      runModuleName =
        'Running "' +
        moduleName.name +
        " > " +
        moduleName.subName +
        '" module task in "' +
        taskPath +
        '"...';
    } else {
      taskPath =
        moduleName.name +
        "/" +
        moduleName.subName +
        "/views/themes/" +
        moduleName.name +
        "/assets/";
      runModuleName =
        'Running "' +
        moduleName.name +
        " > " +
        moduleName.subName +
        '" module task in "' +
        taskPath +
        '"...';
    }
  }
  console.log(runModuleName);
  moduleTasks(moduleName);
} else {
  console.log('Running "admin" module task in "pkg/admin/views/assets/"...');
  adminTasks();
}

// Task for compress js and css vendor assets
gulp.task("combineJavaScriptVendor", function() {
  return gulp
    .src([
      "!pkg/admin/views/assets/javascripts/vendors/jquery.min.js",
      "pkg/admin/views/assets/javascripts/vendors/*.js",
    ])
    .pipe(concat("vendors.js"))
    .pipe(gulp.dest("pkg/admin/views/assets/javascripts"));
});

gulp.task("compressCSSVendor", function() {
  return gulp
    .src("pkg/admin/views/assets/stylesheets/vendors/*.css")
    .pipe(concat("vendors.css"))
    .pipe(gulp.dest("pkg/admin/views/assets/stylesheets"));
});

gulp.task("combineDatetimePicker", function() {
  return gulp
    .src([
      "pkg/admin/views/assets/javascripts/bhojpur/bhojpur-config.js",
      "pkg/admin/views/assets/javascripts/bhojpur/bhojpur-material.js",
      "pkg/admin/views/assets/javascripts/bhojpur/bhojpur-modal.js",
      "pkg/admin/views/assets/javascripts/bhojpur/datepicker.js",
      "pkg/admin/views/assets/javascripts/bhojpur/bhojpur-datepicker.js",
      "pkg/admin/views/assets/javascripts/bhojpur/bhojpur-timepicker.js",
    ])
    .pipe(plumber())
    .pipe(
      eslint({
        configFile: ".eslintrc",
      })
    )
    .pipe(
      babel({
        presets: ["@babel/env"],
      })
    )

    .pipe(concat("datetimepicker.js"))
    .pipe(uglify())
    .pipe(gulp.dest("pkg/admin/views/assets/javascripts"));
});