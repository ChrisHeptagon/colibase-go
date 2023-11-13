import type { Actions } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

let rootURL: string;
rootURL = "http://127.0.0.1:6700";


export const load: PageServerLoad = async () => {
    let UserSchema: any;
    const response = await fetch(`${rootURL}/api/login-schema`);
        const json = await response.json();
        UserSchema = json;
    return {UserSchema};
};

export const actions = {
    init: async (event) => {
        const formData = await event.request.formData();
        const formDataToSend = new FormData();
        for (const [key, value] of formData.entries()) {
            formDataToSend.append(key, value);
        }
        console.log(formDataToSend)
        const fetchResult = await event.fetch(`${rootURL}/api/init-login`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: formDataToSend,
        }
    )
    console.log(await fetchResult.text())
    },
} satisfies Actions;

