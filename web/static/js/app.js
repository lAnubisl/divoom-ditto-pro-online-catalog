document.body.addEventListener('htmx:sendError', function () {
  document.getElementById('status').textContent = 'Connection to the web server was lost. Check your connection and reload the page.';
});
document.body.addEventListener('htmx:timeout', function () {
  document.getElementById('status').textContent = 'The request timed out. The result of the operation is unknown.';
});
document.body.addEventListener('htmx:responseError', function (event) {
  document.getElementById('status').textContent = event.detail.xhr.status === 401 ? 'Your session has ended. Sign in again.' : 'The request failed. Reload the page and try again.';
});

// Keep already loaded cards when a continuation request fails.
['htmx:sendError', 'htmx:timeout', 'htmx:responseError'].forEach(function (name) {
  document.body.addEventListener(name, function (event) {
    var continuation = event.detail.elt;
    if (!continuation || !continuation.matches('.animations-continuation')) return;
    continuation.querySelector('p').textContent = 'Could not load more animations.';
    continuation.querySelector('.animations-retry').hidden = false;
  });
});
document.body.addEventListener('click', function (event) {
  var retry = event.target.closest('.animations-retry');
  if (!retry) return;
  var continuation = retry.closest('.animations-continuation');
  retry.hidden = true;
  continuation.querySelector('p').textContent = 'Loading more animations…';
  htmx.trigger(continuation, 'animations-retry');
});
