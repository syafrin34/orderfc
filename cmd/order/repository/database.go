package repository

import (
	"context"
	"encoding/json"
	"errors"
	"orderfc/infrastructure/constant"
	"orderfc/infrastructure/logger"
	"orderfc/models"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func (r *OrderRepository) WithTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	tx := r.Database.Begin().WithContext(ctx)
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// insert order

func (r *OrderRepository) InsertorderTx(ctx context.Context, tx *gorm.DB, order *models.Order) error {
	err := tx.WithContext(ctx).Table("orders").Create(order).Error
	return err
}

// inser order detail
func (r *OrderRepository) InsertorderDetailTx(ctx context.Context, tx *gorm.DB, orderDetail *models.OrderDetail) error {
	err := tx.WithContext(ctx).Table("order_detail").Create(orderDetail).Error
	return err
}

func (r *OrderRepository) CheckIdempotency(ctx context.Context, token string) (bool, error) {
	var log models.OrderRequestLog
	err := r.Database.WithContext(ctx).Table("order_request_log").First(&log, "idempotency_token=?", token).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}

		return false, err
	}
	return true, nil
}

// save idempotency
func (r *OrderRepository) SaveIdempotency(ctx context.Context, token string) error {
	log := models.OrderRequestLog{
		IdempotencyToken: token,
		CreateTime:       time.Now(),
	}

	err := r.Database.WithContext(ctx).Table("order_request_log").Create(&log).Error
	if err != nil {
		return err
	}
	return nil
}

//order histroy by user id or status

func (r *OrderRepository) GetOrderHistoryByUserID(ctx context.Context, param models.OrderHistoryParam) ([]models.OrderHistoryResponse, error) {
	//join tabel order dan order detail

	//prepare query
	var queryResults []models.OrderHistoryResult
	query := r.Database.WithContext(ctx).Table("orders AS o").
		Select(`o.id, o.amount, o.total_qty, o.status, o.payment_method, o.shipping_address, od.products, od.order_history`).
		Joins("JOIN order_detail od ON o.order_detail_id=od.id").
		Where("o.user_id = ?", param.UserID)

	if param.Status > 0 {
		query = query.Where("o.status = ?", param.Status)
	}

	err := query.Order("o.id DESC").Scan(&queryResults).Error
	if err != nil {
		return nil, err
	}
	var results []models.OrderHistoryResponse
	for _, result := range queryResults {
		var products []models.CheckoutItem
		var orderHistory []models.StatusHistory
		err = json.Unmarshal([]byte(result.Products), &products)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal([]byte(result.OrderHistory), &orderHistory)
		if err != nil {
			return nil, err
		}

		results = append(results, models.OrderHistoryResponse{
			OrderID:         result.ID,
			TotalAmount:     result.Amount,
			TotalQty:        result.TotalQuantity,
			Status:          constant.OrderStatusTranslated[result.Status],
			PaymentMethod:   result.PaymentMethod,
			ShippingAddress: result.ShippingAddress,
			Products:        products,
			History:         orderHistory,
		})
		// var hasils []map[string]interface{}
		// err := r.Database.WithContext(ctx).Table("orders").Find(&hasils).Error
		// if err != nil {
		// 	fmt.Println("gagal ambil data")
		// }

		logger.Logger.WithFields(logrus.Fields{
			"data orders": result,
		}).Infof("some thing wrong %v", result)
	}
	return results, nil

}

// update order status
func (r *OrderRepository) UpdateOrderStatus(ctx context.Context, orderID int64, status int) error {
	err := r.Database.Table("orders").WithContext(ctx).Model(&models.Order{}).
		Updates(map[string]interface{}{
			"status":      status,
			"update_time": time.Now(),
		}).Where("order_id = ? ", orderID).Error
	if err != nil {
		return err
	}
	return nil
}

// get order info
func (r *OrderRepository) GetOrderInfoByOrderID(ctx context.Context, orderID int64) (models.Order, error) {
	var result models.Order
	err := r.Database.Table("orders").WithContext(ctx).Where("order_id = ?", orderID).Find(&result).Error
	if err != nil {
		return models.Order{}, err
	}
	return result, nil
}

// get order detail info
func (r *OrderRepository) GetOrderDetailByOrderDetailID(ctx context.Context, orderDetailID int64) (models.OrderDetail, error) {
	var result models.OrderDetail
	err := r.Database.Table("order_detail").WithContext(ctx).Where("id = ?", orderDetailID).Find(&result).Error
	if err != nil {
		return models.OrderDetail{}, err
	}
	return result, nil
}
