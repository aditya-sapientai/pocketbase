package core_test

//
//import (
//	"errors"
//	"github.com/pocketbase/pocketbase/core"
//	"testing"
//
//	"github.com/pocketbase/dbx"
//	"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/mock"
//)
//
//type MockBaseApp struct {
//	mock.Mock
//}
//
//func (m *MockBaseApp) RecordQuery(collectionNameOrId string) *dbx.SelectQuery {
//	args := m.Called(collectionNameOrId)
//	return args.Get(0).(*dbx.SelectQuery)
//}
//
//func (m *MockBaseApp) Delete(model interface{}) error {
//	args := m.Called(model)
//	return args.Error(0)
//}
//
//type MockSelectQuery struct {
//	mock.Mock
//}
//
//func (m *MockSelectQuery) AndWhere(cond interface{}) *dbx.SelectQuery {
//	args := m.Called(cond)
//	return args.Get(0).(*dbx.SelectQuery)
//}
//
//func (m *MockSelectQuery) OrderBy(columns ...string) *dbx.SelectQuery {
//	args := m.Called(columns)
//	return args.Get(0).(*dbx.SelectQuery)
//}
//
//func (m *MockSelectQuery) Limit(limit int64) *dbx.SelectQuery {
//	args := m.Called(limit)
//	return args.Get(0).(*dbx.SelectQuery)
//}
//
//func (m *MockSelectQuery) One(a interface{}) error {
//	args := m.Called(a)
//	return args.Error(0)
//}
//
//func (m *MockSelectQuery) All(a interface{}) error {
//	args := m.Called(a)
//	return args.Error(0)
//}
//
//func TestFindAllAuthOriginsByRecord(t *testing.T) {
//	mockApp := new(MockBaseApp)
//	mockQuery := new(MockSelectQuery)
//
//	mockApp.On("RecordQuery", core.CollectionNameAuthOrigins).Return(mockQuery)
//	mockQuery.On("AndWhere", mock.Anything).Return(mockQuery)
//	mockQuery.On("OrderBy", "created DESC").Return(mockQuery)
//
//	t.Run("Success", func(t *testing.T) {
//		expectedResult := []*core.AuthOrigin{{}, {}}
//		mockQuery.On("All", mock.AnythingOfType("*[]*core.AuthOrigin")).Run(func(args mock.Arguments) {
//			arg := args.Get(0).(*[]*core.AuthOrigin)
//			*arg = expectedResult
//		}).Return(nil).Once()
//
//		authRecord := &core.Record{}
//		authRecord.SetCollection(&core.Collection{Id: "test_collection"})
//		authRecord.SetId("test_record")
//
//		result, err := mockApp.FindAllAuthOriginsByRecord(authRecord)
//
//		assert.NoError(t, err)
//		assert.Equal(t, expectedResult, result)
//	})
//
//	t.Run("Error", func(t *testing.T) {
//		expectedError := errors.New("database error")
//		mockQuery.On("All", mock.AnythingOfType("*[]*core.AuthOrigin")).Return(expectedError).Once()
//
//		authRecord := &core.Record{}
//		authRecord.SetCollection(&core.Collection{Id: "test_collection"})
//		authRecord.SetId("test_record")
//
//		result, err := mockApp.FindAllAuthOriginsByRecord(authRecord)
//
//		assert.Error(t, err)
//		assert.Equal(t, expectedError, err)
//		assert.Nil(t, result)
//	})
//}
//
//func TestFindAllAuthOriginsByCollection(t *testing.T) {
//	mockApp := new(MockBaseApp)
//	mockQuery := new(MockSelectQuery)
//
//	mockApp.On("RecordQuery", core.CollectionNameAuthOrigins).Return(mockQuery)
//	mockQuery.On("AndWhere", mock.Anything).Return(mockQuery)
//	mockQuery.On("OrderBy", "created DESC").Return(mockQuery)
//
//	t.Run("Success", func(t *testing.T) {
//		expectedResult := []*core.AuthOrigin{{}, {}}
//		mockQuery.On("All", mock.AnythingOfType("*[]*core.AuthOrigin")).Run(func(args mock.Arguments) {
//			arg := args.Get(0).(*[]*core.AuthOrigin)
//			*arg = expectedResult
//		}).Return(nil).Once()
//
//		collection := &core.Collection{Id: "test_collection"}
//
//		result, err := mockApp.FindAllAuthOriginsByCollection(collection)
//
//		assert.NoError(t, err)
//		assert.Equal(t, expectedResult, result)
//	})
//
//	t.Run("Error", func(t *testing.T) {
//		expectedError := errors.New("database error")
//		mockQuery.On("All", mock.AnythingOfType("*[]*core.AuthOrigin")).Return(expectedError).Once()
//
//		collection := &core.Collection{Id: "test_collection"}
//
//		result, err := mockApp.FindAllAuthOriginsByCollection(collection)
//
//		assert.Error(t, err)
//		assert.Equal(t, expectedError, err)
//		assert.Nil(t, result)
//	})
//}
//
//func TestFindAuthOriginById(t *testing.T) {
//	mockApp := new(MockBaseApp)
//	mockQuery := new(MockSelectQuery)
//
//	mockApp.On("RecordQuery", core.CollectionNameAuthOrigins).Return(mockQuery)
//	mockQuery.On("AndWhere", mock.Anything).Return(mockQuery)
//	mockQuery.On("Limit", int64(1)).Return(mockQuery)
//
//	t.Run("Success", func(t *testing.T) {
//		expectedResult := &core.AuthOrigin{Id: "test_id"}
//		mockQuery.On("One", mock.AnythingOfType("*core.AuthOrigin")).Run(func(args mock.Arguments) {
//			arg := args.Get(0).(*core.AuthOrigin)
//			*arg = *expectedResult
//		}).Return(nil).Once()
//
//		result, err := mockApp.FindAuthOriginById("test_id")
//
//		assert.NoError(t, err)
//		assert.Equal(t, expectedResult, result)
//	})
//
//	t.Run("Error", func(t *testing.T) {
//		expectedError := errors.New("database error")
//		mockQuery.On("One", mock.AnythingOfType("*core.AuthOrigin")).Return(expectedError).Once()
//
//		result, err := mockApp.FindAuthOriginById("test_id")
//
//		assert.Error(t, err)
//		assert.Equal(t, expectedError, err)
//		assert.Nil(t, result)
//	})
//}
//
//func TestFindAuthOriginByRecordAndFingerprint(t *testing.T) {
//	mockApp := new(MockBaseApp)
//	mockQuery := new(MockSelectQuery)
//
//	mockApp.On("RecordQuery", core.CollectionNameAuthOrigins).Return(mockQuery)
//	mockQuery.On("AndWhere", mock.Anything).Return(mockQuery)
//	mockQuery.On("Limit", int64(1)).Return(mockQuery)
//
//	t.Run("Success", func(t *testing.T) {
//		expectedResult := &core.AuthOrigin{Id: "test_id"}
//		mockQuery.On("One", mock.AnythingOfType("*core.AuthOrigin")).Run(func(args mock.Arguments) {
//			arg := args.Get(0).(*core.AuthOrigin)
//			*arg = *expectedResult
//		}).Return(nil).Once()
//
//		authRecord := &core.Record{}
//		authRecord.SetCollection(&core.Collection{Id: "test_collection"})
//		authRecord.SetId("test_record")
//
//		result, err := mockApp.FindAuthOriginByRecordAndFingerprint(authRecord, "test_fingerprint")
//
//		assert.NoError(t, err)
//		assert.Equal(t, expectedResult, result)
//	})
//
//	t.Run("Error", func(t *testing.T) {
//		expectedError := errors.New("database error")
//		mockQuery.On("One", mock.AnythingOfType("*core.AuthOrigin")).Return(expectedError).Once()
//
//		authRecord := &core.Record{}
//		authRecord.SetCollection(&core.Collection{Id: "test_collection"})
//		authRecord.SetId("test_record")
//
//		result, err := mockApp.FindAuthOriginByRecordAndFingerprint(authRecord, "test_fingerprint")
//
//		assert.Error(t, err)
//		assert.Equal(t, expectedError, err)
//		assert.Nil(t, result)
//	})
//}
//
//func TestDeleteAllAuthOriginsByRecord(t *testing.T) {
//	mockApp := new(MockBaseApp)
//	mockQuery := new(MockSelectQuery)
//
//	mockApp.On("RecordQuery", core.CollectionNameAuthOrigins).Return(mockQuery)
//	mockQuery.On("AndWhere", mock.Anything).Return(mockQuery)
//	mockQuery.On("OrderBy", "created DESC").Return(mockQuery)
//
//	t.Run("Success", func(t *testing.T) {
//		authOrigins := []*core.AuthOrigin{{Id: "1"}, {Id: "2"}}
//		mockQuery.On("All", mock.AnythingOfType("*[]*core.AuthOrigin")).Run(func(args mock.Arguments) {
//			arg := args.Get(0).(*[]*core.AuthOrigin)
//			*arg = authOrigins
//		}).Return(nil).Once()
//
//		mockApp.On("Delete", mock.AnythingOfType("*core.AuthOrigin")).Return(nil).Times(len(authOrigins))
//
//		authRecord := &core.Record{}
//		authRecord.SetCollection(&core.Collection{Id: "test_collection"})
//		authRecord.SetId("test_record")
//
//		err := mockApp.DeleteAllAuthOriginsByRecord(authRecord)
//
//		assert.NoError(t, err)
//	})
//
//	t.Run("FindError", func(t *testing.T) {
//		expectedError := errors.New("find error")
//		mockQuery.On("All", mock.AnythingOfType("*[]*core.AuthOrigin")).Return(expectedError).Once()
//
//		authRecord := &core.Record{}
//		authRecord.SetCollection(&core.Collection{Id: "test_collection"})
//		authRecord.SetId("test_record")
//
//		err := mockApp.DeleteAllAuthOriginsByRecord(authRecord)
//
//		assert.Error(t, err)
//		assert.Equal(t, expectedError, err)
//	})
//
//	t.Run("PartialDeleteError", func(t *testing.T) {
//		authOrigins := []*core.AuthOrigin{{Id: "1"}, {Id: "2"}}
//		mockQuery.On("All", mock.AnythingOfType("*[]*core.AuthOrigin")).Run(func(args mock.Arguments) {
//			arg := args.Get(0).(*[]*core.AuthOrigin)
//			*arg = authOrigins
//		}).Return(nil).Once()
//
//		deleteError := errors.New("delete error")
//		mockApp.On("Delete", mock.MatchedBy(func(ao *core.AuthOrigin) bool {
//			return ao.Id == "1"
//		})).Return(nil).Once()
//		mockApp.On("Delete", mock.MatchedBy(func(ao *core.AuthOrigin) bool {
//			return ao.Id == "2"
//		})).Return(deleteError).Once()
//
//		authRecord := &core.Record{}
//		authRecord.SetCollection(&core.Collection{Id: "test_collection"})
//		authRecord.SetId("test_record")
//
//		err := mockApp.DeleteAllAuthOriginsByRecord(authRecord)
//
//		assert.Error(t, err)
//		assert.Equal(t, deleteError, err)
//	})
//}
