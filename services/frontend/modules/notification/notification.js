class NoticeManager {
    constructor(notificationSocket) {
        this.socket = notificationSocket;
        this.haveNotice = false;
        this.token = localStorage.getItem("token");
        this.notifications = new Map();
        this.socket.onmessage = (event) => {
            const notification = JSON.parse(event.data);
            document.querySelector(".navbar__content__pages__notice").classList.add("haveNotice");
            this.notifications.set(notification.user_id, notification);
        }   
    }

    render() {
        console.log("render");
        if(this.notifications.size == 0){
            document.querySelector(".navbar__content__pages__notice").classList.remove("haveNotice");
        }
        if (document.querySelector(".notice-container")) {
            if(this.notifications.size == 0) {
                document.querySelector(".notice-container").innerHTML = `
                    <p class="notice-card__message">У вас нет уведомлений</p>
                `
            }
            console.log("return");
            return;
        }
        
        const notice = document.createElement("div");
        console.log("create");
        notice.classList.add("notice-container");
        if(this.notifications.size > 0){
            notice.innerHTML = Array.from(this.notifications.values()).map(note => 
                `
                <div id="request-${note.user_id}" class="notice-card">
                    <p class="notice-card__message">Вам запрос на дружбу от <span class="from-request">${note.username}</span></p>
                    <div class="notice-card__actions">
                    <button class="btn accept-btn" data-id="${note.user_id}" class"notice-card__accept-btn">Принять</button>
                    <button class="btn decline-btn" data-id="${note.user_id}" class="notice-card__decline-btn">Отклонить</button>
                    </div>
                </div>
                `
            ).join("");
            console.log(notice);
            console.log("map");

            notice.querySelectorAll(".accept-btn").forEach(btn => {
                btn.addEventListener("click", (event) => {
                    const userId = event.currentTarget.dataset.id;
                    event.stopPropagation();
                    this.acceptRequest(userId);
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
                <p class="notice-card__message">У вас нет уведомлений</p>
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

    async acceptRequest(userFrom){
        try {
            await api.acceptRequest(userFrom);
            this.notifications.delete(Number(userFrom));
            document.getElementById(`request-${userFrom}`).remove();
            tempNotice.success("Вы приняли запрос на дружбу"); 
            this.render();
        } catch(error) {
            tempNotice.error("Ошибка, попробуйте еще раз через некоторое время");
            console.log(`Error: ${error.message}`);
        }
    }

    async declineRequest(userFrom){
        try {
            await api.declineRequest(userFrom);
            this.notifications.delete(Number(userFrom));
            document.getElementById(`request-${userFrom}`).remove();
            tempNotice.success("Вы отклонили запрос на дружбу"); 
            this.render();
        } catch(error) {
            tempNotice.error("Ошибка, попробуйте еще раз через некоторое время");
            console.log(`Error: ${error.message}`);
        }
    }
}