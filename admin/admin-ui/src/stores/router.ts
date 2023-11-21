import { createRouter } from "@nanostores/router";
import { atom } from "nanostores";

export const $router = createRouter({
    dashboard: "/ui/dashboard",
    login: "/entry/login",
    init: "/entry/init",
    logout: "/entry/logout",
    database: "/ui/database",
    settings: "/ui/settings",
});

export const $prefetchStore = atom<any[]>([]);


export const $pageDataStore = atom<any[]>([]);

export function addPageData(data: { [key: string]: any }) {
    const temp = []
    for (const [key, value] of Object.entries(data)) {
        if (value === undefined) {
            delete data[key];
        }
    }
    for (const item of $pageDataStore.get()) {
        if (item.title === data['title']) {
            return;
        } else {
            temp.push(item);
        }
    }
    const pageData = $pageDataStore.get();
    pageData.push(data);
} 

export function getPageData(title: string): any {
    const pageData = $pageDataStore.get();
    for (const item of pageData) {
        if (item.title === title) {
            return item;
        }
    }
    return null;
}

export function updatePageData(title: string, data: { [key: string]: any }) {
    const pageData = $pageDataStore.get();
    for (const item of pageData) {
        if (item.title === title) {
            for (const [key, value] of Object.entries(data)) {
                if (value === undefined) {
                    delete data[key];
                }
            }
            Object.assign(item, data);
        }
    }
    $pageDataStore.set(pageData);
}

export function removePageData(title: string) {
    const pageData = $pageDataStore.get();
    for (const [index, item] of pageData.entries()) {
        if (item.title === title) {
            pageData.splice(index, 1);
        }
    }
    $pageDataStore.set(pageData);
}

export function clearPageData() {
    $pageDataStore.set([]);
}