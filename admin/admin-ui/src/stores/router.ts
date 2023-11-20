import { createRouter } from "@nanostores/router";
import { persistentAtom } from "@nanostores/persistent";

export const $router = createRouter({
    dashboard: "/ui/dashboard",
    login: "/entry/login",
    init: "/entry/init",
    logout: "/entry/logout",
    database: "/ui/database",
    settings: "/ui/settings",
});