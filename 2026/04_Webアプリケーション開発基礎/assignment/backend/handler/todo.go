package handler

import (
	"errors"
	"net/http"
	"strconv"

	"backend/model"
	"backend/service"

	"github.com/labstack/echo/v4"
)

// TodoHandler はTODO関連のHTTPハンドラー
type TodoHandler struct {
	svc service.TodoService
}

// NewTodoHandler は新しいTodoHandlerを作成する
func NewTodoHandler(svc service.TodoService) *TodoHandler {
	return &TodoHandler{svc: svc}
}

// RegisterRoutes はルーティングを登録する
func (h *TodoHandler) RegisterRoutes(e *echo.Echo) {
	e.GET("/todos", h.List)
	e.POST("/todos", h.Create)
	e.PATCH("/todos/:id", h.Update)
	e.DELETE("/todos/:id", h.Delete)
}

// List は全てのTODOを取得する
func (h *TodoHandler) List(c echo.Context) error {
	todos := h.svc.List()
	return c.JSON(http.StatusOK, NewTodoResponses(todos))
}

// Create は新しいTODOを作成する
func (h *TodoHandler) Create(c echo.Context) error {
	var req CreateTodoRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	todo, err := h.svc.Create(req.Title)
	if err != nil {
		if errors.Is(err, model.ErrEmptyTitle) || errors.Is(err, model.ErrInvalidID) {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, NewTodoResponse(todo))
}

// Update は既存のTODOを更新する
func (h *TodoHandler) Update(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}

	var req UpdateTodoRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	if req.Title == nil && req.Completed == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "no fields to update")
	}

	todo, err := h.svc.Update(id, req.Title, req.Completed)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "todo not found")
		}
		if errors.Is(err, model.ErrEmptyTitle) {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, NewTodoResponse(todo))
}

// Delete は指定IDのTODOを削除する
func (h *TodoHandler) Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}

	if err := h.svc.Delete(id); err != nil {
		if errors.Is(err, service.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "todo not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}
