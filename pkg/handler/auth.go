package handler

import (
	"github.com/gin-gonic/gin"
	todo "github.com/h1lton/mz-course"
)

func (h Handler) signUp(ctx *gin.Context) {
	var input todo.User

	if err := ctx.BindJSON(&input); err != nil {
		NewErrorResponse(ctx, 400, err.Error())
		return
	}

	id, err := h.services.Authorization.CreateUser(input)
	if err != nil {
		NewErrorResponse(ctx, 500, err.Error())
		return
	}

	ctx.JSON(200, map[string]any{
		"id": id,
	})
}

func (h Handler) signIn(c *gin.Context) {

}
