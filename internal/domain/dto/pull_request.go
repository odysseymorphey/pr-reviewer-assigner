package dto

type PRRequest struct {
	ID       string `json:"pull_request_id"`
	Name     string `json:"pull_request_name"`
	AuthorID string `json:"author_id"`
}

type PR struct {
	ID        string   `json:"pull_request_id"`
	Name      string   `json:"pull_request_name"`
	AuthorID  string   `json:"author_id"`
	Status    string   `json:"status"`
	Reviewers []string `json:"assigned_reviewers"`
	CreatedAt string   `json:"createdAt,omitempty"`
	MergedAt  string   `json:"mergedAt,omitempty"`
}

type PRShort struct {
	ID       string `json:"pull_request_id"`
	Name     string `json:"pull_request_name"`
	AuthorID string `json:"author_id"`
	Status   string `json:"status"`
}

type PRResponse struct {
	PR PR `json:"pr"`
}
