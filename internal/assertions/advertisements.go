package assertions

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/advertisements/models"
)

// ContainsAdvertisementByTitle проверяет, что среди объявлений есть объявление с указанным названием
func ContainsAdvertisementByTitle(t *testing.T, items []models.Advertisement, title string) {
	found := slices.ContainsFunc(items, func(item models.Advertisement) bool {
		return item.Title == title
	})
	require.Truef(t, found, "Объявление с названием %q не найдено среди %d результатов", title, len(items))
}
