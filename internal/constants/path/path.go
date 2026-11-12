package path

import "fmt"

const (
	AuthLoginPath        = "/auth/login"
	MyAdvertisementsPath = "/my/advertisements"
	AdvertisementsPath   = "/advertisements"
	AdvertisementPath    = "/advertisement"
	FavoritesPath        = "/favorites"
)

// AdvertisementPhotosPath возвращает путь к фотографиям объявления: /advertisements/{id}/photos
func AdvertisementPhotosPath(advertisementId string) string {
	return fmt.Sprintf("%s/%s/photos", AdvertisementsPath, advertisementId)
}

// FavoritePath возвращает путь к объявлению в избранном: /favorites/{ad_id}
func FavoritePath(advertisementId string) string {
	return fmt.Sprintf("%s/%s", FavoritesPath, advertisementId)
}
