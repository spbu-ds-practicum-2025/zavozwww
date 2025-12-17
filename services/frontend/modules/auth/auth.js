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
                <h1 class="logo_auth">FILMBUDDY</h1>
                <h2>Сперва войдите в свой аккаунт</h2>
                <form id="login-form" class="auth__form">
                    <div class="auth__form__group">
                        <input type="text" id="email" placeholder="Email" required>
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
                <h1 class="logo_auth">FILMBUDDY</h1>
                <h2>Регистрация</h2>
                <form id="reg-form" class="auth__form">
                    <div class="auth__form__group">
                        <input type="text" id="username" placeholder="Логин" required>
                    </div>  
                    <div class="auth__form__group">
                        <input type="email" id="email" placeholder="Email" required>
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

        const email = document.getElementById("email").value;
        const password = document.getElementById("login-password").value;

        if (email && password) {
            //send to server
            try {
                console.log(email);
                let data = await api.login(email, password);
                localStorage.setItem("token", data.access_token);
                this.token = data.access_token;
                document.getElementById("navbar").classList.remove("hidden");
                app.loadPages("search");
                app.noticeManager.start();
            } catch(error){
                tempNotice.error("Ошибка, повторите попытку немного позже");
                console.log(error.message);
            }
        } else {
            tempNotice.error("Пожалуйста, заполните все поля");
        }
    }

    async registration(event){
        event.preventDefault();

        const username = document.getElementById("username").value;
        const email = document.getElementById("email").value;
        const password = document.getElementById("login-password").value;
        
        if (username && email && password) {
            // send to server
            try{
                await api.registration(username, email, password);
                this.renderConfirmCodeForm(email);
            } catch(error) {
                tempNotice.error("Ошибка, повторите попытку немного позже");
                console.log(error.message);
            }
        } else {
            tempNotice.error("Пожалуйста, заполните все поля");
        }
    }

    renderConfirmCodeForm(email) {
        const mainContent = document.getElementById("main-content");

        mainContent.innerHTML = `
        <div class="auth">
            <img src="assets/auth-img.jpg" class="auth__image" />
            <div class="auth__block">
                <h1 class="logo_auth">FILMBUDDY</h1>
                <h2>Мы выслали Вам код на указанный email, пожалуйста введите его для завершения регистрации</h2>
                <form id="confirm-form" class="auth__form">
                    <div class="auth__form__group">
                        <input type="text" id="code" placeholder="Код с почты" required>
                    </div>
                    <button class="auth__form__button" id="apply-code" type="submit">Подтвердить</button>
                    <button class="auth__form__button_sec" id="one-more-time" type="button">Отправить код еще раз</button>
                </form>
            </div>
        </div>
        `

        this.setupConfirmForm(email);
    }

    setupConfirmForm(email) {
        document.getElementById("confirm-form").addEventListener("submit", (event) => {
            event.preventDefault();
            const code = document.getElementById("code").value;
            this.confirm(code, email);
        });

        document.getElementById("one-more-time").addEventListener("click", () => {
            this.again(email);
        });
    }

    async confirm(code, email) {
        if(code) {
            try {
                const data = await api.confirm(code, email);
                if(!data.correct) {
                    this.renderAuthForm();
                } else {
                    throw Error("Uncorrect data");
                }
            } catch(error) {
                tempNotice.error("Ошибка, повторите попытку немного позже");
                console.log(error.message); 
            }
        } else {
            tempNotice.error("Пожалуйста, заполните поле с кодом");
        }
    }

    async again(email) {
        await api.again(email);
        document.getElementById("one-more-time").classList.add("no-active");
        document.getElementById("one-more-time").setAttribute("title", "Пожалуйста подождите 30 сек. перед тем как запросить еще раз")
        setTimeout(() => {
            document.getElementById("one-more-time").classList.remove("no-active");
            document.getElementById("one-more-time").removeAttribute("title");
        }, 30000)
    }
}