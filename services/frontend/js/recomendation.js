class RecomendationManager {
    async render() {
        let data = await api.recomendation();
        const mainContent = document.getElementById("main-content");
        mainContent.innerHTML = `
            <div class="recomendations-page">
                <h1>Рекомедации для Вас</h1>
                <div id="recomendations-results">
                </div>
            </div>
        `

        const results = document.getElementById("recomendations-results");

        // TO DO:
        // добавить обработчик клика на фильм для recomendations-results
        // сделать при этом клике открытие большой статьи (или что-то типа) о фильме

        results.innerHTML = data.map(movie => 
            `
            <div class="movie-card">
                <h3 class="movie-card__title">${movie.title}</h3>
                <img src="${movie.imageSrc}" class="movie-card__image"/>
                <div class="movie-card__raiting">
                    <span class="movie-card__raiting__value">${movie.rating ? movie.rating.toFixed(1) : "Нет оценок"}</span>
                    <span class="movie-card__raiting__count">${movie.countRaitings ? movie.countRaitings + " оценок" : ""}</span>
                </div>

                <button class="movie-card__rate-btn" onclick="searchManager.showRate(${movie.id})">Оценить</button>
            </div>
            `
        ).join("");
    }
}

let recomendationManager = new RecomendationManager()