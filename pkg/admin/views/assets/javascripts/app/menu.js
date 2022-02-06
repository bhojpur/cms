$(function() {
    'use strict';
  
    let menuDatas = [],
      storageName = 'bhojpuradmin_menu_status',
      lastMenuStatus = localStorage.getItem(storageName);
  
    if (lastMenuStatus && lastMenuStatus.length) {
      menuDatas = lastMenuStatus.split(',');
    }
  
    $('.bhojpur-menu-container')
      .on('click', '> ul > li > a', function() {
        let $this = $(this),
          $li = $this.parent(),
          $ul = $this.next('ul'),
          menuName = $li.attr('bhojpur-icon-name');
  
        if (!$ul.length) {
          return;
        }
  
        if ($ul.hasClass('in')) {
          menuDatas.push(menuName);
  
          $li.removeClass('is-expanded');
          $ul
            .one('transitionend', function() {
              $ul.removeClass('collapsing in');
            })
            .addClass('collapsing')
            .height(0);
        } else {
          menuDatas = _.without(menuDatas, menuName);
  
          $li.addClass('is-expanded');
          $ul
            .one('transitionend', function() {
              $ul.removeClass('collapsing');
            })
            .addClass('collapsing in')
            .height($ul.prop('scrollHeight'));
        }
        localStorage.setItem(storageName, menuDatas);
      })
      .find('> ul > li > a')
      .each(function() {
        let $this = $(this),
          $li = $this.parent(),
          $ul = $this.next('ul'),
          menuName = $li.attr('bhojpur-icon-name');
  
        if (!$ul.length) {
          return;
        }
  
        $ul.addClass('collapse');
        $li.addClass('has-menu');
  
        if (menuDatas.indexOf(menuName) != -1) {
          $ul.height(0);
        } else {
          $li.addClass('is-expanded');
          $ul.addClass('in').height($ul.prop('scrollHeight'));
        }
      });
  
    let $pageHeader = $('.bhojpur-page > .bhojpur-page__header'),
      $pageBody = $('.bhojpur-page > .bhojpur-page__body'),
      triggerHeight = $pageHeader.find('.bhojpur-page-subnav__header').length ? 96 : 48;
  
    if ($pageHeader.length) {
      if ($pageHeader.height() > triggerHeight) {
        $pageBody.css('padding-top', $pageHeader.height());
      }
  
      $('.bhojpur-page').addClass('has-header');
      $('header.mdl-layout__header').addClass('has-action');
    }
  });