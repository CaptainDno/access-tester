package main

import (
	"access-tester/src/common"
	"access-tester/src/scraper"
	server2 "access-tester/src/server"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	nethttpmiddleware "github.com/oapi-codegen/nethttp-middleware"
	"github.com/oapi-codegen/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func doTestRequest(t *testing.T, router http.Handler, host string) *httptest.ResponseRecorder {
	var query string
	if host != "" {
		query = fmt.Sprintf("?host=%s", url.QueryEscape(host))
	} else {
		query = ""
	}
	return testutil.NewRequest().Get(fmt.Sprintf("/status%s", query)).GoWithHTTPHandler(t, router).Recorder
}

func getResult(t *testing.T, body *bytes.Buffer) common.Result {
	var result common.Result
	err := json.NewDecoder(body).Decode(&result)
	assert.NoError(t, err, "error unmarshalling response")
	return result
}

func TestAPI(t *testing.T) {
	swagger, err := server2.GetSwagger()
	require.NoError(t, err)
	swagger.Servers = nil

	sv := server2.ServerImpl{Scraper: scraper.NewScraper()}
	t.Log(sv.Scraper.phrases)
	handler := server2.NewStrictHandler(sv, nil)
	router := chi.NewRouter()
	router.Use(nethttpmiddleware.OapiRequestValidator(swagger))
	server2.HandlerFromMux(handler, router)

	t.Run("No hostname provided", func(t *testing.T) {
		assert.Equal(t, http.StatusBadRequest, doTestRequest(t, router, "").Code)
	})

	t.Run("Bad hostname", func(t *testing.T) {
		const host = "http://not+host?.name"
		rr := doTestRequest(t, router, host)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})
	t.Run("Resource with provided hostname cannot respond", func(t *testing.T) {
		const host = "localhost"
		assert.Equal(t, http.StatusInternalServerError, doTestRequest(t, router, host).Code)
	})

	t.Run("Available resource", func(t *testing.T) {
		const host = "google.com"
		rr := doTestRequest(t, router, host)
		assert.Equal(t, http.StatusOK, rr.Code)
		result := getResult(t, rr.Body)
		assert.NotEqual(t, http.StatusForbidden, result.Status)
		assert.Equal(t, host, result.Host)
		assert.Equal(t, false, result.CanBeBlocked)
		assert.Empty(t, result.TriggerPhrasesFound)
	})

	t.Run("Unavailable resource returns 403 code", func(t *testing.T) {
		// Returns 403 in Russia
		// Page contains iframe, that contains text "access denied". We don't render pages,
		const host = "redis.com"
		rr := doTestRequest(t, router, host)
		assert.Equal(t, http.StatusOK, rr.Code)
		result := getResult(t, rr.Body)
		assert.Equal(t, http.StatusForbidden, result.Status)
		assert.Equal(t, host, result.Host)
		assert.Equal(t, true, result.CanBeBlocked)
	})

	t.Run("Unavailable resource returns error page with phrase like 'access denied' with 200 code", func(t *testing.T) {
		const host = "analog.com"
		rr := doTestRequest(t, router, host)
		assert.Equal(t, http.StatusOK, rr.Code)
		result := getResult(t, rr.Body)
		assert.Equal(t, http.StatusOK, result.Status)
		assert.Equal(t, host, result.Host)
		assert.Equal(t, true, result.CanBeBlocked)
		assert.Contains(t, result.TriggerPhrasesFound, strings.ToLower("Access denied under U.S. Export Administration Regulations"), result.TriggerPhrasesFound)
	})
}
