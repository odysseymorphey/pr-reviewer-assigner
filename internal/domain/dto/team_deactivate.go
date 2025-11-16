package dto

type TeamDeactivateRequest struct {
	TeamName string   `json:"team_name"`
	UserIDs  []string `json:"user_ids"`
}

type TeamDeactivateResponse struct {
	TeamName    string   `json:"team_name"`
	Deactivated []string `json:"deactivated_user_ids"`
}
