# План сквозного (E2E) тестирования FilmBuddy

## Введение

Данный файл описывает сценарии сквозного тестирования проекта FilmBuddy.

**Охват**: Аутентификация, Управление профилем, Социальные взаимодействия, Система рекомендаций.

**Предусловия для всех тестов**:
- Все микросервисы (AuthService, SocialService, RecSystem) запущены.
- API Gateway доступен по адресу `http://localhost`.
- База данных PostgreSQL инициализирована и доступны миграции.

---

## 1. Сервис Аутентификации (Auth Service)

### Тест 1.1: Регистрация нового пользователя

**Цель**: Проверить создание новой учетной записи.

**Предусловия**:
- Email и Username уникальны (не существуют в БД).

**Шаги**:
1. Отправить HTTP-запрос:
   ```http
   POST http://localhost/filmbuddy/register
   Content-Type: application/json

   {
     "username": "testuser_01",
     "email": "testuser_01@example.com",
     "password": "password123"
   }
   ```

**Ожидаемый результат**:
- HTTP статус: `201 Created`
- Тело ответа:
  ```json
  {
    "message": "verification email sent"
  }
  ```
- На почту (Mailpit) отправлено письмо с кодом подтверждения.

### Тест 1.2: Верификация Email

**Цель**: Подтвердить аккаунт кодом из письма.

**Предусловия**:
- Пользователь зарегистрирован, но не подтвержден.
- Известен код подтверждения (например, `123456`).

**Шаги**:
1. Отправить HTTP-запрос:
   ```http
   POST http://localhost/filmbuddy/verify
   Content-Type: application/json

   {
     "email": "testuser_01@example.com",
     "code": "123456"
   }
   ```

**Ожидаемый результат**:
- HTTP статус: `200 OK`
- Тело ответа содержит токены:
  ```json
  {
    "access_token": "<jwt_token>",
    "refresh_token": "<uuid>"
  }
  ```

### Тест 1.3: Вход в систему (Login)

**Цель**: Получение токенов доступа по логину и паролю.

**Предусловия**:
- Пользователь подтвержден.

**Шаги**:
1. Отправить HTTP-запрос:
   ```http
   POST http://localhost/filmbuddy/login
   Content-Type: application/json

   {
     "email": "testuser_01@example.com",
     "password": "password123"
   }
   ```

**Ожидаемый результат**:
- HTTP статус: `200 OK`
- Возвращены валидные `access_token` и `refresh_token`.

### Тест 1.4: Обновление токенов (Refresh)

**Цель**: Получение новой пары токенов с помощью refresh-токена.

**Предусловия**:
- Есть валидный `refresh_token`.

**Шаги**:
1. Отправить HTTP-запрос:
   ```http
   POST http://localhost/filmbuddy/refresh
   Content-Type: application/json

   {
     "refresh_token": "<valid_refresh_token>"
   }
   ```

**Ожидаемый результат**:
- HTTP статус: `200 OK`
- Возвращены новые `access_token` и `refresh_token`.

### Тест 1.5: Получение профиля пользователя

**Цель**: Получение данных текущего пользователя.

**Предусловия**:
- Есть валидный `access_token`.

**Шаги**:
1. Отправить HTTP-запрос:
   ```http
   GET http://localhost/filmbuddy/profile
   Authorization: Bearer <access_token>
   ```

**Ожидаемый результат**:
- HTTP статус: `200 OK`
- Тело ответа:
  ```json
  {
    "id": "<uuid>",
    "username": "testuser_01",
    "email": "testuser_01@example.com",
    "first_name": "",
    "last_name": "",
    ...
  }
  ```

### Тест 1.6: Обновление профиля

**Цель**: Заполнение личных данных.

**Предусловия**:
- Пользователь авторизован.

**Шаги**:
1. Отправить HTTP-запрос:
   ```http
   POST http://localhost/filmbuddy/profile
   Authorization: Bearer <access_token>
   Content-Type: application/json

   {
     "username": "testuser_01",
     "first_name": "Ivan",
     "last_name": "Ivanov",
     "age": 25,
     "city": "Moscow",
     "info": "Movie lover"
   }
   ```

**Ожидаемый результат**:
- HTTP статус: `200 OK`
- Сообщение: `profile saved successfully`.

### Тест 1.7: Поиск пользователей (для добавления в друзья)

**Цель**: Найти пользователя по частичному совпадению username.

**Шаги**:
1. Отправить HTTP-запрос:
   ```http
   POST http://localhost/filmbuddy/friends/searchFriends
   Content-Type: application/json

   {
     "username": "testuser"
   }
   ```

**Ожидаемый результат**:
- HTTP статус: `200 OK`
- Список найденных профилей.

### Тест 1.8: Выход (Logout)

**Цель**: Аннулирование refresh-токена.

**Шаги**:
1. Отправить HTTP-запрос:
   ```http
   POST http://localhost/filmbuddy/logout
   Content-Type: application/json

   {
     "refresh_token": "<refresh_token>"
   }
   ```

**Ожидаемый результат**:
- HTTP статус: `200 OK`.

---

## 2. Социальный Сервис (Social Service)

### Тест 2.1: Получение социального профиля

**Цель**: Просмотр статистики (друзья, оценки).

**Шаги**:
1. Отправить HTTP-запрос:
   ```http
   GET http://localhost/social/profile
   Authorization: Bearer <access_token>
   ```

**Ожидаемый результат**:
- HTTP статус: `200 OK`
- JSON с количеством друзей и оценок.

### Тест 2.2: Отправка заявки в друзья

**Цель**: Инициировать дружбу с другим пользователем.

**Предусловия**:
- Пользователь А (отправитель) и Пользователь Б (получатель) существуют.
- Известен `username` получателя.

**Шаги**:
1. Отправить HTTP-запрос (от имени А):
   ```http
   POST http://localhost/social/friends/requests
   Authorization: Bearer <access_token_A>
   Content-Type: application/json

   {
     "from_username": "user_A",
     "to_username": "user_B"
   }
   ```

**Ожидаемый результат**:
- HTTP статус: `200 OK`
- Сообщение: `friend request sent`.

### Тест 2.3: Просмотр входящих заявок

**Цель**: Увидеть список запросов на дружбу.

**Шаги**:
1. Отправить HTTP-запрос (от имени Б):
   ```http
   GET http://localhost/social/friends/requests
   Authorization: Bearer <access_token_B>
   ```

**Ожидаемый результат**:
- HTTP статус: `200 OK`
- Список содержит заявку от `user_A`.

### Тест 2.4: Принятие заявки в друзья

**Цель**: Подтвердить дружбу.

**Шаги**:
1. Отправить HTTP-запрос (от имени Б):
   ```http
   POST http://localhost/social/friends/requests/accept
   Authorization: Bearer <access_token_B>
   Content-Type: application/json

   {
     "from_username": "user_A",
     "to_username": "user_B"
   }
   ```

**Ожидаемый результат**:
- HTTP статус: `200 OK`
- Сообщение: `friend request accepted`.

### Тест 2.5: Просмотр списка друзей

**Цель**: Убедиться, что пользователи стали друзьями.

**Шаги**:
1. Отправить HTTP-запрос:
   ```http
   GET http://localhost/social/friends
   Authorization: Bearer <access_token_A>
   ```

**Ожидаемый результат**:
- HTTP статус: `200 OK`
- В списке друзей присутствует `user_B`.

### Тест 2.6: Добавление оценки фильму

**Цель**: Оценить фильм и оставить отзыв.

**Шаги**:
1. Отправить HTTP-запрос:
   ```http
   POST http://localhost/social/ratings
   Authorization: Bearer <access_token>
   Content-Type: application/json

   {
     "film_id": 1,
     "grade": 9,
     "review": "Great movie!",
     "username": "user_A"
   }
   ```

**Ожидаемый результат**:
- HTTP статус: `201 Created`.

### Тест 2.7: Просмотр своих оценок

**Цель**: Получить историю оценок.

**Шаги**:
1. Отправить HTTP-запрос:
   ```http
   GET http://localhost/social/ratings
   Authorization: Bearer <access_token>
   ```

**Ожидаемый результат**:
- HTTP статус: `200 OK`
- Список содержит добавленную оценку.

---

## 3. Система Рекомендаций (RecSystem)

### Тест 3.1: Поиск фильмов

**Цель**: Найти фильм по названию или жанру.

**Шаги**:
1. Отправить HTTP-запрос:
   ```http
   POST http://localhost/recsys/search
   Content-Type: application/json

   {
     "title": "Star",
     "genre": "Фантастика"
   }
   ```

**Ожидаемый результат**:
- HTTP статус: `200 OK`
- Список фильмов, соответствующих критериям.

### Тест 3.2: Получение деталей фильма

**Цель**: Получить информацию о конкретном фильме.

**Шаги**:
1. Отправить HTTP-запрос:
   ```http
   GET http://localhost/recsys/movie/1
   ```

**Ожидаемый результат**:
- HTTP статус: `200 OK`
- JSON с деталями фильма (название, год, жанр).

### Тест 3.3: Получение рекомендаций

**Цель**: Получить список рекомендованных фильмов на основе истории оценок.

**Шаги**:
1. Отправить HTTP-запрос:
   ```http
   POST http://localhost/recsys/recommend
   Content-Type: application/json

   [
     {
       "film_id": 1,
       "grade": 10
     },
     {
       "film_id": 2,
       "grade": 5
     }
   ]
   ```

**Ожидаемый результат**:
- HTTP статус: `200 OK`
- Список рекомендованных фильмов.

---

## 4. Комплексный сценарий (Happy Path)

### Тест 4.1: Полный жизненный цикл пользователя

**Цель**: Проверить интеграцию всех сервисов в едином сценарии.

**Предусловия**:
- Система чистая или используются уникальные данные.

**Шаги**:
1. **Регистрация**: `POST /filmbuddy/register` -> Получение кода.
2. **Верификация**: `POST /filmbuddy/verify` -> Получение токена.
3. **Профиль**: `POST /filmbuddy/profile` -> Заполнение данных.
4. **Поиск фильма**: `POST /recsys/search` -> Выбор фильма ID=100.
5. **Оценка**: `POST /social/ratings` -> Оценка фильма ID=100 на 8 баллов.
6. **Рекомендации**: `POST /recsys/recommend` (с историей оценок) -> Получение списка.
7. **Социальное**: `POST /filmbuddy/friends/searchFriends` -> Поиск друга.
8. **Заявка**: `POST /social/friends/requests` -> Отправка заявки.
9. **Выход**: `POST /filmbuddy/logout`.

**Ожидаемый результат**:
- Все шаги выполняются успешно с кодом 200/201.
- Данные корректно сохраняются и передаются между сервисами.
