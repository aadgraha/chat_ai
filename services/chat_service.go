package services

import (
	"chat_ai/models"
	"context"
	"encoding/json"

	"github.com/google/generative-ai-go/genai"
	gcache "github.com/patrickmn/go-cache"
	"google.golang.org/api/option"
)

type ChatService struct {
	apiKey    string
	modelName string
	cache     *gcache.Cache
}

func NewChatService(apiKey, modelName string, cache *gcache.Cache) *ChatService {
	return &ChatService{
		apiKey:    apiKey,
		modelName: modelName,
		cache:     cache,
	}
}

func (s *ChatService) ProcessChat(chatRequest models.ChatRequest) (string, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(s.apiKey))
	if err != nil {
		return "", err
	}
	defer client.Close()
	model := client.GenerativeModel(s.modelName)
	aiContext, found := s.cache.Get(chatRequest.ChatID)
	if !found {
		aiContext = "default context"
	}
	prompt := aiContext.(string) + "\nUser: " + chatRequest.Prompt + "\nAI:"
	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", err
	}
	json, err := json.Marshal(resp.Candidates[0].Content.Parts[0])
	if err != nil {
		return "", err
	}
	aiResponse := string(json)
	newAIContext := aiContext.(string) + "\nUser: " + chatRequest.Prompt + "\nAI: " + aiResponse
	s.cache.Set(chatRequest.ChatID, newAIContext, gcache.DefaultExpiration)
	return aiResponse, nil
}
