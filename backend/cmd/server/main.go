package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	"github.com/temanbatin/omnichannel/internal/config"
	"github.com/temanbatin/omnichannel/internal/controllers"
	"github.com/temanbatin/omnichannel/internal/repositories"
	"github.com/temanbatin/omnichannel/internal/services"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Initialize config
	cfg := config.Load()

	// Initialize database
	db, err := repositories.NewPostgresDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	messageRepo := repositories.NewMessageRepository(db)
	contactRepo := repositories.NewContactRepository(db)
	conversationRepo := repositories.NewConversationRepository(db)

	// Initialize services
	messagingSvc := services.NewMessagingService(messageRepo, contactRepo, conversationRepo, cfg)

	// Initialize controllers
	messageCtrl := controllers.NewMessageController(messagingSvc)
	webhookCtrl := controllers.NewWebhookController(messagingSvc, cfg)

	// Setup router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://omni.otomasi.click", "http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// API routes
	r.Route("/api", func(r chi.Router) {
		r.Route("/messages", func(r chi.Router) {
			r.Get("/", messageCtrl.List)
			r.Post("/", messageCtrl.Send)
			r.Get("/{id}", messageCtrl.Get)
		})

		r.Route("/conversations", func(r chi.Router) {
			r.Get("/", messageCtrl.ListConversations)
			r.Get("/{id}", messageCtrl.GetConversation)
		})

		r.Route("/contacts", func(r chi.Router) {
			r.Get("/", messageCtrl.ListContacts)
			r.Post("/", messageCtrl.CreateContact)
		})
	})

	// Webhook routes (Meta verification)
	r.Route("/webhooks", func(r chi.Router) {
		r.Get("/whatsapp", webhookCtrl.VerifyWhatsApp)
		r.Post("/whatsapp", webhookCtrl.HandleWhatsApp)
		r.Get("/instagram", webhookCtrl.VerifyInstagram)
		r.Post("/instagram", webhookCtrl.HandleInstagram)
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
