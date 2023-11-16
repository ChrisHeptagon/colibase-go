import { fail, type Actions } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";

let rootURL: string;
rootURL = "http://127.0.0.1:6700";



export const load: PageServerLoad = async () => {
  let UserSchema: any;
  const response = await fetch(`${rootURL}/api/login-schema`);
  const json = await response.json();
  UserSchema = json;
  console.log(UserSchema);
  return { UserSchema };
};

export const actions = {
  action: async ({ request }) => {
    const formData = await request.formData();
    if (formData) {
      const finalForm = Object.fromEntries(formData.entries());
      const fetchResult = await fetch(`${rootURL}/api/init-login`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(finalForm),
      });
      const json = await fetchResult.json();
      let error = json.error;
      if (json.error) {
        if (Array.isArray(json.error)) {
          let errorArray = error as string[];
          return fail(400, errorArray as any);
        } else {
          return fail(400, json.error);
        }
      }
    }
  },
} satisfies Actions;
