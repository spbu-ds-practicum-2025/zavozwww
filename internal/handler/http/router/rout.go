package rout

import (
	"authServ/internal/services"
	"authServ/pkg/logger"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
)

type ProfileReq struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Age       int    `json:"age"`
	Info      string `json:"info"`
	City      string `json:"city"`
}

// Handler обрабатывает HTTP-запросы, связанные с аутентификацией и управлением пользователями.
type Handler struct {
	userService    services.UserService
	log            *logger.Logger
	profileService services.ProfileService
}

// NewHandler создает новый экземпляр Handler с заданными сервисами и логгером.
func NewHandler(userService services.UserService, profileService services.ProfileService, log *logger.Logger) *Handler {
	return &Handler{
		userService:    userService,
		profileService: profileService,
		log:            log,
	}
}

// logRequestMiddleware логирует детали каждого HTTP-запроса.
func (h *Handler) logRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		start := time.Now()
		defer func() {
			reqID := middleware.GetReqID(r.Context())

			h.log.Info("request completed",
				"request_id", reqID,
				"method", r.Method,
				"path", r.URL.Path,
				"status", ww.Status(),
				"duration", time.Since(start),
				"bytes_written", ww.BytesWritten(),
			)
		}()

		next.ServeHTTP(ww, r)
	})
}

// newRateLimiter создает middleware для ограничения частоты запросов.
func newRateLimiter() func(http.Handler) http.Handler {
	lmt := tollbooth.NewLimiter(1, &limiter.ExpirableOptions{DefaultExpirationTTL: time.Hour})
	lmt.SetMessage("You have reached maximum request limit.")
	lmt.SetMessageContentType("application/json; charset=utf-8")
	lmt.SetOnLimitReached(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
	})
	return func(next http.Handler) http.Handler {
		return tollbooth.LimitFuncHandler(lmt, next.ServeHTTP)
	}
}

func (h *Handler) InitRoutes() *chi.Mux {
	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://127.0.0.1:5500", "http://localhost:5500"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
		Debug:            true,
	}))

	router.Use(middleware.RequestID)
	router.Use(h.logRequestMiddleware)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Use(render.SetContentType(render.ContentTypeJSON))

	// csrfMiddleware := csrf.Protect(
	// 	[]byte("32-byte-long-auth-key"),
	// 	csrf.Secure(false), // Установить в true для продакшена
	// 	csrf.HttpOnly(true),
	// 	csrf.Path("/"),
	// )
	// router.Use(csrfMiddleware)

	router.Group(func(r chi.Router) {
		r.Use(newRateLimiter())
		r.Post("/filmbuddy/register", h.register)
		r.Post("/filmbuddy/login", h.login)
		r.Post("/filmbuddy/resend-verification", h.resendVerificationEmail)
	})

	router.Post("/filmbuddy/verify", h.verifyEmail)
	router.Post("/filmbuddy/refresh", h.refreshTokens)
	router.Post("/filmbuddy/logout", h.logout)
	router.Post("/filmbuddy/profile", h.saveProfile)
	return router
}

// register обрабатывает запросы на регистрацию новых пользователей.
func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	var input services.RegisterUser
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.log.Error("failed to decode request body", "id", middleware.GetReqID(r.Context()), "error", err)
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "invalid request body"})
		return
	}

	err := h.userService.RegisterUser(r.Context(), input)
	if err != nil {
		h.log.Error("failed to register user", "id", middleware.GetReqID(r.Context()), "error", err)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, map[string]string{"message": "verification email sent"})
}

// login обрабатывает запросы на аутентификацию пользователей.
func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	var input services.LoginUser
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.log.Error("failed to decode request body", "id", middleware.GetReqID(r.Context()), "error", err)
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "invalid request body"})
		return
	}

	tokens, err := h.userService.LoginUser(r.Context(), input)
	if err != nil {
		h.log.Error("failed to login user", "id", middleware.GetReqID(r.Context()), "error", err)
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}

	render.JSON(w, r, tokens)
}

// verifyEmail обрабатывает запросы на верификацию электронной почты пользователей.
func (h *Handler) verifyEmail(w http.ResponseWriter, r *http.Request) {
	var input services.CheckTruthEmail
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.log.Error("failed to decode request body", "id", middleware.GetReqID(r.Context()), "error", err)
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "invalid request body"})
		return
	}

	tokens, err := h.userService.CheckTruthEmail(r.Context(), input)
	if err != nil {
		h.log.Error("failed to verify email", "id", middleware.GetReqID(r.Context()), "error", err)
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": err.Error()})

		return
	}

	render.JSON(w, r, tokens)
}

// refreshTokens обрабатывает запросы на обновление токенов аутентификации.
func (h *Handler) refreshTokens(w http.ResponseWriter, r *http.Request) {
	var input services.RefreshInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.log.Error("failed to decode request body", "id", middleware.GetReqID(r.Context()), "error", err)
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "invalid request body"})
		return
	}

	tokens, err := h.userService.RefreshTokens(r.Context(), input)
	if err != nil {
		h.log.Error("failed to refresh tokens", "id", middleware.GetReqID(r.Context()), "error", err)
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}

	render.JSON(w, r, tokens)
}

// logout обрабатывает запросы на выход пользователей из системы.
func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {
	var input services.RefreshInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.log.Error("failed to decode request body", "id", middleware.GetReqID(r.Context()), "error", err)
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "invalid request body"})
		return
	}

	if err := h.userService.Logout(r.Context(), input); err != nil {
		h.log.Error("failed to logout", "id", middleware.GetReqID(r.Context()), "error", err)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]string{"message": "logged out successfully"})
}

// resendVerificationEmail обрабатывает запрос на повторную отправку письма верификации.
func (h *Handler) resendVerificationEmail(w http.ResponseWriter, r *http.Request) {
	var input services.ResendVerificationEmailInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.log.Error("failed to decode request body", "id", middleware.GetReqID(r.Context()), "error", err)
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "invalid request body"})
		return
	}

	err := h.userService.ResendVerificationEmail(r.Context(), input)
	if err != nil {
		h.log.Error("failed to resend verification email", "id", middleware.GetReqID(r.Context()), "error", err)
		if strings.Contains(err.Error(), "please wait") {
			render.Status(r, http.StatusTooManyRequests)
		} else {
			render.Status(r, http.StatusInternalServerError)
		}
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]string{"message": "a new verification email has been sent"})
}

func (h *Handler) saveProfile(w http.ResponseWriter, r *http.Request) {

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, map[string]string{"error": "authorization header is required"})
		return
	}

	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, map[string]string{"error": "invalid authorization header format"})
		return
	}
	tokenString := headerParts[1]

	userID, err := h.userService.ParseAccessToken(r.Context(), tokenString)
	if err != nil {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, map[string]string{"error": "invalid token"})
		return
	}

	var input ProfileReq
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.log.Error("failed to decode request body", "id", middleware.GetReqID(r.Context()), "error", err)
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "invalid request body"})
		return
	}

	profileEntity := &services.ProfileReq{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Age:       input.Age,
		Info:      input.Info,
		City:      input.City,
	}

	err = h.profileService.SaveProfile(r.Context(), *profileEntity, userID)
	if err != nil {
		h.log.Error("failed to save profile", "id", middleware.GetReqID(r.Context()), "error", err)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "failed to save profile"})
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]string{"message": "profile saved successfully"})
}
