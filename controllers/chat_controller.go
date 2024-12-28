package controllers

import (
	"chat_ai/models"
	"chat_ai/services"

	"github.com/gofiber/fiber/v2"
)

type ChatController struct {
	chatService *services.ChatService
}

func NewChatController(chatService *services.ChatService) *ChatController {
	return &ChatController{
		chatService: chatService,
	}
}

func (cc *ChatController) HandleChat(c *fiber.Ctx) error {
	var chatRequest models.ChatRequest
	if err := c.BodyParser(&chatRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	aiResponse, err := cc.chatService.ProcessChat(chatRequest)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to process chat",
		})
	}

	return c.JSON(models.ChatResponse{
		Message: aiResponse,
		ChatID:  chatRequest.ChatID,
	})
}
