package advertisements

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/client/http/advertisements"
	"api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/managers/myAdvertisements"
)

func GetAdvertisementsList(t *testing.T, token, text string) models.AdvertisementsList {
	raw := advertisements.HttpGetAdvertisements(t, token, text)
	defer raw.Body.Close()

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err)

	resp := models.AdvertisementsList{}
	err = json.Unmarshal(body, &resp)
	require.NoError(t, err)

	return resp
}

// CreateAdvertisement создаёт объявление, ожидает 201 и возвращает созданное объявление
func CreateAdvertisement(t *testing.T, token string, request models.CreateAdvertisementRequest) models.Advertisement {
	raw := advertisements.HttpPostAdvertisement(t, token, request)
	defer raw.Body.Close()

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err)
	require.Equalf(t, http.StatusCreated, raw.StatusCode,
		"HTTP status code должен быть %d, тело ответа: %s", http.StatusCreated, string(body))

	resp := models.Advertisement{}
	err = json.Unmarshal(body, &resp)
	require.NoError(t, err)

	return resp
}

// CreateAdvertisementWithError отправляет запрос создания объявления, который должен завершиться ошибкой,
// проверяет ожидаемый статус-код и возвращает тело ошибки
func CreateAdvertisementWithError(t *testing.T, token string, request models.CreateAdvertisementRequest,
	expectedStatusCode int) models.ErrorResponse {
	raw := advertisements.HttpPostAdvertisement(t, token, request)
	defer raw.Body.Close()

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err)
	require.Equalf(t, expectedStatusCode, raw.StatusCode,
		"HTTP status code должен быть %d, тело ответа: %s", expectedStatusCode, string(body))

	resp := models.ErrorResponse{}
	err = json.Unmarshal(body, &resp)
	require.NoError(t, err)

	return resp
}

// GetAdvertisement возвращает объявление по id, ожидает 200
func GetAdvertisement(t *testing.T, id string) models.Advertisement {
	raw := advertisements.HttpGetAdvertisement(t, id)
	defer raw.Body.Close()

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err)
	require.Equalf(t, http.StatusOK, raw.StatusCode,
		"HTTP status code должен быть %d, тело ответа: %s", http.StatusOK, string(body))

	resp := models.Advertisement{}
	err = json.Unmarshal(body, &resp)
	require.NoError(t, err)

	return resp
}

// GetAdvertisementPhotos возвращает список фотографий объявления, ожидает 200
func GetAdvertisementPhotos(t *testing.T, token, id string) []models.Photo {
	raw := advertisements.HttpGetAdvertisementPhotos(t, token, id)
	defer raw.Body.Close()

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err)
	require.Equalf(t, http.StatusOK, raw.StatusCode,
		"HTTP status code должен быть %d, тело ответа: %s", http.StatusOK, string(body))

	var resp []models.Photo
	err = json.Unmarshal(body, &resp)
	require.NoError(t, err)

	return resp
}

// DeleteAdvertisement удаляет объявление по id, ожидает 204
func DeleteAdvertisement(t *testing.T, token, id string) {
	raw := advertisements.HttpDeleteAdvertisement(t, token, id)
	defer raw.Body.Close()

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err)
	require.Equalf(t, http.StatusNoContent, raw.StatusCode,
		"HTTP status code должен быть %d, тело ответа: %s", http.StatusNoContent, string(body))
}

// myAdvertisementsLimit максимальное количество собственных объявлений, которое запрашивается для очистки
const myAdvertisementsLimit = 50

// DeleteMyAdvertisementsByTitlePrefix удаляет объявления пользователя, название которых начинается с префикса,
// и возвращает id удалённых объявлений. Нужна для очистки объявлений, оставшихся от прерванных прогонов тестов
func DeleteMyAdvertisementsByTitlePrefix(t *testing.T, token, titlePrefix string) []string {
	body := myAdvertisements.GetMyAdvertisements(t, token, http.StatusOK, myAdvertisementsLimit)

	myAds := models.AdvertisementsList{}
	err := json.Unmarshal([]byte(body), &myAds)
	require.NoError(t, err)

	var deletedIds []string
	for _, ad := range myAds.Items {
		if strings.HasPrefix(ad.Title, titlePrefix) {
			DeleteAdvertisement(t, token, ad.Id)
			deletedIds = append(deletedIds, ad.Id)
		}
	}

	return deletedIds
}
