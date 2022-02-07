# Bhojpur CMS - Content Management System

The Bhojpur CMS is a high performance content management system and web application development framework applied within the [Bhojpur.NET Platform](https://github.com/bhojpur/platform) ecosystem for building distributed enterprise applications, digital content authoring, and publishing. It features rich web user interfaces for System Administrator. It generates and delivers highly scalable, web applications or services, hosted using a W3C standards compliant HTTP server (e.g., [Bhojpur Web](https://github.com/bhojpur/web) for application delivery).

## Core Modules

* [Admin](https://github.com/bhojpur/cms/pkg/admin) - It is a core component of the [Bhojpur CMS](https://github.com/bhojpur/cms) that generates a system administrator's user interface and RESTful API for you to manage enterprise data

* [Publish](https://github.com/bhojpur/cms/pkg/publish) - It provide a staging environment for all digital content changes to be reviewed before being published to the live systems

* [Transition](https://github.com/bhojpur/cms/pkg/transition) - A configurable `Finite State Machine`: define states, events (e.g., pay order), and validation constraints for state transitions

* [Media Library](https://github.com/bhojpur/cms/pkg/media/media_library) - A digital `Asset Management` with support for several Cloud storage backends and publishing system via content delivery network

* [Worker](https://github.com/bhojpur/cms/pkg/worker) - A feature rich batch job processing scheduler

* [Exchange](https://github.com/bhojpur/cms/pkg/exchange) - A data exchange with other business applications using CSV or Excel data formats

* [Internationalization](https://github.com/bhojpur/cms/pkg/i18n) The `I18N` framework is used for managing and inline editing of natural language translations

* [Localization](https://github.com/bhojpur/cms/pkg/l10n) The `L10N` is ued for managing database-backed models on per-locale basis, with support for defining/editing localizable attributes, and locale-based querying

* [Roles](https://github.com/bhojpur/application/pkg/roles) - An application level Access Control mechanism

* and, learn more [https://github.com/bhojpur](https://github.com/bhojpur)

## Command Line Interface

You can access the following server-side [Bhojpur CMS](https://github.com/bhojpur/cms) tools

```bash
cmssvr [options]
```

You can access the following client-side [Bhojpur CMS](https://github.com/bhojpur/cms) tools

```bash
cmsctl [options]
```

## Web Frontend Development

It requires [Node.js](https://nodejs.org/) and [Gulp](http://gulpjs.com/) for building web frontend files

```bash
npm install && npm install -g gulp
```

- To watch SCSS/JavaScript changes: `gulp`
- To build release files: `gulp release`

## Web Backend Development

It requires [Go](https://go.dev/) >= 1.17 programming language tools to build web application backend files. Each web application requires a set of template files to be able to generate required web user interfaces. You must generate these web interface using following command. For example

```bash
cmsctl template internal/demo
```

Then, the compilable Go source code must integrated within your web application.

## License

Released under the [MIT License](http://opensource.org/licenses/MIT).