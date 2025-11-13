package service

import (
	"errors"

	"github.com/modsynth/e-commerce-api/internal/domain"
	"github.com/modsynth/e-commerce-api/internal/repository"
)

type ProductService interface {
	CreateProduct(req *domain.CreateProductRequest) (*domain.Product, error)
	GetProductByID(id uint) (*domain.Product, error)
	GetProductBySlug(slug string) (*domain.Product, error)
	UpdateProduct(id uint, req *domain.UpdateProductRequest) (*domain.Product, error)
	DeleteProduct(id uint) error
	ListProducts(query *domain.ProductListQuery) ([]*domain.Product, int64, error)
	CheckStock(productID uint, quantity int) (bool, error)
}

type productService struct {
	productRepo repository.ProductRepository
}

func NewProductService(productRepo repository.ProductRepository) ProductService {
	return &productService{
		productRepo: productRepo,
	}
}

func (s *productService) CreateProduct(req *domain.CreateProductRequest) (*domain.Product, error) {
	product := &domain.Product{
		CategoryID:     req.CategoryID,
		Name:           req.Name,
		Slug:           req.Slug,
		Description:    req.Description,
		Price:          req.Price,
		ComparePrice:   req.ComparePrice,
		CostPrice:      req.CostPrice,
		SKU:            req.SKU,
		Barcode:        req.Barcode,
		StockQuantity:  req.StockQuantity,
		TrackInventory: req.TrackInventory,
		Weight:         req.Weight,
		IsActive:       req.IsActive,
		Featured:       req.Featured,
	}

	if err := s.productRepo.Create(product); err != nil {
		return nil, errors.New("failed to create product")
	}

	return product, nil
}

func (s *productService) GetProductByID(id uint) (*domain.Product, error) {
	product, err := s.productRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (s *productService) GetProductBySlug(slug string) (*domain.Product, error) {
	product, err := s.productRepo.FindBySlug(slug)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (s *productService) UpdateProduct(id uint, req *domain.UpdateProductRequest) (*domain.Product, error) {
	// Find existing product
	product, err := s.productRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields
	if req.CategoryID != nil {
		product.CategoryID = req.CategoryID
	}
	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	if req.Price != nil {
		product.Price = *req.Price
	}
	if req.ComparePrice != nil {
		product.ComparePrice = req.ComparePrice
	}
	if req.CostPrice != nil {
		product.CostPrice = req.CostPrice
	}
	if req.SKU != "" {
		product.SKU = req.SKU
	}
	if req.Barcode != "" {
		product.Barcode = req.Barcode
	}
	if req.StockQuantity != nil {
		product.StockQuantity = *req.StockQuantity
	}
	if req.TrackInventory != nil {
		product.TrackInventory = *req.TrackInventory
	}
	if req.Weight != nil {
		product.Weight = req.Weight
	}
	if req.IsActive != nil {
		product.IsActive = *req.IsActive
	}
	if req.Featured != nil {
		product.Featured = *req.Featured
	}

	if err := s.productRepo.Update(product); err != nil {
		return nil, errors.New("failed to update product")
	}

	return product, nil
}

func (s *productService) DeleteProduct(id uint) error {
	// Check if product exists
	_, err := s.productRepo.FindByID(id)
	if err != nil {
		return err
	}

	return s.productRepo.Delete(id)
}

func (s *productService) ListProducts(query *domain.ProductListQuery) ([]*domain.Product, int64, error) {
	return s.productRepo.List(query)
}

func (s *productService) CheckStock(productID uint, quantity int) (bool, error) {
	product, err := s.productRepo.FindByID(productID)
	if err != nil {
		return false, err
	}

	if !product.TrackInventory {
		return true, nil
	}

	return product.StockQuantity >= quantity, nil
}
