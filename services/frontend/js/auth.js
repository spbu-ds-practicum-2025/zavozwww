class AuthManager {
    constructor() {
        this.currentUser = null;
        this.token = localStorage.getItem("token");
        this.renderAuthForm();
    }

    renderAuthForm(){
        const mainContent = document.getElementById("main-content");

        mainContent.innerHTML = `
        <div class="auth">
            <img src="assets/auth-img.jpg" class="auth__image" />
            <div class="auth__block">
                <h1 class="logo">FILMBUDDY</h1>
                <h2>Сперва войдите в свой аккаунт</h2>
                <form id="login-form" class="auth__form">
                    <div class="auth__form__group">
                        <input type="text" id="username" placeholder="Логин" required>
                    </div>
                    <div class="auth__form__group">
                        <input type="password" id="login-password" placeholder="Пароль" required>
                    </div>
                    <button class="auth__form__button" type="submit">Войти</button>
                    <button class="auth__form__button_sec" id="reg-button" type="button">Регистрация</button>
                </form>
            </div>
        </div>
        `

        this.setupAuthForm();
    }

    renderRegForm(){
        const mainContent = document.getElementById("main-content");

        mainContent.innerHTML = `
        <div class="auth">
            <img src="assets/auth-img.jpg" class="auth__image" />
            <div class="auth__block">
                <h1 class="logo">FILMBUDDY</h1>
                <h2>Регистрация</h2>
                <form id="reg-form" class="auth__form">
                    <div class="auth__form__group">
                        <input type="text" id="username" placeholder="Логин" required>
                    </div>  
                    <div class="auth__form__group">
                        <input type="email" id="email-username" placeholder="Почта" required>
                    </div>
                    <div class="auth__form__group">
                        <input type="password" id="login-password" placeholder="Пароль" required>
                    </div>
                    <button class="auth__form__button" id="apply-reg-button" type="submit">Зарегистрироваться</button>
                </form>
            </div>
        </div>
        `

        this.setupAuthForm();
    }

    setupAuthForm(){
        let loginForm = document.getElementById("login-form");
        if(loginForm !== null) {
            loginForm.addEventListener("submit", (event) => {
                this.login(event);
            })
        }

        let regForm = document.getElementById("reg-form");
        if(regForm !== null){
            regForm.addEventListener("submit", (event) => {
                this.registration(event);
            })
        }
        
        let regButton = document.getElementById("reg-button");
        if(regButton !== null) {
            regButton.addEventListener("click", () => {
                this.renderRegForm();
            })
        }
    }

    async login(event) {
        event.preventDefault();

        const username = document.getElementById("username").value;
        const password = document.getElementById("login-password").value;

        if (username && password) {
            //send to server
            try {
                let data = await api.login(username, password);
                localStorage.setItem("token", data.token);
                this.token = data.token;
                this.currentUser = data.user_profile;
                
                document.getElementById("navbar").classList.remove("hidden");
                app.loadPages("search");
            } catch(error){
                console.log(error.message);
            }
        } else {
            alert("Пожалуйста, заполните все поля");
        }
    }

    async registration(event){
        event.preventDefault();

        const username = document.getElementById("username").value;
        const email = document.getElementById("email-username").value;
        const password = document.getElementById("login-password").value;
        
        if (username && email && password) {
            // send to server
            try{
                let data = await api.registration(username, email, password);
                this.token = data.token;

                document.getElementById("navbar").classList.remove("hidden");
                app.loadPages("search");
            } catch(error) {
                console.log(error.message);
            }
        } else {
            alert("Пожалуйста, заполните все поля");
        }
    }
}