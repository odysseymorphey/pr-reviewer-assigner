package dto

type TeamMember struct {
	ID       string `json:"user_id"`
	Name     string `json:"username"`
	IsActive bool   `json:"is_active"`
}

type Team struct {
	Name    string       `json:"team_name"`
	Members []TeamMember `json:"members"`
}
