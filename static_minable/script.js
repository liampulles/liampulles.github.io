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

// --- Insert maybe pages into 404 page ---
function levenshteinDistance(s, t) {
  if (!s.length) return t.length;
  if (!t.length) return s.length;
  const arr = [];
  for (let i = 0; i <= t.length; i++) {
    arr[i] = [i];
    for (let j = 1; j <= s.length; j++) {
      arr[i][j] =
        i === 0
          ? j
          : Math.min(
              arr[i - 1][j] + 1,
              arr[i][j - 1] + 1,
              arr[i - 1][j - 1] + (s[j - 1] === t[i - 1] ? 0 : 1)
            );
    }
  }
  return arr[t.length][s.length];
};

function short(url) {
  const elem = url.split("/")
  const last = elem[elem.length-1]
  const lastElem = last.split(".")
  return lastElem[0]
}

function sortMaybesSimilar() {
  const currentShort = short(window.location.href)
  var copy = [...maybe_pages]
  copy.sort((a,b) => levenshteinDistance(currentShort,a.location) - levenshteinDistance(currentShort,b.location))
  return copy
}

document.getElementById('maybePages').innerHTML = "<ul>"+sortMaybesSimilar().map(p => `<li><a href=${p.location}>${p.title}</a></li>`).join("\n")+"</ul>"