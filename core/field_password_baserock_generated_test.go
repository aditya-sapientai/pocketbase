package core

import (
	"errors"
	"testing"

	"github.com/pocketbase/dbx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAppForFieldPassword struct {
	mock.Mock
}

func (m *MockAppForFieldPassword) RecordQuery(collectionNameOrId string) *dbx.SelectQuery {
	args := m.Called(collectionNameOrId)
	return args.Get(0).(*dbx.SelectQuery)
}

func (m *MockAppForFieldPassword) Delete(model interface{}) error {
	args := m.Called(model)
	return args.Error(0)
}

type MockSelectQuery struct {
	mock.Mock
}

func (m *MockSelectQuery) AndWhere(cond dbx.Expression) *dbx.SelectQuery {
	args := m.Called(cond)
	return args.Get(0).(*dbx.SelectQuery)
}

func (m *MockSelectQuery) OrderBy(cols ...string) *dbx.SelectQuery {
	args := m.Called(cols)
	return args.Get(0).(*dbx.SelectQuery)
}

func (m *MockSelectQuery) All(slice interface{}) error {
	args := m.Called(slice)
	return args.Error(0)
}

func (m *MockSelectQuery) Limit(limit int64) *dbx.SelectQuery {
	args := m.Called(limit)
	return args.Get(0).(*dbx.SelectQuery)
}

func (m *MockSelectQuery) One(pointer interface{}) error {
	args := m.Called(pointer)
	return args.Error(0)
}

func TestFindAllAuthOriginsByRecord(t *testing.T) {
	mockApp := new(MockAppForFieldPassword)
	mockQuery := new(MockSelectQuery)

	mockApp.On("RecordQuery", CollectionNameAuthOrigins).Return(mockQuery)
	mockQuery.On("AndWhere", mock.Anything).Return(mockQuery)
	mockQuery.On("OrderBy", "created DESC").Return(mockQuery)

	testCases := []struct {
		name          string
		authRecord    *Record
		expectedError error
		mockSetup     func()
	}{
		{
			name:       "Success",
			authRecord: &Record{},
			mockSetup: func() {
				mockQuery.On("All", mock.AnythingOfType("*[]*core.AuthOrigin")).Return(nil).Once()
			},
		},
		{
			name:          "Error",
			authRecord:    &Record{},
			expectedError: errors.New("database error"),
			mockSetup: func() {
				mockQuery.On("All", mock.AnythingOfType("*[]*core.AuthOrigin")).Return(errors.New("database error")).Once()
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()

			app := &BaseApp{}
			//app.MockAppForFieldPassword = mockApp

			results, err := app.FindAllAuthOriginsByRecord(tc.authRecord)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedError.Error())
				assert.Nil(t, results)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, results)
			}
		})
	}
}

func TestFindAllAuthOriginsByCollection(t *testing.T) {
	mockApp := new(MockAppForFieldPassword)
	mockQuery := new(MockSelectQuery)

	mockApp.On("RecordQuery", CollectionNameAuthOrigins).Return(mockQuery)
	mockQuery.On("AndWhere", mock.Anything).Return(mockQuery)
	mockQuery.On("OrderBy", "created DESC").Return(mockQuery)

	testCases := []struct {
		name          string
		collection    *Collection
		expectedError error
		mockSetup     func()
	}{
		{
			name:       "Success",
			collection: &Collection{},
			mockSetup: func() {
				mockQuery.On("All", mock.AnythingOfType("*[]*core.AuthOrigin")).Return(nil).Once()
			},
		},
		{
			name:          "Error",
			collection:    &Collection{},
			expectedError: errors.New("database error"),
			mockSetup: func() {
				mockQuery.On("All", mock.AnythingOfType("*[]*core.AuthOrigin")).Return(errors.New("database error")).Once()
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()

			app := &BaseApp{}
			//app.MockAppForFieldPassword = mockApp

			results, err := app.FindAllAuthOriginsByCollection(tc.collection)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedError.Error())
				assert.Nil(t, results)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, results)
			}
		})
	}
}

func TestFindAuthOriginById(t *testing.T) {
	mockApp := new(MockAppForFieldPassword)
	mockQuery := new(MockSelectQuery)

	mockApp.On("RecordQuery", CollectionNameAuthOrigins).Return(mockQuery)
	mockQuery.On("AndWhere", mock.Anything).Return(mockQuery)
	mockQuery.On("Limit", int64(1)).Return(mockQuery)

	testCases := []struct {
		name          string
		id            string
		expectedError error
		mockSetup     func()
	}{
		{
			name: "Success",
			id:   "testId",
			mockSetup: func() {
				mockQuery.On("One", mock.AnythingOfType("*core.AuthOrigin")).Return(nil).Once()
			},
		},
		{
			name:          "Error",
			id:            "testId",
			expectedError: errors.New("database error"),
			mockSetup: func() {
				mockQuery.On("One", mock.AnythingOfType("*core.AuthOrigin")).Return(errors.New("database error")).Once()
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()

			app := &BaseApp{}
			//app.MockAppForFieldPassword = mockApp

			result, err := app.FindAuthOriginById(tc.id)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedError.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

func TestFindAuthOriginByRecordAndFingerprint(t *testing.T) {
	mockApp := new(MockAppForFieldPassword)
	mockQuery := new(MockSelectQuery)

	mockApp.On("RecordQuery", CollectionNameAuthOrigins).Return(mockQuery)
	mockQuery.On("AndWhere", mock.Anything).Return(mockQuery)
	mockQuery.On("Limit", int64(1)).Return(mockQuery)

	testCases := []struct {
		name          string
		authRecord    *Record
		fingerprint   string
		expectedError error
		mockSetup     func()
	}{
		{
			name:        "Success",
			authRecord:  &Record{},
			fingerprint: "testFingerprint",
			mockSetup: func() {
				mockQuery.On("One", mock.AnythingOfType("*core.AuthOrigin")).Return(nil).Once()
			},
		},
		{
			name:          "Error",
			authRecord:    &Record{},
			fingerprint:   "testFingerprint",
			expectedError: errors.New("database error"),
			mockSetup: func() {
				mockQuery.On("One", mock.AnythingOfType("*core.AuthOrigin")).Return(errors.New("database error")).Once()
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()

			app := &BaseApp{}
			//app.MockAppForFieldPassword = mockApp

			result, err := app.FindAuthOriginByRecordAndFingerprint(tc.authRecord, tc.fingerprint)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedError.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

func TestDeleteAllAuthOriginsByRecord(t *testing.T) {
	mockApp := new(MockAppForFieldPassword)
	mockQuery := new(MockSelectQuery)

	mockApp.On("RecordQuery", CollectionNameAuthOrigins).Return(mockQuery)
	mockQuery.On("AndWhere", mock.Anything).Return(mockQuery)
	mockQuery.On("OrderBy", "created DESC").Return(mockQuery)

	testCases := []struct {
		name          string
		authRecord    *Record
		expectedError error
		mockSetup     func()
	}{
		{
			name:       "Success",
			authRecord: &Record{},
			mockSetup: func() {
				mockQuery.On("All", mock.AnythingOfType("*[]*core.AuthOrigin")).Run(func(args mock.Arguments) {
					arg := args.Get(0).(*[]*AuthOrigin)
					*arg = []*AuthOrigin{{}, {}}
				}).Return(nil).Once()
				mockApp.On("Delete", mock.AnythingOfType("*core.AuthOrigin")).Return(nil).Twice()
			},
		},
		{
			name:          "Error in FindAllAuthOriginsByRecord",
			authRecord:    &Record{},
			expectedError: errors.New("database error"),
			mockSetup: func() {
				mockQuery.On("All", mock.AnythingOfType("*[]*core.AuthOrigin")).Return(errors.New("database error")).Once()
			},
		},
		{
			name:          "Error in Delete",
			authRecord:    &Record{},
			expectedError: errors.New("delete error"),
			mockSetup: func() {
				mockQuery.On("All", mock.AnythingOfType("*[]*core.AuthOrigin")).Run(func(args mock.Arguments) {
					arg := args.Get(0).(*[]*AuthOrigin)
					*arg = []*AuthOrigin{{}}
				}).Return(nil).Once()
				mockApp.On("Delete", mock.AnythingOfType("*core.AuthOrigin")).Return(errors.New("delete error")).Once()
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()

			app := &BaseApp{}
			//app.MockAppForFieldPassword = mockApp

			err := app.DeleteAllAuthOriginsByRecord(tc.authRecord)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tc.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
