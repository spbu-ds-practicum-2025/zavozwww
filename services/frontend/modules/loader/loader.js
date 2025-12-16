class LoaderManager {
    constructor() {
        this.loader = document.getElementById("loader");
        this.countRequest = 0;
    }

    show() {
        this.countRequest++;
        if(this.countRequest == 1) {
            this.timeout = setTimeout(() => {
                this.loader.classList.remove("hidden");
            }, 200);
        }
    }

    hide() {
        this.countRequest--;
        if(this.countRequest == 0) {
            clearTimeout(this.timeout);
            this.loader.classList.add("hidden");
        }
    }
}

let loaderManager = new LoaderManager();