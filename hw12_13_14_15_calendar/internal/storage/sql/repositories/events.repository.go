package repositories

import (
	"github.com/jmoiron/sqlx"

	domain "github.com/ToggyO/otus-golang-for-pro/hw12_13_14_15_calendar/internal/domain/repositories"
)

type eventsRepository struct {
	connection *sqlx.DB
}

func NewEventsRepository(connection interface{}) domain.IEventsRepository {
	sqlxConn, ok := connection.(*sqlx.DB)
	if !ok {
		panic("cannot cast object to type `sqlx.DB`")
	}
	return &eventsRepository{connection: sqlxConn}
}
