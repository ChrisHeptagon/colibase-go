import { atom } from "nanostores";



export const $callbackStore = atom<any[]>([]);

export const $prefetchStore = atom<any[]>([]);

export const $visitedPages = atom<string[]>([]);

export function addVisitedPage(page: string) {
  const visitedPages = $visitedPages.get();
  if (visitedPages.indexOf(page) === -1) {
    $visitedPages.set([...visitedPages, page]);
  }
}

export interface pageData extends Object {
  [key: string]: any
  title: string
  path?: string
  content?: Element
  head: HTMLHeadElement
}

export const $pageDataStore = atom<pageData[]>([]);

export function addPageData(data: pageData ) {
  for (const [key, value] of Object.entries(data)) {
    if (value === undefined) {
      delete data[key];
    }
  }
  for (const item of $pageDataStore.get()) {
    if (item.title === data["title"]) {
      return;
    }
  }
  const pageData = $pageDataStore.get()
  $pageDataStore.set([...pageData, data]);
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

export function updatePageData(data: pageData) {
  const pageData = $pageDataStore.get();
  for (const [index, item] of pageData.entries()) {
    if (item.title === data["title"]) {
      pageData[index] = data;
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
