package advertisementSearch

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"api-tests-template/internal/assertions"
	advList "api-tests-template/internal/managers/advertisements"
	advModels "api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/managers/auth"
	"api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/utils"

	base "api-tests-template/tests"
)

const (
	existingAdTitle = "VW Tiguan"
	searchBeginning = "VW"
	searchMiddle    = "Tig"
	searchEnd       = "uan"
)

type TestSuite struct {
	suite.Suite
	loginData models.LoginOkResponse
}

func TestSuiteRun(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

func (s *TestSuite) SetupSuite() {
	base.SetupSuite()

	base.Precondition("Проверяем, что константы заданы")
	require.NotEmpty(s.T(), existingAdTitle, "задайте existingAdTitle")
	require.NotEmpty(s.T(), searchBeginning, "задайте searchBeginning")
	require.NotEmpty(s.T(), searchMiddle, "задайте searchMiddle")
	require.NotEmpty(s.T(), searchEnd, "задайте searchEnd")

	base.Precondition("Авторизация пользователя с кредами из переменных окружения и получение его параметров")
	s.loginData = auth.Login(s.T(), os.Getenv("TEST_LOGIN"), os.Getenv("TEST_PASSWORD"))
}

func (s *TestSuite) TearDownSuite() {
	base.TearDownSuite()
}

// Positive: Проверяем, что поиск по части названия находит существующее объявление
func (s *TestSuite) TestSearchByTitle() {
	var list advModels.AdvertisementsList
	s.Run("Ищем объявления по части названия", func() {
		list = advList.GetAdvertisementsList(s.T(), s.loginData.Token, searchMiddle)
	})
	s.Run("Проверяем, что в ответе есть хотя бы одно объявление", func() {
		require.NotEmpty(s.T(), list.Items, "Поиск вернул пустой список объявлений")
	})
	s.Run("Проверяем, что среди найденных есть наше объявление", func() {
		assertions.ContainsAdvertisementByTitle(s.T(), list.Items, existingAdTitle)
	})
}

// Positive: Проверяем поиск по началу, середине и концу названия и независимость от регистра
func (s *TestSuite) TestSearchByTitle_Substring() {
	cases := []struct {
		name   string
		search string
	}{
		{"начало", searchBeginning},
		{"середина", searchMiddle},
		{"конец", searchEnd},
		{"нижний регистр", strings.ToLower(searchMiddle)},
		{"верхний регистр", strings.ToUpper(searchMiddle)},
		{"смешанный регистр", utils.MixCase(searchMiddle)},
	}
	for _, tc := range cases {
		s.Run(tc.name, func() {
			list := advList.GetAdvertisementsList(s.T(), s.loginData.Token, tc.search)
			assertions.ContainsAdvertisementByTitle(s.T(), list.Items, existingAdTitle)
		})
	}
}
