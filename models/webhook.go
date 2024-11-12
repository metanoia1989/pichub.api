package models

type WebhookPayload struct {
	Ref        string            `json:"ref"`
	Repository WebhookRepository `json:"repository"`
	Commits    []Commit          `json:"commits"`
}

type Commit struct {
	ID        string   `json:"id"`
	Message   string   `json:"message"`
	Added     []string `json:"added"`
	Removed   []string `json:"removed"`
	Modified  []string `json:"modified"`
	Timestamp string   `json:"timestamp"`
}

type WebhookRepository struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	HTMLURL  string `json:"html_url"`
}
