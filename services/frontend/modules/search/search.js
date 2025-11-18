class SearchManager{
    constructor(){
        this.filmManager = new FilmManager();
        this.currentMovies = [];
        this.currentSearchQuery = '';
        this.currentGenre = 'all';
        
        window.addEventListener('popstate', (event) => {
            if (event.state && event.state.page === 'search') {
                this.render(true)
            }
        });
    }

    render(restore = false){
        const mainContent = document.getElementById("main-content");
        mainContent.innerHTML = `
            <div class="search-page">
                <h1 class="search-page__title">Поиск фильмов</h1>
                <div class="search-box">
                    <input type="text" id="film-search" class="search-box__input" placeholder="Название фильма...">
                    <div class="search-box__filters">
                        <select class="filters" id="genres">
                            <option value="all">Все жанры</option>
                            <option value="drama">Драма</option>
                            <option value="comedy">Комедия</option>
                            <option value="action">Боевик</option>
                            <option value="adventure">Приключения</option>
                            <option value="thriller">Триллер</option>
                            <option value="horror">Ужасы</option>
                            <option value="romantic-comedy">Романтическая комедия </option>
                            <option value="western">Вестерн</option>
                            <option value="animated">Мультфильмы</option>
                            <option value="sci-fi">Научная фантастика</option>
                        </select>
                        <svg class="select-icon" focusable="false" width="21" height="12" viewBox="0 0 21 12" fill="none" xmlns="http://www.w3.org/2000/svg">
                            <path d="M1.5 1.5L9.98528 9.98528M10.1465 9.98529L18.6318 1.50001" stroke-width="3" stroke-linecap="round"/>
                        </svg>

                    </div>
                    <button id="search-btn" class="search-box__button">Найти</button>
                </div>
                <div id="films-results" class="results search-page__result">
                    <p>Введите поисковый запрос</p>
                </div>
            </div>
        `

        if (restore) {
            this.restoreSearchState();
        }

        document.getElementById("search-btn").addEventListener("click", () => {
            this.search();
        })

        document.getElementById("film-search").addEventListener("keypress", (event) => {
            if(event.key == "Enter"){
                this.search();
            }
        })

        document.getElementById("genres").addEventListener("change", () => {
            this.search();
        })

        if (this.currentMovies.length > 0 && restore) {
            this.restoreSearchState();
        }
    }

    async search(){
        let searchQuery = document.getElementById("film-search").value;
        //searchQuery может быть пустым, в таком случае осуществляем поиск только по жанру
        let genre = document.getElementById("genres").value;

        this.currentSearchQuery = searchQuery;
        this.currentGenre = genre;

        try{
            let data = await api.searchMovie(searchQuery, genre);
            this.renderMovies(data.movies);
        } catch(error) {
            tempNotice.error("Ошибка поиска, попробуйте еще раз через некоторое время");
            console.log(`Error: ${error.message}`);
        }
    }

    restoreSearchState(){
        document.getElementById('film-search').value = this.currentSearchQuery;
        document.getElementById('genres').value = this.currentGenre;
        this.renderMovies(this.currentMovies);
    }

    renderMovies(movies){
        const results = document.getElementById("films-results");
        
        if(movies.length == 0){
            results.innerHTML = `
            <p>Фильмы не найдены :(</p>
            `

            return;
        }

        // очень надо чтобы жанры с tmdb были преобразованы в одну строку
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
                this.currentSearchQuery = document.getElementById("film-search").value;
                this.currentGenre = document.getElementById("genres").value;
                this.filmManager.render(film, searchManager);
            }
        })

        this.currentMovies = movies;
    }

    showRate(movieId){
        const rateContainer = document.createElement("div");
        rateContainer.classList.add("rate-container");
        rateContainer.innerHTML = `
            <div class="rating">
                <button id="close" class="close">
                    <svg class="cross" viewBox="0 0 20 20" fill="none" xmlns="http://www.w3.org/2000/svg">
<path d="M1.5 18.5L18.5 1.5M1.5 1.5L18.5 18.5"  stroke-width="3" stroke-linecap="round"/>
</svg>

                </button>
                <h1 class="rate-title">Как Вам этот фильм?</h1>
                <div class="stars-group">
                    <div class="active star">
                    <svg class="star-svg" viewBox="0 0 27 26" fill="none" xmlns="http://www.w3.org/2000/svg">
<path d="M12.1901 0.690964C12.4894 -0.230347 13.7928 -0.230344 14.0922 0.690966L16.5088 8.12866C16.6427 8.54068 17.0267 8.81964 17.4599 8.81964H25.2803C26.2491 8.81964 26.6518 10.0593 25.8681 10.6287L19.5412 15.2254C19.1908 15.4801 19.0441 15.9314 19.178 16.3434L21.5946 23.7811C21.894 24.7024 20.8395 25.4686 20.0558 24.8992L13.7289 20.3024C13.3784 20.0478 12.9038 20.0478 12.5533 20.3024L6.22644 24.8992C5.44273 25.4686 4.38825 24.7024 4.68761 23.7811L7.10426 16.3434C7.23813 15.9314 7.09148 15.4801 6.74099 15.2254L0.414104 10.6287C-0.369609 10.0593 0.033169 8.81964 1.00189 8.81964H8.82235C9.25557 8.81964 9.63953 8.54068 9.7734 8.12866L12.1901 0.690964Z"/>
</svg>

                    </div>
                    <div class="star">
                    <svg class="star-svg" viewBox="0 0 27 26" fill="none" xmlns="http://www.w3.org/2000/svg">
<path d="M12.1901 0.690964C12.4894 -0.230347 13.7928 -0.230344 14.0922 0.690966L16.5088 8.12866C16.6427 8.54068 17.0267 8.81964 17.4599 8.81964H25.2803C26.2491 8.81964 26.6518 10.0593 25.8681 10.6287L19.5412 15.2254C19.1908 15.4801 19.0441 15.9314 19.178 16.3434L21.5946 23.7811C21.894 24.7024 20.8395 25.4686 20.0558 24.8992L13.7289 20.3024C13.3784 20.0478 12.9038 20.0478 12.5533 20.3024L6.22644 24.8992C5.44273 25.4686 4.38825 24.7024 4.68761 23.7811L7.10426 16.3434C7.23813 15.9314 7.09148 15.4801 6.74099 15.2254L0.414104 10.6287C-0.369609 10.0593 0.033169 8.81964 1.00189 8.81964H8.82235C9.25557 8.81964 9.63953 8.54068 9.7734 8.12866L12.1901 0.690964Z"/>
</svg>

                    </div>
                    <div class="star">
                    <svg class="star-svg" viewBox="0 0 27 26" fill="none" xmlns="http://www.w3.org/2000/svg">
<path d="M12.1901 0.690964C12.4894 -0.230347 13.7928 -0.230344 14.0922 0.690966L16.5088 8.12866C16.6427 8.54068 17.0267 8.81964 17.4599 8.81964H25.2803C26.2491 8.81964 26.6518 10.0593 25.8681 10.6287L19.5412 15.2254C19.1908 15.4801 19.0441 15.9314 19.178 16.3434L21.5946 23.7811C21.894 24.7024 20.8395 25.4686 20.0558 24.8992L13.7289 20.3024C13.3784 20.0478 12.9038 20.0478 12.5533 20.3024L6.22644 24.8992C5.44273 25.4686 4.38825 24.7024 4.68761 23.7811L7.10426 16.3434C7.23813 15.9314 7.09148 15.4801 6.74099 15.2254L0.414104 10.6287C-0.369609 10.0593 0.033169 8.81964 1.00189 8.81964H8.82235C9.25557 8.81964 9.63953 8.54068 9.7734 8.12866L12.1901 0.690964Z"/>
</svg>

                    </div>
                    <div class="star">
                    <svg class="star-svg" viewBox="0 0 27 26" fill="none" xmlns="http://www.w3.org/2000/svg">
<path d="M12.1901 0.690964C12.4894 -0.230347 13.7928 -0.230344 14.0922 0.690966L16.5088 8.12866C16.6427 8.54068 17.0267 8.81964 17.4599 8.81964H25.2803C26.2491 8.81964 26.6518 10.0593 25.8681 10.6287L19.5412 15.2254C19.1908 15.4801 19.0441 15.9314 19.178 16.3434L21.5946 23.7811C21.894 24.7024 20.8395 25.4686 20.0558 24.8992L13.7289 20.3024C13.3784 20.0478 12.9038 20.0478 12.5533 20.3024L6.22644 24.8992C5.44273 25.4686 4.38825 24.7024 4.68761 23.7811L7.10426 16.3434C7.23813 15.9314 7.09148 15.4801 6.74099 15.2254L0.414104 10.6287C-0.369609 10.0593 0.033169 8.81964 1.00189 8.81964H8.82235C9.25557 8.81964 9.63953 8.54068 9.7734 8.12866L12.1901 0.690964Z"/>
</svg>

                    </div>
                    <div class="star">
                    <svg class="star-svg" viewBox="0 0 27 26" fill="none" xmlns="http://www.w3.org/2000/svg">
<path d="M12.1901 0.690964C12.4894 -0.230347 13.7928 -0.230344 14.0922 0.690966L16.5088 8.12866C16.6427 8.54068 17.0267 8.81964 17.4599 8.81964H25.2803C26.2491 8.81964 26.6518 10.0593 25.8681 10.6287L19.5412 15.2254C19.1908 15.4801 19.0441 15.9314 19.178 16.3434L21.5946 23.7811C21.894 24.7024 20.8395 25.4686 20.0558 24.8992L13.7289 20.3024C13.3784 20.0478 12.9038 20.0478 12.5533 20.3024L6.22644 24.8992C5.44273 25.4686 4.38825 24.7024 4.68761 23.7811L7.10426 16.3434C7.23813 15.9314 7.09148 15.4801 6.74099 15.2254L0.414104 10.6287C-0.369609 10.0593 0.033169 8.81964 1.00189 8.81964H8.82235C9.25557 8.81964 9.63953 8.54068 9.7734 8.12866L12.1901 0.690964Z"/>
</svg>

                    </div>
                </div>
                <textarea class="review" placeholder="Напишите свое мнение о фильме"></textarea>
                <button id="rate" class="rate-btn">Оценить</button>
            </div>
        `

        document.body.appendChild(rateContainer);
        let rating = 1;

        const stars = document.querySelectorAll(".star");

        stars.forEach((item, index) => {
            item.addEventListener("mouseenter", () => {
                updateClass(index + 1);
            })
        })

        stars.forEach((item, index) => {
            item.addEventListener("click", () => {
                rating = index + 1;
                console.log(index + 1);
                updateClass(rating);
            }
            );
        });

        const updateClass = (RATE) => {
            stars.forEach((item, index) => {
                if(index < RATE) {
                    item.classList.add("active");
                } else {
                    item.classList.remove("active");
                }
            });
        }

        document.getElementById("close").addEventListener("click", () => {
            document.body.removeChild(rateContainer);
        })

        document.getElementById("rate").addEventListener("click", async () => {
            try{
                await api.rating(movieId, parseInt(rating), document.querySelector(".review").value);
                document.body.removeChild(rateContainer);
                tempNotice.success("Оценка сохранена!");
            } catch(error) {
                document.body.removeChild(rateContainer);
                tempNotice.error("Ошибка, попробуйте еще раз через некоторое время");
            }
        })
    }
}

const searchManager = new SearchManager();