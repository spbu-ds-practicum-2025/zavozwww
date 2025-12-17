(function () {
    console.log("%c MOCK API ENABLED ", "background: #222; color: #bada55; padding: 4px; border-radius: 4px;");

    // --- БАЗА ДАННЫХ (В ПАМЯТИ) ---
    
    // Текущий пользователь
    let currentUser = {
        id: 1,
        email: "test@mail.com",
        username: "testuser",
        firstname: "Иван",
        secondname: "Тестовый",
        city: "Москва",
        age: 25,
        info: "Разработчик на JS",
        countRateFilms: 5,
        countFriends: 2
    };

    // Фильмы
    const movies = [
        { id: 1, title: "Интерстеллар", imageSrc: "https://via.placeholder.com/300x450?text=Interstellar", genres: "триллер, научная фантастика", year: 2014, rating: 4.8, countRaitings: 1500 },
        { id: 2, title: "Начало", imageSrc: "https://via.placeholder.com/300x450?text=Inception", genres: "триллер, боевик", year: 2010, rating: 4.7, countRaitings: 1200 },
        { id: 3, title: "Матрица", imageSrc: "https://via.placeholder.com/300x450?text=Matrix", genres: "боевик, научная фантастика", year: 1999, rating: 4.6, countRaitings: 2000 },
        { id: 4, title: "Крестный отец", imageSrc: "https://via.placeholder.com/300x450?text=Godfather", genres: "драма, криминал", year: 1972, rating: 4.9, countRaitings: 1800 },
        { id: 5, title: "Зеленая миля", imageSrc: "https://via.placeholder.com/300x450?text=GreenMile", genres: "драма, фэнтези", year: 1999, rating: 4.8, countRaitings: 1600 }
    ];

    // Друзья (изначальный список)
    let friends = [
        { id: 101, username: "ivan_petrov", first_name: "Иван", last_name: "Петров", city: "Москва", age: 30, info: "Люблю кино", countRateFilms: 24, countFriends: 15 },
        { id: 102, username: "maria_sid", first_name: "Мария", last_name: "Сидорова", city: "СПб", age: 27, info: "Киноман", countRateFilms: 42, countFriends: 28 }
    ];

    // Все пользователи (для поиска)
    let allUsers = [
        ...friends,
        { id: 103, username: "dmitry_k", first_name: "Дмитрий", last_name: "Козлов", city: "Казань", age: 22, info: "Новичок", countRateFilms: 1, countFriends: 0 },
        { id: 104, username: "anna_v", first_name: "Анна", last_name: "Волкова", city: "Минск", age: 24, info: "Смотрю ужасы", countRateFilms: 10, countFriends: 5 }
    ];

    // Уведомления
    let notifications = [
        { user_id: 123, username: "ivan_petrov", type: "friend_request", timestamp: new Date().toISOString() }
    ];

    // --- ВСПОМОГАТЕЛЬНЫЕ ФУНКЦИИ ---

    function jsonResponse(body, init = {}) {
        const status = init.status ?? 200;
        return Promise.resolve({
            ok: status >= 200 && status < 300,
            status,
            statusText: init.statusText ?? "OK",
            headers: new Headers({ "Content-Type": "application/json" }),
            json: () => Promise.resolve(body),
            text: () => Promise.resolve(JSON.stringify(body)),
        });
    }

    // --- ГЛАВНАЯ ФУНКЦИЯ MOCK FETCH ---

    async function mockFetch(url, options = {}) {
        // Эмуляция задержки сети (300-600мс)
        await new Promise((r) => setTimeout(r, 300 + Math.random() * 300));

        const method = (options.method || "GET").toUpperCase();
        const body = options.body ? JSON.parse(options.body) : null;
        const u = new URL(url.startsWith("http") ? url : `http://localhost${url}`);
        const path = u.pathname;

        // ================= АВТОРИЗАЦИЯ =================

        // POST /login
        if (path.includes("/login") && method === "POST") {
            const { email, password } = body || {};
            if (email === "test@mail.com" && password === "123456") {
                return jsonResponse({
                    access_token: "mock_access_token_" + Date.now(),
                    refresh_token: "mock_refresh_token_" + Date.now(),
                    user_profile: currentUser
                });
            }
            return jsonResponse({ message: "Неверный логин или пароль" }, { status: 401 });
        }

        // POST /register
        if (path.includes("/register") && method === "POST") {
            const { email, username } = body || {};
            currentUser.email = email;
            currentUser.username = username;
            return jsonResponse({
                access_token: "mock_access_token_" + Date.now(),
                refresh_token: "mock_refresh_token_" + Date.now(),
                user_profile: currentUser
            });
        }

        // POST /verify (confirm)
        if (path.includes("/verify") && method === "POST") {
            return jsonResponse({ message: "Email confirmed" });
        }
        
        // POST /resend-verification
        if (path.includes("/resend-verification")) {
             return jsonResponse({ message: "Code resent" });
        }

        // POST /logout
        if (path.includes("/logout")) {
            return jsonResponse({ message: "Logged out" });
        }

        // POST /refresh
        if (path.includes("/refresh")) {
            return jsonResponse({ 
                access_token: "mock_access_token_refreshed_" + Date.now(),
                refresh_token: "mock_refresh_token_new" 
            });
        }

        // ================= ПРОФИЛЬ =================

        // GET /profile
        if (path.endsWith("/profile") && method === "GET") {
            return jsonResponse(currentUser);
        }

        // POST /profile (Update)
        if (path.endsWith("/profile") && method === "POST") {
            // Обновляем текущего пользователя
            Object.assign(currentUser, body);
            return jsonResponse(currentUser);
        }

        // ================= ФИЛЬМЫ =================

        // GET /recomendations
        if (path.includes("/recomendations")) {
            // Возвращаем случайные фильмы + перемешиваем
            const shuffled = [...movies].sort(() => 0.5 - Math.random());
            return jsonResponse(shuffled);
        }

        // GET /search
        if (path.includes("/search")) {
            const query = (u.searchParams.get("query") || "").toLowerCase();
            const genre = (u.searchParams.get("genre") || "").toLowerCase();
            
            const filtered = movies.filter(m => {
                const matchTitle = m.title.toLowerCase().includes(query);
                const matchGenre = !genre || m.genres.toLowerCase().includes(genre);
                return matchTitle && matchGenre;
            });
            return jsonResponse({ movies: filtered });
        }

        // POST /rating
        if (path.match(/\/movies\/\d+\/rating/)) {
            return jsonResponse({ message: "Оценка сохранена" });
        }

        // GET /ratedFilms
        if (path.includes("/ratedFilms")) {
             // Вернем пару фильмов как "оцененные"
             return jsonResponse(movies.slice(0, 3));
        }

        // ================= ДРУЗЬЯ =================

        // GET /friends
        if (path.endsWith("/friends") && method === "GET") {
            return jsonResponse(friends);
        }

        // POST /friends/searchFriends
        if (path.includes("/friends/searchFriends")) {
            const { username } = body || {};
            const searchName = (username || "").toLowerCase();
            
            // Ищем среди всех юзеров, кроме себя
            const result = allUsers.filter(u => 
                (u.username.toLowerCase().includes(searchName) || 
                 u.first_name.toLowerCase().includes(searchName))
            );
            return jsonResponse(result); // Бэкенд возвращает массив
        }

        // POST /friends/requests/accept (Отправка и Принятие)
        // В вашем API.js методы addFriend и acceptRequest стучатся на один и тот же URL
        if (path.includes("/friends/requests/accept")) {
            if (body.to_name) {
                 return jsonResponse({ message: `Запрос отправлен пользователю ${body.to_name}` });
            }
            if (body.username) {
                 // Добавляем друга в список
                 const newFriend = allUsers.find(u => u.username === body.username);
                 if (newFriend && !friends.find(f => f.username === newFriend.username)) {
                     friends.push(newFriend);
                 }
                 return jsonResponse({ message: `Заявка от ${body.username} принята` });
            }
        }

        // POST /friends/req/request_id (Отклонение)
        if (path.includes("/friends/req/")) {
            return jsonResponse({ message: "Заявка отклонена" });
        }

        // ================= УВЕДОМЛЕНИЯ =================
        if (path.includes("/notifications")) {
             return jsonResponse({ notifications: notifications });
        }

        // Если путь не найден
        console.warn(`[MockAPI] 404 Not Found: ${method} ${path}`);
        return jsonResponse({ message: "Not found" }, { status: 404 });
    }

    // --- ПОДМЕНА ГЛОБАЛЬНОГО FETCH ---
    window.fetch = mockFetch;


    // --- MOCK WEBSOCKET ---
    const OriginalWebSocket = window.WebSocket;
    class MockWebSocket {
        constructor(url) {
            console.log(`[MockWS] Connecting to ${url}`);
            this.readyState = 0; // CONNECTING
            setTimeout(() => {
                this.readyState = 1; // OPEN
                if (this.onopen) this.onopen();
                
                // Шлем тестовое уведомление через 2 секунды
                setTimeout(() => {
                    const msg = {
                        type: "friend_request",
                        username: "mock_friend",
                        user_id: 999
                    };
                    if (this.onmessage) this.onmessage({ data: JSON.stringify(msg) });
                }, 2000);

            }, 500);
        }
        send(data) { console.log("[MockWS] Sent:", data); }
        close() { this.readyState = 3; }
    }
    window.WebSocket = MockWebSocket;

})();
