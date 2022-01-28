# Bhojpur CMS - Content Management

The Bhojpur CMS is a software-as-a-service product used for content authoring and publishing based on Bhojpur.NET Platform for application delivery.
## The modules

* [Admin](https://github.com/bhojpur/cms/pkg/admin) - The core part of Bhojpur CMS that generates a system administrator's interface and RESTful API for you to manage data

* [Publish](https://github.com/bhojpur/cms/pkg/publish) - It provide a staging environment for all content changes to be reviewed before being published to the live system

* [Transition](https://github.com/bhojpur/cms/pkg/transition) - A configurable Finite State Machine: define states, events (e.g., pay order), and validation constraints for state transitions

* [Media Library](https://github.com/bhojpur/cms/pkg/media/media_library) - A digital Asset Management with support for several cloud storage backends and publishing via CDN

* [Worker](https://github.com/bhojpur/cms/pkg/worker) - A batch processing scheduler

* [Exchange](https://github.com/bhojpur/cms/pkg/exchange) - Data Exchange with other business applications using CSV or Excel data

* [Internationalization](https://github.com/bhojpur/i18n/pkg/engine) I18N is used for managing and inline editing of natural language translations

* [Localization](https://github.com/bhojpur/cms/pkg/l10n) L10N is ued for managing database-backed models on per-locale basis, with support for defining/editing localizable attributes, and locale-based querying

* [Roles](https://github.com/bhojpur/application/pkg/roles) - Access Control

* and, more [https://github.com/bhojpur](https://github.com/bhojpur)

## Frontend Development

Requires [Node.js](https://nodejs.org/) and [Gulp](http://gulpjs.com/) for building frontend files

```bash
npm install && npm install -g gulp
```

- Watch SCSS/JavaScript changes: `gulp`
- Build Release files: `gulp release`

## License

Released under the [MIT License](http://opensource.org/licenses/MIT).