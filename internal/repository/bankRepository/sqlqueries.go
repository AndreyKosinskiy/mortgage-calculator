package bankrepository

const create = `INSERT INTO banks(name, rate, max_loan, min_down_payment, loan_term) 
				VALUES ($1,$2,$3,$4,$5) 
				RETURNING *;`
const update = `UPDATE banks 
				SET rate=$2, max_loan=$3,  min_down_payment=$4, loan_term=$5 
				WHERE name=$1
				RETURNING *;`
const delete = `DELETE FROM banks 
				WHERE name=$1 
				RETURNING *;`
const list = `SELECT * FROM banks;`
const byName = `SELECT * FROM banks
				WHERE name=$1;`
