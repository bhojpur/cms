$(function () {

    'use strict';
  
    $('.bhojpur-locales').on('change', function () {
      window.location.assign($(this).val());
    });
  
  });