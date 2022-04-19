CREATE TABLE IF NOT EXISTS banks(
    id serial PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL,
    rate FLOAT NOT NULL,
    max_loan FLOAT NOT NULL,
    min_down_payment FLOAT NOT NULL,
    loan_term INT NOT NULL
);