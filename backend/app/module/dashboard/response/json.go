package response

type DashboardSummaryResponse struct {
	TotalSongs int64  `json:"total_songs"`
	TotalSize  string `json:"total_size"`
	LastUpload string `json:"last_upload"`
}
