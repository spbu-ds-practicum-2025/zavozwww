const originalWebSocket = window.WebSocket;
window.WebSocket = class MockWebSocket {
  constructor(url) {
    this.url = url;
    this.readyState = 0; // CONNECTING
    this.onopen = null;
    this.onmessage = null;
    this.onerror = null;
    this.onclose = null;
    
    // Имитируем успешное подключение
    setTimeout(() => {
      this.readyState = 1; // OPEN
      if (this.onopen) this.onopen(new Event('open'));
      
      // Имитируем получение тестового уведомления через 2 секунды
      setTimeout(() => {
        if (this.onmessage) {
          const mockNotification = {
            user_id: 123,
            username: 'test_friend1',
            type: 'friend_request'
          };
          this.onmessage({ 
            data: JSON.stringify(mockNotification) 
          });
        }
      }, 2000);
      setTimeout(() => {
        if (this.onmessage) {
          const mockNotification = {
            user_id: 321,
            username: 'test_friend2',
            type: 'friend_request'
          };
          this.onmessage({ 
            data: JSON.stringify(mockNotification) 
          });
        }
      }, 5000);
    }, 100);
  }

  send(data) {
    console.log('Mock WebSocket sent:', data);
  }

  close() {
    this.readyState = 3; // CLOSED
    if (this.onclose) this.onclose(new Event('close'));
  }
};

// Переопределяем fetch для мокирования
const originalFetch = window.fetch

window.fetch = async function(...args) {
  const [url, options] = args
  
  // Имитируем задержку сети
  await new Promise(resolve => setTimeout(resolve, 500))
  
  // Мокируем конкретные endpoints
  if (url.includes('/api/login') && options?.method === 'POST') {
    const body = JSON.parse(options.body)
    
    if (body.username === 'testuser' && body.password === '123456') {
      return new Response(JSON.stringify({
        access_token: 'mock_jwt_token',
        user_profile: {
          id: 1,
          username: body.username,
          email: 'test@mail.com',
          friends: [
            {
              name: 'Иван Петров',
              regDate: '15.03.2021',
              countRateFilms: 24,
              countFriends: 15
            },
            {
              name: 'Мария Сидорова', 
              regDate: '02.11.2020',
              countRateFilms: 42,
              countFriends: 28
            }
          ],
          regDate: '12.12.2012',
          countRateFilms: 12,
          countFriends: 2
        }
      }), { status: 200 })
    } else {
      return new Response(JSON.stringify({ 
        error: 'Invalid credentials'
      }), { status: 401,
            statusText: "invalid data"
       })
    }
  }
  
  if (url.includes('/api/register') && options?.method === 'POST') {
    return new Response(JSON.stringify({
      access_token: 'mock_jwt_token',
      user_profile: {
        first_name: "Алексей",
        second_name: "Алексеев",
        city: "Алексеево",
        age: 100,
        info: "В чащах юга жил-был цитрус — да, но фальшивый экземпляръ!"
      }
    }), { status: 200 })
  }
  
  if (url.includes('/api/search')) {
    const mockMovies = [
      {
        id: 1,
        title: "Интерстеллар",
        imageSrc: "https://clck.ru/3Q3i7R",
        genres: "триллер, научная фантастика",
        year: 2014,
        rating: 4.8,
        countRaitings: 1500
      },
      {
        id: 2, 
        title: "Начало",
        year: 2010,
        imageSrc: "https://clck.ru/3Q3i7u",
        genres: "триллер, боевик",
        rating: 4.7,
        countRaitings: 1200
      },
      {
        id: 3,
        title: "Матрица",
        year: 1999,
        imageSrc: "https://clck.ru/3Q3i9s",
        genres: "боевик, научная фантастика, триллер",
        rating: 4.6,
        countRaitings: 2000
      },
      {
        id: 4,
        title: "Крестный отец",
        year: 1972,
        imageSrc: "https://clck.ru/3Q3iAB",
        genres: "боевик",
        rating: 4.9,
        countRaitings: 1800
      }
    ]
    
    return new Response(JSON.stringify({ movies: mockMovies }))
  }

  if (url.includes('/api/recomendations')) {
    const mockRecommendations = [
      {
        id: 3,
        title: "Матрица",
        year: 1999,
        imageSrc: "https://clck.ru/3Q3i9s",
        genres: "боевик, научная фантастика, триллер",
        rating: 4.6,
        countRaitings: 2000
      },
      {
        id: 4,
        title: "Крестный отец",
        year: 1972,
        imageSrc: "https://clck.ru/3Q3iAB",
        genres: "боевик",
        rating: 4.9,
        countRaitings: 1800
      }
    ]
    
    return new Response(JSON.stringify(mockRecommendations), {
      headers: { 'Content-Type': 'application/json' }
    })
  }

  if (url.includes('/api/friends/search')) {
    const urlObj = new URL(url, window.location.origin)
    const searchName = urlObj.searchParams.get('name')
    
    let mockFriends = []
    
    if (searchName && searchName.trim() !== '') {
      mockFriends = [
        {
          name: 'Алексей Козлов',
          regDate: '08.07.2019',
          countRateFilms: 18,
          countFriends: 9
        },
        {
          name: 'Александр Новиков',
          regDate: '22.12.2022',
          countRateFilms: 7,
          countFriends: 3
        }
      ].filter(friend => 
        friend.name.toLowerCase().includes(searchName.toLowerCase())
      )
    }
    
    return new Response(JSON.stringify({ friends: mockFriends }), {
      status: 200,
      headers: { 'Content-Type': 'application/json' }
    })
  }

  if (url.includes('/api/friends/request') && options?.method === 'POST') {
    return new Response(JSON.stringify({ 
      message: 'Friend request sent successfully'
    }), {
      status: 200,
      headers: { 'Content-Type': 'application/json' }
    })
  }

  if (url.includes('/api/movies/') && url.includes('/rating') && options?.method === 'POST') {
    return new Response(JSON.stringify({ 
      message: 'Rating saved successfully'
    }), {
      status: 200,
      headers: { 'Content-Type': 'application/json' }
    })
  }
  
   if (url.includes('/api/friends/accept/request_id') && options?.method === 'POST') {
    const body = JSON.parse(options.body)
    return new Response(JSON.stringify({ 
      message: `Friend request from user ${body.user_id} accepted successfully` 
    }), {
      status: 200,
      headers: { 'Content-Type': 'application/json' }
    })
  }

  if (url.includes('/api/friends/decline/request_id') && options?.method === 'POST') {
    const body = JSON.parse(options.body)
    return new Response(JSON.stringify({ 
      message: `Friend request from user ${body.user_id} declined` 
    }), {
      status: 200,
      headers: { 'Content-Type': 'application/json' }
    })
  }

  // Мок для получения списка уведомлений
  if (url.includes('/api/notifications')) {
    const mockNotifications = [
      {
        user_id: 123,
        username: 'Иван Петров',
        type: 'friend_request',
        timestamp: new Date().toISOString()
      },
      {
        user_id: 124,
        username: 'Мария Сидорова',
        type: 'friend_request', 
        timestamp: new Date(Date.now() - 3600000).toISOString()
      }
    ]
    
    return new Response(JSON.stringify({ notifications: mockNotifications }), {
      status: 200,
      headers: { 'Content-Type': 'application/json' }
    })
  }

  // Мок для отправки тестового уведомления (для отладки)
  if (url.includes('/api/test/notification') && options?.method === 'POST') {
    // Триггерим mock WebSocket сообщение
    if (window.triggerMockNotification) {
      window.triggerMockNotification();
    }
    
    return new Response(JSON.stringify({ 
      message: 'Test notification sent' 
    }), {
      status: 200,
      headers: { 'Content-Type': 'application/json' }
    })
  }

  // Для остальных запросов используем оригинальный fetch
  return originalFetch.apply(this, args)
}


// Глобальная функция для ручного триггера уведомлений (для тестирования)
window.triggerMockNotification = function() {
  // Находим активные WebSocket соединения и отправляем тестовое уведомление
  if (window.activeMockSockets) {
    window.activeMockSockets.forEach(socket => {
      if (socket.onmessage && socket.readyState === 1) {
        const testNotification = {
          user_id: Math.floor(Math.random() * 1000),
          username: 'Test User ' + Math.floor(Math.random() * 100),
          type: 'friend_request'
        };
        socket.onmessage({ 
          data: JSON.stringify(testNotification) 
        });
      }
    });
  }
};

// Глобальная функция для ручного триггера уведомлений (для тестирования)
window.triggerMockNotification = function() {
  // Находим активные WebSocket соединения и отправляем тестовое уведомление
  if (window.activeMockSockets) {
    window.activeMockSockets.forEach(socket => {
      if (socket.onmessage && socket.readyState === 1) {
        const testNotification = {
          user_id: Math.floor(Math.random() * 1000),
          username: 'Test User ' + Math.floor(Math.random() * 100),
          type: 'friend_request'
        };
        socket.onmessage({ 
          data: JSON.stringify(testNotification) 
        });
      }
    });
  }
};

// Сохраняем ссылки на созданные WebSocket соединения
window.activeMockSockets = [];
const originalMockWebSocket = window.WebSocket;
window.WebSocket = function(...args) {
  const socket = new originalMockWebSocket(...args);
  window.activeMockSockets.push(socket);
  
  // Очищаем при закрытии соединения
  const originalClose = socket.close;
  socket.close = function() {
    const index = window.activeMockSockets.indexOf(socket);
    if (index > -1) {
      window.activeMockSockets.splice(index, 1);
    }
    return originalClose.apply(this, arguments);
  };
  
  return socket;
};