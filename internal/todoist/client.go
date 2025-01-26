package todoist

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

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

const apiUrlBase = "https://api.todoist.com/rest/v2"

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

func (c *TodoistClient) makeGetRequest(endpoint string) (*http.Response, error) {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s%s", apiUrlBase, endpoint), nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	return http.DefaultClient.Do(req)
}
