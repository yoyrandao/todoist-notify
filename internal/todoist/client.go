package todoist

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

var _ json.Unmarshaler = (*DateOnlyTime)(nil)

type DateOnlyTime struct {
	time.Time
}

func (t *DateOnlyTime) UnmarshalJSON(b []byte) (err error) {
	date, err := time.Parse(time.DateOnly, strings.Trim(string(b), `"`))
	if err != nil {
		return err
	}

	t.Time = date
	return
}

type TodoistClient struct {
	apiKey string
}

type Task struct {
	Id          string    `json:"id"`
	ProjectId   string    `json:"project_id"`
	Content     string    `json:"content"`
	Description string    `json:"description"`
	IsCompleted bool      `json:"is_completed"`
	Labels      []string  `json:"labels"`
	CreatedAt   time.Time `json:"created_at"`
	Due         struct {
		Date DateOnlyTime `json:"date"`
	} `json:"due"`
}

type OAuthAccessToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

const (
	API_URL_BASE       = "https://api.todoist.com/rest/v2"
	AUTH_URL_BASE      = "https://todoist.com/oauth/authorize"
	TOKEN_EXCHANGE_URL = "https://todoist.com/oauth/access_token"
)

func NewTodoistClient(apiKey string) *TodoistClient {
	return &TodoistClient{apiKey}
}

func (c *TodoistClient) GetTasks() ([]*Task, error) {
	response, err := c.makeGetRequest("/tasks")
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var tasks []*Task
	if err := json.Unmarshal(body, &tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (c *TodoistClient) GetTask(id string) (*Task, error) {
	response, err := c.makeGetRequest(fmt.Sprintf("/tasks/%s", id))
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var task *Task
	if err := json.Unmarshal(body, &task); err != nil {
		return nil, err
	}

	return task, nil
}

func GetAccessToken(grantingCode, oauthClientId, oauthClientSecret string) (*OAuthAccessToken, error) {
	// building request parameters
	jsonData := []byte(fmt.Sprintf(`{
			"client_id": "%s",
			"client_secret": "%s",
			"code": "%s"
		}`, oauthClientId, oauthClientSecret, grantingCode))

	// building request to get access token from granting code
	request, _ := http.NewRequest(http.MethodPost, TOKEN_EXCHANGE_URL, bytes.NewBuffer(jsonData))
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
	var accessToken OAuthAccessToken
	if err := json.Unmarshal(payload, &accessToken); err != nil {
		return nil, err
	}

	slog.Debug("got access token", "access_token", accessToken.AccessToken)

	return &accessToken, nil
}

func (c *TodoistClient) makeGetRequest(endpoint string) (*http.Response, error) {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s%s", API_URL_BASE, endpoint), nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	return http.DefaultClient.Do(req)
}
