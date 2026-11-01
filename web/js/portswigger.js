import { AppMode } from "./mode.js";
import { getPortSwiggerTopic, PORTSWIGGER_TOPICS } from "./portswigger-data.js";
import { bindProductSearchDemo, productSearchDemoMarkup } from "./lab-demo.js";

const content = document.getElementById("portswigger-content");
const esc = (value) => String(value ?? "").replace(/[&<>"']/g, (char) => ({ "&": "&amp;", "<": "&lt;", ">": "&gt;", '"': "&quot;", "'": "&#039;" })[char]);

function topicCard(topic) {
    return `<a class="lab-card portswigger-card" href="#/portswigger/${topic.slug}">
        <div class="lab-card__top"><span class="lab-id">PortSwigger topic</span><span class="status status--implemented">${topic.labs} Labs</span></div>
        <h3>${esc(topic.title)}</h3>
        <p class="portswigger-card__description">${esc(topic.description)}</p>
        <div class="lab-card__meta"><span>${topic.labs} hands-on labs</span><span class="secondary-button">Go to topic</span></div>
    </a>`;
}

function renderTopics(filter = "") {
    const query = filter.trim().toLowerCase();
    const topics = PORTSWIGGER_TOPICS.filter((topic) => `${topic.title} ${topic.description}`.toLowerCase().includes(query));
    const grid = document.getElementById("portswigger-grid");
    const count = document.getElementById("portswigger-count");
    const empty = document.getElementById("portswigger-empty");
    grid.innerHTML = topics.map(topicCard).join("");
    count.textContent = `${topics.length} ${topics.length === 1 ? "topic" : "topics"}`;
    empty.classList.toggle("hidden", topics.length !== 0);
}

function topicHeader(topic) {
    return `<div class="lab-heading"><div><p class="eyebrow">PortSwigger Web Security Academy / ${esc(topic.category)}</p><h1>${esc(topic.title)}</h1><p class="lead">${esc(topic.description)}</p></div><a class="secondary-button" href="#/portswigger">All PortSwigger topics</a></div>`;
}

function topicInformation(topic) {
    return `<div class="info-grid"><article class="info-panel"><h3>How it works</h3><p>${esc(topic.howItWorks)}</p></article><article class="info-panel"><h3>Potential consequences</h3><p>${esc(topic.impact)}</p></article><article class="info-panel"><h3>Typical mitigation</h3><p>${esc(topic.mitigation)}</p></article></div>`;
}

function requestPanel(topic) {
    const endpoint = AppMode.isVulnerable() ? `/api/v1/vulnerable${topic.endpoint}` : `/api/v1${topic.endpoint}`;
    return `<section class="demo-panel"><div class="demo-panel__heading"><div><p class="eyebrow">Practical demonstration</p><h2>Backend integration point</h2></div><span class="status status--pending">Backend Integration Pending</span></div>
        <p>The page is ready for the project owner’s Go implementation. Requests are disabled until the endpoint exists; no exploit result or HTTP response is simulated.</p>
        <div class="request-summary"><div><span class="request-summary__label">HTTP method</span><code>${esc(topic.method)}</code></div><div><span class="request-summary__label">Endpoint</span><code>${esc(endpoint)}</code></div></div>
        <div class="request-example"><span class="request-summary__label">Example request</span><pre>${esc(topic.requestExample)}</pre></div>
        <form class="pending-form"><label for="pending-request-editor">Request editor</label><textarea id="pending-request-editor" class="request-editor" disabled>${esc(topic.requestExample)}</textarea><button class="primary-button" type="submit" disabled>Send Request when endpoint is available</button></form>
        <div class="response-viewer"><span class="muted">HTTP status and response body will appear after backend integration.</span></div>
        <div class="contract-note"><span>Expected contract is documented in <code>web/docs/OWASP_API_CONTRACT.md</code>.</span></div></section>`;
}

function renderDetail(topic) {
    const demo = topic.demo === "product-search"
        ? productSearchDemoMarkup({ prefix: "portswigger-sql-demo" })
        : requestPanel(topic);
    content.innerHTML = `<div class="lab-detail">${topicHeader(topic)}<div class="topic-meta"><span class="status status--implemented">${topic.labs} official Labs</span><span class="topic-meta__category">Category: ${esc(topic.category)}</span></div>${topicInformation(topic)}${demo}<div class="official-resource"><div><p class="eyebrow">Official PortSwigger Resources</p><h2>Continue your study</h2></div><a class="primary-button" href="${esc(topic.url)}" target="_blank" rel="noopener noreferrer">Open PortSwigger Academy ↗</a></div></div>`;
    if (topic.demo === "product-search") bindProductSearchDemo({ prefix: "portswigger-sql-demo" });
}

export function renderPortSwigger(slug = null) {
    if (!content) return;
    const topic = slug ? getPortSwiggerTopic(slug) : null;
    if (topic) {
        renderDetail(topic);
        return;
    }
    content.innerHTML = `<div class="lab-overview portswigger-overview">
        <div class="lab-heading"><div><p class="eyebrow">External practice / PortSwigger Labs</p><h1>PortSwigger Web Security Academy</h1><p class="lead">Explore web security topics, understand vulnerability mechanics and practice with hands-on labs.</p></div><a class="secondary-button" href="https://portswigger.net/web-security" target="_blank" rel="noopener noreferrer">Official website</a></div>
        <div class="portswigger-toolbar"><label for="portswigger-search">Search topics</label><input id="portswigger-search" type="search" placeholder="Search by topic..." autocomplete="off"><span id="portswigger-count" class="products-count"></span></div>
        <div id="portswigger-grid" class="labs-grid"></div><p id="portswigger-empty" class="products-state hidden">No topics match your search.</p>
    </div>`;
    const search = document.getElementById("portswigger-search");
    search.addEventListener("input", () => renderTopics(search.value));
    renderTopics();
}
