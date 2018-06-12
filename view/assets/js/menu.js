function openNav() {
  document.getElementById('side-nav').style.width = '250px';
  document.getElementById('shy').style.opacity = '0.5';
}

function closeNav() {
  document.getElementById('side-nav').style.width = '0';
  document.getElementById('shy').style.opacity = '1';
}

$(document).on('ready', function() {
  $('.content-wrapper').on('scroll', function() {
    var winHeight = $(window).height(),
      docHeight = $('.content-wrapper').prop('scrollHeight'),
      progressBar = $('progress'),
      max,
      value;

    /* Set the max scrollable area */
    max = docHeight - winHeight;
    progressBar.attr('max', max);
    value = $('.content-wrapper').scrollTop();
    progressBar.attr('value', value);
  });

  // Disable Safari bounce scroll
  $('body').on({
    'touchstart': _onTouchStart,
    'touchmove': _onTouchMove,
    'touchend': _onTouchEnd
});

  function _onTouchStart(e) {

      e.stopPropagation();
  }

  function _onTouchMove(e) {

      e.stopPropagation();
  }

  function _onTouchEnd(e) {

      e.stopPropagation();
  }
  document.ontouchmove = function(event){
    event.preventDefault();
}

});
