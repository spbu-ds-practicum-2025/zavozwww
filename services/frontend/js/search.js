class SearchManager{
    render(){
        const mainContent = document.getElementById("main-content");
        mainContent.innerHTML = `
            <div class="films-page">
                <h1>Поиск фильмов</h1>
                <div class="search-box">
                    <input type="text" id="film-search" placeholder="Название фильма...">
                    <button id="search-btn">Найти</button>
                    <div class="filters">
                        <select id="genres">
                            <option value="all">Все жанры</option>
                            <option value="drama">Драма</option>
                            <option value="comedy">Комедия</option>
                            <option value="action">Боевик</option>
                            <option value="sci-fi">Научная фантастика</option>
                        </select>
                    </div>
                </div>
                <div id="films-results">
                    <p>Введите поисковый запрос</p>
                </div>
            </div>
        `

        document.getElementById("search-btn").addEventListener("click", () => {
            this.search();
        })

        document.getElementById("film-search").addEventListener("keypress", (event) => {
            if(event.key == "Enter"){
                this.search();
            }
        })
    }

    async search(){
        let searchQuery = document.getElementById("film-search").value;
        let genre = document.getElementById("genres").value;

        try{
            let data = await api.searchMovie(searchQuery, genre);
            this.renderMovies(data.movies);
        } catch(error) {
            console.log(`Error ${error.message}`);
        }
    }

    renderMovies(movies){
        const results = document.getElementById("films-results");
        
        if(movies.length == 0){
            results.innerHTML = `
            <p>Фильмы не найдены :(</p>
            `
        }

        //отрисовали карточки фильмов, попавших в поиск
        results.innerHTML = movies.map(movie => 
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

    //TO DO:
    // заменить promt на модальное окно с звездами и text area
    // заменить логирование на поп апы

    async showRate(movieId){
        let rating = prompt("Пожалуйста оценить фильм");
        let message = prompt("Ваше мнение о фильме")

        if(rating && rating >= 1 && rating <= 10){
            try{
                await api.rating(movieId, parseInt(rating), message);
                alert("Оценка сохранена!");
            } catch(error) {
                console.log(`Error ${error.message}`);
            }
        }
    }
}

const searchManager = new SearchManager();