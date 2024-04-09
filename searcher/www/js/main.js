const searchInput = document.querySelector("#searchInput");
const searchBtn = document.querySelector("#searchBtn");

const placeHolder = "Search for something...";
let i = 0;
const interval = setInterval(() => {
  if (searchInput.placeholder == placeHolder) {
    clearInterval(interval);
    return;
  }
  searchInput.placeholder = searchInput.placeholder + placeHolder[i];
  i++;
}, 40);

const search = () => {
  const searchValue = searchInput.value.trim();
  if (!searchValue) {
    return;
  }
  // alert("/search?q=" + searchValue);
  window.location.href = "/search?q=" + searchValue;
};

searchInput.addEventListener("keypress", (e) => {
  if (e.key === "Enter") {
    search();
  }
});

searchBtn.addEventListener("click", search);
