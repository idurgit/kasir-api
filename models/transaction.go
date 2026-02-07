package models

type Transaction struct {
	ID          int                 `json:"id"`
	TotalAmount int                 `json:"total_amount"`
	Details     []TransactionDetail `json:"details"`
}

type TransactionDetail struct {
	ID            int    `json:"id"`
	TransactionID int    `json:"transaction_id"`
	ProductID     int    `json:"product_id"`
	ProductName   string `json:"product_name"`
	Quantity      int    `json:"quantity"`
	Subtotal      int    `json:"subtotal"`
}

type CheckoutRequest struct {
	Items []CheckoutItem `json:"items"`
}

type CheckoutItem struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

type DailySalesSummary struct {
	TotalTransaction   int                `json:"total_transaction"`
	TotalRevenue       int                `json:"total_revenue"`
	MostSellingProduct MostSellingProduct `json:"mostselling_product"`
}

type MostSellingProduct struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
}
