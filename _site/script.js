if(window.CSS&&CSS.supports("color","var(--primary)")){var toggleColorMode=function toggleColorMode(e){if(e.currentTarget.classList.contains("light--hidden")){document.documentElement.setAttribute("color-mode","light");localStorage.setItem("color-mode","light");}else{document.documentElement.setAttribute("color-mode","dark");localStorage.setItem("color-mode","dark");}};var toggleColorButtons=document.querySelectorAll(".color-mode__btn");toggleColorButtons.forEach(function(btn){btn.addEventListener("click",toggleColorMode);});}else{var btnContainer=document.querySelector(".color-mode__header");btnContainer.style.display="none";}
const summaries=document.querySelectorAll('summary');summaries.forEach((summary)=>{summary.addEventListener('click',startDetailsImageLoad);summary.addEventListener('mouseover',startDetailsImageLoad);summary.addEventListener('click',detailsExpand);});function detailsExpand(){summaries.forEach((summary)=>{let detail=summary.parentNode;if(detail!=this.parentNode){detail.removeAttribute('open');}});var el=this.closest("summary");el.scrollIntoView(true);}
function startDetailsImageLoad(){var details=this.closest("details");var imgs=details.querySelectorAll("img");imgs.forEach(function(i){i.removeAttribute("loading");});}