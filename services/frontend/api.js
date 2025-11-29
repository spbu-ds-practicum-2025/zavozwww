class Api {
    constructor(){
        this.apiURL = "http://localhost:8080/api";
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
            let data = await response.json();

            if(response.ok){
                return data;
            } else {
                throw response;
            }
        } catch(error){
            if(error.status != 401){
                throw new Error("Network error");
            }
            throw new Error(error.statusText);
        }
    }

    async login(Username, Password){
        return this.request("/login", {
            method: "POST",
            body: JSON.stringify({ username: Username, password: Password })
        });
    }

    async registration(Username, Email, Password){
        return this.request("/register", {
            method: "POST",
            body: JSON.stringify({ email: Email, username: Username, password: Password })
        });
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
        })
    }
}

const api = new Api();