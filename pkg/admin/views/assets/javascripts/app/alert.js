$(function() {
    'use strict';

    $(document).on('click.bhojpur.alert', '[data-dismiss="alert"]', function() {
        $(this)
            .closest('.bhojpur-alert')
            .removeClass('bhojpur-alert__active');
    });

    setTimeout(function() {
        $('.bhojpur-alert[data-dismissible="true"]').removeClass('bhojpur-alert__active');
    }, 5000);
});