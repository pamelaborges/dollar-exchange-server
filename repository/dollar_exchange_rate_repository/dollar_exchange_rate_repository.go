package dollarexchangeraterepository

import (
	"database/sql"

	"github.com/pamelaborges/dollar-exchange-server/models/dollarexchangerate"

	_ "github.com/go-sql-driver/mysql"
)

type Repository struct {
	db sql.DB
}

func (db *Repository) Create(entity *dollarexchangerate.DollarExchangeRate) {

}
