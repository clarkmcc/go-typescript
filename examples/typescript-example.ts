class Person {
    name: string;

    constructor(name: string) {
        this.name = name
    }

    public greet(): string {
        return `Hello ${this.name}!`
    }
}

const me = new Person("John Doe")
me.greet()