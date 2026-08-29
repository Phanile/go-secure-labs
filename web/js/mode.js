const MODE_STORAGE_KEY = "appleshop-mode";

export const MODES = Object.freeze({
    SECURE: "secure",
    VULNERABLE: "vulnerable",
});

const listeners = new Set();
const readStoredMode = () => localStorage.getItem(MODE_STORAGE_KEY) === MODES.VULNERABLE
    ? MODES.VULNERABLE
    : MODES.SECURE;

let currentMode = readStoredMode();

function notify() {
    document.body.dataset.mode = currentMode;
    listeners.forEach((listener) => listener(currentMode));
}

export const AppMode = {
    get current() {
        return currentMode;
    },
    isVulnerable() {
        return currentMode === MODES.VULNERABLE;
    },
    set(mode) {
        if (!Object.values(MODES).includes(mode) || mode === currentMode) return;
        currentMode = mode;
        localStorage.setItem(MODE_STORAGE_KEY, currentMode);
        notify();
    },
    toggle() {
        this.set(this.isVulnerable() ? MODES.SECURE : MODES.VULNERABLE);
    },
    subscribe(listener) {
        listeners.add(listener);
        return () => listeners.delete(listener);
    },
};

document.body.dataset.mode = currentMode;
