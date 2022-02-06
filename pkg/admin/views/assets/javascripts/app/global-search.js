$(function () {

    'use strict';
  
    var modal = (
      '<div class="bhojpur-dialog bhojpur-dialog--global-search" tabindex="-1" role="dialog" aria-hidden="true">' +
        '<div class="bhojpur-dialog-content">' +
          '<form action=[[actionUrl]]>' +
            '<div class="mdl-textfield mdl-js-textfield" id="global-search-textfield">' +
              '<input class="mdl-textfield__input ignore-dirtyform" name="keyword" id="globalSearch" value="" type="text" placeholder="" />' +
              '<label class="mdl-textfield__label" for="globalSearch">[[placeholder]]</label>' +
            '</div>' +
          '</form>' +
        '</div>' +
      '</div>'
    );
  
    $(document).on('click', '.bhojpur-dialog--global-search', function(e){
      e.stopPropagation();
      if (!$(e.target).parents('.bhojpur-dialog-content').length && !$(e.target).is('.bhojpur-dialog-content')){
        $('.bhojpur-dialog--global-search').remove();
      }
    });
  
    $(document).on('click', '.bhojpur-global-search--show', function(e){
        e.preventDefault();
  
        var data = $(this).data();
        var modalHTML = window.Mustache.render(modal, data);
  
        $('body').append(modalHTML);
        window.componentHandler.upgradeElement(document.getElementById('global-search-textfield'));
        $('#globalSearch').focus();
  
    });
  });