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

        const themeBtn = document.getElementById("theme-toggle");

        if (!themeBtn.dataset.hasListener) {
            if (localStorage.getItem("theme") == "light") {
                document.body.classList.add("light-theme");
                document.querySelector(".sun").classList.add("hidden"); 
                document.querySelector(".moon").classList.remove("hidden");
            } else {
                document.body.classList.remove("light-theme");
                document.querySelector(".sun").classList.remove("hidden"); 
                document.querySelector(".moon").classList.add("hidden");
            }

            themeBtn.addEventListener("click", () => {
                if (!document.body.classList.contains("light-theme")) {
                    document.body.classList.add("light-theme");
                    document.querySelector(".sun").classList.add("hidden"); 
                    document.querySelector(".moon").classList.remove("hidden");
                } else {
                    document.body.classList.remove("light-theme");
                    document.querySelector(".sun").classList.remove("hidden"); 
                    document.querySelector(".moon").classList.add("hidden");
                }

                const isLight = document.body.classList.contains("light-theme");
                localStorage.setItem("theme", isLight ? "light" : "dark");
            });
            
            themeBtn.dataset.hasListener = "true";
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

