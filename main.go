package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	gcache "github.com/patrickmn/go-cache"
	"google.golang.org/api/option"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	app := fiber.New()
	port := os.Getenv("PORT")
	cache := gcache.New(5*time.Minute, 10*time.Minute)
	if err != nil {
		log.Fatal(err)
	}
	app.Post("/chat", func(c *fiber.Ctx) error {
		var chatrequest ChatRequest
		err := c.BodyParser(&chatrequest)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
		}
		aiResponse := TextPrompt(chatrequest, cache)
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
	ChatID string `json:"chat_id"`
}

func TextPrompt(chatRequest ChatRequest, cache *gcache.Cache) string {
	geminiApiKey := os.Getenv("GEMINI_API_KEY")
	modelName := os.Getenv("MODEL_NAME")
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(geminiApiKey))
	if err != nil {
		log.Fatal(err)
	}
	model := client.GenerativeModel(modelName)
	aiContext, found := cache.Get(chatRequest.ChatID)
	if !found {
		aiContext = "default context"
	}
	prompt := aiContext.(string) + "\nUser: " + chatRequest.Prompt + "\nAI:"
	defer client.Close()
	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		log.Fatal(err)
	}
	json, err := json.Marshal(resp.Candidates[0].Content.Parts[0])
	aiResponse := string(json)
	newAIContext := aiContext.(string) + "\nUser: " + chatRequest.Prompt + "\nAI: " + aiResponse
	cache.Set(chatRequest.ChatID, newAIContext, gcache.DefaultExpiration)
	if err != nil {
		log.Fatal(err)
	}
	return aiResponse
}
