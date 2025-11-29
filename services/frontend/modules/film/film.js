class FilmManager {
    render(movie, fromManager) {
        const mainContent = document.getElementById("main-content");
        mainContent.innerHTML = `
        <div class="film-page">
            <a id="back" class="back-btn">← Назад</a>
            <div class="film-article">
                <img class="film-article__poster" src="${movie.imageSrc}" />
                <div class="film-info">
                    <h1 class="film-article__title">${movie.title}</h1> 
                    <span class="film-article__year">Год: ${movie.year}</span>
                    <span class="film-article__rating">Рейтинг: ${movie.rating ? movie.rating.toFixed(1) : 'Нет оценок'}</span>
                    <span class="film-article__description">Описание: ${movie.description || 'Описание отсутствует'}</span>
                    <button class="movie-card__rate-btn film-article__rate-btn" onclick="searchManager.showRate(${movie.id})">Оценить</button>
                </div>
            </div>
        </div>
        `;

        document.getElementById('back').addEventListener('click', () => {
            fromManager.render(true);
        });
    }
}