package createAdvertisement

import (
	"fmt"
	"net/http"
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"api-tests-template/internal/assertions"
	advManager "api-tests-template/internal/managers/advertisements"
	advModels "api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/managers/auth"
	authModels "api-tests-template/internal/managers/auth/models"

	base "api-tests-template/tests"
)

// Валидные данные для создания объявления
const (
	titlePrefix      = "ByAutotest"
	validDescription = "Объявление было создано автотестом"
	validPrice       = 1687.4
	validQuantity    = 3
	incorrectToken   = "incorrect_token"
)

const photoFilePath = "testdata/photo.jpg"

// Ограничения на параметры объявления
const (
	maxTitleLength       = 50
	maxDescriptionLength = 500
	maxPrice             = 1000000
	maxQuantity          = 100
	maxPhotosCount       = 3
)

// Ожидаемые значения поля error в ответе
const (
	errorInvalidRequest = "invalid_request"
	errorUnauthorized   = "unauthorized"
)

// Названия полей, которые должны упоминаться в сообщении об ошибке
const (
	fieldTitle       = "title"
	fieldDescription = "description"
	fieldPrice       = "price"
	fieldQuantity    = "quantity"
	fieldPhotos      = "photo"
)

type TestSuite struct {
	suite.Suite
	loginData    authModels.LoginOkResponse
	createdAdIds []string
}

func TestSuiteRun(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

func (s *TestSuite) SetupSuite() {
	base.SetupSuite()

	base.Precondition("Авторизация пользователя с кредами из переменных окружения и получение его параметров")
	s.loginData = auth.Login(s.T(), os.Getenv("TEST_LOGIN"), os.Getenv("TEST_PASSWORD"))

	base.Precondition("Проверяем, что тестовая фотография для загрузки существует")
	require.FileExists(s.T(), photoFilePath, "Не найден тестовый файл фотографии")
}

func (s *TestSuite) TearDownSuite() {
	base.TearDownSuite()
}

func (s *TestSuite) SetupTest() {
	base.Precondition("Очищаем список объявлений, созданных в тесте")
	s.createdAdIds = nil
}

func (s *TestSuite) TearDownTest() {
	for _, id := range s.createdAdIds {
		base.Postcondition(fmt.Sprintf("Удаляем созданное в тесте объявление %s", id))
		advManager.DeleteAdvertisement(s.T(), s.loginData.Token, id)
	}
}

// validCreateRequest возвращает запрос на создание объявления со всеми параметрами и уникальным названием
func (s *TestSuite) validCreateRequest() advModels.CreateAdvertisementRequest {
	title := fmt.Sprintf("%s %s", titlePrefix, uuid.NewString())
	description := validDescription
	price := validPrice
	quantity := validQuantity

	return advModels.CreateAdvertisementRequest{
		Title:       &title,
		Description: &description,
		Price:       &price,
		Quantity:    &quantity,
		PhotoPaths:  []string{photoFilePath},
	}
}

// Positive: Проверяем создание объявления со всеми параметрами: оно доступно по id, находится в поиске и имеет загруженные фото
func (s *TestSuite) TestCreateAdvertisementPositive() {
	request := s.validCreateRequest()

	var createdAd advModels.Advertisement
	s.Run("Создаём объявление со всеми возможными параметрами", func() {
		createdAd = advManager.CreateAdvertisement(s.T(), s.loginData.Token, request)
		require.NotEmpty(s.T(), createdAd.Id, "id созданного объявления пустой")
		s.createdAdIds = append(s.createdAdIds, createdAd.Id)
	})

	s.Run("Проверяем, что объявление доступно по id и возвращает те же параметры", func() {
		ad := advManager.GetAdvertisement(s.T(), createdAd.Id)
		require.Equal(s.T(), createdAd.Id, ad.Id, "id объявления не совпадает")
		require.Equal(s.T(), s.loginData.User.Id, ad.UserId, "user_id не совпадает с текущим пользователем")
		assertions.AdvertisementMatchesRequest(s.T(), request, ad)
	})

	s.Run("Проверяем, что объявление находится в поиске по полному названию", func() {
		list := advManager.GetAdvertisementsList(s.T(), s.loginData.Token, *request.Title)
		assertions.ContainsAdvertisementById(s.T(), list.Items, createdAd.Id)
	})

	s.Run("Проверяем, что фотография загружена и доступна по ручке фотографий объявления", func() {
		photos := advManager.GetAdvertisementPhotos(s.T(), s.loginData.Token, createdAd.Id)
		require.Lenf(s.T(), photos, len(request.PhotoPaths), "Количество фотографий должно быть %d", len(request.PhotoPaths))
		require.Lenf(s.T(), createdAd.Photos, len(request.PhotoPaths),
			"Количество фотографий в ответе на создание должно быть %d", len(request.PhotoPaths))
		for i, photo := range photos {
			require.Equal(s.T(), createdAd.Photos[i].Id, photo.Id, "id фотографии не совпадает с ответом на создание")
			require.NotEmpty(s.T(), photo.Url, "url фотографии пустой")
		}
	})
}

// Negative: Проверяем 400 ошибку при создании объявления без обязательных или с некорректными параметрами
func (s *TestSuite) TestCreateAdvertisementInvalidParams() {
	tooLongTitle := strings.Repeat("a", maxTitleLength+1)
	tooLongDescription := strings.Repeat("d", maxDescriptionLength+1)
	tooManyPhotos := slices.Repeat([]string{photoFilePath}, maxPhotosCount+1)

	cases := []struct {
		name          string
		modify        requestModifier
		expectedField string
	}{
		{"Без title", withoutTitle(), fieldTitle},
		{"Пустой title", withTitle(""), fieldTitle},
		{"Слишком длинный title", withTitle(tooLongTitle), fieldTitle},
		{"Без description", withoutDescription(), fieldDescription},
		{"Пустой description", withDescription(""), fieldDescription},
		{"Слишком длинный description", withDescription(tooLongDescription), fieldDescription},
		{"Отрицательная цена", withPrice(-1), fieldPrice},
		{"Цена больше максимальной", withPrice(maxPrice + 1), fieldPrice},
		{"Отрицательное количество", withQuantity(-1), fieldQuantity},
		{"Количество больше максимального", withQuantity(maxQuantity + 1), fieldQuantity},
		{"Без фотографий", withPhotos(nil), fieldPhotos},
		{"Фотографий больше максимального количества", withPhotos(tooManyPhotos), fieldPhotos},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			request := s.validCreateRequest()
			tc.modify(&request)

			var errorResponse advModels.ErrorResponse
			s.Run("Отправляем запрос на создание объявления с некорректными параметрами", func() {
				errorResponse = advManager.CreateAdvertisementWithError(s.T(), s.loginData.Token, request, http.StatusBadRequest)
			})

			s.Run("Проверяем, что ошибка указывает на некорректное поле", func() {
				assertions.ErrorMessageContainsField(s.T(), errorResponse, errorInvalidRequest, tc.expectedField)
			})
		})
	}
}

// Negative: Проверяем 401 ошибку при создании объявления без токена и с неверным токеном
func (s *TestSuite) TestCreateAdvertisementUnauthorized() {
	cases := []struct {
		name  string
		token string
	}{
		{"Без авторизационного токена", ""},
		{"С неверным авторизационным токеном", incorrectToken},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			var errorResponse advModels.ErrorResponse
			s.Run("Отправляем запрос на создание объявления", func() {
				errorResponse = advManager.CreateAdvertisementWithError(s.T(), tc.token, s.validCreateRequest(), http.StatusUnauthorized)
			})

			s.Run("Проверяем, что показывается ошибка авторизации", func() {
				require.Equal(s.T(), errorUnauthorized, errorResponse.Error, "Некорректное значение error параметра")
				require.NotEmpty(s.T(), errorResponse.Message, "Сообщение об ошибке пустое")
			})
		})
	}
}

// requestModifier изменяет запрос на создание объявления для негативного сценария
type requestModifier func(request *advModels.CreateAdvertisementRequest)

// withTitle подменяет название объявления в запросе
func withTitle(title string) requestModifier {
	return func(request *advModels.CreateAdvertisementRequest) { request.Title = &title }
}

// withoutTitle убирает поле title из запроса
func withoutTitle() requestModifier {
	return func(request *advModels.CreateAdvertisementRequest) { request.Title = nil }
}

// withDescription подменяет описание объявления в запросе
func withDescription(description string) requestModifier {
	return func(request *advModels.CreateAdvertisementRequest) { request.Description = &description }
}

// withoutDescription убирает поле description из запроса
func withoutDescription() requestModifier {
	return func(request *advModels.CreateAdvertisementRequest) { request.Description = nil }
}

// withPrice подменяет цену объявления в запросе
func withPrice(price float64) requestModifier {
	return func(request *advModels.CreateAdvertisementRequest) { request.Price = &price }
}

// withQuantity подменяет количество в запросе
func withQuantity(quantity int) requestModifier {
	return func(request *advModels.CreateAdvertisementRequest) { request.Quantity = &quantity }
}

// withPhotos подменяет список фотографий в запросе; nil — запрос без фотографий
func withPhotos(photoPaths []string) requestModifier {
	return func(request *advModels.CreateAdvertisementRequest) { request.PhotoPaths = photoPaths }
}
