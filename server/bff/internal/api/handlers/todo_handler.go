package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/tadasy/mytodo202507/server/bff/internal/api/middleware"
	"github.com/tadasy/mytodo202507/server/bff/internal/clients"
	"github.com/tadasy/mytodo202507/server/bff/internal/models"
)

type TodoHandler struct {
	todoClient *clients.TodoServiceClient
}

func NewTodoHandler(todoClient *clients.TodoServiceClient) *TodoHandler {
	return &TodoHandler{
		todoClient: todoClient,
	}
}

func (h *TodoHandler) CreateTodo(c echo.Context) error {
	userID := middleware.GetUserIDFromContext(c)
	if userID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not authenticated")
	}

	var req models.CreateTodoRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	todo, err := h.todoClient.CreateTodo(c.Request().Context(), userID, req.Title, req.Description)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, todo)
}

func (h *TodoHandler) GetTodo(c echo.Context) error {
	userID := middleware.GetUserIDFromContext(c)
	if userID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not authenticated")
	}

	todoID := c.Param("id")
	if todoID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "todo ID is required")
	}

	todo, err := h.todoClient.GetTodo(c.Request().Context(), todoID, userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, todo)
}

func (h *TodoHandler) ListTodos(c echo.Context) error {
	userID := middleware.GetUserIDFromContext(c)
	if userID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not authenticated")
	}

	// Check if we should list completed todos only
	completedOnly, _ := strconv.ParseBool(c.QueryParam("completed"))

	var todos []*models.Todo
	var err error

	if completedOnly {
		todos, err = h.todoClient.ListCompletedTodos(c.Request().Context(), userID)
	} else {
		todos, err = h.todoClient.ListTodos(c.Request().Context(), userID)
	}

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, todos)
}

func (h *TodoHandler) UpdateTodo(c echo.Context) error {
	userID := middleware.GetUserIDFromContext(c)
	if userID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not authenticated")
	}

	todoID := c.Param("id")
	if todoID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "todo ID is required")
	}

	var req models.UpdateTodoRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	todo, err := h.todoClient.UpdateTodo(c.Request().Context(), todoID, userID, req.Title, req.Description)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, todo)
}

func (h *TodoHandler) MarkTodoComplete(c echo.Context) error {
	userID := middleware.GetUserIDFromContext(c)
	if userID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not authenticated")
	}

	todoID := c.Param("id")
	if todoID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "todo ID is required")
	}

	var req models.MarkTodoCompleteRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	todo, err := h.todoClient.MarkTodoComplete(c.Request().Context(), todoID, userID, req.Completed)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, todo)
}

func (h *TodoHandler) DeleteTodo(c echo.Context) error {
	userID := middleware.GetUserIDFromContext(c)
	if userID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not authenticated")
	}

	todoID := c.Param("id")
	if todoID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "todo ID is required")
	}

	err := h.todoClient.DeleteTodo(c.Request().Context(), todoID, userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "todo deleted successfully"})
}
