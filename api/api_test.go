package api_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/snyk/snyk-code-review-exercise/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//comment: this is a nice black-box test that tests the app end-to-end, it would be nice to have:
// additional black box unit tests, and maybe the ability to cache and control the HTTP responses from the registry.
func TestPackageHandler(t *testing.T) {
	handler := api.New()
	server := httptest.NewServer(handler)
	defer server.Close()

	resp, err := server.Client().Get(server.URL + "/package/react/16.13.0")
	require.Nil(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	require.Nil(t, err)

	var data api.NpmPackageVersion
	err = json.Unmarshal(body, &data)
	require.Nil(t, err)

	assert.Equal(t, "react", data.Name)
	assert.Equal(t, "16.13.0", data.Version)

	fixture, err := os.Open(filepath.Join("testdata", "react-16.13.0.json"))
	require.Nil(t, err)
	var fixtureObj api.NpmPackageVersion
	require.Nil(t, json.NewDecoder(fixture).Decode(&fixtureObj))

	assert.Equal(t, fixtureObj, data)
}
