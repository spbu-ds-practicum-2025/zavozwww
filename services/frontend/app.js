class App {
    constructor(){
        this.currentPage = localStorage.getItem("currentPage") || "search";
        this.token = localStorage.getItem("token");
        const userData = localStorage.getItem("currentUser");
        this.currentUser = userData ? JSON.parse(userData) : null;
        this.notificationSocket = null;
        this.profileShowen = false;
        this.currentPage = "recomendation";
        this.init();
    }

    async connectNotification() {
        try {
            const ws = new WebSocket(`/users/notifications?token=${this.token}`);
            ws.onopen = () => {
                console.log('WebSocket connected');
            };
            ws.onerror = (error) => {
                console.error('WebSocket error:', error);
            };
            
            return ws;
        } catch (error) {
            console.error('Failed to connect to notifications:', error);
            return null;
        }
    }


    async init(){   
        localStorage.setItem('currentPage', this.currentPage);
        if (this.token) {
            this.notificationSocket = await this.connectNotification();
            if (this.notificationSocket) {
                this.noticeManager = new NoticeManager(this.notificationSocket);
            }
        }
        document.addEventListener("click", (event) => {
            if(event.target.closest("[data-page]")){
                event.preventDefault();
                this.loadPages(event.target.dataset.page);
            }
        });

        if (this.token) {
            this.loadPages(this.currentPage);
        }
    }

    loadPages(name){
        console.log(name);
        if(name != "profile" && name != "notice"){
            this.currentPage = name;
            localStorage.setItem('currentPage', name);
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
                console.log(this.profileShowen);
                if(!this.profileShowen){
                    this.profileShowen = true;
                    profileManager.render();
                } else {
                    document.querySelector(".profile-container").classList.remove("show-profile");
                    document.querySelector(".profile-container").classList.add("hide-profile");
                    this.profileShowen = false;
                }
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