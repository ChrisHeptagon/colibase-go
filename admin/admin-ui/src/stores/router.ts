import { createRouter } from "@nanostores/router";

export const $router = createRouter({
    dashboard: "/ui/dashboard",
    login: "/entry/login",
    init: "/entry/init",
    logout: "/entry/logout",
    database: "/ui/database",
    settings: "/ui/settings",
});