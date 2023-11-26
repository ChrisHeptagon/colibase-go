import { atom } from "nanostores";

export interface stats {
    memory: Object[]
    cpu: Object[]
}

export const $statsStore = atom<stats[]>([]);
