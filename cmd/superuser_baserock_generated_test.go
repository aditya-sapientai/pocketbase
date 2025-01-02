package cmd_test

import (
	"errors"
	"testing"

	"github.com/pocketbase/pocketbase/cmd"
	"github.com/pocketbase/pocketbase/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockApp struct {
	mock.Mock
}

func (m *mockApp) FindCachedCollectionByNameOrId(nameOrId string) (*core.Collection, error) {
	args := m.Called(nameOrId)
	return args.Get(0).(*core.Collection), args.Error(1)
}

func (m *mockApp) FindAuthRecordByEmail(collection any, email string) (*core.Record, error) {
	args := m.Called(collection, email)
	return args.Get(0).(*core.Record), args.Error(1)
}

func (m *mockApp) Save(record *core.Record) error {
	args := m.Called(record)
	return args.Error(0)
}

func (m *mockApp) Delete(record *core.Record) error {
	args := m.Called(record)
	return args.Error(0)
}

func TestNewSuperuserCommand(t *testing.T) {
	app := &mockApp{}
	cmd := cmd.NewSuperuserCommand(app)

	assert.NotNil(t, cmd)
	assert.Equal(t, "superuser", cmd.Use)
	assert.Equal(t, "Manage superusers", cmd.Short)
	assert.Len(t, cmd.Commands(), 5)
}

func TestSuperuserUpsertCommand(t *testing.T) {
	app := &mockApp{}
	command := cmd.NewSuperuserCommand(app).Commands()[0]

	t.Run("MissingArguments", func(t *testing.T) {
		err := command.Execute()
		assert.Error(t, err)
		assert.Equal(t, "Missing email and password arguments.", err.Error())
	})

	t.Run("InvalidEmail", func(t *testing.T) {
		command.SetArgs([]string{"invalid-email", "password"})
		err := command.Execute()
		assert.Error(t, err)
		assert.Equal(t, "Missing or invalid email address.", err.Error())
	})

	t.Run("SuccessfulUpsert", func(t *testing.T) {
		email := "test@example.com"
		password := "password123"

		collection := &core.Collection{}
		app.On("FindCachedCollectionByNameOrId", core.CollectionNameSuperusers).Return(collection, nil)

		record := &core.Record{}
		app.On("FindAuthRecordByEmail", collection, email).Return(record, nil)
		app.On("Save", record).Return(nil)

		command.SetArgs([]string{email, password})
		err := command.Execute()

		assert.NoError(t, err)
		app.AssertExpectations(t)
	})
}

func TestSuperuserCreateCommand(t *testing.T) {
	app := &mockApp{}
	command := cmd.NewSuperuserCommand(app).Commands()[1]

	t.Run("MissingArguments", func(t *testing.T) {
		err := command.Execute()
		assert.Error(t, err)
		assert.Equal(t, "Missing email and password arguments.", err.Error())
	})

	t.Run("InvalidEmail", func(t *testing.T) {
		command.SetArgs([]string{"invalid-email", "password"})
		err := command.Execute()
		assert.Error(t, err)
		assert.Equal(t, "Missing or invalid email address.", err.Error())
	})

	t.Run("SuccessfulCreate", func(t *testing.T) {
		email := "test@example.com"
		password := "password123"

		collection := &core.Collection{}
		app.On("FindCachedCollectionByNameOrId", core.CollectionNameSuperusers).Return(collection, nil)

		record := &core.Record{}
		app.On("Save", mock.AnythingOfType("*core.Record")).Return(nil)

		command.SetArgs([]string{email, password})
		err := command.Execute()

		assert.NoError(t, err)
		app.AssertExpectations(t)
	})
}

func TestSuperuserUpdateCommand(t *testing.T) {
	app := &mockApp{}
	command := cmd.NewSuperuserCommand(app).Commands()[2]

	t.Run("MissingArguments", func(t *testing.T) {
		err := command.Execute()
		assert.Error(t, err)
		assert.Equal(t, "Missing email and password arguments.", err.Error())
	})

	t.Run("InvalidEmail", func(t *testing.T) {
		command.SetArgs([]string{"invalid-email", "password"})
		err := command.Execute()
		assert.Error(t, err)
		assert.Equal(t, "Missing or invalid email address.", err.Error())
	})

	t.Run("SuperuserNotFound", func(t *testing.T) {
		email := "test@example.com"
		password := "password123"

		app.On("FindAuthRecordByEmail", core.CollectionNameSuperusers, email).Return(&core.Record{}, errors.New("not found"))

		command.SetArgs([]string{email, password})
		err := command.Execute()

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "doesn't exist")
		app.AssertExpectations(t)
	})

	t.Run("SuccessfulUpdate", func(t *testing.T) {
		email := "test@example.com"
		password := "password123"

		record := &core.Record{}
		app.On("FindAuthRecordByEmail", core.CollectionNameSuperusers, email).Return(record, nil)
		app.On("Save", record).Return(nil)

		command.SetArgs([]string{email, password})
		err := command.Execute()

		assert.NoError(t, err)
		app.AssertExpectations(t)
	})
}

func TestSuperuserDeleteCommand(t *testing.T) {
	app := &mockApp{}
	command := cmd.NewSuperuserCommand(app).Commands()[3]

	t.Run("MissingEmail", func(t *testing.T) {
		err := command.Execute()
		assert.Error(t, err)
		assert.Equal(t, "Invalid or missing email address.", err.Error())
	})

	t.Run("InvalidEmail", func(t *testing.T) {
		command.SetArgs([]string{"invalid-email"})
		err := command.Execute()
		assert.Error(t, err)
		assert.Equal(t, "Invalid or missing email address.", err.Error())
	})

	t.Run("SuperuserNotFound", func(t *testing.T) {
		email := "test@example.com"

		app.On("FindAuthRecordByEmail", core.CollectionNameSuperusers, email).Return(&core.Record{}, errors.New("not found"))

		command.SetArgs([]string{email})
		err := command.Execute()

		assert.NoError(t, err) // This command doesn't return an error for missing superuser
		app.AssertExpectations(t)
	})

	t.Run("SuccessfulDelete", func(t *testing.T) {
		email := "test@example.com"

		record := &core.Record{}
		app.On("FindAuthRecordByEmail", core.CollectionNameSuperusers, email).Return(record, nil)
		app.On("Delete", record).Return(nil)

		command.SetArgs([]string{email})
		err := command.Execute()

		assert.NoError(t, err)
		app.AssertExpectations(t)
	})
}

func TestSuperuserOTPCommand(t *testing.T) {
	app := &mockApp{}
	command := cmd.NewSuperuserCommand(app).Commands()[4]

	t.Run("MissingEmail", func(t *testing.T) {
		err := command.Execute()
		assert.Error(t, err)
		assert.Equal(t, "Invalid or missing email address.", err.Error())
	})

	t.Run("InvalidEmail", func(t *testing.T) {
		command.SetArgs([]string{"invalid-email"})
		err := command.Execute()
		assert.Error(t, err)
		assert.Equal(t, "Invalid or missing email address.", err.Error())
	})

	t.Run("SuperuserNotFound", func(t *testing.T) {
		email := "test@example.com"

		app.On("FindAuthRecordByEmail", core.CollectionNameSuperusers, email).Return(&core.Record{}, errors.New("not found"))

		command.SetArgs([]string{email})
		err := command.Execute()

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "doesn't exist")
		app.AssertExpectations(t)
	})

	t.Run("OTPNotEnabled", func(t *testing.T) {
		email := "test@example.com"

		collection := &core.Collection{OTP: &core.OTPConfig{Enabled: false}}
		record := &core.Record{}
		record.setCollection(collection)

		app.On("FindAuthRecordByEmail", core.CollectionNameSuperusers, email).Return(record, nil)

		command.SetArgs([]string{email})
		err := command.Execute()

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "OTP is not enabled")
		app.AssertExpectations(t)
	})

	// Note: Testing the successful OTP creation scenario would require mocking more complex behaviors
	// and possibly adjusting the source code to allow for better testability. This might involve
	// dependency injection for the OTP creation process.
}
