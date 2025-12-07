class ProfileManager {
    constructor() {
        this.currentUser = this.getUser();
    }

    async getUser() {
        return await api.getProfile();
    }
    render() {
        if (document.querySelector('.profile-container')) {
            return;
        }
        
        /*
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
        */

        const profile = document.createElement("div");
        profile.classList.add("profile-container");
        profile.innerHTML= `
            <div class="profile-content">
                <div class="profile">
                    <p class="profile__title">Имя: <span class="profile__data">${this.currentUser.firstname}</span></p>
                    <p class="profile__title"> Фамилия: <span class="profile__data">${this.currentUser.secondname}</span></p>
                    <p class="profile__title">Город: <span class="profile__data">${this.currentUser.city}</span></p>
                    <p class="profile__title">Возвраст: <span class="profile__data">${this.currentUser.age}</span></p>
                    <p class="profile__title">О Вас: <span class="profile__data">${this.currentUser.info}</span></p>
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
            this.logout();
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
    
    async logout() {
        const profile = document.querySelector('.profile-container');
        profile.classList.remove("show-profile");
        profile.classList.add("hide-profile");
        app.profileShowen = false;
        profile.remove();
        localStorage.removeItem("token");
        await api.logout();
        location.reload();
    }

    showProfileForm() {
        if(document.querySelector("setProfile-container")) {
            return;
        }
        const profileFormContainer = document.createElement("div");
        profileFormContainer.classList.add("setProfile-container");
        profileFormContainer.innerHTML = `
            <div class="rating setProfile">
                <h1 class="save-title">Сперва заполните свой профиль</h1>
                <div class="auth__form__group">
                    <input type="text" id="firstname" placeholder="Имя" required>
                </div>
                <div class="auth__form__group">
                    <input type="text" id="secondname" placeholder="Фамилия" required>
                </div>
                <div class="auth__form__group">
                    <input type="text" id="age" placeholder="Возвраст" required>
                </div>
                <div class="auth__form__group">
                    <input type="text" id="city" placeholder="Город" required>
                </div>
                <textarea class="about" placeholder="Расскажите о себе"></textarea>
                <button id="save" class="save-btn">Сохранить</button>
            </div>
        `

        document.body.appendChild(profileFormContainer);

        document.getElementById("save").addEventListener("click", () => {
            console.log("CLICK");
            this.setProfile();
        });
    }

    async setProfile() {
        try{
            const firstname = document.getElementById("firstname").value; 
            const secondname = document.getElementById("secondname").value; 
            const age = document.getElementById("age").value; 
            const city = document.getElementById("city").value; 
            const about = document.querySelector(".about").value;
            if(firstname && secondname && age && city) {
                await api.setProfile(firstname, secondname, age, city, about);
                document.body.removeChild(document.querySelector(".setProfile-container"));
                tempNotice.success("Данные профиля сохранены!");
            } else {
                console.log("else ALDLAMKSDMLASD");
                tempNotice.error("Заполните все поля");
            }
        } catch(error) {
            tempNotice.error("Ошибка, попробуйте еще раз через некоторое время");
        }
    }
}

const profileManager = new ProfileManager();