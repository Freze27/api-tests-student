package favorites

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/client/http/favorites"
	advModels "api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/managers/favorites/models"
)

// GetFavorites возвращает список избранного пользователя, ожидает 200
func GetFavorites(t *testing.T, token string) models.FavoritesList {
	raw := favorites.HttpGetFavorites(t, token)
	defer raw.Body.Close()

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err)
	require.Equalf(t, http.StatusOK, raw.StatusCode,
		"HTTP status code должен быть %d, тело ответа: %s", http.StatusOK, string(body))

	resp := models.FavoritesList{}
	err = json.Unmarshal(body, &resp)
	require.NoError(t, err)

	return resp
}

// GetFavoritesWithError запрашивает список избранного, который должен завершиться ошибкой,
// проверяет ожидаемый статус-код и возвращает тело ошибки
func GetFavoritesWithError(t *testing.T, token string, expectedStatusCode int) advModels.ErrorResponse {
	raw := favorites.HttpGetFavorites(t, token)
	defer raw.Body.Close()

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err)
	require.Equalf(t, expectedStatusCode, raw.StatusCode,
		"HTTP status code должен быть %d, тело ответа: %s", expectedStatusCode, string(body))

	resp := advModels.ErrorResponse{}
	err = json.Unmarshal(body, &resp)
	require.NoError(t, err)

	return resp
}

// AddFavorite добавляет объявление в избранное, ожидает 201 и возвращает добавленную запись
func AddFavorite(t *testing.T, token, advertisementId string) models.FavoriteItem {
	raw := favorites.HttpPostFavorite(t, token, models.AddFavoriteRequest{AdvertisementId: advertisementId})
	defer raw.Body.Close()

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err)
	require.Equalf(t, http.StatusCreated, raw.StatusCode,
		"HTTP status code должен быть %d, тело ответа: %s", http.StatusCreated, string(body))

	resp := models.FavoriteItem{}
	err = json.Unmarshal(body, &resp)
	require.NoError(t, err)

	return resp
}

// AddFavoriteWithError отправляет запрос добавления в избранное, который должен завершиться ошибкой,
// проверяет ожидаемый статус-код и возвращает тело ошибки
func AddFavoriteWithError(t *testing.T, token string, request models.AddFavoriteRequest,
	expectedStatusCode int) advModels.ErrorResponse {
	raw := favorites.HttpPostFavorite(t, token, request)
	defer raw.Body.Close()

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err)
	require.Equalf(t, expectedStatusCode, raw.StatusCode,
		"HTTP status code должен быть %d, тело ответа: %s", expectedStatusCode, string(body))

	resp := advModels.ErrorResponse{}
	err = json.Unmarshal(body, &resp)
	require.NoError(t, err)

	return resp
}

// RemoveFavorite удаляет объявление из избранного, ожидает 204
func RemoveFavorite(t *testing.T, token, advertisementId string) {
	raw := favorites.HttpDeleteFavorite(t, token, advertisementId)
	defer raw.Body.Close()

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err)
	require.Equalf(t, http.StatusNoContent, raw.StatusCode,
		"HTTP status code должен быть %d, тело ответа: %s", http.StatusNoContent, string(body))
}

// RemoveFavoriteWithError отправляет запрос удаления из избранного, который должен завершиться ошибкой,
// проверяет ожидаемый статус-код и возвращает тело ошибки
func RemoveFavoriteWithError(t *testing.T, token, advertisementId string,
	expectedStatusCode int) advModels.ErrorResponse {
	raw := favorites.HttpDeleteFavorite(t, token, advertisementId)
	defer raw.Body.Close()

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err)
	require.Equalf(t, expectedStatusCode, raw.StatusCode,
		"HTTP status code должен быть %d, тело ответа: %s", expectedStatusCode, string(body))

	resp := advModels.ErrorResponse{}
	err = json.Unmarshal(body, &resp)
	require.NoError(t, err)

	return resp
}
