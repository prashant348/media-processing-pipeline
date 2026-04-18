export abstract class Component {
    constructor(protected params: any = {}) {}
    // every component must have its own html
    abstract render(): string;

    // logic to run after the html is injected into the DOM
    abstract mount(): void;
}
