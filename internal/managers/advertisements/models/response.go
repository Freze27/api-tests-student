package models

type AdvertisementsList struct {
	Items []Advertisement `json:"items"`
}

type Advertisement struct {
	Id          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Price       *float64 `json:"price"`
	Quantity    *int     `json:"quantity"`
	UserId      string   `json:"user_id"`
	Photos      []Photo  `json:"photos"`
	CreatedAt   string   `json:"created_at"`
	UpdatedAt   string   `json:"updated_at"`
}

type Photo struct {
	Id        string `json:"id"`
	Url       string `json:"url"`
	SortOrder int    `json:"sort_order"`
	CreatedAt string `json:"created_at"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}
