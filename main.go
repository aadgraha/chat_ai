package main

import (
	"context"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	app := fiber.New()
	ctx := context.Background()
	geminiApiKey := os.Getenv("GEMINI_API_KEY")
	port := os.Getenv("PORT")
	client, err := genai.NewClient(ctx, option.WithAPIKey(geminiApiKey))
	if err != nil {
		log.Fatal(err)
	}
	app.Post("/chat", func(c *fiber.Ctx) error {
		var chatrequest ChatRequest
		err := c.BodyParser(&chatrequest)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
		}
		aiResponse := TextPrompt(chatrequest, client, ctx)

		response := fiber.Map{"message": aiResponse}
		return c.JSON(response)
	})

	err = app.Listen(port)
	if err != nil {
		panic(err)
	}
}

type ChatRequest struct {
	Prompt string `json:"prompt"`
}

func TextPrompt(chatRequest ChatRequest, client *genai.Client, ctx context.Context) *genai.GenerateContentResponse {
	modelName := os.Getenv("MODEL_NAME")
	model := client.GenerativeModel(modelName)
	resp, err := model.GenerateContent(ctx, genai.Text(chatRequest.Prompt))
	if err != nil {
		log.Fatal(err)
	}
	return resp
}
