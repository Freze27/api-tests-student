package models

type CreateAdvertisementRequest struct {
	Title       *string
	Description *string
	Price       *string
	Quantity    *string
	PhotoPaths  []string
}
