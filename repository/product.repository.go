package repository

import (
	"golang-rest-api-jwt/entity"

	"gorm.io/gorm"
)

type ProductRepository interface {
	InsertProduct(product entity.Product) (entity.Product, error)

	All(userID string) ([]entity.Product, error)
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{
		db: db,
	}
}

func (r *productRepository) InsertProduct(product entity.Product) (entity.Product, error) {
	r.db.Save(&product)
	r.db.Preload("User").Find(&product)
	return product, nil
}

func (r *productRepository) All(userID string) ([]entity.Product, error) {
	products := []entity.Product{}
	r.db.Preload("User").
		Where("user_id = ?", userID).
		Find(&products)
	return products, nil
}
