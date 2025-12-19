class NoticeManager {
    constructor() {
        // this.socket = notificationSocket;
        this.notifications = [];
        this.waitRequest = false;
        /*
        this.socket.onmessage = (event) => {
            const notification = JSON.parse(event.data);
            document.querySelector(".navbar__content__pages__notice").classList.add("haveNotice");
            this.notifications.set(notification.username, notification);
        }   
        */
    }

    start() {
        if(this.waitRequest) {
            return;
        }

        this.getNotifications();

        this.waitRequest = setInterval(() => {
            this.getNotifications();
        }, 30000);
    }

    stop() {
        if (this.waitRequest) {
            clearInterval(this.waitRequest);
            this.waitRequest = null;
        }
    }

    async getNotifications() {
        try {
            const data = await api.getNotice(); 
            const newNotifications = data || [];

            if (JSON.stringify(newNotifications) !== JSON.stringify(this.notifications)) {
                this.notifications = newNotifications;
                document.querySelector(".navbar__content__pages__notice").classList.add("haveNotice");
            }
        } catch (error) {
            console.error("Error:", error);
        }
    }

    async render(username) {
        this.username = username;

        await this.getNotifications();
        console.log("render");
        if(this.notifications.length == 0){
            document.querySelector(".navbar__content__pages__notice").classList.remove("haveNotice");
        }
        if (document.querySelector(".notice-container")) {
            if(this.notifications.length == 0) {
                document.querySelector(".notice-container").innerHTML = `
                    <p class="notice-card__message">У вас нет уведомлений!</p>
                `
            }
            console.log("return");
            return;
        }
        
        const notice = document.createElement("div");
        console.log("create");
        notice.classList.add("notice-container");
        if(this.notifications.length > 0){
            notice.innerHTML = this.notifications.map(note => 
                `
                <div id="request-${note.request_id}" class="notice-card">
                    <p class="notice-card__message">Вам запрос на дружбу от <span class="from-request">${note.from_username} </span><span class="date-request">${new Date(note.created_at).toLocaleDateString()}</span></p>
                    <div class="notice-card__actions">
                    <button class="btn accept-btn" data-id="${note.request_id}" data-from="${note.from_username}">Принять</button>
                    <button class="btn decline-btn" data-id="${note.request_id}">Отклонить</button>
                    </div>
                </div>
                `
            ).join("");
            console.log(notice);

            notice.querySelectorAll(".accept-btn").forEach(btn => {
                btn.addEventListener("click", (event) => {
                    const userId = event.currentTarget.dataset.id;
                    const fromUsername = event.currentTarget.dataset.from;
                    event.stopPropagation();
                    this.acceptRequest(userId, fromUsername);
                });
            });
            notice.querySelectorAll(".decline-btn").forEach(btn => {
                btn.addEventListener("click", (event) => {
                    const userId = event.target.dataset.id;
                    event.stopPropagation();
                    this.declineRequest(userId);
                });
            });
        } else {
            notice.innerHTML = `
                <p class="notice-card__message">У вас нет уведомлений!</p>
            `
        }

        document.getElementById("navbar").after(notice);
        console.log("add");
        setTimeout(() => {
            notice.classList.add("show-notice");
        }, 10);

        const closeHandler = (event) => {
            if (!event.target.closest(".notice-container")) {
                notice.classList.remove("show-notice");
                notice.classList.add("hide-notice");
                
                setTimeout(() => {
                    notice.remove();
                    document.removeEventListener("click", closeHandler); 
                }, 300);
            }
        };

        setTimeout(() => {
            document.addEventListener("click", closeHandler);
        }, 0);
    }

    async acceptRequest(requestID, fromUsername){
        try {
            await api.acceptRequest(fromUsername, this.username);
            this.notifications = this.notifications.filter(note => note.request_id !== requestID);
            document.getElementById(`request-${requestID}`).remove();
            tempNotice.success("Вы приняли запрос на дружбу"); 
            this.render();
        } catch(error) {
            tempNotice.error("Ошибка, попробуйте еще раз через некоторое время");
            console.log(`Error: ${error.message}`);
        }
    }

    async declineRequest(requestID){
        try {
            await api.declineRequest(requestID);
            this.notifications = this.notifications.filter(note => note.request_id !== requestID);
            document.getElementById(`request-${requestID}`).remove();
            tempNotice.success("Вы отклонили запрос на дружбу"); 
            this.render();
        } catch(error) {
            tempNotice.error("Ошибка, попробуйте еще раз через некоторое время");
            console.log(`Error: ${error.message}`);
        }
    }
}