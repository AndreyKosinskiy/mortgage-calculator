package testrepository

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/AndreyKosinskiy/mortgage-calculator/internal/models"
)

// bs - Bank slice
var bs = []*models.Bank{
	&models.Bank{
		Name:           "PrivatBank",
		Rate:           10,
		MaxLoan:        10000,
		MinDownPayment: 1000,
		LoanTerm:       6,
	},
	&models.Bank{
		Name:           "Monobank",
		Rate:           12,
		MaxLoan:        15000,
		MinDownPayment: 2000,
		LoanTerm:       12,
	},
	&models.Bank{
		Name:           "MyBank",
		Rate:           15,
		MaxLoan:        30000,
		MinDownPayment: 5000,
		LoanTerm:       24,
	},
	&models.Bank{
		Name:           "UkraineBank",
		Rate:           17.5,
		MaxLoan:        50000,
		MinDownPayment: 10000,
		LoanTerm:       32,
	},
}

type TestRepository struct {
	db     []*models.Bank
	logger *log.Logger
}

func New() *TestRepository {
	return &TestRepository{bs, log.New(os.Stdout, "test:", 0)}
}

func (r *TestRepository) Create(ctx context.Context, b *models.Bank) (*models.Bank, error) {
	r.db = append(r.db, b)
	fmt.Println(" r.db len: ", len(r.db))
	return b, nil
}
func (r *TestRepository) Update(ctx context.Context, b *models.Bank) (*models.Bank, error) {
	return &models.Bank{}, nil
}
func (r *TestRepository) Delete(ctx context.Context, name string) error {
	return nil
}

func (r *TestRepository) BankByName(ctx context.Context, name string) (*models.Bank, error) {
	var b *models.Bank
	for _, bi := range r.db {
		if bi.Name == name {
			b = bi
			break
		}
	}
	return b, nil
}

func (r *TestRepository) BankList(ctx context.Context) ([]*models.Bank, error) {
	return r.db, nil
}
