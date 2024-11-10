package controllers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"pichub.api/models"
	"pichub.api/services"
)

// GithubWebhook 处理 GitHub webhook 请求
func GithubWebhook(c *gin.Context) {
	// 读取请求体
	payload, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}

	// 验证签名
	signature := c.GetHeader("X-Hub-Signature-256")
	if !services.WebhookService.ValidateSignature(payload, signature) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid signature"})
		return
	}

	// 检查事件类型
	event := c.GetHeader("X-GitHub-Event")
	if event != "push" {
		c.JSON(http.StatusOK, gin.H{"message": "Event ignored"})
		return
	}

	// 解析 payload
	var webhookPayload models.WebhookPayload
	if err := json.Unmarshal(payload, &webhookPayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload format"})
		return
	}

	// 处理 push 事件
	if err := services.WebhookService.HandlePush(&webhookPayload); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Webhook processed successfully"})
}
