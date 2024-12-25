package main

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func main() {
	app := fiber.New()

	app.Post("/chat", func(c *fiber.Ctx) error {
		var chatrequest ChatRequest
		err := c.BodyParser(&chatrequest)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
		}
		aiResponse := TextPrompt(chatrequest)

		response := fiber.Map{"message": aiResponse}
		return c.JSON(response)
	})

	err := app.Listen(":")
	if err != nil {
		panic(err)
	}
}

type ChatRequest struct {
	Prompt string `json:"prompt"`
}

func TextPrompt(chatRequest ChatRequest) *genai.GenerateContentResponse {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(""))
	if err != nil {
		log.Fatal(err)
	}
	// defer client.Close()

	model := client.GenerativeModel("gemini-1.5-flash")
	resp, err := model.GenerateContent(ctx, genai.Text(chatRequest.Prompt))
	if err != nil {
		log.Fatal(err)
	}
	return resp
}
