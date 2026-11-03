import { AppMode, MODES } from "./mode.js";
import { loadProducts } from "./products.js";
import { renderLabs } from "./labs.js";
import { renderPortSwigger } from "./portswigger.js";
import { AuthState } from "./auth.js";
import { renderUserRoute } from "./account.js";

const modeToggle = document.getElementById("mode-toggle");
const modeLabel = document.getElementById("mode-label");
const searchInput = document.getElementById("products-search");
const searchButton = document.getElementById("products-search-button");
const views = [...document.querySelectorAll("[data-view]")];
const routeLinks = [...document.querySelectorAll("[data-route-link]")];
const guestAuthLinks = document.getElementById("guest-auth-links");
const userAuthLinks = document.getElementById("user-auth-links");
const logoutButton = document.getElementById("logout-button");

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

function updateAuthNavigation(authState) {
    guestAuthLinks.classList.toggle("hidden", authState.authenticated);
    userAuthLinks.classList.toggle("hidden", !authState.authenticated);
}

AuthState.subscribe((authState) => {
    updateAuthNavigation(authState);
    const route = getRoute();
    if (["login", "register", "account", "orders", "order", "balance"].includes(route.section)) renderUserRoute(route);
});

logoutButton.addEventListener("click", async () => {
    logoutButton.disabled = true;
    try {
        await AuthState.logout();
        window.location.hash = "#/store";
    } catch (error) {
        console.error("Logout failed", error);
    } finally {
        logoutButton.disabled = false;
    }
});

updateModeControls(AppMode.current);
updateAuthNavigation(AuthState.current);

function getRoute() {
    const parts = window.location.hash.replace(/^#\/?/, "").split("/").filter(Boolean);
    const rawSection = parts[0] || "store";
    return {
        section: rawSection === "orders" && parts[1] ? "order" : rawSection,
        labId: parts[1]?.toUpperCase() || null,
        topicSlug: parts[1]?.toLowerCase() || null,
        id: rawSection === "orders" && parts[1] ? parts[1] : parts[2] || null,
    };
}

function navigate() {
    const route = getRoute();
    const userSections = ["login", "register", "account", "orders", "order", "balance"];
    const section = route.section === "products" ? "products" : route.section === "labs" ? "labs" : route.section === "portswigger" ? "portswigger" : route.section === "about" ? "about" : userSections.includes(route.section) ? route.section : "store";
    views.forEach((view) => view.classList.toggle("hidden", !view.dataset.view.split(" ").includes(section)));
    routeLinks.forEach((link) => link.classList.toggle("nav__link--active", link.dataset.routeLink === section || (section === "products" && link.dataset.routeLink === "store") || (section === "order" && link.dataset.routeLink === "orders")));
    const showSearch = section === "store" || section === "products";
    document.querySelector(".products-search").classList.toggle("hidden", !showSearch);
    if (section === "labs") renderLabs(route.labId);
    if (section === "portswigger") renderPortSwigger(route.topicSlug);
    if (userSections.includes(section)) renderUserRoute(route);
    if (showSearch) loadProducts(searchInput.value);
}

window.addEventListener("hashchange", navigate);
navigate();
AuthState.restoreSession();
