package path

import "fmt"

const (
	AuthLoginPath        = "/auth/login"
	MyAdvertisementsPath = "/my/advertisements"
	AdvertisementsPath   = "/advertisements"
	AdvertisementPath    = "/advertisement"
)

// AdvertisementPhotosPath возвращает путь к фотографиям объявления: /advertisements/{id}/photos
func AdvertisementPhotosPath(advertisementId string) string {
	return fmt.Sprintf("%s/%s/photos", AdvertisementsPath, advertisementId)
}
