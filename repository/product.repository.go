package repository

import (
	"golang-rest-api-jwt/entity"

	"gorm.io/gorm"
)

type ProductRepository interface {
	InsertProduct(product entity.Product) (entity.Product, error)

	All(userID string) ([]entity.Product, error)

	FindOneProductByID(productID string) (entity.Product, error)

	UpdateProduct(product entity.Product) (entity.Product, error)

	DeleteProduct(productID string) error
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

func (r *productRepository) FindOneProductByID(productID string) (entity.Product, error) {
	var product entity.Product
	result := r.db.Preload("User").Where("id = ?", productID).Take(&product)
	if result.Error != nil {
		return product, result.Error
	}
	return product, nil
}

func (r *productRepository) UpdateProduct(product entity.Product) (entity.Product, error) {
	r.db.Save(&product)
	r.db.Preload("User").Find(&product)
	return product, nil
}

func (r *productRepository) DeleteProduct(productID string) error {
	var product entity.Product
	result := r.db.Preload("User").Where("id = ?", productID).Take(&product)
	if result.Error != nil {
		return result.Error
	}
	r.db.Delete(&product)
	return nil
}
