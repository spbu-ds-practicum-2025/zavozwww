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
          email: 'test@mail.com'
        }
      }), { status: 200 })
    } else {
      return new Response(JSON.stringify({ 
        error: 'Invalid credentials' 
      }), { status: 401 })
    }
  }
  
  if (url.includes('/api/search')) {
    const mockMovies = [
      {
        id: 1,
        title: "Интерстеллар",
        imageSrc: "https://via.placeholder.com/300x450/3498db/ffffff?text=Interstellar",
        rating: 4.8,
        countRaitings: 1500
      },
      {
        id: 2, 
        title: "Начало",
        imageSrc: "https://via.placeholder.com/300x450/e74c3c/ffffff?text=Inception",
        rating: 4.7,
        countRaitings: 1200
      }
    ]
    
    return new Response(JSON.stringify({ movies: mockMovies }))
  }
  
  // Для остальных запросов используем оригинальный fetch
  return originalFetch.apply(this, args)
}