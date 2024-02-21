package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
)

func (h *Handler) userIdentity(ctx *gin.Context) {
	header := ctx.GetHeader(authorizationHeader)
	if header == "" {
		NewErrorResponse(ctx, http.StatusUnauthorized, "empty auth header")
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		NewErrorResponse(ctx, http.StatusUnauthorized, "invalid auth header")
		return
	}

	userId, err := h.services.Authorization.ParseToken(headerParts[1])
	if err != nil {
		NewErrorResponse(ctx, http.StatusUnauthorized, err.Error())
	}

	ctx.Set(userCtx, userId)
}

func getUserId(ctx *gin.Context) (int, error) {
	id, ok := ctx.Get(userCtx)
	if !ok {
		errorMessage := "user id not found"
		NewErrorResponse(ctx, http.StatusInternalServerError, errorMessage)
		return 0, errors.New(errorMessage)
	}

	idInt, ok := id.(int)
	if !ok {
		errorMessage := "user id is of invalid type"
		NewErrorResponse(ctx, http.StatusInternalServerError, errorMessage)
		return 0, errors.New(errorMessage)
	}

	return idInt, nil
}
