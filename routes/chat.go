package routes

import (
	"chat_ai/controllers"
	"chat_ai/services"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/patrickmn/go-cache"
)

func SetupChatRoutes(router fiber.Router) {
	cache := cache.New(5*time.Minute, 10*time.Minute)
	chatService := services.NewChatService(os.Getenv("GEMINI_API_KEY"), os.Getenv("MODEL_NAME"), cache)
	chatController := controllers.NewChatController(chatService)
	chat := router.Group("/chat")
	chat.Post("/", chatController.HandleChat)
}
