import { AppMode } from "./mode.js";
import { runProductSearch } from "./api.js";

const escapeHtml = (value) => String(value ?? "").replace(/[&<>"']/g, (char) => ({ "&": "&amp;", "<": "&lt;", ">": "&gt;", '"': "&quot;", "'": "&#039;" })[char]);

export function productSearchDemoMarkup({ prefix = "sql-demo" } = {}) {
    const formId = `${prefix}-form`;
    const inputId = `${prefix}-input`;
    const resultId = `${prefix}-result`;
    return `<section class="demo-panel"><div class="demo-panel__heading"><div><p class="eyebrow">Practical demonstration / live endpoint</p><h2>Product search</h2></div><span class="status status--implemented">Backend connected</span></div>
        <p>Submit a search value to the current global mode. The response below is the backend’s actual response; the frontend never labels a payload as successful unless the server responds.</p>
        <form id="${formId}" class="demo-form"><label for="${inputId}">Search input</label><div class="input-row"><input id="${inputId}" placeholder="Try a normal product search" autocomplete="off"><button class="primary-button" type="submit">Send request</button></div></form>
        <div id="${resultId}" class="response-viewer" aria-live="polite"><span class="muted">No request sent yet.</span></div>
        <div class="request-summary"><div><span class="request-summary__label">HTTP method</span><code>GET</code></div><div><span class="request-summary__label">Endpoint</span><code>${escapeHtml(AppMode.isVulnerable() ? "/api/v1/vulnerable/products?search=..." : "/api/v1/products?search=...")}</code></div></div><div class="contract-note"><span>Secure mode uses a parameterized query; vulnerable mode routes to the intentionally unsafe repository.</span></div></section>`;
}

export function bindProductSearchDemo({ prefix = "sql-demo" } = {}) {
    const form = document.getElementById(`${prefix}-form`);
    if (!form) return;
    form.addEventListener("submit", async (event) => {
        event.preventDefault();
        const input = document.getElementById(`${prefix}-input`);
        const result = document.getElementById(`${prefix}-result`);
        result.innerHTML = `<span class="loading-dot">Sending request in ${AppMode.current} mode…</span>`;
        try {
            const response = await runProductSearch(input.value);
            const success = response.status >= 200 && response.status < 300;
            result.innerHTML = `<div class="response-meta"><span class="result-indicator ${success ? "result-indicator--success" : "result-indicator--error"}">${success ? "HTTP success" : "HTTP error"}</span><span>Status ${response.status}</span></div><pre>${escapeHtml(JSON.stringify(response.body, null, 2))}</pre>`;
        } catch (error) {
            result.innerHTML = `<div class="result-indicator result-indicator--error">Request failed</div><pre>${escapeHtml(error.message)}</pre>`;
        }
    });
}
