package assertions

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/require"

	favModels "api-tests-template/internal/managers/favorites/models"
)

// FavoritesContainAdvertisement проверяет, что в списке избранного есть объявление с указанным id
func FavoritesContainAdvertisement(t *testing.T, items []favModels.FavoriteItem, advertisementId string) {
	require.Truef(t, favoritesHaveAdvertisement(items, advertisementId),
		"Объявление с id %q не найдено в избранном среди %d записей", advertisementId, len(items))
}

// FavoritesNotContainAdvertisement проверяет, что в списке избранного нет объявления с указанным id
func FavoritesNotContainAdvertisement(t *testing.T, items []favModels.FavoriteItem, advertisementId string) {
	require.Falsef(t, favoritesHaveAdvertisement(items, advertisementId),
		"Объявление с id %q всё ещё находится в избранном", advertisementId)
}

// favoritesHaveAdvertisement возвращает true, если в списке избранного есть объявление с указанным id
func favoritesHaveAdvertisement(items []favModels.FavoriteItem, advertisementId string) bool {
	return slices.ContainsFunc(items, func(item favModels.FavoriteItem) bool {
		return item.Advertisement.Id == advertisementId
	})
}
