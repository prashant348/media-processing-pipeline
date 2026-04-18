export class State<T> {

    private value: T;
    private render: (val: T) => void;

    constructor(initialValue: T, render: (val: T) => void) {
        this.value = initialValue;
        this.render = render;
        try {
            this.render(this.value);
        } catch (err) {}
    }

    get val() {
        return this.value;
    }

    set val(newValue: T) {
        this.value = newValue
        this.render(this.value)
    }
}