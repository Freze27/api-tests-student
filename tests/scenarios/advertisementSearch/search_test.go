package advertisementSearch

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"api-tests-template/internal/assertions"
	advManager "api-tests-template/internal/managers/advertisements"
	advModels "api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/managers/auth"
	"api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/utils"

	base "api-tests-template/tests"
)

const (
	titlePrefix      = "VW Tiguan"
	validDescription = "Объявление для проверки поиска, создано автотестом"
	photoFilePath    = "testdata/tiguan.jpg"
)

type TestSuite struct {
	suite.Suite
	loginData models.LoginOkResponse
	adId      string
	adTitle   string
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
	s.adId = ""

	base.Precondition("Создаём объявление для проверки поиска")
	ad := advManager.CreateAdvertisement(s.T(), s.loginData.Token, newAdvertisementRequest())
	s.adId = ad.Id
	s.adTitle = ad.Title
	require.NotEmpty(s.T(), s.adId, "id созданного объявления пустой")
	require.GreaterOrEqual(s.T(), len([]rune(s.adTitle)), 3, "Название слишком короткое для поиска по трём частям")
}

func (s *TestSuite) TearDownTest() {
	if s.adId != "" {
		base.Postcondition(fmt.Sprintf("Удаляем объявление для проверки поиска %s", s.adId))
		advManager.DeleteAdvertisement(s.T(), s.loginData.Token, s.adId)
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

func (s *TestSuite) titleParts() (string, string, string) {
	title := []rune(s.adTitle)
	partLength := len(title) / 3
	return string(title[:partLength]), string(title[partLength : 2*partLength]), string(title[2*partLength:])
}

// Positive: Проверяем поиск по началу, середине и концу названия и независимость от регистра
func (s *TestSuite) TestSearchByTitle_Substring() {
	searchBeginning, searchMiddle, searchEnd := s.titleParts()
	cases := []struct {
		name   string
		search string
	}{
		{"начало", searchBeginning},
		{"середина", searchMiddle},
		{"конец", searchEnd},
		{"нижний регистр", strings.ToLower(searchBeginning)},
		{"верхний регистр", strings.ToUpper(searchBeginning)},
		{"смешанный регистр", utils.MixCase(searchBeginning)},
	}
	for _, tc := range cases {
		s.Run(tc.name, func() {
			list := advManager.GetAdvertisementsList(s.T(), s.loginData.Token, tc.search)
			assertions.ContainsAdvertisementById(s.T(), list.Items, s.adId)
		})
	}
}
