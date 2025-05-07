package handler

import (
	"net/http"
	"orderfc/cmd/order/usecase"
	"orderfc/infrastructure/logger"
	"orderfc/models"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type OrderHandler struct {
	OrderUsecase usecase.OrderUsecase
}

func NewOrderHandler(oh usecase.OrderUsecase) *OrderHandler {
	return &OrderHandler{
		OrderUsecase: oh,
	}
}

func (h *OrderHandler) CheckOutOrder(c *gin.Context) {
	var param models.CheckoutRequest

	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":        "invalid request",
			"error detail": err.Error()})
		return
	}

	// auth session Login -->vcek user id
	usrIDStr, isExists := c.Get("user_id")
	if isExists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"erro": "unauthorized",
		})
		return
	}

	userID, ok := usrIDStr.(float64)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid user id",
		})
		return
	}

	if len(param.Items) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid parameter",
		})
		return
	}
	param.UserID = int64(userID)

	orderID, err := h.OrderUsecase.CheckOutOrder(c.Request.Context(), &param)
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"param": param,
		}).Errorf("h.orderUseCae.CekecoutOrder got error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "order created",
		"order_id": orderID,
	})
	return

}
