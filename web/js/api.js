import { AppMode } from "./mode.js";

const API_PATH = "/api/v1";
const PRODUCTS_PATH = "/products";
const VULNERABLE_PREFIX = "/vulnerable";

export function getApiBaseUrl() {
    const modePrefix = AppMode.isVulnerable() ? VULNERABLE_PREFIX : "";
    return `${API_PATH}${modePrefix}`;
}

export function getApiRoot() {
    return API_PATH;
}

export class ApiError extends Error {
    constructor(message, status, body = null) {
        super(message);
        this.name = "ApiError";
        this.status = status;
        this.body = body;
    }
}

async function readResponseBody(response) {
    const rawBody = await response.text();
    if (rawBody.trim() === "") return null;

    try {
        return JSON.parse(rawBody);
    } catch {
        return rawBody;
    }
}

export async function requestApi(path, options = {}, { modeAware = true } = {}) {
    const baseUrl = modeAware ? getApiBaseUrl() : getApiRoot();
    const response = await fetch(`${baseUrl}${path}`, {
        ...options,
        credentials: "same-origin",
        headers: { "Accept": "application/json", ...(options.headers ?? {}) },
    });
    const body = await readResponseBody(response);
    if (!response.ok) throw new ApiError(`API request failed: ${response.status}`, response.status, body);
    return { status: response.status, headers: response.headers, body };
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

export async function runProductSearch(search) {
    const url = new URL(`${getApiBaseUrl()}${PRODUCTS_PATH}`, window.location.origin);
    url.searchParams.set("search", search);
    const response = await fetch(url);
    const contentType = response.headers.get("content-type") ?? "";
    const body = contentType.includes("application/json")
        ? await response.json()
        : await response.text();
    return { status: response.status, headers: response.headers, body };
}
