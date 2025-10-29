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
        token: 'mock_jwt_token',
        user_profile: {
          id: 1,
          username: body.username,
          email: 'test@mail.com',
          friends: []
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
  
  if (url.includes('/api/search')) {
    const mockMovies = [
      {
        id: 1,
        title: "Интерстеллар",
        imageSrc: "https://clck.ru/3Q3i7R",
        rating: 4.8,
        countRaitings: 1500
      },
      {
        id: 2, 
        title: "Начало",
        imageSrc: "https://clck.ru/3Q3i7u",
        rating: 4.7,
        countRaitings: 1200
      }
    ]
    
    return new Response(JSON.stringify({ movies: mockMovies }))
  }

  if (url.includes('/api/recomendations')) {
    const mockRecommendations = [
      {
        id: 3,
        title: "Матрица",
        imageSrc: "https://clck.ru/3Q3i9s",
        rating: 4.6,
        countRaitings: 2000
      },
      {
        id: 4,
        title: "Крестный отец",
        imageSrc: "https://clck.ru/3Q3iAB",
        rating: 4.9,
        countRaitings: 1800
      }
    ]
    
    return new Response(JSON.stringify(mockRecommendations), {
      headers: { 'Content-Type': 'application/json' }
    })
  }

  if (url.includes('/api/friends/request') && options?.method === 'POST') {
    const body = JSON.parse(options.body)
    const mockFriends = [
      {
        name: body.target_username,
        avatarSrc: "https://clck.ru/3Q3iBp"
      }
    ]
    
    return new Response(JSON.stringify({ friends: mockFriends }), {
      headers: { 'Content-Type': 'application/json' }
    })
  }

  if (url.includes('/api/friends/add') && options?.method === 'POST') {
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
  // Для остальных запросов используем оригинальный fetch
  return originalFetch.apply(this, args)
}