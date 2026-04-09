package models

import "time"

type ProductStatus string

const (
	ProductOnSale  ProductStatus = "on_sale"
	ProductOffSale ProductStatus = "off_sale"
	ProductSoldOut ProductStatus = "sold_out"
)

type Product struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	SKU         string        `json:"sku,omitempty"`
	Price       float64       `json:"price"`
	Status      ProductStatus `json:"status"`
	CategoryID  string        `json:"category_id,omitempty"`
	Category    string        `json:"category,omitempty"`
	Description string        `json:"description,omitempty"`
	ImageURL    string        `json:"image_url,omitempty"`
	CreatedAt   time.Time     `json:"created_at,omitempty"`
	UpdatedAt   time.Time     `json:"updated_at,omitempty"`
}

type Category struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	ParentID  string `json:"parent_id,omitempty"`
	SortOrder int    `json:"sort_order,omitempty"`
}
