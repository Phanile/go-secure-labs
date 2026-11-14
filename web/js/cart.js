import { requestApi } from "./api.js";
import { AuthState, getApiErrorMessage } from "./auth.js";

const KEY = "appleshop.cart.v1";
const content = document.getElementById("cart-content");
const esc = (value) => String(value ?? "").replace(/[&<>"']/g, (c) => ({ "&": "&amp;", "<": "&lt;", ">": "&gt;", '"': "&quot;", "'": "&#039;" })[c]);
const price = (n) => new Intl.NumberFormat("ru-RU", { style: "currency", currency: "RUB", maximumFractionDigits: 2 }).format(Number(n || 0));
function readCart() { try { const value = JSON.parse(localStorage.getItem(KEY) || "[]"); return Array.isArray(value) ? value : []; } catch { return []; } }
function saveCart(items) { localStorage.setItem(KEY, JSON.stringify(items)); updateCartCount(); }
export function addToCart(product) {
    const cart = readCart(); const item = cart.find((entry) => entry.id === product.id);
    if (item) item.quantity = Math.min(item.quantity + 1, Number(product.stock_quantity) || 99);
    else cart.push({ id: product.id, name: product.name, price: product.price, image_url: product.image_url, quantity: 1, stock_quantity: product.stock_quantity });
    saveCart(cart); window.dispatchEvent(new Event("cart-updated"));
}
export function updateCartCount() { const count = readCart().reduce((sum, item) => sum + item.quantity, 0); const node = document.getElementById("cart-count"); if (node) node.textContent = String(count); }
function render() {
    const cart = readCart();
    if (!cart.length) { content.innerHTML = `<div class="user-page"><p class="eyebrow">AppleShop</p><h1>Корзина</h1><div class="empty-state">Корзина пуста. <a href="#/store">Перейти в магазин</a></div></div>`; return; }
    const total = cart.reduce((sum, item) => sum + Number(item.price) * item.quantity, 0);
    content.innerHTML = `<div class="user-page"><p class="eyebrow">AppleShop</p><h1>Корзина</h1><div class="order-items">${cart.map((item) => `<div class="order-item"><span><strong>${esc(item.name)}</strong><br>${price(item.price)} × <button type="button" data-dec="${esc(item.id)}" aria-label="Уменьшить количество">−</button> ${item.quantity} <button type="button" data-inc="${esc(item.id)}" aria-label="Увеличить количество">+</button> <button type="button" data-remove="${esc(item.id)}">Удалить</button></span><strong>${price(Number(item.price) * item.quantity)}</strong></div>`).join("")}</div><p class="cart-total">Итого: <strong>${price(total)}</strong></p><form id="checkout-form" class="account-form"><h2>Оформление заказа</h2><label>Страна<input name="country" required value="Россия"></label><label>Город<input name="city" required></label><label>Улица<input name="street" required></label><label>Дом<input name="house" required></label><label>Квартира<input name="apartment"></label><label>Индекс<input name="postal_code" autocomplete="postal-code"></label><div id="checkout-message" class="form-message hidden" role="alert"></div><button class="primary-button" type="submit">Оформить заказ</button></form></div>`;
    content.querySelectorAll("[data-inc],[data-dec],[data-remove]").forEach((button) => button.addEventListener("click", () => { const id = Number(button.dataset.inc ?? button.dataset.dec ?? button.dataset.remove); const next = readCart().map((item) => item.id === id ? { ...item, quantity: item.quantity + (button.dataset.inc ? 1 : -1) } : item).filter((item) => item.id !== id || (!button.dataset.remove && item.quantity > 0)); saveCart(next); render(); }));
    content.querySelector("#checkout-form").addEventListener("submit", async (event) => { event.preventDefault(); const form = event.currentTarget; const button = form.querySelector("button"); const message = content.querySelector("#checkout-message"); if (!AuthState.current.authenticated) { window.location.hash = "#/login"; return; } button.disabled = true; try { const address = Object.fromEntries(["country", "city", "street", "house", "apartment", "postal_code"].map((name) => [name, form.elements[name].value.trim()])); const result = await requestApi("/orders", { method: "POST", body: JSON.stringify({ address, items: readCart().map(({ id, quantity }) => ({ product_id: id, quantity })) }), headers: { "Content-Type": "application/json" } }, { modeAware: false }); saveCart([]); window.location.hash = `#/orders/${result.body.id}`; } catch (error) { message.textContent = getApiErrorMessage(error); message.className = "form-message form-message--error"; } finally { button.disabled = false; } });
}
export function renderCart() { if (content) render(); }
window.addEventListener("cart-updated", updateCartCount); updateCartCount();
