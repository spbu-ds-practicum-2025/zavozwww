class RatingList {

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
            const results = document.querySelector(".rated-movies__results");
            if(ratedMovies.length === 0) {
                results.innerHTML = `
                    <p class="rated-movies__results__without">У Вас пока нет оцененных фильмов</p>
                `
            } else {
                results.innerHTML = ratedMovies.map((film) => 
                    `
                    <div class="rated-movies__results__film">
                        <p class="rated-movies__results__film__title>${film.title}</p>
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