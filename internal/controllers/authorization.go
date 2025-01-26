package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
	"todoist-tg/internal/container"
	"todoist-tg/internal/messages"
	"todoist-tg/internal/storage"
	"todoist-tg/internal/todoist"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/gin-gonic/gin"
)

var (
	TODOIST_CLIENT_ID     = os.Getenv("TODOIST_CLIENT_ID")
	TODOIST_CLIENT_SECRET = os.Getenv("TODOIST_CLIENT_SECRET")

	TELEGRAM_BOT_URL = os.Getenv("TELEGRAM_BOT_URL")
)

type AuthorizationController struct {
	container *container.Container
}

func NewAuthorizationController(container *container.Container) *AuthorizationController {
	return &AuthorizationController{container}
}

func (c *AuthorizationController) Authorize(ctx *gin.Context) {
	scope := []string{"data:read"}
	state := ctx.Query("chat_id")
	if state == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "chat_id is required"})
		return
	}

	redirectUrl, err := buildAuthorizationUrl(TODOIST_CLIENT_ID, strings.Join(scope, ","), state)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.Redirect(http.StatusMovedPermanently, redirectUrl)
}

func (c *AuthorizationController) AuthorizeCallback(ctx *gin.Context) {
	grantingCode := ctx.Query("code")
	state := ctx.Query("state")

	accessToken, err := getAccessToken(grantingCode)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	encryptedAccessToken, err := c.container.Encryptor.Encrypt(accessToken.AccessToken)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	chatId, _ := strconv.ParseInt(state, 10, 64)
	// call repository to store token
	if _, err := c.container.UserRepository.CreateOrUpdate(ctx, &storage.User{
		ChatId:                      chatId,
		EncryptedTodoistAccessToken: encryptedAccessToken,
	}); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	message := telegram.NewMessage(chatId, messages.AuthorizationSuccessful)
	if _, err := c.container.Api.Send(message); err != nil {
		slog.Warn("cannot send successful authorization message", "error", err.Error())
	}

	ctx.Redirect(http.StatusMovedPermanently, TELEGRAM_BOT_URL)
}

func buildAuthorizationUrl(clientId string, scope string, state string) (string, error) {
	request, err := http.NewRequest(http.MethodGet, todoist.AUTH_URL_BASE, nil)
	if err != nil {
		return "", err
	}

	queryParams := request.URL.Query()
	queryParams.Add("client_id", clientId)
	queryParams.Add("scope", scope)
	queryParams.Add("state", state)

	request.URL.RawQuery = queryParams.Encode()

	return request.URL.String(), nil
}

func getAccessToken(grantingCode string) (*todoist.OAuthAccessToken, error) {
	// building request parameters
	jsonData := []byte(fmt.Sprintf(`{
		"client_id": "%s",
		"client_secret": "%s",
		"code": "%s"
	}`, TODOIST_CLIENT_ID, TODOIST_CLIENT_SECRET, grantingCode))

	// building request to get access token from granting code
	request, _ := http.NewRequest(http.MethodPost, todoist.TOKEN_EXCHANGE_URL, bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// reading response body
	payload, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	// getting access token from responses
	var accessToken todoist.OAuthAccessToken
	if err := json.Unmarshal(payload, &accessToken); err != nil {
		return nil, err
	}

	return &accessToken, nil
}
