class Api {
    constructor(){
        this.apiURL = "http://localhost/filmbuddy";
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

    async request(endpoint, options = {}, urlService = this.apiURL){
        loaderManager.show();

        let url = `${urlService}${endpoint}`;
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
        
        console.log(`Request: ${fetchInit.method || 'GET'} ${url}`, fetchInit);

        try {
            let response = await fetch(url, fetchInit);

            if(response.status === 401 && !endpoint.includes("/login") && !endpoint.includes("/register")){
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

            const text = await response.text();
            let data = {};

            try {
                data = text ? JSON.parse(text) : {};
            } catch (error) {
                throw new Error("Server returned non-JSON response");
            }

            if (!response.ok) {
                throw new Error(data.message || response.statusText);
            }
            
            return data;
        } catch(error){
            throw new Error(error.message);
        } finally {
            loaderManager.hide();
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
        return this.request("/verify", {
            method: "POST",
            body: JSON.stringify({code: userCode, email: Email}),
        })
    }

    async again(Email) {
        return this.request("/resend-verification", {
            method: "POST",
            body: JSON.stringify({email: Email}),
        })
    }

    async searchMovie(searchQuery, genre) {
        return this.request("/search", {
            method: "POST",
            body: JSON.stringify({
                title: searchQuery,
                genre: genre === "all" ? "" : genre
            })
        }, "http://localhost/recsys");
    }

    async rating(ID, Rating, Message){
        return this.request("/ratings", {
                    method: "POST",
                    body: JSON.stringify({ film_id: ID, grade: Rating, review: Message, username: "" }),
                }, "http://localhost/social");
    }

    async recomendation() {
        try {
            const ratings = await this.getRatedMovies();
            const grades = ratings.map(r => ({
                film_id: r.film_id,
                grade: r.grade
            }));

            return this.request("/recommend", {
                method: "POST",
                body: JSON.stringify(grades)
            }, "http://localhost/recsys");
        } catch (error) {
            console.error("Failed to get recommendations:", error);
            return [];
        }
    }

    async searchFriend(Name) {
        return this.request(`/friends/searchFriends`, {
            method: "POST",
            body: JSON.stringify({ username: Name }),
        });
    }

    async addFriend(FromName, ToName) {
        return this.request("/friends/requests", {
           method: "POST",
           body: JSON.stringify({from_username: FromName, to_username: ToName}),
       }, "http://localhost/social");
    }

    async getNotice() {
        const data = await api.request("/friends/requests", {}, "http://localhost/social");
        return data;
    }

    async acceptRequest(FromName, ToName) {
        return this.request("/friends/requests/accept", {
            method: "POST",
            body: JSON.stringify({from_username: FromName, to_username: ToName}),
        }, "http://localhost/social");
    }

    async declineRequest(Request_id) {
        return this.request("/friends/requests/reject", {
            method: "POST",
            body: JSON.stringify({request_id: Request_id}),
        }, "http://localhost/social");
    }

    async logout(refreshToken){
        try {
            await this.request("/logout", {
                method: "POST",
                body: JSON.stringify({refresh_token: refreshToken})
            });
        } finally {
            localStorage.removeItem("token");
            localStorage.removeItem("refresh_token");
        }
    }

    async getProfile() {
        const data = await this.request("/profile");
        return data;
    }

    async getFriends() {
        const data = await this.request("/friends", {}, "http://localhost/social");
        return data;
    }

    async setProfile(Username, Firstname, Lastname, Age, City, Info) {
        return this.request("/profile", {
            method: "POST",
            body: JSON.stringify({username: Username, first_name: Firstname, last_name: Lastname, age: Age, city: City, info: Info}),
        });
    }

    async getRatedMovies() {
        const data = await this.request("/ratings", {}, "http://localhost/social");
        return data;
    }
    async getFilm(filmID) {
        const data = await this.request(`/movie/${filmID}`, {}, "http://localhost/recsys");
        console.log("FILM DATA: ", data);
        return data.movie_title;
    }

    async getMoviesBatch(filmIDs) {
        const data = await this.request("/movies/batch", {
            method: "POST",
            body: JSON.stringify({ movie_ids: filmIDs })
        }, "http://localhost/recsys");
        return data;
    }
}

const api = new Api();