package schema

type Track struct {
	ID               uint64  `gorm:"primary_key;column:id" json:"id"`
	UserID           uint64  `gorm:"column:user_id;not null" json:"user_id"`
	Title            string  `gorm:"column:title;not null;index:idx_title" json:"title"`
	Artist           string  `gorm:"column:artist;default:'Unknown Artist';index:idx_artist" json:"artist"`
	Album            *string `gorm:"column:album" json:"album"`
	Duration         int     `gorm:"column:duration;default:0" json:"duration"`
	StorageFilename  string  `gorm:"column:storage_filename;not null" json:"storage_filename"`
	OriginalFilename string  `gorm:"column:original_filename;not null" json:"original_filename"`
	FileSize         int64   `gorm:"column:file_size;default:0" json:"file_size"`
	MimeType         string  `gorm:"column:mime_type;default:'audio/mpeg'" json:"mime_type"`
	Base

	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
}
