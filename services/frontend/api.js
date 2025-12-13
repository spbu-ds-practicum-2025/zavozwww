class Api {
    constructor(){
        this.apiURL = "http://localhost:8080/api";
    }

    async refresh(){
        try {
            const refreshToken = localStorage.getItem("refresh_token");
            if (!refreshToken) return false;

            const response = await fetch(`${this.apiURL}/refresh`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json"
                },
                body: JSON.stringify({refresh_token: refreshToken})
            }); 
            if (response.ok) {
                const data = await response.json();
                localStorage.setItem("token", data.access_token);
                localStorage.setItem("refresh_token", data.refresh_token);
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

            if(response.status === 401 && !endpoint.includes("/login")){
                const refreshed = await this.refresh();
                if (refreshed) {
                    const newToken = localStorage.getItem("token");
                    fetchInit.headers["Authorization"] = `Bearer ${newToken}`;
                    return this.request(endpoint, options);
                } else {
                    localStorage.removeItem("token");
                    localStorage.removeItem("refresh_token");
                    authManager.renderAuthForm();
                    throw new Error("Session expired");
                }
            }

            let data = {};
            const text = await response.text();
            if(!response.ok) {
                throw new Error(data.message || response.statusText);
            }
            try {
                data = text ? JSON.parse(text) : {};
            } catch (error) {
                throw new Error("Server returned non-JSON response");
            }
            
            return data;
        } catch(error){
            throw new Error(error.message);
        }
    }

    async login(Email, Password){
        const data = await this.request("/login", {
            method: "POST",
            body: JSON.stringify({  email: Email,
                                    password: Password,}),
        });

        if (data.access_token) localStorage.setItem("token", data.access_token);
        if (data.refresh_token) localStorage.setItem("refresh_token", data.refresh_token);
        
        return data;
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
        localStorage.removeItem("token");
        localStorage.removeItem("refresh_token");
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

    async setProfile(Firstname, Lastname, Age, City, Info) {
        return this.request("/profile", {
            method: "POST",
            body: JSON.stringify({first_name: Firstname, last_name: Lastname, age: Age, city: City, info: Info}),
        });
    }

}

const api = new Api();