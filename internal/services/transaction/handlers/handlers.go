package handlers

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lalizita/shard-banking/internal/services/transaction/model"
	"github.com/lalizita/shard-banking/internal/services/transaction/service"
)

type Handler struct {
	echo               *echo.Echo
	TransactionService service.ITransactionService
}

func NewTransactionHandler(e *echo.Echo, service service.ITransactionService) *Handler {
	return &Handler{
		echo:               e,
		TransactionService: service,
	}
}

func (h *Handler) RegisterRoutes() {
	h.echo.POST("/transactions/credit", h.CreateCreditTransaction)
}

func (h *Handler) CreateCreditTransaction(c echo.Context) error {
	slog.Info("Creating credit transaction...")

	var req model.CreateTransactionRequest
	if err := c.Bind(&req); err != nil {
		slog.Error("Create transaction failed", slog.Any("error", err))
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "failed to bind request"})
	}

	tx := model.Transaction{
		ClientID:  req.ClientID,
		Amount:    req.Amount,
		EntryType: model.EntryTypeCredit,
	}

	if err := h.TransactionService.CreateTransaction(c.Request().Context(), tx); err != nil {
		slog.Error("Create transaction failed", slog.Any("error", err))
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "failed to create transaction"})
	}

	slog.Info("transaction created successfully")
	return c.JSON(http.StatusCreated, echo.Map{"message": "transaction created successfully"})
}
