package advertisements

import (
	"net/http"
	"strconv"
	"testing"

	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
	advModels "api-tests-template/internal/managers/advertisements/models"
)

func HttpGetAdvertisements(t *testing.T, token, textForSearch string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().
		Get(path.AdvertisementsPath).
		Query("search", textForSearch).
		Expect(t).
		Status(http.StatusOK).
		End().Response
}

// HttpPostAdvertisement выполняет POST-запрос создания объявления
func HttpPostAdvertisement(t *testing.T, token string, request advModels.CreateAdvertisementRequest) *http.Response {
	req := apiRunner.GetRunner().Auth(token).Create().
		Post(path.AdvertisementPath)

	if request.Title != nil {
		req = req.MultipartFormData("title", *request.Title)
	}
	if request.Description != nil {
		req = req.MultipartFormData("description", *request.Description)
	}
	if request.Price != nil {
		req = req.MultipartFormData("price", strconv.FormatFloat(*request.Price, 'f', -1, 64))
	}
	if request.Quantity != nil {
		req = req.MultipartFormData("quantity", strconv.Itoa(*request.Quantity))
	}
	for _, photoPath := range request.PhotoPaths {
		req = req.MultipartFile("photos", photoPath)
	}

	return req.Expect(t).End().Response
}

// HttpGetAdvertisement выполняет GET-запрос получения объявления по его id
func HttpGetAdvertisement(t *testing.T, id string) *http.Response {
	return apiRunner.GetRunner().Create().
		Get(path.AdvertisementPath).
		Query("id", id).
		Expect(t).
		End().Response
}

// HttpGetAdvertisementPhotos выполняет GET-запрос получения фотографий объявления
func HttpGetAdvertisementPhotos(t *testing.T, token, id string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().
		Get(path.AdvertisementPhotosPath(id)).
		Expect(t).
		End().Response
}

// HttpDeleteAdvertisement выполняет DELETE-запрос удаления объявления по его id
func HttpDeleteAdvertisement(t *testing.T, token, id string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().
		Delete(path.AdvertisementPath).
		Query("id", id).
		Expect(t).
		End().Response
}
