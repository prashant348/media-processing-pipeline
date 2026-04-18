import { Component } from "./Component";
import { match, type MatchFunction } from "path-to-regexp";

export type Route = {
    path: string;
    view: new (params: any) => Component;
}

type CompiledRoute = {
    matcher: MatchFunction<any>;
    view: new (params: any) => Component;
}

export class Router {

    private container: HTMLElement;
    private routes: Route[];
    private compiledRoutes: CompiledRoute[];

    constructor(container: HTMLElement, routes: Route[]) {
        this.container = container;
        this.routes = routes;
        this.compiledRoutes = this.routes.map(route => ({
            matcher: match(route.path, { decode: decodeURIComponent }),
            view: route.view
        }));
        window.addEventListener("popstate", () => this.handleRoute());
        this.handleRoute();
    }

    public navigate(url: string) {
        history.pushState(null, '', url);
        this.handleRoute();
    }


    private handleRoute() {
        const path = window.location.pathname;
        let activeRoute = null;
        let params = {};

        for (const route of this.compiledRoutes) {
            const result = route.matcher(path);
            if (result) {
                activeRoute = route;
                params = result.params;
                break;
            };
        };

        if (!activeRoute) {
            this.container.innerHTML = "<h1>404 Not Found</h1>";
            return;
        }

        const instance = new activeRoute.view(params);
        this.container.innerHTML = instance.render();
        instance.mount();
    }
}