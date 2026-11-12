package favorites

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"api-tests-template/internal/assertions"
	advManager "api-tests-template/internal/managers/advertisements"
	advModels "api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/managers/auth"
	authModels "api-tests-template/internal/managers/auth/models"
	favManager "api-tests-template/internal/managers/favorites"
	favModels "api-tests-template/internal/managers/favorites/models"

	base "api-tests-template/tests"
)

// Данные объявлений, создаваемых для тестов избранного
const (
	titlePrefix      = "FavAutotest"
	validDescription = "Объявление для проверки избранного, создано автотестом"
	photoFilePath    = "testdata/photo.jpg"
)

const incorrectToken = "incorrect_token"

const (
	errorUnauthorized   = "unauthorized"
	errorInvalidRequest = "invalid_request"
	errorNotFound       = "not_found"
	errorForbidden      = "forbidden"
	errorAlreadyExists  = "already_exists"
)

type TestSuite struct {
	suite.Suite
	loginData       authModels.LoginOkResponse
	sellerLoginData authModels.LoginOkResponse
	sellerAdId      string
	ownAdIds        []string
}

func TestSuiteRun(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

func (s *TestSuite) SetupSuite() {
	base.SetupSuite()

	base.Precondition("Авторизация основного пользователя, который работает с избранным")
	s.loginData = auth.Login(s.T(), os.Getenv("TEST_LOGIN"), os.Getenv("TEST_PASSWORD"))

	base.Precondition("Авторизация второго пользователя, владельца объявлений для избранного")
	s.sellerLoginData = auth.Login(s.T(), os.Getenv("TEST_SECOND_LOGIN"), os.Getenv("TEST_SECOND_PASSWORD"))

	base.Precondition("Удаляем объявления обоих пользователей, оставшиеся от прерванных прогонов")
	for _, token := range []string{s.loginData.Token, s.sellerLoginData.Token} {
		deletedIds := advManager.DeleteMyAdvertisementsByTitlePrefix(s.T(), token, titlePrefix+" ")
		base.Precondition(fmt.Sprintf("Удалено оставшихся объявлений: %d", len(deletedIds)))
	}

	base.Precondition("Проверяем, что тестовая фотография для загрузки существует")
	require.FileExists(s.T(), photoFilePath, "Не найден тестовый файл фотографии")
}

func (s *TestSuite) TearDownSuite() {
	base.TearDownSuite()
}

func (s *TestSuite) SetupTest() {
	s.sellerAdId = ""
	s.ownAdIds = nil

	base.Precondition("Второй пользователь создаёт объявление для добавления в избранное")
	sellerAd := advManager.CreateAdvertisement(s.T(), s.sellerLoginData.Token, newAdvertisementRequest())
	s.sellerAdId = sellerAd.Id
}

func (s *TestSuite) TearDownTest() {
	if s.sellerAdId != "" {
		base.Postcondition(fmt.Sprintf("Удаляем объявление второго пользователя %s", s.sellerAdId))
		advManager.DeleteAdvertisement(s.T(), s.sellerLoginData.Token, s.sellerAdId)
	}

	for _, id := range s.ownAdIds {
		base.Postcondition(fmt.Sprintf("Удаляем объявление основного пользователя %s", id))
		advManager.DeleteAdvertisement(s.T(), s.loginData.Token, id)
	}
}

// newAdvertisementRequest возвращает запрос на создание объявления с обязательными полями и уникальным названием
func newAdvertisementRequest() advModels.CreateAdvertisementRequest {
	title := fmt.Sprintf("%s %s", titlePrefix, uuid.NewString())
	description := validDescription

	return advModels.CreateAdvertisementRequest{
		Title:       &title,
		Description: &description,
		PhotoPaths:  []string{photoFilePath},
	}
}

// Positive: Проверяем полный цикл избранного: добавление, удаление и исчезновение после удаления самого объявления
func (s *TestSuite) TestFavoritesFullCycle() {
	s.Run("Добавляем объявление второго пользователя в избранное", func() {
		favorite := favManager.AddFavorite(s.T(), s.loginData.Token, s.sellerAdId)
		require.Equal(s.T(), s.sellerAdId, favorite.Advertisement.Id, "В избранное добавлено не то объявление")
	})

	s.Run("Проверяем, что объявление появилось в списке избранного", func() {
		list := favManager.GetFavorites(s.T(), s.loginData.Token)
		assertions.FavoritesContainAdvertisement(s.T(), list.Items, s.sellerAdId)
	})

	s.Run("Удаляем объявление из избранного", func() {
		favManager.RemoveFavorite(s.T(), s.loginData.Token, s.sellerAdId)
	})

	s.Run("Проверяем, что объявление пропало из списка избранного", func() {
		list := favManager.GetFavorites(s.T(), s.loginData.Token)
		assertions.FavoritesNotContainAdvertisement(s.T(), list.Items, s.sellerAdId)
	})

	s.Run("Снова добавляем объявление в избранное", func() {
		favManager.AddFavorite(s.T(), s.loginData.Token, s.sellerAdId)
	})

	deletedAdId := s.sellerAdId
	s.Run("Второй пользователь удаляет само объявление", func() {
		advManager.DeleteAdvertisement(s.T(), s.sellerLoginData.Token, deletedAdId)
		s.sellerAdId = ""
	})

	s.Run("Проверяем, что удалённое объявление пропало из избранного", func() {
		list := favManager.GetFavorites(s.T(), s.loginData.Token)
		assertions.FavoritesNotContainAdvertisement(s.T(), list.Items, deletedAdId)
	})
}

// Negative: Проверяем ошибки добавления в избранное: несуществующее, без id, собственное и уже добавленное объявление
func (s *TestSuite) TestAddFavoriteErrors() {
	cases := []struct {
		name               string
		prepareAdId        func() string
		expectedStatusCode int
		expectedError      string
	}{
		{"Несуществующее объявление", uuid.NewString, http.StatusNotFound, errorNotFound},
		{"Без advertisement_id", func() string { return "" }, http.StatusBadRequest, errorInvalidRequest},
		{"Собственное объявление", s.createOwnAdvertisement, http.StatusForbidden, errorForbidden},
		{"Объявление уже в избранном", s.addSellerAdToFavorites, http.StatusConflict, errorAlreadyExists},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			request := favModels.AddFavoriteRequest{AdvertisementId: tc.prepareAdId()}

			var errorResponse advModels.ErrorResponse
			s.Run("Добавляем объявление в избранное", func() {
				errorResponse = favManager.AddFavoriteWithError(s.T(), s.loginData.Token, request, tc.expectedStatusCode)
			})

			s.Run("Проверяем, что показывается правильная ошибка", func() {
				assertions.ErrorResponseHasCode(s.T(), errorResponse, tc.expectedError)
			})
		})
	}
}

// Negative: Проверяем 404 ошибку при удалении из избранного несуществующего объявления
func (s *TestSuite) TestRemoveNonexistentAdvertisementFromFavorites() {
	var errorResponse advModels.ErrorResponse
	s.Run("Удаляем из избранного объявление с несуществующим id", func() {
		errorResponse = favManager.RemoveFavoriteWithError(s.T(), s.loginData.Token, uuid.NewString(), http.StatusNotFound)
	})

	s.Run("Проверяем, что показывается правильная ошибка", func() {
		assertions.ErrorResponseHasCode(s.T(), errorResponse, errorNotFound)
	})
}

// Negative: Проверяем 401 ошибку на всех ручках избранного без токена и с неверным токеном
func (s *TestSuite) TestFavoritesUnauthorized() {
	tokens := []struct {
		name  string
		token string
	}{
		{"Без авторизационного токена", ""},
		{"С неверным авторизационным токеном", incorrectToken},
	}

	endpoints := []struct {
		name string
		send func(token string) advModels.ErrorResponse
	}{
		{"GET /favorites", s.getFavoritesUnauthorized},
		{"POST /favorites", s.addFavoriteUnauthorized},
		{"DELETE /favorites/{ad_id}", s.removeFavoriteUnauthorized},
	}

	for _, endpoint := range endpoints {
		for _, tc := range tokens {
			s.Run(fmt.Sprintf("%s: %s", endpoint.name, tc.name), func() {
				var errorResponse advModels.ErrorResponse
				s.Run("Отправляем запрос", func() {
					errorResponse = endpoint.send(tc.token)
				})

				s.Run("Проверяем, что показывается ошибка авторизации", func() {
					assertions.ErrorResponseHasCode(s.T(), errorResponse, errorUnauthorized)
				})
			})
		}
	}
}

// createOwnAdvertisement создаёт объявление от лица основного пользователя и возвращает его id
func (s *TestSuite) createOwnAdvertisement() string {
	var ownAdId string
	s.Run("Создаём объявление от лица основного пользователя", func() {
		ownAd := advManager.CreateAdvertisement(s.T(), s.loginData.Token, newAdvertisementRequest())
		ownAdId = ownAd.Id
		s.ownAdIds = append(s.ownAdIds, ownAdId)
	})
	return ownAdId
}

// addSellerAdToFavorites добавляет объявление второго пользователя в избранное и возвращает его id
func (s *TestSuite) addSellerAdToFavorites() string {
	s.Run("Заранее добавляем объявление второго пользователя в избранное", func() {
		favManager.AddFavorite(s.T(), s.loginData.Token, s.sellerAdId)
	})
	return s.sellerAdId
}

// getFavoritesUnauthorized запрашивает список избранного с переданным токеном, ожидает 401
func (s *TestSuite) getFavoritesUnauthorized(token string) advModels.ErrorResponse {
	return favManager.GetFavoritesWithError(s.T(), token, http.StatusUnauthorized)
}

// addFavoriteUnauthorized добавляет объявление второго пользователя в избранное с переданным токеном, ожидает 401
func (s *TestSuite) addFavoriteUnauthorized(token string) advModels.ErrorResponse {
	request := favModels.AddFavoriteRequest{AdvertisementId: s.sellerAdId}
	return favManager.AddFavoriteWithError(s.T(), token, request, http.StatusUnauthorized)
}

// removeFavoriteUnauthorized удаляет объявление второго пользователя из избранного с переданным токеном, ожидает 401
func (s *TestSuite) removeFavoriteUnauthorized(token string) advModels.ErrorResponse {
	return favManager.RemoveFavoriteWithError(s.T(), token, s.sellerAdId, http.StatusUnauthorized)
}
