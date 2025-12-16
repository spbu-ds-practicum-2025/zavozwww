class RecomendationManager {
    constructor(){
        this.filmManager = new FilmManager();
        this.currentMovies = [];
        
        window.addEventListener('popstate', (event) => {
            if (event.state && event.state.page === 'recomendation') {
                this.render()
            }
        });
    }

    async render() {
        let movies = await api.recomendation();
        const mainContent = document.getElementById("main-content");
        mainContent.innerHTML = `
            <div class="recomendations-page">
                <h1 class="recomendations-page__title">Рекомедации</h1>
                <div class="results" id="recomendations-results">
                </div>
            </div>
        `

        const results = document.getElementById("recomendations-results");

        results.innerHTML = movies.map(movie => 
            `
            <div class="movie-card" data-movie-id="${movie.id}">
                <img src="${movie.imageSrc}" class="movie-card__image"/>
                <div class="movie-card-additional">
                    <div class="movie-card-info">
                        <h3 class="movie-card__title">${movie.title}</h3>
                        <span class="movie-card__year">${movie.year} год</span>
                        <span class="movie-card__genres">${movie.genres}</span>
                    </div>
                    <div class="movie-card-rate">
                        <span class="movie-card__raiting"> Оценка: ${movie.rating ? movie.rating.toFixed(1): "нет оценок"}</span><br>
                        <button class="movie-card__rate-btn" onclick="searchManager.showRate(${movie.id})">Оценить</button>
                    </div>
                </div>
            </div>
            `
        ).join("");

        results.addEventListener("click", (event) => {
            let movie = event.target.closest(".movie-card");
            if(movie && !event.target.classList.contains("movie-card__rate-btn")){
                let film = movies.find((item) => item.id === parseInt(movie.dataset.movieId));
                this.filmManager.render(film, recomendationManager);
            }
        })

        this.currentMovies = movies;
    }
}

const recomendationManager = new RecomendationManager();