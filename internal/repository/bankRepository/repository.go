package bankrepository

import (
	"context"

	"github.com/AndreyKosinskiy/mortgage-calculator/internal/models"
)

type Repository interface {
	Create(ctx context.Context, b models.Bank)
	Update(ctx context.Context, b models.Bank)
	Delete(ctx context.Context, name string)
	BankByName(ctx context.Context, name string)
}
