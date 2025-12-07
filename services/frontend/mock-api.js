// mock-api.js
(function () {
  // Сохраняем оригинальные реализации
  const originalFetch = window.fetch ? window.fetch.bind(window) : null;
  const OriginalWebSocket = window.WebSocket;

  // Простейшее состояние мока
  let currentToken = "mock_jwt_token_initial";
  let currentUser = {
    id: 1,
    email: "test@mail.com",
    username: "testuser",
    firstname: "j",  // Пусто
    secondname: "", // Пусто
    city: "",       // Пусто
    age: null,
    info: ""
  };

  const movies = [
    {
      id: 1,
      title: "Интерстеллар",
      imageSrc: "https://clck.ru/3Q3i7R",
      genres: "триллер, научная фантастика",
      year: 2014,
      rating: 4.8,
      countRaitings: 1500,
    },
    {
      id: 2,
      title: "Начало",
      imageSrc: "https://clck.ru/3Q3i7u",
      genres: "триллер, боевик",
      year: 2010,
      rating: 4.7,
      countRaitings: 1200,
    },
    {
      id: 3,
      title: "Матрица",
      imageSrc: "https://clck.ru/3Q3i9s",
      genres: "боевик, научная фантастика, триллер",
      year: 1999,
      rating: 4.6,
      countRaitings: 2000,
    },
    {
      id: 4,
      title: "Крестный отец",
      imageSrc: "https://clck.ru/3Q3iAB",
      genres: "боевик",
      year: 1972,
      rating: 4.9,
      countRaitings: 1800,
    },
  ];

  const friends = [
    {
      id: 101,
      name: "Иван Петров",
      firstname: "Иван",
      secondname: "Петров",
      city: "Москва",
      age: 30,
      info: "Люблю кино и программирование",
      regDate: "15.03.2021",
      countRateFilms: 24,
      countFriends: 15,
    },
    {
      id: 102,
      name: "Мария Сидорова",
      firstname: "Мария",
      secondname: "Сидорова",
      city: "Санкт-Петербург",
      age: 27,
      info: "Киноман",
      regDate: "02.11.2020",
      countRateFilms: 42,
      countFriends: 28,
    },
  ];

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

  async function mockFetch(url, options = {}) {
    // Если вдруг передали Request — отдаем в оригинал
    if (!(typeof url === "string")) {
      return originalFetch ? originalFetch(url, options) : jsonResponse({});
    }

    const method = (options.method || "GET").toUpperCase();
    const body = options.body ? JSON.parse(options.body) : null;

    console.log("MockFetch:", { url, method, body });
    // Эмуляция задержки сети
    await new Promise((r) => setTimeout(r, 300));

    if (!url.includes("/api")) {
      return originalFetch ? originalFetch(url, options) : jsonResponse({});
    }

    const u = new URL(url, window.location.origin);
    const path = u.pathname; // например: /api/login

    // ====== АУТЕНТИФИКАЦИЯ ======

    // POST /api/login
    if (path.includes("/api/login") && method === "POST") {
      const { email, password } = body || {};
      if (email === "test@mail.com" && password === "123456") {
        currentToken = "mock_jwt_token_" + Date.now();
        return jsonResponse({
          access_token: currentToken,
          user_profile: currentUser,
        });
      }
      return jsonResponse(
        { error: "Invalid credentials" },
        { status: 401, statusText: "invalid data" }
      );
    }

    // POST /api/register
    if (path.includes("/api/register") && method === "POST") {
      const { email, username, password } = body || {};
      currentUser = {
        ...currentUser,
        email: email || "new@mail.com",
        username: username || "newuser",
      };
      currentToken = "mock_jwt_token_" + Date.now();
      return jsonResponse({
        access_token: currentToken,
        user_profile: currentUser,
      });
    }

    // POST /api/confirm
    if (path.includes("/api/confirm") && method === "POST") {
      return jsonResponse({ message: "Code confirmed", status: 200 });
    }

    // POST /api/again
    if (path.includes("/api/again") && method === "POST") {
      return jsonResponse({ message: "Code sent again" });
    }

    // POST /api/refresh
    if (path.includes("/api/refresh") && method === "POST") {
      currentToken = "mock_jwt_token_" + Date.now();
      return jsonResponse({ access_token: currentToken });
    }

    // POST /api/logout
    if (path.includes("/api/logout") && method === "POST") {
      currentToken = null;
      try {
        localStorage.removeItem("token");
      } catch {}
      return jsonResponse({ message: "Logged out" });
    }

    // ====== ПРОФИЛЬ ======

    // GET /api/profile
    if (path.includes("/api/profile") && method === "GET") {
      return jsonResponse(currentUser);
    }

    // POST /api/profile
    if (path.includes("/api/profile") && method === "POST") {
      const {
        firstname,
        secondname,
        age,
        city,
        about, // см. setProfile в api.js
      } = body || {};
      currentUser = {
        ...currentUser,
        firstname: firstname ?? currentUser.firstname,
        secondname: secondname ?? currentUser.secondname,
        age: age ?? currentUser.age,
        city: city ?? currentUser.city,
        info: about ?? currentUser.info,
      };
      return jsonResponse(currentUser);
    }

    // ====== ФИЛЬМЫ ======

    // GET /api/search?query=&genre=
    if (path.includes("/api/search") && method === "GET") {
      const query = (u.searchParams.get("query") || "").toLowerCase();
      const genre = (u.searchParams.get("genre") || "").toLowerCase();

      const filtered = movies.filter((m) => {
        const byTitle = m.title.toLowerCase().includes(query);
        const byGenre =
          !genre || m.genres.toLowerCase().includes(genre.toLowerCase());
        return byTitle && byGenre;
      });

      return jsonResponse({ movies: filtered });
    }

    // POST /api/movies/:id/rating
    if (path.match(/\/api\/movies\/\d+\/rating$/) && method === "POST") {
      const match = path.match(/\/api\/movies\/(\d+)\/rating$/);
      const movieId = match ? Number(match[1]) : null;
      const { rating, review } = body || {};
      console.log("Mock rating:", { movieId, rating, review });
      return jsonResponse({ message: "Rating saved successfully" });
    }

    // GET /api/recomendations
    if (path.includes("/api/recomendations") && method === "GET") {
      return jsonResponse(movies  , {
        headers: { "Content-Type": "application/json" },
      });
    }

    // ====== ДРУЗЬЯ ======

    // GET /api/friends
    if (path.includes("/api/friends") && method === "GET") {
      return jsonResponse(friends);
    }

    // /api/friends/search?name=...
    if (path.includes("/api/friends/search")) {
      const searchName = (u.searchParams.get("name") || "").toLowerCase();
      let result = [];
      if (searchName.trim() !== "") {
        result = friends.filter((f) =>
          f.name.toLowerCase().includes(searchName)
        );
      }
      return jsonResponse(
        { friends: result },
        { headers: { "Content-Type": "application/json" } }
      );
    }

    // POST /api/friends/request
    if (path.includes("/api/friends/request") && method === "POST") {
      // Можно было бы добавлять в список, но для теста достаточно сообщения
      return jsonResponse(
        { message: "Friend request sent successfully" },
        { headers: { "Content-Type": "application/json" } }
      );
    }

    // POST /api/friends/accept/request_id
    if (path.includes("/api/friends/accept/request_id") && method === "POST") {
      const { user_id } = body || {};
      return jsonResponse(
        { message: `Friend request from user ${user_id} accepted successfully` },
        { headers: { "Content-Type": "application/json" } }
      );
    }

    // POST /api/friends/decline/request_id
    if (
      path.includes("/api/friends/decline/request_id") &&
      method === "POST"
    ) {
      const { user_id } = body || {};
      return jsonResponse(
        { message: `Friend request from user ${user_id} declined` },
        { headers: { "Content-Type": "application/json" } }
      );
    }

    // ====== УВЕДОМЛЕНИЯ ======

    // GET /api/notifications
    if (path.includes("/api/notifications") && method === "GET") {
      const now = Date.now();
      const mockNotifications = [
        {
          user_id: 123,
          username: "Иван Петров",
          type: "friend_request",
          timestamp: new Date(now).toISOString(),
        },
        {
          user_id: 124,
          username: "Мария Сидорова",
          type: "friend_request",
          timestamp: new Date(now - 3600000).toISOString(),
        },
      ];
      return jsonResponse(
        { notifications: mockNotifications },
        { headers: { "Content-Type": "application/json" } }
      );
    }

    // POST /api/test/notification — ручной триггер
    if (path.includes("/api/test/notification") && method === "POST") {
      if (window.triggerMockNotification) {
        window.triggerMockNotification();
      }
      return jsonResponse(
        { message: "Test notification sent" },
        { headers: { "Content-Type": "application/json" } }
      );
    }

    // Всё остальное — в реальный бэкенд, если он есть
    return originalFetch ? originalFetch(url, options) : jsonResponse({});
  }

  // Переопределяем fetch
  if (originalFetch) {
    window.fetch = mockFetch;
  } else {
    window.fetch = mockFetch;
  }

  // ====== Mock WebSocket для уведомлений ======

  class MockWebSocket {
    constructor(url) {
      this.url = url;
      this.readyState = 0; // CONNECTING
      this.onopen = null;
      this.onmessage = null;
      this.onerror = null;
      this.onclose = null;

      if (!window.activeMockSockets) {
        window.activeMockSockets = [];
      }
      window.activeMockSockets.push(this);

      // Имитируем подключение
      setTimeout(() => {
        this.readyState = 1; // OPEN
        if (this.onopen) {
          this.onopen(new Event("open"));
        }

        // Пара тестовых уведомлений
        this._sendMockNotification(123, "test_friend1");
        setTimeout(
          () => this._sendMockNotification(321, "test_friend2"),
          3000
        );
      }, 100);
    }

    _sendMockNotification(id, username) {
      if (this.onmessage && this.readyState === 1) {
        const notification = {
          user_id: id,
          username,
          type: "friend_request",
          timestamp: new Date().toISOString(),
        };
        this.onmessage({ data: JSON.stringify(notification) });
      }
    }

    send(data) {
      console.log("Mock WebSocket sent:", data);
    }

    close() {
      this.readyState = 3; // CLOSED
      if (this.onclose) {
        this.onclose(new Event("close"));
      }
      if (window.activeMockSockets) {
        const idx = window.activeMockSockets.indexOf(this);
        if (idx > -1) {
          window.activeMockSockets.splice(idx, 1);
        }
      }
    }
  }

  // Подменяем WebSocket
  if (OriginalWebSocket) {
    window.WebSocket = MockWebSocket;
  } else {
    window.WebSocket = MockWebSocket;
  }

  // Глобальная функция для ручного триггера уведомлений
  window.triggerMockNotification = function () {
    if (!window.activeMockSockets) return;
    window.activeMockSockets.forEach((socket) => {
      if (socket.onmessage && socket.readyState === 1) {
        const testNotification = {
          user_id: Math.floor(Math.random() * 1000),
          username: "Test User " + Math.floor(Math.random() * 100),
          type: "friend_request",
          timestamp: new Date().toISOString(),
        };
        socket.onmessage({ data: JSON.stringify(testNotification) });
      }
    });
  };
})();
