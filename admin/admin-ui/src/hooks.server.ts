import type { Handle } from "@sveltejs/kit";

export const handle: Handle = async ({event, resolve }) => {
    const response = await resolve(event);
    if (event.url.pathname.includes("/admin/ui")) {
        event.cookies.get("colibase");
        console.log(event.cookies.get("colibase"));
    }
    if (event.url.pathname.includes("/admin/entry")) {

    }
    return response;
}