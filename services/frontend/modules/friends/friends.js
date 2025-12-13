class FriendsManager {
    render() {
        const mainContent = document.getElementById("main-content");
        mainContent.innerHTML = `
            <div class="friends-page">
                <h1 class="friends-page__title">Поиск пользователей</h1>
                <div class="search-box">
                    <input type="text" id="friends-search" class="search-box__input" placeholder="Логин пользователя...">
                    <button id="search-btn" class="search-box__button">Найти</button>
                </div>
                <div id="your-friends-results">
                    <h2 class="friends-page__title_small">Ваши друзья</h2>
                    <div id="your-friends" class="friends-page__result friends-page__your-friends"></div>
                </div>
                <div id="search-friends-results" class="results friends-page__results"></div>
            </div>
        `

        this.renderYourFriends();

        document.getElementById("search-btn").addEventListener("click", () => {
            this.search();
        })

        document.getElementById("friends-search").addEventListener("keypress", (event) => {
            if(event.key == "Enter"){
                this.search();
            }
        })
    }

    async search() {
        let searchQuery = document.getElementById("friends-search").value;
        try{
            let data = await api.searchFriend(searchQuery);
            this.renderSearchFriends(data.friends);
        } catch(error) {
            tempNotice.error("Ошибка, попробуйте еще раз через некоторое время");
            console.log(`Error: ${error.message}`);
        }
    }

    async renderYourFriends(){
        let friends = await api.getFriends();
        const yourFriends = document.getElementById("your-friends");
        if(friends.length == 0){
            yourFriends.innerHTML = `
            <p class="friends-page__without-friends">У Вас пока нет друзей. Скорее найдите их и отправте запрос на дружбу!</p>
            `
            return;
        }

        /*
            <p class="friend-card__title">Имя: <span class="friend-card__data">${friend.name}</span></p>
            <p class="friend-card__title">Дата регистрации: <span class="friend-card__data">${friend.regDate}</span></p>
            <div class="profile-activity profile-acticity-friendsPage">
                <div class="profile-stat profile-stat-friendsPage">
                    <span class="categorie">Оцененных фильмов:</span>
                    <span class="number">${friend.countRateFilms}</span>
                </div>
                <div class="profile-stat profile-stat-friendsPage">
                    <span class="categorie">Количество друзей:</span>
                    <span class="number">${friend.countFriends}</span>
                </div>
            </div>
        */

        yourFriends.innerHTML = friends.map(friend => 
            `
            <div class="friend-card">
                <p class="friend-card__title">Имя: <span class="friend-card__data">${friend.name}</span></p>
                <p class="friend-card__title">Имя: <span class="friend-card__data">${friend.firstname}</span></p>
                <p class="friend-card__title"> Фамилия: <span class="friend-card__data">${friend.secondname}</span></p>
                <p class="friend-card__title">Город: <span class="friend-card__data">${friend.city}</span></p>
                <p class="friend-card__title">Возвраст: <span class="friend-card__data">${friend.age}</span></p>
                <p class="friend-card__title">О Вас: <span class="friend-card__data">${friend.info}</span></p>
            </div>
            `
        ).join("");
    }

    renderSearchFriends(friendsArray) {
        const results = document.getElementById("search-friends-results");
        
        if(friendsArray.length == 0){
            results.innerHTML = `
            <p class="friends-page__without-friends">Просстите, но мы не можем найти такого пользователя</p>
            `
            return;
        }

        results.innerHTML = friendsArray.map(friend => 
            `
            <div class="friend-card">
                <p class="friend-card__title">Имя: <span class="friend-card__data">${friend.firstname}</span></p>
                <p class="friend-card__title"> Фамилия: <span class="friend-card__data">${friend.secondname}</span></p>
                <p class="friend-card__title">Город: <span class="friend-card__data">${friend.city}</span></p>
                <p class="friend-card__title">Возвраст: <span class="friend-card__data">${friend.age}</span></p>
                <p class="friend-card__title">О Вас: <span class="friend-card__data">${friend.info}</span></p>
                <button class="friend-card__add-btn" onclick="friendsManager.sendRequest(${friend.name})">Добавить в друзья</button>
            </div>
            `
        ).join("");
    }

    async sendRequest(name) {
        try {
            await api.addFriend(name);
            tempNotice.success("Заявка отправлена");
        } catch(error) {
            tempNotice.error("Ошибка, попробуйте еще раз через некоторое время");
            console.log(`Error: ${error.message}`);
        }
        
    }
}

let friendsManager = new FriendsManager();