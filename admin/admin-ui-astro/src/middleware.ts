import { defineMiddleware } from 'astro:middleware';

export const onRequest = defineMiddleware(async (context, next) => {
    if (context.url.pathname.includes('ui')) {
        if (context.cookies.get("colibase")) {
            try {
            const res = await fetch('http://0.0.0.0:6700/api/auth-check', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    cookie: context.cookies.get("colibase")
                })
            });
            const data = await res.json();
            if (data.status === 200 && data.message === 'Authorized') {
                return next();
            } else {
                return next(), context.redirect('/entry/login');
            }
        } catch (e) {
            console.log(e);
            return next(),  context.redirect('/entry/login');
        }

        } else {
            console.log(context.url.pathname);
            return next(), context.redirect('/entry/login') ;
        }
    } 
    if (context.url.pathname.includes('entry')) {
        if (context.request.method === 'GET') {
            const res = await fetch('http://0.0.0.0:6700/api/user-initialization-status');
        if (context.url.pathname.includes('init')) {
            if (res.status === 200) {
                return next(), context.redirect('/entry/login');
            } else if (res.status === 500) {
                return next();
            }
        } else if (context.url.pathname.includes('login')) {
            if (res.status === 200) {
                return next();
            } else if (res.status === 500) {
                return next(), context.redirect('/entry/init');
            }
        }}
    }
    return next();
})

