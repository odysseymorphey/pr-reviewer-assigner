package dto

type User struct {
	ID       string `json:"user_id"`
	Name     string `json:"username"`
	Team     string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}

type UserPR struct {
	ID  string    `json:"user_id"`
	PRs []PRShort `json:"pull_requests"`
}

type UserResponse struct {
	User User `json:"user"`
}
