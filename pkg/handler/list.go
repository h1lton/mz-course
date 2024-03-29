package handler

import (
	"github.com/gin-gonic/gin"
	todo "github.com/h1lton/mz-course"
	"net/http"
	"strconv"
)

func (h *Handler) createList(ctx *gin.Context) {
	userId, err := getUserId(ctx)
	if err != nil {
		return
	}

	var input todo.TodoList
	if err = ctx.BindJSON(&input); err != nil {
		NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.TodoList.Create(userId, input)
	if err != nil {
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, map[string]any{"id": id})
}

type getAllListsResponse struct {
	Data []todo.TodoList `json:"data"`
}

func (h *Handler) getAllLists(ctx *gin.Context) {
	userId, err := getUserId(ctx)
	if err != nil {
		return
	}

	lists, err := h.services.TodoList.GetAll(userId)
	if err != nil {
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, getAllListsResponse{lists})
}

func (h *Handler) getListById(ctx *gin.Context) {
	userId, err := getUserId(ctx)
	if err != nil {
		return
	}

	listId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	list, err := h.services.TodoList.GetById(userId, listId)
	if err != nil {
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, list)
}

func (h *Handler) updateList(ctx *gin.Context) {
	userId, err := getUserId(ctx)
	if err != nil {
		return
	}

	listId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	var input todo.UpdateListInput
	if err = ctx.BindJSON(&input); err != nil {
		NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err = h.services.TodoList.Update(userId, listId, input); err != nil {
		NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, statusResponse{"ok"})
}

func (h *Handler) deleteList(ctx *gin.Context) {
	userId, err := getUserId(ctx)
	if err != nil {
		return
	}

	listId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	err = h.services.TodoList.Delete(userId, listId)
	if err != nil {
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, statusResponse{"ok"})
}
