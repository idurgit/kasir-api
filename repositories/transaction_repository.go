package repositories

import (
	"database/sql"
	"fmt"
	"kasir-api/models"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (repo *TransactionRepository) CreateTransaction(items []models.CheckoutItem) (*models.Transaction, error) {
	var (
		res *models.Transaction
	)

	tx, err := repo.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// inisialisasi subtotal -> jumlah total transaksi keseluruhan
	totalAmount := 0
	// inisialisasi modeling transactionDetails -> nanti kita insert ke db
	details := make([]models.TransactionDetail, 0)
	// loop setiap item
	for _, item := range items {
		var productName string
		var productID, price, stock int
		// get product dapet pricing
		err := tx.QueryRow("SELECT id, name, price, stock FROM products WHERE id=$1", item.ProductID).Scan(&productID, &productName, &price, &stock)
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product id %d not found", item.ProductID)
		}

		if err != nil {
			return nil, err
		}

		// hitung current total = quantity * pricing
		// ditambahin ke dalam subtotal
		subtotal := item.Quantity * price
		totalAmount += subtotal

		// kurangi jumlah stok
		_, err = tx.Exec("UPDATE products SET stock = stock - $1 WHERE id = $2", item.Quantity, productID)
		if err != nil {
			return nil, err
		}

		// item nya dimasukkin ke transactionDetails
		details = append(details, models.TransactionDetail{
			ProductID:   productID,
			ProductName: productName,
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		})
	}

	// insert transaction
	var transactionID int
	err = tx.QueryRow("INSERT INTO transactions (total_amount) VALUES ($1) RETURNING ID", totalAmount).Scan(&transactionID)
	if err != nil {
		return nil, err
	}

	// insert transaction details
	for i, detail := range details {
		details[i].TransactionID = transactionID
		err := tx.QueryRow("INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal) VALUES ($1, $2, $3, $4) RETURNING id", transactionID, detail.ProductID, detail.Quantity, detail.Subtotal).Scan(&details[i].ID)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	res = &models.Transaction{
		ID:          transactionID,
		TotalAmount: totalAmount,
		Details:     details,
	}

	return res, nil
}

func (repo *TransactionRepository) GetTodaySalesSummary() (*models.DailySalesSummary, error) {
	// Get total revenue and transaction count for today
	var totalRevenue, totalTransaction int
	err := repo.db.QueryRow(`
		SELECT COALESCE(SUM(total_amount), 0), COUNT(*)
		FROM transactions
		WHERE DATE(created_at) = CURRENT_DATE
	`).Scan(&totalRevenue, &totalTransaction)

	if err != nil {
		return nil, err
	}

	// Get most selling product today
	var mostSelling models.MostSellingProduct
	var name sql.NullString
	var quantity sql.NullInt64

	err = repo.db.QueryRow(`
		SELECT p.name, COALESCE(SUM(td.quantity), 0) as total_qty
		FROM transaction_details td
		JOIN transactions t ON td.transaction_id = t.id
		JOIN products p ON td.product_id = p.id
		WHERE DATE(t.created_at) = CURRENT_DATE
		GROUP BY p.id, p.name
		ORDER BY total_qty DESC
		LIMIT 1
	`).Scan(&name, &quantity)

	if name.Valid {
		mostSelling.Name = name.String
	}
	if quantity.Valid {
		mostSelling.Quantity = int(quantity.Int64)
	}

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return &models.DailySalesSummary{
		TotalTransaction:   totalTransaction,
		TotalRevenue:       totalRevenue,
		MostSellingProduct: mostSelling,
	}, nil
}

func (repo *TransactionRepository) GetSalesSummaryByDateRange(startDate, endDate string) (*models.DailySalesSummary, error) {
	// Get total revenue and transaction count for date range
	var totalRevenue, totalTransaction int
	err := repo.db.QueryRow(`
		SELECT COALESCE(SUM(total_amount), 0), COUNT(*)
		FROM transactions
		WHERE DATE(created_at) BETWEEN $1 AND $2
	`, startDate, endDate).Scan(&totalRevenue, &totalTransaction)

	if err != nil {
		return nil, err
	}

	// Get most selling product in date range
	var mostSelling models.MostSellingProduct
	var name sql.NullString
	var quantity sql.NullInt64

	err = repo.db.QueryRow(`
		SELECT p.name, COALESCE(SUM(td.quantity), 0) as total_qty
		FROM transaction_details td
		JOIN transactions t ON td.transaction_id = t.id
		JOIN products p ON td.product_id = p.id
		WHERE DATE(t.created_at) BETWEEN $1 AND $2
		GROUP BY p.id, p.name
		ORDER BY total_qty DESC
		LIMIT 1
	`, startDate, endDate).Scan(&name, &quantity)

	if name.Valid {
		mostSelling.Name = name.String
	}
	if quantity.Valid {
		mostSelling.Quantity = int(quantity.Int64)
	}

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return &models.DailySalesSummary{
		TotalTransaction:   totalTransaction,
		TotalRevenue:       totalRevenue,
		MostSellingProduct: mostSelling,
	}, nil
}
