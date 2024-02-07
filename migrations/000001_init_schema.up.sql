-- Создание таблицы кошельков
CREATE TABLE IF NOT EXISTS wallets
(
    id      SERIAL PRIMARY KEY,
    balance NUMERIC(10, 2) NOT NULL
);

-- Создание таблицы транзакций
CREATE TABLE IF NOT EXISTS transactions
(
    id             SERIAL PRIMARY KEY,
    from_wallet_id INT            NOT NULL REFERENCES wallets (id),
    to_wallet_id   INT            NOT NULL REFERENCES wallets (id),
    amount         NUMERIC(10, 2) NOT NULL,
    time           TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'utc')
);
