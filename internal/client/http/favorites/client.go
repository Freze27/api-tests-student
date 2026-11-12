package favorites

import (
	"net/http"
	"strconv"
	"testing"

	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
	favModels "api-tests-template/internal/managers/favorites/models"
)

// favoritesPageSize максимальный размер страницы избранного
const favoritesPageSize = 50

// HttpGetFavorites выполняет GET-запрос получения списка избранного авторизованного пользователя
func HttpGetFavorites(t *testing.T, token string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().
		Get(path.FavoritesPath).
		Query("page_size", strconv.Itoa(favoritesPageSize)).
		Expect(t).
		End().Response
}

// HttpPostFavorite выполняет POST-запрос добавления объявления в избранное
func HttpPostFavorite(t *testing.T, token string, request favModels.AddFavoriteRequest) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().
		Post(path.FavoritesPath).
		ContentType("application/json").
		JSON(request).
		Expect(t).
		End().Response
}

// HttpDeleteFavorite выполняет DELETE-запрос удаления объявления из избранного
func HttpDeleteFavorite(t *testing.T, token, advertisementId string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().
		Delete(path.FavoritePath(advertisementId)).
		Expect(t).
		End().Response
}
