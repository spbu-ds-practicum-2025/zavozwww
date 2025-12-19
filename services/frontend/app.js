class App {
    constructor(){
        if(!localStorage.getItem("currentPage")) {
            localStorage.setItem("currentPage", "search");
        }
        this.currentPage = localStorage.getItem("currentPage");
        this.token = localStorage.getItem("token");
        // this.notificationSocket = null;
        // this.currentPage = "search";
        this.noticeManager = new NoticeManager();
        this.init();
    }

    /*
    async connectNotification() {
        try {
            const ws = new WebSocket(`/users/notifications?token=${this.token}`);
            ws.onopen = () => {
                console.log("WebSocket connected");
            };
            ws.onerror = (error) => {
                console.error("WebSocket error:", error);
            };
            
            return ws;
        } catch (error) {
            console.error("Failed to connect to notifications:", error);
            return null;
        }
    }
    */

    async init(){   
        localStorage.setItem("currentPage", this.currentPage);
        
        /*
        if (this.token) {
            this.notificationSocket = await this.connectNotification();
            if (this.notificationSocket) {
                this.noticeManager.socket = this.notificationSocket; 
            }
        }
        */

        document.addEventListener("click", (event) => {
            const clickedPage = event.target.closest("[data-page]");
            if(clickedPage){
                event.preventDefault();
                this.loadPages(clickedPage.dataset.page);
            }
        });

        document.getElementById("theme-toggle").addEventListener('click', () => {
            document.body.classList.toggle("light-theme");
            document.querySelector(".sun").toggle("hidden");
            document.querySelector(".moon").toggle("hidden");
            const isDark = document.body.classList.contains("light-theme");
            localStorage.setItem("theme", isDark ? "light" : "dark");
        });

        if (localStorage.getItem("theme") == "light") {
            document.body.classList.add("light-theme");
        }

        if (this.token) {
            this.loadPages(this.currentPage);
            this.noticeManager.start();
        }
    }

    async loadPages(name){
        let data = await api.getProfile();
        if(data.first_name === "Unknown" || data.last_name == "Unknown" || !data.city || !data.info){
            profileManager.showProfileForm(data.username);
        }
        if(name != "profile" && name != "notice"){
            this.currentPage = name;
            localStorage.setItem("currentPage", name);
        } 
        switch(name){
            case "search":
                searchManager.render(); 
                break;
            case "recomendation":
                recomendationManager.render();
                break;
            case "notice":
                this.noticeManager.render(data.username);
                break;
            case "profile":
                profileManager.render(data);
                break;
            case "friends":
                friendsManager.render(data.username);
                break;
            case "ratingList":
                ratingListManager.render();
                break;
        }
    }
}

const app = new App();

