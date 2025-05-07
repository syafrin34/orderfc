package repository

import (
	"context"
	"errors"
	"orderfc/models"
	"time"

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
