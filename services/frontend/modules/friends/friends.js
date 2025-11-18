class FriendsManager {
    render() {
        const mainContent = document.getElementById("main-content");
        mainContent.innerHTML = `
            <div class="friends-page">
                <h2>Поиск пользователей</h2>
                <div class="search-box">
                    <input type="text" id="friends-search" placeholder="Логин пользователя...">
                    <button id="search-btn">Найти</button>
                <div id="friends-results"></div>
                <h1>Ваши друзья</h1>
                <div id="friends"></div>
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
            console.log(`Error: ${error.message}`);
        }
    }

    renderYourFriends(){
        this.currentUser = JSON.parse(localStorage.getItem("currentUser"));
        let friends = this.currentUser.friends;

        const yourFriends = document.getElementById("friends");

        if(yourFriends.length == 0){
            yourFriends.innerHTML = `
            <p>У Вас пока нет друзей. Скорее найдите их и отправте запрос на дружбу!</p>
            `
        }

        yourFriends.innerHTML = friends.map(friend => 
            `
            <div class="friend-card">
                <img src="${friend.avatarSrc}" class="friend-card__image"/>
                <h3 class="friend-card__name">${friend.name}</h3>
            </div>
            `
        ).join("");
    }

    renderSearchFriends(friendsArray) {
        const results = document.getElementById("friends-results");
        
        if(friendsArray.length == 0){
            results.innerHTML = `
            <p>Просстите, но мы не можем найти такого пользователя</p>
            `
        }

        results.innerHTML = friendsArray.map(friend => 
            `
            <div class="friend-card">
                <img src="${friend.avatarSrc}" class="friend-card__image"/>
                <h3 class="friend-card__name">${friend.name}</h3>
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