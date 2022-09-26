package suites

import (
	"context"
	sqlstorage "github.com/ToggyO/otus-golang-for-pro/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/ToggyO/otus-golang-for-pro/hw12_13_14_15_calendar/pkg/migrations"
	"github.com/stretchr/testify/suite"
	"log"

	"github.com/ToggyO/otus-golang-for-pro/hw12_13_14_15_calendar/pkg/configuration"
	"github.com/ToggyO/otus-golang-for-pro/hw12_13_14_15_calendar/pkg/shared"
)

type ServiceFixtureSuite struct {
	suite.Suite
	repositories

	conf            configuration.Configuration
	client          shared.IDbClient
	migrationRunner migrations.IMigrationRunner

	ctx    context.Context
	cancel context.CancelFunc

	isInitiated bool
}

func (sf *ServiceFixtureSuite) Init() {
	var err error
	sf.ctx, sf.cancel = context.WithCancel(context.Background())

	sf.conf = configuration.NewConfiguration("../configs/config.stage.toml")

	client := sqlstorage.NewDbClient(sf.conf.Storage)
	if err = client.Connect(sf.ctx); err != nil {
		log.Fatal(err)
	}

	sf.client = client
	sf.migrationRunner = migrations.NewMigrationRunner(sf.conf.Storage.Dialect, shared.CreatePgConnectionString(sf.conf.Storage))
	if err = sf.migrationRunner.MigrateUp(sf.ctx); err != nil {
		log.Fatal(err)
	}

	sf.setupRepositories()

	sf.isInitiated = true
}

func (sf *ServiceFixtureSuite) TearDownSuite() {
	var err error
	if err = sf.migrationRunner.MigrateDown(sf.ctx); err != nil {
		log.Fatal(err)
	}

	if err = sf.client.Close(sf.ctx); err != nil {
		log.Fatal(err)
	}
}

func (sf *ServiceFixtureSuite) setupRepositories() {
	sf.repositories = repositories{
		eventsRepository: sqlstorage.NewEventsRepository(sf.client),
	}
}
