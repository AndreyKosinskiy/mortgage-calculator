package bankrepository

import (
	"database/sql"
	"log"
)

type BankRepository struct {
	db     *sql.DB
	logger log.Logger
}

func New(db *sql.DB, logger log.Logger) *BankRepository {
	return &BankRepository{db, logger}
}

func (r *BankRepository) Create() {

}
func (r *BankRepository) Update() {

}
func (r *BankRepository) Delete() {

}
func (r *BankRepository) BankByName() {

}
