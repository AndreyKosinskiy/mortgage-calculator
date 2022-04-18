package bankrepository

import (
	"context"
	"database/sql"
	"log"

	"github.com/AndreyKosinskiy/mortgage-calculator/internal/models"
)

type BankRepository struct {
	db     *sql.DB
	logger log.Logger
}

func New(db *sql.DB, logger log.Logger) *BankRepository {
	return &BankRepository{db, logger}
}

func (r *BankRepository) Create(ctx context.Context, b *models.Bank) (*models.Bank, error) {
	rows, err := r.db.QueryContext(ctx, create, b.Name, b.Rate, b.MaxLoan, b.MinDownPayment, b.LoanTerm)
	if err != nil {
		r.logger.Printf("bankRepository.Create.QueryContext : %s", err)
		return nil, err
	}
	defer rows.Close()

	nb := &models.Bank{}
	for rows.Next() {
		err := rows.Scan(&nb.Name, &nb.Name, &nb.Rate, &nb.MaxLoan, &nb.MinDownPayment, &nb.LoanTerm)
		if err != nil {
			r.logger.Printf("bankRepository.Create.Scan: %s", err)
			return nil, err
		}
	}

	return nb, nil
}
func (r *BankRepository) Update(ctx context.Context, b *models.Bank) (*models.Bank, error) {
	return nil, nil
}
func (r *BankRepository) Delete(ctx context.Context, b *models.Bank) error {
	return nil
}
func (r *BankRepository) BankByName(ctx context.Context, n string) (*models.Bank, error) {
	rows, err := r.db.QueryContext(ctx, byName, n)
	if err != nil {
		r.logger.Printf("bankRepository.BankByName.QueryContext : %s", err)
		return nil, err
	}
	defer rows.Close()

	b := &models.Bank{}
	for rows.Next() {
		err := rows.Scan(&b.Name, &b.Name, &b.Rate, &b.MaxLoan, &b.MinDownPayment, &b.LoanTerm)
		if err != nil {
			r.logger.Printf("bankRepository.BankByName.Scan: %s", err)
			return nil, err
		}
	}

	return b, nil
}
func (r *BankRepository) BankList(ctx context.Context) ([]*models.Bank, error) {
	rows, err := r.db.QueryContext(ctx, list)
	if err != nil {
		r.logger.Printf("bankRepository.BankList.QueryContext : %s", err)
		return nil, err
	}
	defer rows.Close()

	var sb []*models.Bank
	b := &models.Bank{}
	var id int
	for rows.Next() {
		err := rows.Scan(&id, &b.Name, &b.Rate, &b.MaxLoan, &b.MinDownPayment, &b.LoanTerm)
		if err != nil {
			r.logger.Printf("bankRepository.BankList.Scan: %s", err)
			return nil, err
		}
		sb = append(sb, b)
	}

	return sb, nil
}
