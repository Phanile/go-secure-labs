import { AppMode, MODES } from "./mode.js";
import { loadProducts } from "./products.js";

const modeToggle = document.getElementById("mode-toggle");
const modeLabel = document.getElementById("mode-label");
const searchInput = document.getElementById("products-search");
const searchButton = document.getElementById("products-search-button");

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
});

updateModeControls(AppMode.current);
loadProducts();
