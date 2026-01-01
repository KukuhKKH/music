package response

import (
	"git.dev.siap.id/kukuhkkh/app-music/app/database/schema"
)

type TrackResponse struct {
	ID        uint64      `json:"id"`
	Title     string      `json:"title"`
	Artist    string      `json:"artist"`
	Album     *string     `json:"album"`
	Duration  int         `json:"duration"`
	FileSize  int64       `json:"file_size"`
	MimeType  string      `json:"mime_type"`
	PublicURL string      `json:"public_url"`
	CreatedAt string      `json:"created_at"`
	User      schema.User `json:"user,omitempty"`
}

func FromTrackSchema(track schema.Track, publicURL string) TrackResponse {
	return TrackResponse{
		ID:        track.ID,
		Title:     track.Title,
		Artist:    track.Artist,
		Album:     track.Album,
		Duration:  track.Duration,
		FileSize:  track.FileSize,
		MimeType:  track.MimeType,
		PublicURL: publicURL,
		CreatedAt: track.CreatedAt.Format("2006-01-02 15:04:05"),
		User:      track.User,
	}
}

func FromTrackListSchema(tracks []schema.Track, storage interface{ GetURL(string) string }) []TrackResponse {
	var res []TrackResponse
	for _, t := range tracks {
		res = append(res, FromTrackSchema(t, storage.GetURL(t.StorageFilename)))
	}

	return res
}
