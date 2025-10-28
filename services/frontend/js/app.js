class App {
    constructor(){
        this.currentPage = "films";
        this.token = localStorage.getItem("token");
        this.init();
    }

    init(){
        document.addEventListener("click", (event) => {
            if(event.target.matches("[data-page]")){
                event.preventDefault()
                this.loadPages(event.target.dataset.page);
            }
        });
    }

    loadPages(name){
        switch(name){
            case "search":
                searchManager.render();
                break;
            case "recomendation":
                //recomendationManager.render();
                this.renderRecommendationsPage();
                break;
            case "notice":
                //noticeManager.render();
                this.renderNotificationsPage()
                break;
            case "profile":
                //profileManager.render();
                this.renderProfilePage();
                break;
            case "friends":
                //friendsManager.render();
                this.renderFriendsPage();
                break;
        }
    }

    renderRecommendationsPage() {
        document.getElementById("main-content").innerHTML = `
            <h1>Рекомендации</h1>
            <p>Страница рекомендаций в разработке</p>
        `;
    }

    renderNotificationsPage() {
        document.getElementById("main-content").innerHTML = `
            <h1>Уведомления</h1>
            <p>Страница уведомлений в разработке</p>
        `;
    }

    renderProfilePage() {
        document.getElementById("main-content").innerHTML = `
            <h1>Профиль</h1>
            <p>Страница профиля в разработке</p>
        `;
    }

    renderFriendsPage() {
        document.getElementById("main-content").innerHTML = `
            <h1>Друзья</h1>
            <p>Страница друзей в разработке</p>
        `;
    }
}

const app = new App();