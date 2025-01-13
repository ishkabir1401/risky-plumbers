package routes

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"risky-plumbers/settings"
	"sync"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestRouter() *gin.Engine {
	settings.InitializeRiskStore()

	gin.SetMode(gin.TestMode)

	router := gin.Default()
	SetupRoutes(router)

	return router
}

func TestListRisks(t *testing.T) {
	router := setupTestRouter()

	req, _ := http.NewRequest(http.MethodGet, "/v1/risks", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"data": []}`, w.Body.String())
}

func TestCreateRisk(t *testing.T) {
	router := setupTestRouter()

	risk := map[string]string{
		"state":       "open",
		"title":       "Test Risk",
		"description": "This is a test risk",
	}
	body, _ := json.Marshal(risk)

	req, _ := http.NewRequest(http.MethodPost, "/v1/risks", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var createdRisk map[string]map[string]string
	_ = json.Unmarshal(w.Body.Bytes(), &createdRisk)
	assert.Equal(t, "Test Risk", createdRisk["data"]["title"])
	assert.NotEmpty(t, createdRisk["data"]["id"])
}

func TestGetRisk(t *testing.T) {
	router := setupTestRouter()

	risk := map[string]string{
		"state":       "open",
		"title":       "Test Risk",
		"description": "This is a test risk",
	}
	body, _ := json.Marshal(risk)

	req, _ := http.NewRequest(http.MethodPost, "/v1/risks", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var createdRisk map[string]map[string]string
	_ = json.Unmarshal(w.Body.Bytes(), &createdRisk)
	riskID := createdRisk["data"]["id"]

	req, _ = http.NewRequest(http.MethodGet, "/v1/risks/"+riskID, nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var fetchedRisk map[string]map[string]string
	_ = json.Unmarshal(w.Body.Bytes(), &fetchedRisk)
	assert.Equal(t, "Test Risk", fetchedRisk["data"]["title"])
	assert.Equal(t, riskID, fetchedRisk["data"]["id"])
}

func TestCreateRiskWithInvalidState(t *testing.T) {
	router := setupTestRouter()

	risk := map[string]string{
		"state":       "invalid_state",
		"title":       "Invalid State Risk",
		"description": "This risk has an invalid state",
	}
	body, _ := json.Marshal(risk)

	req, _ := http.NewRequest(http.MethodPost, "/v1/risks", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, `{"error": "Invalid state value"}`, w.Body.String())
}

func TestCreateRiskWithEmptyTitle(t *testing.T) {
	router := setupTestRouter()

	risk := map[string]string{
		"state":       "open",
		"title":       "",
		"description": "Risk with empty title",
	}
	body, _ := json.Marshal(risk)

	req, _ := http.NewRequest(http.MethodPost, "/v1/risks", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, `{"error": "Missing required fields"}`, w.Body.String())
}

func TestGetNonExistentRisk(t *testing.T) {
	router := setupTestRouter()

	req, _ := http.NewRequest(http.MethodGet, "/v1/risks/non-existent-id", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.JSONEq(t, `{"error": "Risk not found"}`, w.Body.String())
}

func TestConcurrentRiskCreation(t *testing.T) {
	router := setupTestRouter()

	var wg sync.WaitGroup
	numRequests := 10
	risk := map[string]string{
		"state":       "open",
		"title":       "Concurrent Risk",
		"description": "Testing concurrency",
	}
	body, _ := json.Marshal(risk)

	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			req, _ := http.NewRequest(http.MethodPost, "/v1/risks", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assert.Equal(t, http.StatusOK, w.Code)
		}()
	}
	wg.Wait()

	req, _ := http.NewRequest(http.MethodGet, "/v1/risks", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var risks map[string][]map[string]string
	_ = json.Unmarshal(w.Body.Bytes(), &risks)
	assert.Len(t, risks["data"], numRequests)
}
