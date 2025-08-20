package api

import (
	commonErr "github.com/Base-111/backend/pkg/errors"
	"github.com/Base-111/backend/pkg/logs"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
)

func WriteError(ctx *gin.Context, err error) {
	logs.Error(ctx.Request, err)

	var apiErr Error
	var notFoundErr *commonErr.NotFoundError

	if errors.As(err, &apiErr) {
		ctx.JSON(apiErr.StatusCode(), apiErr.Error())

		return
	}

	if errors.As(err, &notFoundErr) {
		ctx.JSON(http.StatusNotFound, notFoundErr.Error())

		return
	}

	ctx.AbortWithStatus(http.StatusInternalServerError)
}
