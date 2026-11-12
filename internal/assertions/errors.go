package assertions

import (
	"testing"

	"github.com/stretchr/testify/require"

	advModels "api-tests-template/internal/managers/advertisements/models"
)

// ErrorResponseHasCode проверяет код ошибки (поле error) и то, что в ответе есть сообщение об ошибке
func ErrorResponseHasCode(t *testing.T, actual advModels.ErrorResponse, expectedError string) {
	require.Equal(t, expectedError, actual.Error, "Некорректное значение error параметра")
	require.NotEmpty(t, actual.Message, "В ответе нет сообщения об ошибке (поле message)")
}
