package dto

type ReassignRequest struct {
	PullRequestID string `json:"pull_request_id"`
	OldUserID     string `json:"old_user_id"`
}

type ReassignResponse struct {
	PR         PR     `json:"pr"`
	ReplacedBy string `json:"replaced_by"`
}
