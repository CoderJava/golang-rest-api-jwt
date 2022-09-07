package service

import (
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
