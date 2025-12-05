package models

type Album struct {
	ID          string `json:"id" bson:"_id"`
	UserID      string `json:"user_id" bson:"user_id" binding:"required"`
	Name 	  	string `json:"name" bson:"name" binding:"required,min=2,max=200"`
	Description string `json:"description,omitempty" bson:"description,omitempty"`
	CreatedAt   int64  `json:"created_at" bson:"created_at"`
	UpdatedAt   int64  `json:"updated_at" bson:"updated_at"`
}

type UpdateAlbum struct {
	Name        string `json:"name,omitempty" bson:"name,omitempty" binding:"omitempty,min=2,max=200"`
	Description string `json:"description,omitempty" bson:"description,omitempty"`
}

type CreateAlbumInput struct {
	Name        string `json:"name" binding:"required,min=2,max=200"`
	Description string `json:"description,omitempty"`
}

type UpdateAlbumInput struct {
	Name        string `json:"name,omitempty" binding:"omitempty,min=2,max=200"`
	Description string `json:"description,omitempty"`
}
