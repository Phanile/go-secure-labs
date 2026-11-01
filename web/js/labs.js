import { AppMode } from "./mode.js";
import { getLab, LABS } from "./lab-data.js";
import { bindProductSearchDemo, productSearchDemoMarkup } from "./lab-demo.js";

const content = document.getElementById("labs-content");
const esc = (value) => String(value ?? "").replace(/[&<>"']/g, (char) => ({ "&": "&amp;", "<": "&lt;", ">": "&gt;", '"': "&quot;", "'": "&#039;" })[char]);

const statusText = { implemented: "Implemented", pending: "Backend integration pending" };

function labCard(lab) {
    return `<a class="lab-card" href="#/labs/${lab.id}">
        <div class="lab-card__top"><span class="lab-id">${lab.id}:2025</span><span class="status status--${lab.status}">${statusText[lab.status]}</span></div>
        <h3>${esc(lab.title)}</h3><p>${esc(lab.short)}</p>
        <div class="lab-card__meta"><span>${lab.demos} demonstration${lab.demos === 1 ? "" : "s"}</span><span>Open lab →</span></div>
    </a>`;
}

function header(lab) {
    return `<div class="lab-heading"><div><p class="eyebrow">OWASP Top 10:2025 / ${lab.id}</p><h1>${esc(lab.title)}</h1><p class="lead">${esc(lab.description)}</p></div><a class="secondary-button" href="#/labs">All labs</a></div>`;
}

function infoPanel(lab) {
    return `<div class="info-grid"><article class="info-panel"><h3>Common causes</h3><p>${esc(lab.causes)}</p></article><article class="info-panel"><h3>Potential consequences</h3><p>${esc(lab.consequences)}</p></article><article class="info-panel"><h3>Mitigation</h3><p>${esc(lab.mitigation)}</p></article></div>`;
}

function pendingDemo(lab) {
    return `<section class="demo-panel"><div class="demo-panel__heading"><div><p class="eyebrow">Practical demonstration</p><h2>Integration point</h2></div><span class="status status--pending">Backend integration pending</span></div>
        <p>This interface is prepared for the project owner’s Go implementation. It does not fabricate a result or simulate a successful exploit.</p>
        <form class="pending-form" data-pending-form><div class="field-grid">${lab.fields.map((field) => `<label>${esc(field.label)}<input name="${esc(field.name)}" placeholder="${esc(field.placeholder)}" disabled></label>`).join("")}</div><button class="primary-button" type="submit" disabled>Run when endpoint is available</button></form>
        <div class="contract-note"><code>${esc(AppMode.isVulnerable() ? `/api/v1/vulnerable${lab.endpoint}` : `/api/v1${lab.endpoint}`)}</code><span>Expected contract is documented in <code>web/docs/OWASP_API_CONTRACT.md</code>.</span></div></section>`;
}

function sqlDemo() {
    return productSearchDemoMarkup({ prefix: "owasp-sql-demo" });
}

function renderOverview() {
    content.innerHTML = `<div class="lab-overview"><div class="lab-heading"><div><p class="eyebrow">AppleShop / Security Labs</p><h1>OWASP Top 10:2025</h1><p class="lead">Изучайте реальные границы безопасности приложения в контролируемой лаборатории.</p></div><span class="mode-chip">${AppMode.isVulnerable() ? "Vulnerable mode" : "Secure mode"}</span></div>
        <div class="warning-panel"><strong>Учебное предупреждение</strong><p>Уязвимые реализации предназначены только для локального контролируемого окружения. Не публикуйте vulnerable mode в интернете и не используйте реальные персональные или конфиденциальные данные.</p></div>
        <div class="labs-grid">${LABS.map(labCard).join("")}</div></div>`;
}

function renderDetail(lab) {
    content.innerHTML = `<div class="lab-detail">${header(lab)}${infoPanel(lab)}${lab.id === "A05" ? sqlDemo() : pendingDemo(lab)}</div>`;
    if (lab.id === "A05") bindProductSearchDemo({ prefix: "owasp-sql-demo" });
}

export function renderLabs(labId = null) {
    if (!content) return;
    const lab = labId ? getLab(labId) : null;
    lab ? renderDetail(lab) : renderOverview();
}
