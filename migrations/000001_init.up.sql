CREATE TABLE IF NOT EXISTS accounts
(
    accountID         SERIAL PRIMARY KEY,
    balance    NUMERIC(10, 2) DEFAULT 0 CHECK (balance >= 0),
    currency   VARCHAR(16)    DEFAULT 'TMT' NOT NULL,
    is_locked  BOOLEAN        DEFAULT FALSE,
    created_at TIMESTAMP(0)   DEFAULT now()::TIMESTAMP(0),
    deleted_at TIMESTAMP(0)
);


CREATE TYPE transaction_type_enum AS ENUM ('deposit', 'credit');
CREATE TABLE IF NOT EXISTS transactions
(
    accountID               SERIAL PRIMARY KEY,
    account_id       INT REFERENCES accounts (accountID) NOT NULL,
    amount           NUMERIC(10, 2)               NOT NULL,
    transaction_type transaction_type_enum        NOT NULL,
    created_at       TIMESTAMP(0) DEFAULT now(),
    deleted_at       TIMESTAMP(0)
);
