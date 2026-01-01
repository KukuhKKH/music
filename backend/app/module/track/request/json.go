package request

type TrackPaginationRequest struct {
	Search string `query:"search"`
	Page   int    `query:"page"`
	Limit  int    `query:"limit"`
}
type CreateTrackRequest struct {
	Title    string `form:"title" validate:"required"`
	Artist   string `form:"artist"`
	Album    string `form:"album"`
	Duration int    `form:"duration"`
}
