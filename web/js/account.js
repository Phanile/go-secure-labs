import { requestApi } from "./api.js";
import { AuthState, getApiErrorMessage } from "./auth.js";

const content = document.getElementById("user-content");
const esc = (value) => String(value ?? "").replace(/[&<>"']/g, (char) => ({ "&": "&amp;", "<": "&lt;", ">": "&gt;", '"': "&quot;", "'": "&#039;" })[char]);
const formatDate = (value) => value ? new Intl.DateTimeFormat("ru-RU", { dateStyle: "medium" }).format(new Date(value)) : "—";
const formatMoney = (value) => new Intl.NumberFormat("ru-RU", { style: "currency", currency: "RUB", maximumFractionDigits: 2 }).format(Number(value ?? 0));

function shell(title, subtitle, body) {
    return `<div class="user-page"><p class="eyebrow">AppleShop / Account</p><h1>${esc(title)}</h1>${subtitle ? `<p class="lead">${esc(subtitle)}</p>` : ""}${body}</div>`;
}

function authForm(kind) {
    const register = kind === "register";
    return shell(register ? "Регистрация" : "Войти", register ? "Создайте аккаунт AppleShop." : "Войдите, чтобы открыть аккаунт и заказы.", `<div class="auth-card"><form id="${kind}-form" class="account-form" novalidate>
        <div id="${kind}-message" class="form-message hidden" role="alert"></div>
        <label for="${kind}-email">Email<input id="${kind}-email" name="email" type="email" autocomplete="email" required></label>
        <label for="${kind}-password">Пароль<input id="${kind}-password" name="password" type="password" autocomplete="${register ? "new-password" : "current-password"}" required></label>
        ${register ? `<label for="${kind}-confirmation">Подтверждение пароля<input id="${kind}-confirmation" name="password_confirmation" type="password" autocomplete="new-password" required></label>` : ""}
        <button class="primary-button" type="submit">${register ? "Зарегистрироваться" : "Войти"}</button>
    </form><p class="account-form__switch">${register ? "Уже есть аккаунт?" : "Нет аккаунта?"} <a href="#/${register ? "login" : "register"}">${register ? "Войти" : "Зарегистрироваться"}</a></p></div>`);
}

function showMessage(id, message, error = true) {
    const element = document.getElementById(id);
    element.textContent = message;
    element.className = `form-message ${error ? "form-message--error" : "form-message--success"}`;
}

function bindAuthForm(kind) {
    const form = document.getElementById(`${kind}-form`);
    if (!form) return;
    form.addEventListener("submit", async (event) => {
        event.preventDefault();
        const email = form.elements.email.value.trim();
        const password = form.elements.password.value;
        const confirmation = form.elements.password_confirmation?.value ?? "";
        if (!email || !email.includes("@")) return showMessage(`${kind}-message`, "Введите корректный email.");
        if (!password) return showMessage(`${kind}-message`, "Введите пароль.");
        if (kind === "register" && password !== confirmation) return showMessage(`${kind}-message`, "Пароли не совпадают.");
        const button = form.querySelector("button");
        button.disabled = true;
        button.textContent = "Загрузка…";
        try {
            if (kind === "register") await AuthState.register(email, password, confirmation);
            else await AuthState.login(email, password);
            window.location.hash = "#/account";
        } catch (error) {
            showMessage(`${kind}-message`, getApiErrorMessage(error));
        } finally {
            button.disabled = false;
            button.textContent = kind === "register" ? "Зарегистрироваться" : "Войти";
        }
    });
}

function loadingState() {
    return `<div class="account-state">Загрузка данных аккаунта…</div>`;
}

function errorState(error) {
    return `<div class="account-state account-state--error"><strong>Не удалось загрузить данные</strong><p>${esc(getApiErrorMessage(error))}</p></div>`;
}

function pendingState() {
    return `<div class="account-state"><strong>Backend integration pending</strong><p>Этот раздел готов к подключению backend endpoint, но данные пока недоступны.</p></div>`;
}

function accountDashboard(data) {
    const user = data.user ?? data;
    const orders = data.recent_orders ?? data.orders ?? [];
    return shell("Мой аккаунт", "Личная информация и последние заказы.", `<div class="account-summary"><article class="summary-card"><span>Email</span><strong>${esc(user.email ?? "Доступен после интеграции")}</strong></article><article class="summary-card"><span>Баланс</span><strong>${formatMoney(data.balance ?? user.balance)}</strong></article><article class="summary-card"><span>Заказов</span><strong>${data.orders_count ?? orders.length}</strong></article></div><div class="account-section"><div class="section-header"><h2 class="section-title">Последние заказы</h2><a class="secondary-button" href="#/orders">Все заказы</a></div>${orders.length ? orderList(orders.slice(0, 5)) : `<div class="empty-state">Заказов пока нет.</div>`}</div>`);
}

function orderList(orders) {
    return `<div class="order-list">${orders.map((order) => `<a class="order-row" href="#/orders/${encodeURIComponent(order.id)}"><span><strong>#${esc(order.id)}</strong><small>${esc(formatDate(order.created_at ?? order.date))}</small></span><span><strong>${formatMoney(order.total ?? order.total_price)}</strong><small>${esc(order.status ?? "—")}</small></span></a>`).join("")}</div>`;
}

async function renderAccount() {
    content.innerHTML = shell("Мой аккаунт", "Загрузка данных аккаунта.", loadingState());
    try {
        const response = await requestApi("/account", {}, { modeAware: false });
        content.innerHTML = accountDashboard(response.body);
    } catch (error) {
        content.innerHTML = error instanceof Error && error.status === 404 ? shell("Мой аккаунт", "", pendingState()) : shell("Мой аккаунт", "", errorState(error));
    }
}

async function renderOrders() {
    content.innerHTML = shell("Мои заказы", "История заказов аккаунта.", loadingState());
    try {
        const response = await requestApi("/orders", {}, { modeAware: false });
        const orders = Array.isArray(response.body) ? response.body : response.body.orders ?? [];
        content.innerHTML = shell("Мои заказы", "История заказов аккаунта.", orders.length ? orderList(orders) : `<div class="empty-state">Заказов пока нет.</div>`);
    } catch (error) {
        content.innerHTML = error instanceof Error && error.status === 404 ? shell("Мои заказы", "", pendingState()) : shell("Мои заказы", "", errorState(error));
    }
}

async function renderOrder(orderId) {
    content.innerHTML = shell(`Заказ #${orderId}`, "Загрузка деталей заказа.", loadingState());
    try {
        const response = await requestApi(`/orders/${encodeURIComponent(orderId)}`, {}, { modeAware: false });
        const order = response.body.order ?? response.body;
        const items = order.items ?? order.products ?? [];
        content.innerHTML = shell(`Заказ #${order.id ?? orderId}`, "", `<div class="order-detail"><div class="order-detail__meta"><span>Дата<strong>${esc(formatDate(order.created_at ?? order.date))}</strong></span><span>Статус<strong>${esc(order.status ?? "—")}</strong></span><span>Итого<strong>${formatMoney(order.total ?? order.total_price)}</strong></span></div><h2>Товары</h2>${items.length ? `<div class="order-items">${items.map((item) => `<div class="order-item"><span>${esc(item.name ?? `Product #${item.product_id}`)}</span><span>${esc(item.quantity)} × ${formatMoney(item.unit_price ?? item.price)}</span></div>`).join("")}</div>` : `<div class="empty-state">В заказе нет товаров.</div>`}</div>`);
    } catch (error) {
        content.innerHTML = error instanceof Error && error.status === 404 ? shell(`Заказ #${orderId}`, "", pendingState()) : shell(`Заказ #${orderId}`, "", errorState(error));
    }
}

async function renderBalance() {
    content.innerHTML = shell("Баланс", "Текущий баланс аккаунта.", loadingState());
    try {
        const response = await requestApi("/account/balance", {}, { modeAware: false });
        const balance = response.body.balance ?? response.body.amount;
        content.innerHTML = shell("Баланс", "Текущий баланс аккаунта.", `<div class="balance-card"><span>Текущий баланс</span><strong>${formatMoney(balance)}</strong></div><p class="muted account-note">Пополнение и списание средств пока не реализуются в frontend.</p>`);
    } catch (error) {
        content.innerHTML = error instanceof Error && error.status === 404 ? shell("Баланс", "", pendingState()) : shell("Баланс", "", errorState(error));
    }
}

export function renderUserRoute(route) {
    if (!content) return;
    if (["login", "register"].includes(route.section)) {
        content.innerHTML = authForm(route.section);
        bindAuthForm(route.section);
        return;
    }
    if (!AuthState.current.authenticated) {
        content.innerHTML = shell("Требуется вход", "Войдите, чтобы открыть этот раздел.", `<div class="auth-card"><a class="primary-button" href="#/login">Войти</a></div>`);
        return;
    }
    if (route.section === "account") renderAccount();
    else if (route.section === "orders") renderOrders();
    else if (route.section === "order") renderOrder(route.id);
    else if (route.section === "balance") renderBalance();
}
