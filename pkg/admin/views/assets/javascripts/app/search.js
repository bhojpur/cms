$(function () {

    'use strict';
  
    var location = window.location;
  
    $('.bhojpur-search').each(function () {
      var $this = $(this);
      var $input = $this.find('.bhojpur-search__input');
      var $clear = $this.find('.bhojpur-search__clear');
      var isSearched = !!$input.val();
  
      var emptySearch = function () {
        var search = location.search.replace(new RegExp($input.attr('name') + '\\=?\\w*'), '');
        if (search == '?'){
          location.href = location.href.split('?')[0];
        } else {
          location.search = location.search.replace(new RegExp($input.attr('name') + '\\=?\\w*'), '');
        }
      };
  
      $this.closest('.bhojpur-page__header').addClass('has-search');
      $('header.mdl-layout__header').addClass('has-search');
  
      $clear.on('click', function () {
        if ($input.val() || isSearched) {
          emptySearch();
        } else {
          $this.removeClass('is-dirty');
        }
      });
    });
  });