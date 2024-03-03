// --- Theme ---
if (window.CSS && CSS.supports("color", "var(--primary)")) {
    // Setup button
    var toggleColorMode = function toggleColorMode(e) {
      if (e.currentTarget.classList.contains("light--hidden")) {
        // Set light
        document.documentElement.setAttribute("color-mode", "light");
        localStorage.setItem("color-mode", "light");
      } else {
        // Set dark
        document.documentElement.setAttribute("color-mode", "dark");
        localStorage.setItem("color-mode", "dark");
      }
    };
    var toggleColorButtons = document.querySelectorAll(".color-mode__btn");
    toggleColorButtons.forEach(function(btn) {
      btn.addEventListener("click", toggleColorMode);
    });

  } else {
    // If the feature isn't supported, then we hide the toggle buttons
    var btnContainer = document.querySelector(".color-mode__header");
    btnContainer.style.display = "none";
  }

// --- Handle details expansion ---
const summaries = document.querySelectorAll('summary');

summaries.forEach((summary) => {
  summary.addEventListener('mouseover', startDetailsImageLoad);
  summary.addEventListener('click', detailsExpand);
});

function detailsExpand() {
  // Close others
  summaries.forEach((summary) => {
    let detail = summary.parentNode;
      if (detail != this.parentNode) {
        detail.removeAttribute('open');
      }
    });
  
  // Scroll to the summary
  var el = this.closest("summary");
  el.scrollIntoView(true);
}

function startDetailsImageLoad() {
  // Start loading any contained images
  var details = this.closest("details");
  var imgs = details.querySelectorAll("img");
  imgs.forEach(function(i) {
    i.removeAttribute("loading");
  });
}