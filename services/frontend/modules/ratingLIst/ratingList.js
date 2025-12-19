class RatingList {
    async getFilmName(filmID){
        try {
            const filmName = await api.getFilm(filmID);
            return filmName;
        } catch (error) {
            console.log(error.message);
            tempNotice.error("Ошибка, попробуйте еще раз через некоторое время");
        }
    }

    async render() {
        const mainContent = document.getElementById("main-content");
        mainContent.innerHTML = `
        <div class="main-page">
            <h1 class="rated-movies__title">Ваши оценки</h1>
            <div class="rated-movies__results"></div>
        </div>
        `
        try {
            const ratedMovies = await api.getRatedMovies();
            console.log("DATA: ", ratedMovies);
            const results = document.querySelector(".rated-movies__results");
            if(ratedMovies.length === 0) {
                results.innerHTML = `
                    <p class="rated-movies__results__without">У Вас пока нет оцененных фильмов</p>
                `
            } else {
                const filmIds = ratedMovies.map(f => f.film_id);
                const moviesInfo = await api.getMoviesBatch(filmIds);
                
                const moviesMap = {};
                moviesInfo.forEach(m => {
                    moviesMap[m.movie_id] = m.movie_title;
                });

                const moviesWithNames = ratedMovies.map(film => ({
                    ...film,
                    name: moviesMap[film.film_id] || "Unknown Title"
                }));

                results.innerHTML = moviesWithNames.map((film) => 
                    `
                    <div class="rated-movies__results__film">
                        <p class="rated-movies__results__film__title">${film.name}</p>
                        <p class="rated-movies__results__film__createdAt">${new Date(film.created_at).toLocaleDateString()}</p>
                        <p class="rated-movies__results__film__grade">&#9733; ${film.grade}</p>
                        <p class="rated-movies__results__film__review">${film.review}</p>
                    </div>
                    `
                ).join("");
            }
        } catch(error) {
            console.log(error.message);
            tempNotice.error("Ошибка, попробуйте еще раз через некоторое время");
        }
        
    }
}

let ratingListManager = new RatingList();