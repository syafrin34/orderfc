package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"orderfc/models"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type OrderRepository struct {
	Database    *gorm.DB
	Redis       *redis.Client
	ProductHost string
}

func NewOrderRepository(db *gorm.DB, redis *redis.Client, productHost string) *OrderRepository {
	return &OrderRepository{
		Database:    db,
		Redis:       redis,
		ProductHost: productHost,
	}
}

func (r *OrderRepository) GetProductInfo(ctx context.Context, productID int64) (models.Product, error) {
	var product models.Product

	//baca ke config --> utk set product host
	url := fmt.Sprintf("%s/v1/product/%d", r.ProductHost, productID)
	// logger.Logger.WithFields(logrus.Fields{
	// 	"info url": url,
	// })
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return models.Product{}, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return models.Product{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return models.Product{}, fmt.Errorf("invalid response - get product info")
	}

	//
	var response models.GetProductInfo
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return models.Product{}, err
	}
	product = response.Product
	return product, nil

}
