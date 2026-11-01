import { AppMode, MODES } from "./mode.js";
import { loadProducts } from "./products.js";
import { renderLabs } from "./labs.js";
import { renderPortSwigger } from "./portswigger.js";

const modeToggle = document.getElementById("mode-toggle");
const modeLabel = document.getElementById("mode-label");
const searchInput = document.getElementById("products-search");
const searchButton = document.getElementById("products-search-button");
const views = [...document.querySelectorAll("[data-view]")];
const routeLinks = [...document.querySelectorAll("[data-route-link]")];

function updateModeControls(mode) {
    const vulnerable = mode === MODES.VULNERABLE;
    modeToggle.checked = vulnerable;
    modeLabel.textContent = vulnerable ? "Vulnerable" : "Secure";
    modeToggle.setAttribute("aria-label", `Switch to ${vulnerable ? "secure" : "vulnerable"} mode`);
}

modeToggle.addEventListener("change", () => AppMode.toggle());
searchButton.addEventListener("click", () => loadProducts(searchInput.value));
searchInput.addEventListener("keydown", (event) => {
    if (event.key === "Enter") loadProducts(searchInput.value);
});

AppMode.subscribe((mode) => {
    updateModeControls(mode);
    loadProducts(searchInput.value);
    if (getRoute().section === "labs") renderLabs(getRoute().labId);
    if (getRoute().section === "portswigger") renderPortSwigger(getRoute().topicSlug);
});

updateModeControls(AppMode.current);

function getRoute() {
    const parts = window.location.hash.replace(/^#\/?/, "").split("/").filter(Boolean);
    return {
        section: parts[0] || "store",
        labId: parts[1]?.toUpperCase() || null,
        topicSlug: parts[1]?.toLowerCase() || null,
    };
}

function navigate() {
    const route = getRoute();
    const section = route.section === "products" ? "products" : route.section === "labs" ? "labs" : route.section === "portswigger" ? "portswigger" : route.section === "about" ? "about" : "store";
    views.forEach((view) => view.classList.toggle("hidden", !view.dataset.view.split(" ").includes(section)));
    routeLinks.forEach((link) => link.classList.toggle("nav__link--active", link.dataset.routeLink === section || (section === "products" && link.dataset.routeLink === "store")));
    const showSearch = section === "store" || section === "products";
    document.querySelector(".products-search").classList.toggle("hidden", !showSearch);
    if (section === "labs") renderLabs(route.labId);
    if (section === "portswigger") renderPortSwigger(route.topicSlug);
    if (showSearch) loadProducts(searchInput.value);
}

window.addEventListener("hashchange", navigate);
navigate();
