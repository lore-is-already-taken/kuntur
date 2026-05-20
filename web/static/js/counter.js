(function () {
  var el = document.getElementById("counter");
  if (!el) return;
  var counter = 0;
  var set = function (n) {
    counter = n;
    el.innerHTML = "count is " + counter;
  };
  el.addEventListener("click", function () { set(counter + 1); });
  set(0);
})();
