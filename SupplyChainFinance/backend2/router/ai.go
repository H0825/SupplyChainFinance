package router

import (
	"backend/appconfig"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type assistantChatReq struct {
	Message string                 `json:"message" binding:"required"`
	Scene   string                 `json:"scene"`
	Guard   map[string]interface{} `json:"guard"`
}

type openAIChatReq struct {
	Model       string                 `json:"model"`
	Messages    []map[string]string    `json:"messages"`
	Temperature float64                `json:"temperature,omitempty"`
	MaxTokens   int                    `json:"max_tokens,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

type openAIChatResp struct {
	Model   string `json:"model"`
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Usage *struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage,omitempty"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

func systemPromptByScene(scene string) string {
	switch strings.ToLower(strings.TrimSpace(scene)) {
	case "business":
		return "You are a supply-chain-finance copilot for enterprise operators. Focus on receivable lifecycle, document quality, settlement consistency, and risk hints. Keep answers practical and concise."
	case "finance":
		return "You are a supply-chain-finance copilot for financial institutions. Focus on credit review, fraud/overdue risk signals, privacy-preserving data usage, and compliance-aware recommendations."
	default:
		return "You are a supply-chain-finance copilot. Provide clear, practical, compliance-aware suggestions."
	}
}

// ChatAssistant proxies chat requests to configured AI provider.
func ChatAssistant(c *gin.Context) {
	var req assistantChatReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cfg, err := appconfig.Load("config.toml")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "load app config failed"})
		return
	}
	if strings.TrimSpace(cfg.AI.APIKey) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ai api key not configured (set [AI].api_key or AI_API_KEY)"})
		return
	}

	modelReq := openAIChatReq{
		Model: cfg.AI.Model,
		Messages: []map[string]string{
			{"role": "system", "content": systemPromptByScene(req.Scene)},
			{"role": "user", "content": strings.TrimSpace(req.Message)},
		},
		Temperature: cfg.AI.Temperature,
		MaxTokens:   900,
		Metadata: map[string]interface{}{
			"scene": req.Scene,
			"guard": req.Guard,
		},
	}
	b, _ := json.Marshal(modelReq)

	base := strings.TrimRight(strings.TrimSpace(cfg.AI.BaseURL), "/")
	url := base + "/chat/completions"
	httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "build ai request failed"})
		return
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+cfg.AI.APIKey)

	timeout := time.Duration(cfg.AI.TimeoutSec) * time.Second
	client := &http.Client{Timeout: timeout}
	resp, err := client.Do(httpReq)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "ai upstream request failed", "detail": err.Error()})
		return
	}
	defer resp.Body.Close()

	raw, _ := io.ReadAll(resp.Body)
	var upstream openAIChatResp
	_ = json.Unmarshal(raw, &upstream)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		msg := strings.TrimSpace(string(raw))
		if upstream.Error != nil && strings.TrimSpace(upstream.Error.Message) != "" {
			msg = upstream.Error.Message
		}
		c.JSON(http.StatusBadGateway, gin.H{"error": fmt.Sprintf("ai upstream failed: %s", msg)})
		return
	}

	if len(upstream.Choices) == 0 || strings.TrimSpace(upstream.Choices[0].Message.Content) == "" {
		c.JSON(http.StatusBadGateway, gin.H{"error": "ai upstream returned empty content"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"reply": upstream.Choices[0].Message.Content,
		"meta": gin.H{
			"model": strings.TrimSpace(firstNonEmpty(upstream.Model, cfg.AI.Model)),
			"scene": req.Scene,
			"usage": upstream.Usage,
		},
	})
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}
