package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"server/internal/appcontext"
	"server/internal/model"
)

func HandlePackageRequest(requestContext *gin.Context, appContext *appcontext.AppContext) {
	var req model.ProductsPackageRequest

	if err := requestContext.ShouldBindJSON(&req); err != nil {
		requestContext.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := appContext.PackingService.PackItems(req.NumberOfItems)
	if err != nil {
		var emptyPacksConfigError *model.EmptyPacksConfig
		if errors.As(err, &emptyPacksConfigError) {
			requestContext.JSON(http.StatusBadRequest, gin.H{"error": emptyPacksConfigError.Error()})
		} else {
			requestContext.JSON(http.StatusInternalServerError, gin.H{"error": "failed to pack items"})
		}
		return
	}

	requestContext.JSON(http.StatusOK, response)
}

func HandleGetPacksRequest(requestContext *gin.Context, appContext *appcontext.AppContext) {
	response, err := appContext.PacksService.GetPacks()
	if err != nil {
		requestContext.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get packs"})
		return
	}

	requestContext.JSON(http.StatusOK, response)
}

func HandlePacksSyncRequest(requestContext *gin.Context, appContext *appcontext.AppContext) {
	var req model.PacksSyncRequest

	if err := requestContext.ShouldBindJSON(&req); err != nil {
		requestContext.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := appContext.PacksService.SyncPacks(req.Packs); err != nil {
		requestContext.JSON(http.StatusInternalServerError, gin.H{"error": "failed to sync packs"})
		return
	}

	requestContext.JSON(http.StatusOK, map[string]string{"status": "OK"})
}
