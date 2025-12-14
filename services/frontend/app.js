class App {
    constructor(){
        if(!localStorage.getItem("currentPage")) {
            localStorage.setItem("currentPage", "search");
        }
        this.currentPage = localStorage.getItem("currentPage");
        this.token = localStorage.getItem("token");
        this.notificationSocket = null;
        this.currentPage = "search";
        this.init();
    }

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


    async init(){   
        localStorage.setItem("currentPage", this.currentPage);
        if (this.token) {
            this.notificationSocket = await this.connectNotification();
            if (this.notificationSocket) {
                this.noticeManager = new NoticeManager(this.notificationSocket);
            }
        }
        document.addEventListener("click", (event) => {
            const clickedPage = event.target.closest("[data-page]");
            if(clickedPage){
                event.preventDefault();
                this.loadPages(clickedPage.dataset.page);
            }
        });

        if (this.token) {
            this.loadPages(this.currentPage);
        }
    }

    async loadPages(name){
        let data = await api.getProfile();
        if(!data.first_name || !data.last_name || !data.age || !data.city || !data.info){
            profileManager.showProfileForm();
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
                this.noticeManager.render();
                break;
            case "profile":
                profileManager.render(data);
                break;
            case "friends":
                friendsManager.render();
                break;
        }
    }

    renderNotificationsPage() {
        document.getElementById("main-content").innerHTML = `
            <h1>Уведомления</h1>
            <p>Страница уведомлений в разработке</p>
        `;
    }

}

const app = new App();