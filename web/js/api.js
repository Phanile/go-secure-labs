import { AppMode } from "./mode.js";

const API_PATH = "/api/v1";
const PRODUCTS_PATH = "/products";
const VULNERABLE_PREFIX = "/vulnerable";

export function getApiBaseUrl() {
    const modePrefix = AppMode.isVulnerable() ? VULNERABLE_PREFIX : "";
    return `${API_PATH}${modePrefix}`;
}

export async function fetchProducts(search = "") {
    const url = new URL(
        `${getApiBaseUrl()}${PRODUCTS_PATH}`,
        window.location.origin,
    );
    const normalizedSearch = search.trim();
    if (normalizedSearch !== "") url.searchParams.set("search", normalizedSearch);

    const response = await fetch(url);
    if (!response.ok) throw new Error(`Failed to fetch products: ${response.status}`);
    return response.json();
}
