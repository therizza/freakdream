package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"golang.org/x/time/rate"

	"github.com/therizza/freakdream/backend/internal/config"
	"github.com/therizza/freakdream/backend/internal/handlers"
	"github.com/therizza/freakdream/backend/internal/middleware"
	"github.com/therizza/freakdream/backend/internal/repository"
)

func main() {
	cfg := config.Load()

	db, err := config.NewDB(cfg)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	defer db.Close()
	log.Println("database connected")

	if err := os.MkdirAll(cfg.UploadDir, 0755); err != nil {
		log.Fatalf("failed to create upload dir: %v", err)
	}

	// Repositories
	userRepo := repository.NewUserRepository(db)
	postRepo := repository.NewPostRepository(db)
	photoRepo := repository.NewPhotoRepository(db)
	followRepo := repository.NewFollowRepository(db)
	friendRepo := repository.NewFriendRepository(db)
	notifRepo := repository.NewNotificationRepository(db)
	chatRepo := repository.NewChatRepository(db)

	// Handlers
	authH := handlers.NewAuthHandler(userRepo, cfg.JWTSecret, cfg.JWTExpiry)
	userH := handlers.NewUserHandler(userRepo, photoRepo)
	postH := handlers.NewPostHandler(postRepo, photoRepo, userRepo, cfg.UploadDir)
	followH := handlers.NewFollowHandler(followRepo)
	friendH := handlers.NewFriendHandler(friendRepo)
	notifH := handlers.NewNotificationHandler(notifRepo)
	chatHub := handlers.NewChatHub(chatRepo)

	// Rate limiter for auth endpoints (5 req/min per IP)
	loginLimiter := middleware.NewRateLimiter(rate.Every(time.Minute/5), 5)

	r := chi.NewRouter()

	// Global middleware
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.Timeout(30 * time.Second))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://localhost:3000", os.Getenv("FRONTEND_ORIGIN")},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Static file serving for uploads
	uploadsPath := cfg.UploadDir
	r.Handle("/uploads/*", http.StripPrefix("/uploads/",
		http.FileServer(http.Dir(uploadsPath))))

	// Auth routes (rate limited)
	r.Group(func(r chi.Router) {
		r.Use(loginLimiter.Limit)
		r.Post("/api/auth/login", authH.Login)
		r.Post("/api/auth/register", authH.Register)
	})

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(middleware.JWTAuth(cfg.JWTSecret))

		r.Post("/api/auth/logout", authH.Logout)

		// Users
		r.Get("/api/users/me", userH.Me)
		r.Get("/api/users/search", userH.Search)
		r.Get("/api/users/{id}", userH.GetProfile)
		r.Get("/api/users/{id}/posts", postH.GetUserPosts)

		// Posts / Feed
		r.Get("/api/feed", postH.Feed)
		r.Post("/api/posts", postH.Create)
		r.Post("/api/posts/photo", postH.UploadPhoto)
		r.Post("/api/posts/share", postH.Share)
		r.Delete("/api/posts/{id}", postH.Delete)

		// Follow
		r.Post("/api/follow/{id}", followH.Follow)
		r.Delete("/api/follow/{id}", followH.Unfollow)
		r.Get("/api/follow/{id}/count", followH.Count)
		r.Get("/api/follow/{id}/status", followH.IsFollowing)

		// Friends
		r.Post("/api/friends/request/{id}", friendH.SendRequest)
		r.Post("/api/friends/accept/{id}", friendH.Accept)
		r.Delete("/api/friends/{id}", friendH.Remove)
		r.Get("/api/friends/requests", friendH.ListRequests)

		// Notifications
		r.Get("/api/notifications", notifH.List)
		r.Delete("/api/notifications/{id}", notifH.Remove)
		r.Patch("/api/notifications/{id}/read", notifH.MarkRead)

		// Chat history (REST)
		r.Get("/api/chat/{id}/history", chatHub.History)

		// WebSocket chat
		r.Get("/ws/chat", chatHub.ServeWS)
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{"status":"ok"}`)
	})

	absUploadDir, _ := filepath.Abs(cfg.UploadDir)
	log.Printf("upload dir: %s", absUploadDir)
	log.Printf("server listening on :%s", cfg.ServerPort)
	if err := http.ListenAndServe(":"+cfg.ServerPort, r); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
