package handlers

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lalizita/shard-banking/internal/services/account/model"
	"github.com/lalizita/shard-banking/internal/services/account/service"
)

type Handler struct {
	echo           *echo.Echo
	AccountService service.IAccountService
}

func NewAccountHandler(e *echo.Echo, service service.IAccountService) *Handler {
	return &Handler{
		echo:           e,
		AccountService: service,
	}
}

func (h *Handler) RegisterRoutes() {
	h.echo.GET("/accounts", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "ok",
		})
	})

	h.echo.POST("/accounts", h.CreateAccount)
}

func (h *Handler) CreateAccount(c echo.Context) error {
	slog.Info("Creating account...")
	//@TODO: Validate fields
	var newAccount model.CreateAccountRequest
	var err error
	var ctx = c.Request().Context()

	if err = c.Bind(&newAccount); err != nil {
		return c.JSON(http.StatusInternalServerError, []byte(`{"message": "failed bind request"}`))
	}

	newAccModel := model.Account{
		Name:  newAccount.Name,
		Email: newAccount.Email,
	}

	created, err := h.AccountService.CreateAccount(ctx, newAccModel)
	if err != nil {
		slog.Error("Create account failed", slog.Any("error", err))
		return c.JSON(http.StatusInternalServerError, []byte(`{"message": "failed create account"}`))
	}

	slog.Info("account created successfully!!!")
	return c.JSON(http.StatusCreated, created)
}
