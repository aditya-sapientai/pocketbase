package core

import (
	"errors"
	"testing"
	"time"

	"github.com/pocketbase/dbx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAppForMigrationsRunner struct {
	mock.Mock
}

func (m *MockAppForMigrationsRunner) DB() *dbx.DB {
	args := m.Called()
	return args.Get(0).(*dbx.DB)
}

func (m *MockAppForMigrationsRunner) AuxRunInTransaction(fn func(App) error) error {
	args := m.Called(fn)
	return args.Error(0)
}

func (m *MockAppForMigrationsRunner) RunInTransaction(fn func(App) error) error {
	args := m.Called(fn)
	return args.Error(0)
}

type MockDBX struct {
	mock.Mock
}

func (m *MockDBX) NewQuery(query string) *dbx.Query {
	args := m.Called(query)
	return args.Get(0).(*dbx.Query)
}

func (m *MockDBX) Select(cols ...string) *dbx.SelectQuery {
	args := m.Called(cols)
	return args.Get(0).(*dbx.SelectQuery)
}

func (m *MockDBX) Insert(table string, params dbx.Params) *dbx.Query {
	args := m.Called(table, params)
	return args.Get(0).(*dbx.Query)
}

func (m *MockDBX) Delete(table string, condition interface{}) *dbx.Query {
	args := m.Called(table, condition)
	return args.Get(0).(*dbx.Query)
}

type MockQuery struct {
	mock.Mock
}

func (m *MockQuery) Execute() (int64, error) {
	args := m.Called()
	return int64(args.Int(0)), args.Error(1)
}

func (m *MockQuery) Row(a interface{}) error {
	args := m.Called(a)
	return args.Error(0)
}

func (m *MockQuery) From(tables ...string) *dbx.SelectQuery {
	args := m.Called(tables)
	return args.Get(0).(*dbx.SelectQuery)
}

func (m *MockQuery) Where(condition interface{}) *dbx.SelectQuery {
	args := m.Called(condition)
	return args.Get(0).(*dbx.SelectQuery)
}

func (m *MockQuery) AndWhere(condition interface{}) *dbx.SelectQuery {
	args := m.Called(condition)
	return args.Get(0).(*dbx.SelectQuery)
}

func (m *MockQuery) OrderBy(expression string) *dbx.SelectQuery {
	args := m.Called(expression)
	return args.Get(0).(*dbx.SelectQuery)
}

func (m *MockQuery) AndOrderBy(expression string) *dbx.SelectQuery {
	args := m.Called(expression)
	return args.Get(0).(*dbx.SelectQuery)
}

func (m *MockQuery) Limit(limit int64) *dbx.SelectQuery {
	args := m.Called(limit)
	return args.Get(0).(*dbx.SelectQuery)
}

func (m *MockQuery) Column(a interface{}) error {
	args := m.Called(a)
	return args.Error(0)
}

func TestNewMigrationsRunner(t *testing.T) {
	mockApp := new(MockAppForMigrationsRunner)
	migrationsList := MigrationsList{}

	runner := NewMigrationsRunner(mockApp, migrationsList)

	assert.NotNil(t, runner)
	assert.Equal(t, mockApp, runner.app)
	assert.Equal(t, migrationsList, runner.migrationsList)
	assert.Equal(t, DefaultMigrationsTable, runner.tableName)
}

func TestMigrationsRunner_Run(t *testing.T) {
	mockApp := new(MockAppForMigrationsRunner)
	mockDB := new(MockDBX)
	mockQuery := new(MockQuery)

	mockApp.On("DB").Return(mockDB)
	mockDB.On("NewQuery", mock.Anything).Return(mockQuery)
	mockQuery.On("Execute").Return(0, nil)

	runner := NewMigrationsRunner(mockApp, MigrationsList{})

	t.Run("Up command", func(t *testing.T) {
		mockApp.On("AuxRunInTransaction", mock.AnythingOfType("func(App) error")).Return(nil).Once()
		err := runner.Run("up")
		assert.NoError(t, err)
	})

	t.Run("Down command", func(t *testing.T) {
		mockDB.On("Select", mock.Anything).Return(mockQuery)
		mockQuery.On("From", mock.Anything).Return(mockQuery)
		mockQuery.On("Where", mock.Anything).Return(mockQuery)
		mockQuery.On("AndWhere", mock.Anything).Return(mockQuery)
		mockQuery.On("OrderBy", mock.Anything).Return(mockQuery)
		mockQuery.On("AndOrderBy", mock.Anything).Return(mockQuery)
		mockQuery.On("Limit", mock.Anything).Return(mockQuery)
		mockQuery.On("Column", mock.Anything).Return(nil)

		err := runner.Run("down")
		assert.NoError(t, err)
	})

	t.Run("History-sync command", func(t *testing.T) {
		mockDB.On("Delete", mock.Anything, mock.Anything).Return(mockQuery)
		mockQuery.On("Execute").Return(int64(0), nil)

		err := runner.Run("history-sync")
		assert.NoError(t, err)
	})

	t.Run("Unsupported command", func(t *testing.T) {
		err := runner.Run("unsupported")
		assert.Error(t, err)
	})
}

func TestMigrationsRunner_Up(t *testing.T) {
	mockApp := new(MockAppForMigrationsRunner)
	mockDB := new(MockDBX)
	mockQuery := new(MockQuery)

	mockApp.On("DB").Return(mockDB)
	mockDB.On("NewQuery", mock.Anything).Return(mockQuery)
	mockQuery.On("Execute").Return(0, nil)

	runner := NewMigrationsRunner(mockApp, MigrationsList{})

	t.Run("No migrations to apply", func(t *testing.T) {
		mockApp.On("AuxRunInTransaction", mock.AnythingOfType("func(App) error")).Return(nil).Once()
		applied, err := runner.Up()
		assert.NoError(t, err)
		assert.Empty(t, applied)
	})

	t.Run("Apply migrations", func(t *testing.T) {
		migrations := MigrationsList{
			Items: func() []*Migration {
				return []*Migration{
					{File: "migration1.go", Up: func(app App) error { return nil }},
					{File: "migration2.go", Up: func(app App) error { return nil }},
				}
			},
		}
		runner.migrationsList = migrations

		mockApp.On("AuxRunInTransaction", mock.AnythingOfType("func(App) error")).Return(nil).Once()
		mockApp.On("RunInTransaction", mock.AnythingOfType("func(App) error")).Return(nil).Once()
		mockDB.On("Select", mock.Anything).Return(mockQuery)
		mockQuery.On("From", mock.Anything).Return(mockQuery)
		mockQuery.On("Where", mock.Anything).Return(mockQuery)
		mockQuery.On("Limit", mock.Anything).Return(mockQuery)
		mockQuery.On("Row", mock.Anything).Return(nil)
		mockDB.On("Insert", mock.Anything, mock.Anything).Return(mockQuery)

		applied, err := runner.Up()
		assert.NoError(t, err)
		assert.Equal(t, []string{"migration1.go", "migration2.go"}, applied)
	})
}

func TestMigrationsRunner_Down(t *testing.T) {
	mockApp := new(MockAppForMigrationsRunner)
	mockDB := new(MockDBX)
	mockQuery := new(MockQuery)

	mockApp.On("DB").Return(mockDB)
	mockDB.On("NewQuery", mock.Anything).Return(mockQuery)
	mockQuery.On("Execute").Return(0, nil)

	runner := NewMigrationsRunner(mockApp, MigrationsList{})

	t.Run("No migrations to revert", func(t *testing.T) {
		mockDB.On("Select", mock.Anything).Return(mockQuery)
		mockQuery.On("From", mock.Anything).Return(mockQuery)
		mockQuery.On("Where", mock.Anything).Return(mockQuery)
		mockQuery.On("AndWhere", mock.Anything).Return(mockQuery)
		mockQuery.On("OrderBy", mock.Anything).Return(mockQuery)
		mockQuery.On("AndOrderBy", mock.Anything).Return(mockQuery)
		mockQuery.On("Limit", mock.Anything).Return(mockQuery)
		mockQuery.On("Column", mock.Anything).Return(nil)

		reverted, err := runner.Down(1)
		assert.NoError(t, err)
		assert.Empty(t, reverted)
	})

	t.Run("Revert migrations", func(t *testing.T) {
		migrations := MigrationsList{
			Items: func() []*Migration {
				return []*Migration{
					{File: "migration1.go", Down: func(app App) error { return nil }},
					{File: "migration2.go", Down: func(app App) error { return nil }},
				}
			},
		}
		runner.migrationsList = migrations

		mockDB.On("Select", mock.Anything).Return(mockQuery)
		mockQuery.On("From", mock.Anything).Return(mockQuery)
		mockQuery.On("Where", mock.Anything).Return(mockQuery)
		mockQuery.On("AndWhere", mock.Anything).Return(mockQuery)
		mockQuery.On("OrderBy", mock.Anything).Return(mockQuery)
		mockQuery.On("AndOrderBy", mock.Anything).Return(mockQuery)
		mockQuery.On("Limit", mock.Anything).Return(mockQuery)
		mockQuery.On("Column", mock.Anything).Run(func(args mock.Arguments) {
			files := args.Get(0).(*[]string)
			*files = []string{"migration2.go", "migration1.go"}
		}).Return(nil)

		mockApp.On("AuxRunInTransaction", mock.AnythingOfType("func(App) error")).Return(nil).Once()
		mockApp.On("RunInTransaction", mock.AnythingOfType("func(App) error")).Return(nil).Once()
		mockDB.On("Delete", mock.Anything, mock.Anything).Return(mockQuery)

		reverted, err := runner.Down(2)
		assert.NoError(t, err)
		assert.Equal(t, []string{"migration2.go", "migration1.go"}, reverted)
	})
}

func TestMigrationsRunner_RemoveMissingAppliedMigrations(t *testing.T) {
	mockApp := new(MockAppForMigrationsRunner)
	mockDB := new(MockDBX)
	mockQuery := new(MockQuery)

	mockApp.On("DB").Return(mockDB)
	mockDB.On("Delete", mock.Anything, mock.Anything).Return(mockQuery)
	mockQuery.On("Execute").Return(int64(0), nil)

	runner := NewMigrationsRunner(mockApp, MigrationsList{})

	err := runner.RemoveMissingAppliedMigrations()
	assert.NoError(t, err)
}

func TestMigrationsRunner_initMigrationsTable(t *testing.T) {
	mockApp := new(MockAppForMigrationsRunner)
	mockDB := new(MockDBX)
	mockQuery := new(MockQuery)

	mockApp.On("DB").Return(mockDB)
	mockDB.On("NewQuery", mock.Anything).Return(mockQuery)
	mockQuery.On("Execute").Return(int64(0), nil)

	runner := NewMigrationsRunner(mockApp, MigrationsList{})

	err := runner.initMigrationsTable()
	assert.NoError(t, err)
	assert.True(t, runner.inited)
}

func TestMigrationsRunner_isMigrationApplied(t *testing.T) {
	mockApp := new(MockAppForMigrationsRunner)
	mockDB := new(MockDBX)
	mockQuery := new(MockQuery)

	mockApp.On("DB").Return(mockDB)
	mockDB.On("Select", mock.Anything).Return(mockQuery)
	mockQuery.On("From", mock.Anything).Return(mockQuery)
	mockQuery.On("Where", mock.Anything).Return(mockQuery)
	mockQuery.On("Limit", mock.Anything).Return(mockQuery)
	mockQuery.On("Row", mock.Anything).Return(nil)

	runner := NewMigrationsRunner(mockApp, MigrationsList{})

	applied := runner.isMigrationApplied(mockApp, "test_migration.go")
	assert.True(t, applied)
}

func TestMigrationsRunner_saveAppliedMigration(t *testing.T) {
	mockApp := new(MockAppForMigrationsRunner)
	mockDB := new(MockDBX)
	mockQuery := new(MockQuery)

	mockApp.On("DB").Return(mockDB)
	mockDB.On("Insert", mock.Anything, mock.Anything).Return(mockQuery)
	mockQuery.On("Execute").Return(int64(1), nil)

	runner := NewMigrationsRunner(mockApp, MigrationsList{})

	err := runner.saveAppliedMigration(mockApp, "test_migration.go")
	assert.NoError(t, err)
}

func TestMigrationsRunner_saveRevertedMigration(t *testing.T) {
	mockApp := new(MockAppForMigrationsRunner)
	mockDB := new(MockDBX)
	mockQuery := new(MockQuery)

	mockApp.On("DB").Return(mockDB)
	mockDB.On("Delete", mock.Anything, mock.Anything).Return(mockQuery)
	mockQuery.On("Execute").Return(int64(1), nil)

	runner := NewMigrationsRunner(mockApp, MigrationsList{})

	err := runner.saveRevertedMigration(mockApp, "test_migration.go")
	assert.NoError(t, err)
}

func TestMigrationsRunner_lastAppliedMigrations(t *testing.T) {
	mockApp := new(MockAppForMigrationsRunner)
	mockDB := new(MockDBX)
	mockQuery := new(MockQuery)

	mockApp.On("DB").Return(mockDB)
	mockDB.On("Select", mock.Anything).Return(mockQuery)
	mockQuery.On("From", mock.Anything).Return(mockQuery)
	mockQuery.On("Where", mock.Anything).Return(mockQuery)
	mockQuery.On("AndWhere", mock.Anything).Return(mockQuery)
	mockQuery.On("OrderBy", mock.Anything).Return(mockQuery)
	mockQuery.On("AndOrderBy", mock.Anything).Return(mockQuery)
	mockQuery.On("Limit", mock.Anything).Return(mockQuery)
	mockQuery.On("Column", mock.Anything).Run(func(args mock.Arguments) {
		files := args.Get(0).(*[]string)
		*files = []string{"migration2.go", "migration1.go"}
	}).Return(nil)

	runner := NewMigrationsRunner(mockApp, MigrationsList{})

	files, err := runner.lastAppliedMigrations(2)
	assert.NoError(t, err)
	assert.Equal(t, []string{"migration2.go", "migration1.go"}, files)
}
