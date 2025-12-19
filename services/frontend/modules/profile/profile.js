class ProfileManager {
    constructor() {
        this.currentUser = undefined;
        console.log(this.currentUser);
    }

    async getUser() {
        return await api.getProfile();
    }
    render(data) {
        if (document.querySelector('.profile-container')) {
            return;
        }
       
        this.currentUser = data;
        console.log(this.currentUser);
        const profile = document.createElement("div");
        profile.classList.add("profile-container");
        profile.innerHTML= `
            <div class="profile-content">
                <div class="profile">
                    <p class="profile__title">Никнейм: <span class="profile__data">${this.currentUser.username}</span></p>
                    <p class="profile__title">Имя: <span class="profile__data">${this.currentUser.first_name}</span></p>
                    <p class="profile__title"> Фамилия: <span class="profile__data">${this.currentUser.last_name}</span></p>
                    <p class="profile__title">Город: <span class="profile__data">${this.currentUser.city}</span></p>
                    <p class="profile__title">Возвраст: <span class="profile__data">${this.currentUser.age}</span></p>
                    <p class="profile__title">О Вас: <span class="profile__data">${this.currentUser.info}</span></p>
                </div>
                <button id="change" class="change-btn">Изменить</button>
                <button id="quit" class="quit-btn">Выйти</button>
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

        document.getElementById("change").addEventListener("click", () => {
            this.showProfileForm(this.currentUser.username);
        })

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
        await api.logout(localStorage.getItem("refresh_token"));
        location.reload();
    }

    showProfileForm(username) {
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
            this.setProfile(username);
        });
    }

    async setProfile(username) {
        try{
            const firstname = document.getElementById("firstname").value; 
            const secondname = document.getElementById("secondname").value; 
            const age = Number(document.getElementById("age").value); 
            const city = document.getElementById("city").value; 
            const about = document.querySelector(".about").value;

            console.log("Данные формы:", { firstname, secondname, age, city, about });

            if(firstname && secondname && age && city) {
                console.log("Отправка запроса setProfile...");
                await api.setProfile(username, firstname, secondname, age, city, about);
                console.log("Запрос setProfile выполнен успешно");
                
                document.body.removeChild(document.querySelector(".setProfile-container"));
                app.loadPages("search");
                tempNotice.success("Данные профиля сохранены!");
            } else {
                console.log("Проверка полей не пройдена");
                tempNotice.error("Заполните обязательные поля (Имя, Фамилия, Возраст, Город)");
            }
        } catch(error) {
            console.error("Ошибка в setProfile:", error);
            tempNotice.error("Ошибка, попробуйте еще раз через некоторое время");
        }
    }
}

const profileManager = new ProfileManager();