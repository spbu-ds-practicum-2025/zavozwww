class FilmManager {
    render(movie, fromManager) {
        const mainContent = document.getElementById("main-content");
        mainContent.innerHTML = `
            <div class="film-article">
                <button id="back" class="film-article__back">Назад</button>
                <h1 class="film-article__title">${movie.title}</h1> 
                <img class="film-article__poster" src="${movie.imageSrc}" />
                <span class="film-article__year">Год: ${movie.year}</span>
                <span class="film-article__rating">Рейтинг: ${movie.rating ? movie.rating.toFixed(1) : 'Нет оценок'}</span>
                <span film-article__description>Описание: ${movie.description || 'Описание отсутствует'}</span>
                <button class="movie-card__rate-btn film-article__rate-btn" onclick="searchManager.showRate(${movie.id})">Оценить</button>
            </div>
        `;

        document.getElementById('back').addEventListener('click', () => {
            fromManager.render(true);
        });
    }
}