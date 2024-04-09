const counterElement = document.getElementById("counter");
const resultsElement = document.getElementById("results");
const loadMoreButton = document.getElementById("loadMore");

const query = new URLSearchParams(window.location.search).get("q");

if (!query) {
  window.location.href = "/";
}

document.title = `Search results for "${query}"`;

searchInput.value = query;

let currentPage = 1;

const extractDomain = (url) => url.split("/")[2];

const faviconUrl = (url) => {
  const parts = url.split("/");
  return parts[0] + "//" + parts[2] + "/favicon.ico";
};

const onImageError = (event) => {
  console.log("error loading image", event.target.src);
  event.target.src = "/default-favicon.ico";
};

const itemSkeleton = /*html*/ `
  <div class="flex flex-col gap-4 max-w-[50rem] border-b pb-4 animate-pulse">
  <div class="flex gap-4 items-center">
    <div class="w-8 h-8 bg-gray-200 rounded-lg"></div>
    <span>
      <div class="h-4 bg-gray-200 rounded w-1/4"></div>
      <div class="h-4 bg-gray-200 rounded w-full mt-1"></div>
    </span>
  </div>
  <div class="flex flex-col gap-2">
    <div class="h-4 bg-gray-200 rounded w-full"></div>
    <div class="h-4 bg-gray-200 rounded w-2/3 mt-1"></div>
  </div>
</div>
`;

// itemSkeleton * 5
const loadingTemplate = Array(5).fill(itemSkeleton).join("");

const renderResults = (results) => {
  if (currentPage === 1) {
    resultsElement.innerHTML = "";
  }
  results.forEach(({ url, abstract, title }) => {
    const resultElement = document.createElement("div");
    resultElement.className = "result";
    resultElement.innerHTML = /*html*/ `
      <div class="flex flex-col gap-4 max-w-[50rem] border-b pb-4">
        <div class="flex gap-4 items-center">
        <img onError="this.onerror=null;this.src='images/default-favicon.png';" src="${faviconUrl(
          url
        )}" class="w-8 h-8 border rounded-lg" />
          <span>
            <p class="font-medium">${extractDomain(url)}</p>
            <p title="${url}" class="text-gray-600 truncate max-w-[35rem]">${url}</p>
          </span>
        </div>
        <div class="flex flex-col gap-2">
          <a href="${url}" class="text-lg text-gray-800 text-blue-500 cursor-pointer max-w-[50rem] truncate">${title}</a>
          <p class="text-gray-600">
           ${abstract}
          </p>
        </div>
      </div>
    `;
    resultsElement.appendChild(resultElement);
  });
};

const timeToStr = (time) => {
  if (time < 1) {
    return `less than a millisecond`;
  }

  if (time < 1000) {
    return `${time} milliseconds`;
  }

  if (time < 1000 * 60) {
    return `${(time / 1000).toFixed(1)} seconds`;
  }

  return `${(time / 1000 / 60).toFixed(2)} minutes`;
};

const performSearch = async () => {
  const searchValue = searchInput.value;
  if (searchValue?.trim()) {
    if (currentPage === 1) {
      resultsElement.innerHTML = loadingTemplate;
    } else {
      loadMoreButton.disabled = true;
      loadMoreButton.innerHTML = "Loading...";
    }

    const res = await fetch(`/api/search?q=${searchValue}&page=${currentPage}`);
    const data = await res.json();
    const { pages, totalResults, time, page, pagesCount } = data;
    renderResults(pages);
    counterElement.innerHTML = `Found ${totalResults.toLocaleString()} results in ${timeToStr(
      time
    )}`;

    if (page >= pagesCount) {
      loadMoreButton.style.display = "none";
    }
  }
};

performSearch(1);

loadMoreButton.addEventListener("click", async () => {
  currentPage++;
  performSearch(currentPage);
});
