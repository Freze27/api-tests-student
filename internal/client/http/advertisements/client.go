package advertisements

import (
	"net/http"
	"testing"

	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
)

func HttpGetAdvertisements(t *testing.T, token, textForSearch string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().
		Get(path.AdvertisementsPath).
		Query("search", textForSearch).
		Expect(t).
		Status(http.StatusOK).
		End().Response
}
