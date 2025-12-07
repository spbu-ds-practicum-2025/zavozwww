class Api {
    constructor(){
        this.apiURL = "http://localhost:8080/api";
    }

    async refresh(){
        try {
            const response = await fetch(`${this.apiURL}/refresh`, {
                method: "POST",
                credentials: "include",
            });
            if (response.ok) {
                const data = await response.json();
                localStorage.setItem("token", data.access_token);
                return true;
            }
            return false;
        } catch {
            return false;
        }
    }

    async request(endpoint, options = {}){
        let url = `${this.apiURL}${endpoint}`;
        let fetchInit = {
            headers: {
                "Content-Type": "application/json",
                "credentials": "include",
                ...options.headers
            },
            ...options
        };

        const token = localStorage.getItem("token");
        if(token){
            fetchInit.headers["Authorization"] = `Bearer ${token}`;
        }

        try {
            let response = await fetch(url, fetchInit);
            let data = await response.json();

            if(response.status === 401){
                const refreshed = await this.refresh();
                if (refreshed) {
                    return this.request(endpoint, options);
                } else {
                    authManager.renderAuthForm();
                    throw new Error("Session expired");
                }
            }
            
            return data;
        } catch(error){
            throw new Error(error.statusText);
        }
    }

    async login(Email, Password){
        return this.request("/login", {
            method: "POST",
            body: JSON.stringify({  email: Email,
                                    password: Password,}),
        });
    }

    async registration(Username, Email, Password){
        return this.request("/register", {
            method: "POST",
            body: JSON.stringify({ email: Email, username: Username, password: Password }),
        });
    }

    async confirm(userCode, Email) {
        return this.request("/confirm", {
            method: "POST",
            body: JSON.stringify({code: userCode, email: Email}),
        })
    }

    async again(Email) {
        return this.request("/again", {
            method: "POST",
            body: JSON.stringify({email: Email}),
        })
    }

    async searchMovie(searchQuery, genre) {
        return this.request(`/search?query=${encodeURIComponent(searchQuery)}&genre=${genre}`);
    }

    async rating(ID, Rating, Message){
        return this.request(`/movies/${ID}/rating`, {
                    method: "POST",
                    body: JSON.stringify({ rating: Rating, review: Message }),
                });
    }

    async recomendation() {
        return this.request("/recomendations");
    }

    async searchFriend(Name) {
        return this.request(`/friends/search?name=${encodeURIComponent(Name)}`, {
            method: "POST",
            body: JSON.stringify({ target_user_id: Name }),
        });
    }

    async addFriend(Name) {
        return this.request("/friends/request", {
           method: "POST",
           body: JSON.stringify({target_user_id: Name}),
       });
    }

    async acceptRequest(userFrom) {
        return this.request("/friends/accept/request_id", {
            method: "POST",
            body: JSON.stringify({user_id: userFrom}),
        });
    }

    async declineRequest(userFrom) {
        return this.request("/friends/decline/request_id", {
            method: "POST",
            body: JSON.stringify({user_id: userFrom}),
        });
    }

    // что отправлять на эндпоинт?
    async logout(){
        return this.request("/logout", {
            method: "POST",
        });
    }

    async getProfile() {
        return this.request("/profile");
    }

    async getFriends() {
        return this.request("/friends");
    }

    async setProfile(Firstname, Secondname, Age, City, About) {
        return this.request("/profile", {
            method: "POST",
            body: JSON.stringify({firstname: Firstname, secondname: Secondname, age: Age, city: City, about: About}),
        });
    }

}

const api = new Api();