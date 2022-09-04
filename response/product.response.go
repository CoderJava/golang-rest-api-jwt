package response

import "golang-rest-api-jwt/entity"

type ProductResponse struct {
	ID          int64        `json:"id"`
	ProductName string       `json:"product_name"`
	Price       uint64       `json:"price"`
	User        UserResponse `json:"user,omitempty"`
}

func NewProductResponse(product entity.Product) *ProductResponse {
	return &ProductResponse{
		ID:          product.ID,
		ProductName: product.Name,
		Price:       product.Price,
		User:        *NewUserResponse(product.User),
	}
}

func NewProductArrayResponse(product []entity.Product) *[]ProductResponse {
	listProductResponse := []ProductResponse{}
	for _, element := range product {
		productResponse := ProductResponse{
			ID:          element.ID,
			ProductName: element.Name,
			Price:       element.Price,
			User:        *NewUserResponse(element.User),
		}
		listProductResponse = append(listProductResponse, productResponse)
	}
	return &listProductResponse
}
