package core_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAppForMigrationsRunner is a mock of the App interface for testing MigrationsRunner
type MockAppForMigrationsRunner struct {
	mock.Mock
}

func (m *MockAppForMigrationsRunner) DB() *dbx.DB {
	args := m.Called()
	return args.Get(0).(*dbx.DB)
}

func (m *MockAppForMigrationsRunner) AuxRunInTransaction(fn func(app core.App) error) error {
	args := m.Called(fn)
	return args.Error(0)
}

func (m *MockAppForMigrationsRunner) RunInTransaction(fn func(app core.App) error) error {
	args := m.Called(fn)
	return args.Error(0)
}

func TestNewMigrationsRunner(t *testing.T) {
	mockApp := new(MockAppForMigrationsRunner)
	migrationsList := core.MigrationsList{}

	runner := core.NewMigrationsRunner(mockApp, migrationsList)

	assert.NotNil(t, runner)
	assert.Equal(t, mockApp, runner.app)
	assert.Equal(t, migrationsList, runner.migrationsList)
	assert.Equal(t, core.DefaultMigrationsTable, runner.tableName)
}

func TestMigrationsRunner_Run(t *testing.T) {
	mockApp := new(MockAppForMigrationsRunner)
	migrationsList := core.MigrationsList{}
	runner := core.NewMigrationsRunner(mockApp, migrationsList)

	mockDB := &dbx.DB{}
	mockApp.On("DB").Return(mockDB)

	mockQuery := &dbx.Query{}
	mockDB.On("NewQuery", mock.Anything).Return(mockQuery)
	mockQuery.On("Execute").Return(nil, nil)

	t.Run("Up command", func(t *testing.T) {
		err := runner.Run("up")
		assert.NoError(t, err)
	})

	t.Run("Down command", func(t *testing.T) {
		mockApp.On("DB").Return(mockDB)
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
		mockApp.On("DB").Return(mockDB)
		mockDB.On("Delete", mock.Anything, mock.Anything).Return(mockQuery)
		mockQuery.On("Execute").Return(nil, nil)

		err := runner.Run("history-sync")
		assert.NoError(t, err)
	})

	t.Run("Invalid command", func(t *testing.T) {
		err := runner.Run("invalid")
		assert.Error(t, err)
	})
}

func TestMigrationsRunner_Up(t *testing.T) {
	mockApp := new(MockAppForMigrationsRunner)
	migrationsList := core.MigrationsList{}
	migrationsList.Register(core.Migration{
		File: "test_migration",
		Up: func(app core.App) error {
			return nil
		},
	})
	runner := core.NewMigrationsRunner(mockApp, migrationsList)

	mockDB := &dbx.DB{}
	mockApp.On("DB").Return(mockDB)

	mockQuery := &dbx.Query{}
	mockDB.On("NewQuery", mock.Anything).Return(mockQuery)
	mockQuery.On("Execute").Return(nil, nil)

	mockApp.On("AuxRunInTransaction", mock.AnythingOfType("func(core.App) error")).Return(nil)
	mockApp.On("RunInTransaction", mock.AnythingOfType("func(core.App) error")).Return(nil)

	applied, err := runner.Up()

	assert.NoError(t, err)
	assert.Equal(t, []string{"test_migration"}, applied)
}

func TestMigrationsRunner_Down(t *testing.T) {
	mockApp := new(MockAppForMigrationsRunner)
	migrationsList := core.MigrationsList{}
	migrationsList.Register(core.Migration{
		File: "test_migration",
		Down: func(app core.App) error {
			return nil
		},
	})
	runner := core.NewMigrationsRunner(mockApp, migrationsList)

	mockDB := &dbx.DB{}
	mockApp.On("DB").Return(mockDB)

	mockQuery := &dbx.Query{}
	mockDB.On("NewQuery", mock.Anything).Return(mockQuery)
	mockQuery.On("Execute").Return(nil, nil)

	mockDB.On("Select", mock.Anything).Return(mockQuery)
	mockQuery.On("From", mock.Anything).Return(mockQuery)
	mockQuery.On("Where", mock.Anything).Return(mockQuery)
	mockQuery.On("AndWhere", mock.Anything).Return(mockQuery)
	mockQuery.On("OrderBy", mock.Anything).Return(mockQuery)
	mockQuery.On("AndOrderBy", mock.Anything).Return(mockQuery)
	mockQuery.On("Limit", mock.Anything).Return(mockQuery)
	mockQuery.On("Column", mock.Anything).Run(func(args mock.Arguments) {
		files := args.Get(0).(*[]string)
		*files = append(*files, "test_migration")
	}).Return(nil)

	mockApp.On("AuxRunInTransaction", mock.AnythingOfType("func(core.App) error")).Return(nil)
	mockApp.On("RunInTransaction", mock.AnythingOfType("func(core.App) error")).Return(nil)

	reverted, err := runner.Down(1)

	assert.NoError(t, err)
	assert.Equal(t, []string{"test_migration"}, reverted)
}

func TestMigrationsRunner_RemoveMissingAppliedMigrations(t *testing.T) {
	mockApp := new(MockAppForMigrationsRunner)
	migrationsList := core.MigrationsList{}
	runner := core.NewMigrationsRunner(mockApp, migrationsList)

	mockDB := &dbx.DB{}
	mockApp.On("DB").Return(mockDB)

	mockQuery := &dbx.Query{}
	mockDB.On("Delete", mock.Anything, mock.Anything).Return(mockQuery)
	mockQuery.On("Execute").Return(nil, nil)

	err := runner.RemoveMissingAppliedMigrations()

	assert.NoError(t, err)
}

func TestMigrationsRunner_initMigrationsTable(t *testing.T) {
	mockApp := new(MockAppForMigrationsRunner)
	migrationsList := core.MigrationsList{}
	runner := core.NewMigrationsRunner(mockApp, migrationsList)

	mockDB := &dbx.DB{}
	mockApp.On("DB").Return(mockDB)

	mockQuery := &dbx.Query{}
	mockDB.On("NewQuery", mock.Anything).Return(mockQuery)
	mockQuery.On("Execute").Return(nil, nil)

	err := runner.initMigrationsTable()

	assert.NoError(t, err)
	assert.True(t, runner.inited)
}

func TestMigrationsRunner_isMigrationApplied(t *testing.T) {
	mockApp := new(MockAppForMigrationsRunner)
	migrationsList := core.MigrationsList{}
	runner := core.NewMigrationsRunner(mockApp, migrationsList)

	mockDB := &dbx.DB{}
	mockApp.On("DB").Return(mockDB)

	mockQuery := &dbx.Query{}
	mockDB.On("Select", mock.Anything).Return(mockQuery)
	mockQuery.On("From", mock.Anything).Return(mockQuery)
	mockQuery.On("Where", mock.Anything).Return(mockQuery)
	mockQuery.On("Limit", mock.Anything).Return(mockQuery)
	mockQuery.On("Row", mock.Anything).Run(func(args mock.Arguments) {
		exists := args.Get(0).(*bool)
		*exists = true
	}).Return(nil)

	applied := runner.isMigrationApplied(mockApp, "test_migration")

	assert.True(t, applied)
}

func TestMigrationsRunner_saveAppliedMigration(t *testing.T) {
	mockApp := new(MockAppForMigrationsRunner)
	migrationsList := core.MigrationsList{}
	runner := core.NewMigrationsRunner(mockApp, migrationsList)

	mockDB := &dbx.DB{}
	mockApp.On("DB").Return(mockDB)

	mockQuery := &dbx.Query{}
	mockDB.On("Insert", mock.Anything, mock.Anything).Return(mockQuery)
	mockQuery.On("Execute").Return(nil, nil)

	err := runner.saveAppliedMigration(mockApp, "test_migration")

	assert.NoError(t, err)
}

func TestMigrationsRunner_saveRevertedMigration(t *testing.T) {
	mockApp := new(MockAppForMigrationsRunner)
	migrationsList := core.MigrationsList{}
	runner := core.NewMigrationsRunner(mockApp, migrationsList)

	mockDB := &dbx.DB{}
	mockApp.On("DB").Return(mockDB)

	mockQuery := &dbx.Query{}
	mockDB.On("Delete", mock.Anything, mock.Anything).Return(mockQuery)
	mockQuery.On("Execute").Return(nil, nil)

	err := runner.saveRevertedMigration(mockApp, "test_migration")

	assert.NoError(t, err)
}

func TestMigrationsRunner_lastAppliedMigrations(t *testing.T) {
	mockApp := new(MockAppForMigrationsRunner)
	migrationsList := core.MigrationsList{}
	migrationsList.Register(core.Migration{File: "test_migration"})
	runner := core.NewMigrationsRunner(mockApp, migrationsList)

	mockDB := &dbx.DB{}
	mockApp.On("DB").Return(mockDB)

	mockQuery := &dbx.Query{}
	mockDB.On("Select", mock.Anything).Return(mockQuery)
	mockQuery.On("From", mock.Anything).Return(mockQuery)
	mockQuery.On("Where", mock.Anything).Return(mockQuery)
	mockQuery.On("AndWhere", mock.Anything).Return(mockQuery)
	mockQuery.On("OrderBy", mock.Anything).Return(mockQuery)
	mockQuery.On("AndOrderBy", mock.Anything).Return(mockQuery)
	mockQuery.On("Limit", mock.Anything).Return(mockQuery)
	mockQuery.On("Column", mock.Anything).Run(func(args mock.Arguments) {
		files := args.Get(0).(*[]string)
		*files = append(*files, "test_migration")
	}).Return(nil)

	files, err := runner.lastAppliedMigrations(1)

	assert.NoError(t, err)
	assert.Equal(t, []string{"test_migration"}, files)
}
