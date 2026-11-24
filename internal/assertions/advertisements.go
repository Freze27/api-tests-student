package assertions

import (
	"slices"
	"strconv"
	"strings"
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

// ContainsAdvertisementById проверяет, что среди объявлений есть объявление с указанным id
func ContainsAdvertisementById(t *testing.T, items []models.Advertisement, id string) {
	found := slices.ContainsFunc(items, func(item models.Advertisement) bool {
		return item.Id == id
	})
	require.Truef(t, found, "Объявление с id %q не найдено среди %d результатов", id, len(items))
}

// AdvertisementMatchesRequest проверяет, что параметры объявления совпадают с параметрами запроса на создание
func AdvertisementMatchesRequest(t *testing.T, request models.CreateAdvertisementRequest, actual models.Advertisement) {
	require.NotNil(t, request.Title, "В запросе на создание не задан title")
	require.NotNil(t, request.Description, "В запросе на создание не задан description")
	require.Equal(t, *request.Title, actual.Title, "title не совпадает с переданным при создании")
	require.Equal(t, *request.Description, actual.Description, "description не совпадает с переданным при создании")

	var expectedPrice *float64
	if request.Price != nil {
		price, err := strconv.ParseFloat(*request.Price, 64)
		require.NoError(t, err, "price в запросе должен быть числом")
		expectedPrice = &price
	}
	require.Equal(t, expectedPrice, actual.Price, "price не совпадает с переданным при создании")

	var expectedQuantity *int
	if request.Quantity != nil {
		quantity, err := strconv.Atoi(*request.Quantity)
		require.NoError(t, err, "quantity в запросе должен быть целым числом")
		expectedQuantity = &quantity
	}
	require.Equal(t, expectedQuantity, actual.Quantity, "quantity не совпадает с переданным при создании")

	require.Lenf(t, actual.Photos, len(request.PhotoPaths),
		"Количество фотографий должно быть %d", len(request.PhotoPaths))
}

// ErrorMessageContainsField проверяет код ошибки и то, что сообщение об ошибке указывает на некорректное поле
func ErrorMessageContainsField(t *testing.T, actual models.ErrorResponse, expectedError, expectedField string) {
	require.Equal(t, expectedError, actual.Error, "Некорректное значение error параметра")
	require.Containsf(t, strings.ToLower(actual.Message), expectedField,
		"Сообщение об ошибке %q должно указывать на поле %q", actual.Message, expectedField)
}
