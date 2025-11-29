class ProfileManager {
    constructor() {
        this.currentUser = JSON.parse(localStorage.getItem("currentUser"));
    }

    render() {
        if (document.querySelector('.profile-container')) {
            return;
        }
        
        const profile = document.createElement("div");
        profile.classList.add("profile-container");
        profile.innerHTML= `
            <div class="profile-content">
                <div class="profile">
                    <p class="profile__title">Ваше имя: <span class="profile__data">${this.currentUser.username}</span></p>
                    <p class="profile__title"> Email: <span class="profile__data">${this.currentUser.email}</span></p>
                    <p class="profile__title">Дата регистрации: <span class="profile__data">${this.currentUser.regDate}</span></p>
                </div>
                <div class="profile-activity">
                    <div class="profile-stat profile-stat-profilePage">
                        <span class="categorie">Оцененных фильмов:</span>
                        <span class="number">${this.currentUser.countRateFilms}</span>
                    </div>
                    <div class="profile-stat profile-stat-profilePage">
                        <span class="categorie">Количество друзей:</span>
                        <span class="number">${this.currentUser.countFriends}</span>
                    </div>
                </div>
                <button id="quit" class="quit">Выйти</button>
            </div>
        `

        document.getElementById("navbar").after(profile);
        console.log("add");
        setTimeout(() => {
            profile.classList.add("show-profile");
        }, 10);

        document.getElementById("quit").addEventListener("click", () => {
            profile.classList.remove("show-profile");
            profile.classList.add("hide-profile");
            app.profileShowen = false;
            profile.remove();
            localStorage.removeItem("token");
            location.reload();

        });

        const closeHandler = (event) => {
            if(!event.target.closest(".profile-container")) {
                profile.classList.remove("show-profile");
                profile.classList.add("hide-profile");
                app.profileShowen = false;
                setTimeout(() => {
                    profile.remove();
                    document.removeEventListener("click", closeHandler);
                }, 300);
            }
        }
        document.addEventListener("click", closeHandler);
    }    
}

const profileManager = new ProfileManager();