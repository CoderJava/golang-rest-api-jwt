package service

import (
	"errors"
	"fmt"
	"golang-rest-api-jwt/dto"
	"golang-rest-api-jwt/entity"
	"golang-rest-api-jwt/repository"
	"golang-rest-api-jwt/response"
	"strconv"
)

type ProductService interface {
	CreateProduct(productRequest dto.CreateProductRequest, userID string) (*response.ProductResponse, error)

	All(userID string) (*[]response.ProductResponse, error)

	FindOneProductByID(productID string) (*response.ProductResponse, error)

	UpdateProduct(updateProductRequest dto.UpdateProductRequest, userID string) (*response.ProductResponse, error)

	DeleteProduct(productID string, userID string) error
}

type productService struct {
	productRepository repository.ProductRepository
}

func NewProductService(productRepository repository.ProductRepository) ProductService {
	return &productService{
		productRepository: productRepository,
	}
}

func (s *productService) CreateProduct(
	productRequest dto.CreateProductRequest,
	userID string,
) (*response.ProductResponse, error) {
	product := entity.Product{
		Name:  productRequest.Name,
		Price: productRequest.Price,
	}

	id, err := strconv.Atoi(userID)
	if err != nil {
		return nil, err
	}
	product.UserID = int64(id)

	resultInsertProduct, err := s.productRepository.InsertProduct(product)
	if err != nil {
		return nil, err
	}

	response := response.NewProductResponse(resultInsertProduct)
	return response, nil
}

func (s *productService) All(userID string) (*[]response.ProductResponse, error) {
	allProductEntity, err := s.productRepository.All(userID)
	if err != nil {
		return nil, err
	}

	allProductResponse := response.NewProductArrayResponse(allProductEntity)
	return allProductResponse, nil
}

func (s *productService) FindOneProductByID(productID string) (*response.ProductResponse, error) {
	productEntity, err := s.productRepository.FindOneProductByID(productID)
	if err != nil {
		return nil, err
	}

	result := response.NewProductResponse(productEntity)
	return result, nil
}

func (s *productService) UpdateProduct(
	updateProductRequest dto.UpdateProductRequest,
	userID string,
) (*response.ProductResponse, error) {
	strProductID := strconv.Itoa(int(updateProductRequest.ID))
	productEntity, err := s.productRepository.FindOneProductByID(strProductID)
	if err != nil {
		return nil, err
	}

	intUserID, err := strconv.ParseInt(userID, 0, 64)
	if err != nil {
		return nil, err
	} else if intUserID != productEntity.UserID {
		return nil, errors.New("this product is not yours")
	}

	updateProductEntity := entity.Product{
		ID:     updateProductRequest.ID,
		Name:   updateProductRequest.Name,
		Price:  updateProductRequest.Price,
		UserID: intUserID,
	}
	updateProductEntity, err = s.productRepository.UpdateProduct(updateProductEntity)
	if err != nil {
		return nil, err
	}

	result := response.NewProductResponse(updateProductEntity)
	return result, nil
}

func (s *productService) DeleteProduct(productID string, userID string) error {
	productEntity, err := s.productRepository.FindOneProductByID(productID)
	if err != nil {
		return err
	}
	if fmt.Sprintf("%d", productEntity.UserID) != userID {
		return errors.New("this product is not yours")
	}

	s.productRepository.DeleteProduct(productID)
	return nil
}
