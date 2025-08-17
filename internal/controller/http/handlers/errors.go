package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/Base-111/backend/internal/common/errors/api"
	"github.com/Base-111/backend/internal/common/errors/database"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"log/slog"
	"net/http"
)

func WriteError(ctx *gin.Context, err error) {
	slog.ErrorContext(ctx, fmt.Sprintf("%s : %v", err.Error(), err))

	var dbNotFound *database.NotFoundError
	if errors.As(err, &dbNotFound) {
		ctx.JSON(http.StatusNotFound, api.ErrorResponse{
			Message: "not found",
		})

		return
	}

	var dbRowExists *database.RowExistsError
	if errors.As(err, &dbRowExists) {
		ctx.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Message: "saving entity error",
		})

		return
	}

	var unmarshalError *json.UnmarshalTypeError
	if errors.As(err, &unmarshalError) {
		ctx.JSON(http.StatusBadRequest, api.ErrorResponse{
			Message: "invalid parsing",
		})

		return
	}

	var validateError *api.ValidationError
	if errors.As(err, &validateError) {
		var errs *api.ValidationError
		errors.As(err, &errs)
		ctx.JSON(http.StatusBadRequest, api.ErrorResponse{
			Message: "validation error",
			Data:    errs.GetValidationErrors(),
		})

		return
	}

	ctx.JSON(http.StatusInternalServerError, api.ErrorResponse{
		Message: "unknown error",
	})
}
