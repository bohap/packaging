package itest

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"server/internal/appcontext"
	"server/internal/controller"
	"server/internal/service"
	"server/test/stub"
	"testing"
)

func TestPacksGet(t *testing.T) {
	// given
	gin.SetMode(gin.TestMode)

	appContext := appcontext.BuildAppContext()
	// insert initial date
	initSql := `INSERT INTO packs(size) VALUES (100), (200), (1000);`
	if err := appContext.DB.Exec(initSql).Error; err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := cleanupDb(appContext.DB); err != nil {
			t.Fatal(err)
		}
	}()

	router := controller.SetupRouter(appContext)

	expectedResponse := []int{100, 200, 1000}
	expectedResponseJson, _ := json.Marshal(expectedResponse)

	// when
	response := executeRequest(router, "GET", "/packs", nil)

	// then
	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, string(expectedResponseJson), response.Body.String())
}

func TestPackGet(t *testing.T) {
	// given
	gin.SetMode(gin.TestMode)

	repoStub := stub.PacksRepositoryStub{Error: errors.New("error")}
	appContext := &appcontext.AppContext{
		PacksService: service.NewPacksService(repoStub),
	}

	router := controller.SetupRouter(appContext)

	expectedResponse := map[string]interface{}{
		"error": "failed to get packs",
	}
	expectedResponseJson, _ := json.Marshal(expectedResponse)

	// when
	response := executeRequest(router, "GET", "/packs", nil)

	// then
	assert.Equal(t, http.StatusInternalServerError, response.Code)
	assert.Equal(t, string(expectedResponseJson), response.Body.String())
}

func TestPacksSync_ValidRequest_EmptyState(t *testing.T) {
	// given
	gin.SetMode(gin.TestMode)

	appContext := appcontext.BuildAppContext()
	defer func() {
		if err := cleanupDb(appContext.DB); err != nil {
			t.Fatal(err)
		}
	}()

	router := controller.SetupRouter(appContext)

	requestBody := map[string]interface{}{
		"packs": []int{100, 200, 1000},
	}

	// when
	response := executeRequest(router, "POST", "/packs", requestBody)

	// then
	assert.Equal(t, http.StatusOK, response.Code)

	sizes, err := appContext.PacksService.GetPacks()
	assert.Nil(t, err)
	assert.Equal(t, []int{100, 200, 1000}, sizes)
}

func TestPacksSync_ValidRequest_NotEmptyState(t *testing.T) {
	// given
	gin.SetMode(gin.TestMode)

	appContext := appcontext.BuildAppContext()

	// insert initial date
	initSql := `INSERT INTO packs(size) VALUES (100), (200), (1000);`
	if err := appContext.DB.Exec(initSql).Error; err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := cleanupDb(appContext.DB); err != nil {
			t.Fatal(err)
		}
	}()

	router := controller.SetupRouter(appContext)

	requestBody := map[string]interface{}{
		"packs": []int{100, 201, 1000},
	}

	// when
	response := executeRequest(router, "POST", "/packs", requestBody)

	// then
	assert.Equal(t, http.StatusOK, response.Code)

	sizes, err := appContext.PacksService.GetPacks()
	assert.Nil(t, err)
	assert.Equal(t, []int{100, 201, 1000}, sizes)
}

func TestPacksSync_ValidRequest_EmptyList(t *testing.T) {
	// given
	gin.SetMode(gin.TestMode)

	appContext := appcontext.BuildAppContext()
	// insert initial date
	initSql := `INSERT INTO packs(size) VALUES (100), (200), (1000);`
	if err := appContext.DB.Exec(initSql).Error; err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := cleanupDb(appContext.DB); err != nil {
			t.Fatal(err)
		}
	}()

	router := controller.SetupRouter(appContext)

	requestBody := map[string]interface{}{
		"packs": []int{},
	}

	// when
	response := executeRequest(router, "POST", "/packs", requestBody)

	// then
	assert.Equal(t, http.StatusOK, response.Code)

	sizes, err := appContext.PacksService.GetPacks()
	assert.Nil(t, err)
	assert.Equal(t, 0, len(sizes))
}

func TestPacksSync_ValidRequest_PacksSaveError(t *testing.T) {
	// given
	gin.SetMode(gin.TestMode)

	repoStub := stub.PacksRepositoryStub{Error: errors.New("error")}
	appContext := &appcontext.AppContext{
		PacksService: service.NewPacksService(repoStub),
	}

	router := controller.SetupRouter(appContext)

	requestBody := map[string]interface{}{
		"packs": []int{100, 200, 1000},
	}

	expectedResponse := map[string]interface{}{
		"error": "failed to sync packs",
	}
	expectedResponseJson, _ := json.Marshal(expectedResponse)

	// when
	response := executeRequest(router, "POST", "/packs", requestBody)

	// then
	assert.Equal(t, http.StatusInternalServerError, response.Code)
	assert.Equal(t, string(expectedResponseJson), response.Body.String())
}

func TestPacksSync_InvalidRequest_Missing(t *testing.T) {
	// given
	gin.SetMode(gin.TestMode)

	appContext := appcontext.BuildAppContext()
	router := controller.SetupRouter(appContext)

	requestBody := map[string]interface{}{
		"test": "test",
	}

	// when
	response := executeRequest(router, "POST", "/packs", requestBody)

	// then
	assert.Equal(t, http.StatusBadRequest, response.Code)
}

func TestPackItems_ValidRequest(t *testing.T) {
	// given
	gin.SetMode(gin.TestMode)

	appContext := appcontext.BuildAppContext()
	// insert packs data
	initSql := `INSERT INTO packs(size) VALUES (100), (200), (1000);`
	if err := appContext.DB.Exec(initSql).Error; err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := cleanupDb(appContext.DB); err != nil {
			t.Fatal(err)
		}
	}()

	router := controller.SetupRouter(appContext)

	requestBody := map[string]interface{}{
		"numberOfItems": 1001,
	}

	expectedResponse := map[string]interface{}{
		"1000": 1, "100": 1,
	}
	expectedResponseJson, _ := json.Marshal(expectedResponse)

	// when
	response := executeRequest(router, "POST", "/package", requestBody)

	// then
	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, string(expectedResponseJson), response.Body.String())
}

func TestPackItems_ValidRequest_EmptyPacksConfig(t *testing.T) {
	// given
	gin.SetMode(gin.TestMode)

	appContext := appcontext.BuildAppContext()
	router := controller.SetupRouter(appContext)

	requestBody := map[string]interface{}{
		"numberOfItems": 1001,
	}

	expectedResponse := map[string]interface{}{
		"error": "no packs configured",
	}
	expectedResponseJson, _ := json.Marshal(expectedResponse)

	// when
	response := executeRequest(router, "POST", "/package", requestBody)

	// then
	assert.Equal(t, http.StatusBadRequest, response.Code)
	assert.Equal(t, string(expectedResponseJson), response.Body.String())
}

func TestPackItems_ValidRequest_RepoError(t *testing.T) {
	// given
	gin.SetMode(gin.TestMode)

	repoStub := stub.PacksRepositoryStub{Error: errors.New("error")}
	appContext := &appcontext.AppContext{
		PackingService: service.NewPackagingService(service.NewPacksService(repoStub)),
	}

	router := controller.SetupRouter(appContext)

	requestBody := map[string]interface{}{
		"numberOfItems": 1001,
	}

	expectedResponse := map[string]interface{}{
		"error": "failed to pack items",
	}
	expectedResponseJson, _ := json.Marshal(expectedResponse)

	// when
	response := executeRequest(router, "POST", "/package", requestBody)

	// then
	assert.Equal(t, http.StatusInternalServerError, response.Code)
	assert.Equal(t, string(expectedResponseJson), response.Body.String())
}

func TestPackItems_InvalidRequest(t *testing.T) {
	// given
	gin.SetMode(gin.TestMode)

	appContext := appcontext.BuildAppContext()

	router := controller.SetupRouter(appContext)

	requestBody := map[string]interface{}{
		"numberOfItems": 0,
	}

	// when
	response := executeRequest(router, "POST", "/package", requestBody)

	// then
	assert.Equal(t, http.StatusBadRequest, response.Code)
}

func executeRequest(
	router *gin.Engine, method string, url string, body map[string]interface{},
) *httptest.ResponseRecorder {
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest(method, "/api"+url, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	return w
}
