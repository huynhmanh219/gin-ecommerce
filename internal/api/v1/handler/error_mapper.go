package handler

import (
	stdErrors "errors"
	"net/http"

	"github.com/gin-gonic/gin"
	appErrors "huynhmanh.com/gin/internal/common/errors"
)

func writeError(c *gin.Context, err error) {
	var appErr *appErrors.AppError
	if stdErrors.As(err, &appErr) {
		switch appErr.Code {
		case "VALIDATION_ERROR":
			c.JSON(http.StatusBadRequest,appErr)
		case "UNAUTHORIZED":
			c.JSON(http.StatusUnauthorized,appErr)
		case "NOT_FOUND":
			c.JSON(http.StatusNotFound,appErr)
		case "CONFLICT":
			c.JSON(http.StatusInternalServerError,appErr)
		}
		return
	}
	c.JSON(http.StatusInternalServerError,gin.H{
		"code": "INTERNAL_ERROR",
		"message": "ERROR NOT IDENTIFY",
	})
}

