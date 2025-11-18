class ProfileManager {
    constructor() {
        this.currentUser = JSON.parse(localStorage.getItem("currentUser"));
    }

    render() {
        const mainContent = document.getElementById("main-content");
        mainContent.innerHTML = `
            <div class="profile-container">
                <span class="profile__name">Ваше имя: ${this.currentUser.username}</span>
                <span class="profile__email"> Email: ${this.currentUser.email}</span>
                <span class="reg-date">Дата регистрации: ${this.currentUser.regDate}</span>
            </div>
            <div class="profile-activity">
                <div class="profile-stat">
                    <span class="categorie">Количесвто оцененных фильмов:</span>
                    <span class="number">${this.currentUser.countRateFilms}</span>
                </div>
                <div class="profile-stat">
                    <span class="categorie">Количесвто друзей:</span>
                    <span class="number">${this.currentUser.countFriends}</span>
                </div>
            </div>
        `
    }    
}

const profileManager = new ProfileManager();