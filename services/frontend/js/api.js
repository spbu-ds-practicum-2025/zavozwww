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
                throw new Error(`Error: ${response.status}, ${response.statusText}` || "Server error");
            }
        } catch(error){
            throw new Error(error.message || "Network error");
        }
    }

    async login(username, password){
        return this.request("/login", {
            method: "POST",
            body: JSON.stringify({ username, password })
        });
    }

    async registration(username, email, password){
        return this.request("/register", {
            method: "POST",
            body: JSON.stringify({ email, username, password })
        });
    }

    async searchMovie(searchQuery, genre) {
        return this.request(`/search?query=${encodeURIComponent(searchQuery)}&genre=${genre}`);
    }

    async rating(rating, message){
        return this.request(`/movies/${id}/rating`, {
                    method: "POST",
                    body: JSON.stringify({rating, message}),
                    headers: { "Content-Type": "application/json" }
                });
    }
}

const api = new Api();