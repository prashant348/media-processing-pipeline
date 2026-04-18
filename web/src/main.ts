import './style.css'
import { Router, type Route } from './libs/Router';
import HomePage from './pages/Home';
import VideoPage from './pages/Video';

const appContainer = document.querySelector<HTMLDivElement>('#app')!;

const routes: Route[] = [
    { path: "/", view: HomePage },
    { path: "/video/:id", view: VideoPage },
];

const appRouter = new Router(appContainer, routes);

export const navigate = (url: string) => appRouter.navigate(url);

document.addEventListener("click", (e) => {
    const target = e.target as HTMLElement;

    if (target.matches("[data-link]")) {
        e.preventDefault();
        navigate((target as HTMLAnchorElement).href);
    }
})