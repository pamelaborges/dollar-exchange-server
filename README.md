# dollar-exchange-server
Desafio Pós Graduação GoExpert


# Script de Banco de Dados

CREATE TABLE dollar_exchange_rate (
    id SERIAL PRIMARY KEY,
    code VARCHAR(10) NOT NULL,
    codein VARCHAR(10) NOT NULL,
    name VARCHAR(255) NOT NULL,
    high FLOAT NOT NULL,
    low FLOAT NOT NULL,
    varBid FLOAT NOT NULL,
    pctChange FLOAT NOT NULL,
    bid FLOAT NOT NULL,
    ask FLOAT NOT NULL
);