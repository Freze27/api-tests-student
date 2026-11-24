package myAdvertisement

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/tidwall/gjson"

	"api-tests-template/internal/assertions"
	advManager "api-tests-template/internal/managers/advertisements"
	advModels "api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/managers/auth"
	"api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/managers/myAdvertisements"

	base "api-tests-template/tests"
)

const (
	titlePrefix      = "MyAdsAutotest"
	validDescription = "Объявление для проверки своих объявлении, создано автотестом"
	photoFilePath    = "testdata/photo.jpg"
)

type TestSuite struct {
	suite.Suite
	loginData models.LoginOkResponse
	ownAdId   string
}

func TestSuiteRun(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

func (s *TestSuite) SetupSuite() {
	base.SetupSuite()

	base.Precondition("Авторизация пользователя с кредами из переменных окружения и получение его параметров")
	s.loginData = auth.Login(s.T(), os.Getenv("TEST_LOGIN"), os.Getenv("TEST_PASSWORD"))
}

func (s *TestSuite) TearDownSuite() {
	base.TearDownSuite()
}

func (s *TestSuite) SetupTest() {
	s.ownAdId = ""

	base.Precondition("Создаём объявление от основного пользователя")
	ad := advManager.CreateAdvertisement(s.T(), s.loginData.Token, newAdvertisementRequest())
	s.ownAdId = ad.Id
	require.NotEmpty(s.T(), s.ownAdId, "id созданного объявления пустой")
}

func (s *TestSuite) TearDownTest() {
	if s.ownAdId != "" {
		base.Postcondition(fmt.Sprintf("Удаляем объявление основного пользователя %s", s.ownAdId))
		advManager.DeleteAdvertisement(s.T(), s.loginData.Token, s.ownAdId)
	}
}

func newAdvertisementRequest() advModels.CreateAdvertisementRequest {
	title := fmt.Sprintf("%s %s", titlePrefix, uuid.NewString())
	description := validDescription

	return advModels.CreateAdvertisementRequest{
		Title:       &title,
		Description: &description,
		PhotoPaths:  []string{photoFilePath},
	}
}

// Positive: Проверяем получение собственных объявлений в профиле
func (s *TestSuite) TestGetMyAdvertisementsPositive() {
	var list advModels.AdvertisementsList
	s.Run("Получаем список собственных объявлений", func() {
		body := myAdvertisements.GetMyAdvertisements(s.T(), s.loginData.Token, http.StatusOK, 50)
		require.NoError(s.T(), json.Unmarshal([]byte(body), &list))
	})

	s.Run("Проверяем, что созданное объявление есть в списке", func() {
		assertions.ContainsAdvertisementById(s.T(), list.Items, s.ownAdId)
	})

	s.Run("Проверяем, что наши объявления имеют принадлежность к нашему пользователю и объявления имеют НЕ пустые параметры", func() {
		for _, item := range list.Items {
			require.Equal(s.T(), s.loginData.User.Id, item.UserId, "user_id не совпадает")

			require.NotEmpty(s.T(), item.Id, "id пустое")
			require.NotEmpty(s.T(), item.CreatedAt, "created_at пустое")
			require.NotEmpty(s.T(), item.UpdatedAt, "updated_at пустое")
			require.NotEmpty(s.T(), item.Title, "Название пустое")
			require.NotEmpty(s.T(), item.Description, "Описание пустое")
			require.NotEmpty(s.T(), item.Photos, "Список фотографий пуст")
		}
	})
}

// Negative: Проверяем 401 ошибку при попытке получить собственные объявления без токена
func (s *TestSuite) TestGetMyAdvertisementsIncorrectToken() {
	var advertisementsBody string
	s.Run("Получаем список собственных объявлений с неправильным авторизационным токеном", func() {
		advertisementsBody = myAdvertisements.GetMyAdvertisements(s.T(), "incorrect_token", http.StatusUnauthorized, 50)
	})

	s.Run("Проверяем, что показывается правильная ошибка", func() {
		require.Equal(s.T(), "unauthorized", gjson.Get(advertisementsBody, "error").String(),
			"Некорректное значение error параметра")
		require.Equal(s.T(), "Invalid or expired token", gjson.Get(advertisementsBody, "message").String(),
			"Некорректное значение message параметра")
	})
}
