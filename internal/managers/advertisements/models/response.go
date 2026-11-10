package models

type AdvertisementsList struct {
	Items []Advertisement `json:"items"`
}

type Advertisement struct {
	Title string `json:"title"`
}
