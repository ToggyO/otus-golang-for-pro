package sqlstorage

import (
	"context"
	"fmt"
	"github.com/ToggyO/otus-golang-for-pro/hw12_13_14_15_calendar/pkg/shared"
	"reflect"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"

	"github.com/ToggyO/otus-golang-for-pro/hw12_13_14_15_calendar/internal/domain/models"
	domain "github.com/ToggyO/otus-golang-for-pro/hw12_13_14_15_calendar/internal/domain/repositories"
	"github.com/ToggyO/otus-golang-for-pro/hw12_13_14_15_calendar/internal/storage/dbconverter"
	"github.com/ToggyO/otus-golang-for-pro/hw12_13_14_15_calendar/internal/storage/sql/entities"
	"github.com/ToggyO/otus-golang-for-pro/hw12_13_14_15_calendar/internal/storage/sql/table"
)

type eventsRepository struct {
	connection           *sqlx.DB
	primaryKeyColumnName string
	commonColumns        []string
}

func NewEventsRepository(client shared.IDbClient) domain.IEventsRepository {
	connection, ok := client.GetConnection().(*sqlx.DB)
	if !ok {
		panic(fmt.Errorf("invalid connection type provided< must be a %s", reflect.TypeOf(sqlx.DB{}).Name()))
	}
	return &eventsRepository{
		connection:           connection,
		primaryKeyColumnName: "id",
		commonColumns:        []string{"title", "start_date", "end_date", "description", "owner_id", "notification_date"},
	}
}

func (e *eventsRepository) CreateEvent(ctx context.Context, eventInfo *models.EventInfo) (*models.Event, error) {
	dbEventInfo, err := dbconverter.FromEventInfo(eventInfo)
	if err != nil {
		return nil, err
	}

	sql, args, err := squirrel.Insert(table.EventTableName).
		Columns(e.commonColumns...).
		PlaceholderFormat(squirrel.Dollar).
		Values(
			dbEventInfo.Title,
			dbEventInfo.StartDate,
			dbEventInfo.EndDate,
			dbEventInfo.Description,
			dbEventInfo.OwnerId,
			dbEventInfo.NotificationDate,
		).ToSql()

	if err != nil {
		return nil, err
	}

	var lastInsertedId int64
	row := e.connection.QueryRowContext(ctx, e.setReturningLastInsertedId(sql, "id"), args...)
	err = row.Scan(&lastInsertedId)
	if err != nil {
		return nil, err
	}

	event := &models.Event{
		ID:        lastInsertedId,
		EventInfo: eventInfo,
	}

	return event, nil
}

func (e *eventsRepository) UpdateEvent(ctx context.Context, id int64, eventInfo *models.EventInfo) (*models.Event, error) {
	dbEventInfo, err := dbconverter.FromEventInfo(eventInfo)
	if err != nil {
		return nil, err
	}

	builder := squirrel.Update(table.EventTableName).Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		Set("updated_at", time.Now()).
		Set("title", dbEventInfo.Title).
		Set("start_date", dbEventInfo.StartDate).
		Set("end_date", dbEventInfo.EndDate).
		Set("description", dbEventInfo.Description).
		Set("owner_id", dbEventInfo.OwnerId).
		Set("notification_date", dbEventInfo.NotificationDate)

	sql, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	_, err = e.connection.ExecContext(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	return e.GetEventById(ctx, id)
}

func (e *eventsRepository) DeleteEvent(ctx context.Context, id int64) error {
	sql, args, err := squirrel.Delete(table.EventTableName).
		PlaceholderFormat(squirrel.Dollar).
		Where(squirrel.Eq{"id": id}).
		ToSql()

	if err != nil {
		return err
	}

	_, err = e.connection.ExecContext(ctx, sql, args...)
	return err
}

func (e *eventsRepository) GetEventsList(ctx context.Context, filter *models.EventsFilter) ([]models.Event, error) {
	builder := squirrel.Select(e.appendPrimaryKeyToColumns()...).
		PlaceholderFormat(squirrel.Dollar).
		From(table.EventTableName).
		Offset(filter.Page).
		Limit(filter.PageSize)

	isValidStartDate := filter.StartDate.IsZero()
	isValidEndDate := filter.EndDate.IsZero()
	if isValidStartDate && !isValidEndDate {
		builder.Where(squirrel.GtOrEq{"start_date": filter.StartDate})
	}

	if !isValidStartDate && isValidEndDate {
		builder.Where(squirrel.LtOrEq{"end_date": filter.EndDate})
	}

	if isValidStartDate && isValidEndDate {
		builder.Where(squirrel.And{
			squirrel.GtOrEq{"start_date": filter.StartDate},
			squirrel.LtOrEq{"end_date": filter.EndDate},
		})
	}

	sql, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	var eventDbModels []*entities.EventDbModel
	err = e.connection.SelectContext(ctx, &eventDbModels, sql, args...)
	if err != nil {
		return nil, err
	}

	events := make([]models.Event, 0, len(eventDbModels))
	for _, dbEv := range eventDbModels {
		ev, err := dbconverter.ToEvent(dbEv)
		if err != nil {
			return events[0:0], err
		}
		events = append(events, *ev)
	}

	return events, nil
}

func (e *eventsRepository) GetEventById(ctx context.Context, id int64) (*models.Event, error) {
	sql, args, err := squirrel.Select(e.appendPrimaryKeyToColumns()...).
		PlaceholderFormat(squirrel.Dollar).
		From(table.EventTableName).
		Where(squirrel.Eq{"id": id}).
		ToSql()

	if err != nil {
		return nil, err
	}

	eventDbModel := &entities.EventDbModel{}
	err = e.connection.GetContext(ctx, eventDbModel, sql, args...)
	if err != nil {
		return nil, err
	}

	return dbconverter.ToEvent(eventDbModel)
}

func (e *eventsRepository) setReturningLastInsertedId(sql, idColumnName string) string {
	return sql + fmt.Sprintf("RETURNING %s", idColumnName)
}

func (e *eventsRepository) appendPrimaryKeyToColumns() []string {
	return append(e.commonColumns, e.primaryKeyColumnName)
}
