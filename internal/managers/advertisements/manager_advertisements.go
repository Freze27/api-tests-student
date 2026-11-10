package advertisements

import (
	"api-tests-template/internal/client/http/advertisements"
	"api-tests-template/internal/managers/advertisements/models"
	"encoding/json"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
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
