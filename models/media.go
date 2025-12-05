package models

type Media struct {
	ID          string `json:"id" bson:"_id"`
	UserID      string `json:"user_id" bson:"user_id" binding:"required"`
	AlbumID     string `json:"album_id,omitempty" bson:"album_id,omitempty"`
	Type        string `json:"type" bson:"type" binding:"required,oneof=image video"`
	URL         string `json:"url" bson:"url" binding:"required,url"`
	Title       string `json:"title,omitempty" bson:"title,omitempty"`
	Description string `json:"description,omitempty" bson:"description,omitempty"`
	CreatedAt   int64  `json:"created_at" bson:"created_at"`
	UpdatedAt   int64  `json:"updated_at" bson:"updated_at"`
}

type CreatedImageInput struct {
	URL         string `json:"url" binding:"required,url"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	AlbumID     string `json:"album_id,omitempty"`
	Location   	string `json:"location,omitempty"`
	DeviceInfo  string `json:"device_info,omitempty"`
	CaptureDate string `json:"capture_date,omitempty"`
}

type CreatedVideoInput struct {
	URL         string `json:"url" binding:"required,url"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	AlbumID     string `json:"album_id,omitempty"`
}
