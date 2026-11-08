import { ApiError, requestApi } from "./api.js";

const listeners = new Set();
const AUTH_ENDPOINTS = Object.freeze({
    login: "/login",
    register: "/register",
    logout: "/logout",
    currentUser: "/me",
});

let state = { authenticated: false, user: null, loading: false };

function notify() {
    listeners.forEach((listener) => listener(state));
}

function normalizeUser(body) {
    if (!body || typeof body !== "object") return null;
    return body.user ?? body.data?.user ?? body.data ?? body;
}

export const AuthState = {
    get current() { return state; },
    subscribe(listener) {
        listeners.add(listener);
        return () => listeners.delete(listener);
    },
    async login(email, password) {
        state = { ...state, loading: true };
        notify();
        try {
            const response = await requestApi(AUTH_ENDPOINTS.login, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ email, password }),
            }, { modeAware: false });
            state = { authenticated: true, user: normalizeUser(response.body), loading: false };
            notify();
            return response;
        } catch (error) {
            state = { ...state, loading: false };
            notify();
            throw error;
        }
    },
    async register(email, password, passwordConfirmation) {
        state = { ...state, loading: true };
        notify();
        try {
            const response = await requestApi(AUTH_ENDPOINTS.register, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ email, password, password_confirmation: passwordConfirmation }),
            }, { modeAware: false });
            state = { authenticated: true, user: normalizeUser(response.body), loading: false };
            notify();
            return response;
        } catch (error) {
            state = { ...state, loading: false };
            notify();
            throw error;
        }
    },
    async logout() {
        await requestApi(AUTH_ENDPOINTS.logout, { method: "POST" }, { modeAware: false });
        state = { authenticated: false, user: null, loading: false };
        notify();
    },
    async restoreSession() {
        try {
            const response = await requestApi(AUTH_ENDPOINTS.currentUser, {}, { modeAware: false });
            state = { authenticated: true, user: normalizeUser(response.body), loading: false };
            notify();
        } catch (error) {
            if (!(error instanceof ApiError) || ![401, 404].includes(error.status)) console.warn("Unable to restore session", error);
            state = { authenticated: false, user: null, loading: false };
            notify();
        }
    },
};

export function getApiErrorMessage(error) {
    if (error instanceof ApiError && error.status === 404) return "Этот backend endpoint пока не реализован.";
    if (error instanceof ApiError && error.status === 401) return "Неверные учетные данные или истекшая сессия.";
    if (error instanceof ApiError && error.body && typeof error.body === "object") return error.body.message ?? error.body.error ?? "Backend отклонил запрос.";
    return error?.message ?? "Не удалось выполнить запрос.";
}
