package core

//
//import (
//	"context"
//	"github.com/pocketbase/pocketbase/tools/hook"
//	"github.com/pocketbase/pocketbase/tools/types"
//	"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/mock"
//	"log/slog"
//	"testing"
//)
//
//type MockApp struct {
//	mock.Mock
//}
//
//func (m *MockApp) FindCachedCollectionByNameOrId(nameOrId string) (*Collection, error) {
//	args := m.Called(nameOrId)
//	return args.Get(0).(*Collection), args.Error(1)
//}
//
//func (m *MockApp) UnsafeWithoutHooks() bool {
//	args := m.Called()
//	return args.Bool(0)
//}
//
//func (m *MockApp) AppLogger() *slog.Logger {
//	args := m.Called()
//	return args.Get(0).(*slog.Logger)
//}
//
//func (m *MockApp) IsBootstrapped() bool {
//	args := m.Called()
//	return args.Bool(0)
//}
//
//func (m *MockApp) Logger() *slog.Logger {
//	args := m.Called()
//	return args.Get(0).(*slog.Logger)
//}
//
//func (m *MockApp) IsTransactional() bool {
//	args := m.Called()
//	return args.Bool(0)
//}
//
//func (m *MockApp) Bootstrap() error {
//	args := m.Called()
//	return args.Error(0)
//}
//
//func TestNewAuthOrigin(t *testing.T) {
//	mockApp := new(MockApp)
//	mockCollection := NewCollection(&SchemaField{})
//	mockCollection.Name = CollectionNameAuthOrigins
//
//	mockApp.On(" FindCachedCollectionByNameOrId", CollectionNameAuthOrigins).Return(mockCollection, nil)
//
//	authOrigin := NewAuthOrigin(mockApp)
//
//	assert.NotNil(t, authOrigin)
//	assert.Equal(t, CollectionNameAuthOrigins, authOrigin.Record.Collection().Name)
//}
//
//func TestAuthOrigin_PreValidate(t *testing.T) {
//	mockApp := new(MockApp)
//	authOrigin := NewAuthOrigin(mockApp)
//
//	err := authOrigin.PreValidate(context.Background(), mockApp)
//	assert.NoError(t, err)
//
//	invalidAuthOrigin := &AuthOrigin{}
//	err = invalidAuthOrigin.PreValidate(context.Background(), mockApp)
//	assert.Error(t, err)
//	assert.Equal(t, " missing or invalid AuthOrigin ProxyRecord", err.Error())
//}
//
//func TestAuthOrigin_ProxyRecord(t *testing.T) {
//	mockApp := new(MockApp)
//	authOrigin := NewAuthOrigin(mockApp)
//
//	proxyRecord := authOrigin.ProxyRecord()
//	assert.NotNil(t, proxyRecord)
//	assert.Equal(t, authOrigin.Record, proxyRecord)
//}
//
//func TestAuthOrigin_SetProxyRecord(t *testing.T) {
//	mockApp := new(MockApp)
//	authOrigin := NewAuthOrigin(mockApp)
//
//	newRecord := NewRecord(NewCollection(&SchemaField{}))
//	authOrigin.SetProxyRecord(newRecord)
//
//	assert.Equal(t, newRecord, authOrigin.Record)
//}
//
//func TestAuthOrigin_CollectionRef(t *testing.T) {
//	mockApp := new(MockApp)
//	authOrigin := NewAuthOrigin(mockApp)
//
//	collectionRef := " testCollectionRef"
//	authOrigin.SetCollectionRef(collectionRef)
//
//	assert.Equal(t, collectionRef, authOrigin.CollectionRef())
//}
//
//func TestAuthOrigin_RecordRef(t *testing.T) {
//	mockApp := new(MockApp)
//	authOrigin := NewAuthOrigin(mockApp)
//
//	recordRef := " testRecordRef"
//	authOrigin.SetRecordRef(recordRef)
//
//	assert.Equal(t, recordRef, authOrigin.RecordRef())
//}
//
//func TestAuthOrigin_Fingerprint(t *testing.T) {
//	mockApp := new(MockApp)
//	authOrigin := NewAuthOrigin(mockApp)
//
//	fingerprint := " testFingerprint"
//	authOrigin.SetFingerprint(fingerprint)
//
//	assert.Equal(t, fingerprint, authOrigin.Fingerprint())
//}
//
//func TestAuthOrigin_Created(t *testing.T) {
//	mockApp := new(MockApp)
//	authOrigin := NewAuthOrigin(mockApp)
//
//	now := types.NowDateTime()
//	authOrigin.Set(" created", now)
//
//	assert.Equal(t, now, authOrigin.Created())
//}
//
//func TestAuthOrigin_Updated(t *testing.T) {
//	mockApp := new(MockApp)
//	authOrigin := NewAuthOrigin(mockApp)
//
//	now := types.NowDateTime()
//	authOrigin.Set(" updated", now)
//
//	assert.Equal(t, now, authOrigin.Updated())
//}
//
//type MockBaseApp struct {
//	mock.Mock
//}
//
//func (m *MockBaseApp) OnRecordUpdate() *hook.Hook[*RecordEvent] {
//	args := m.Called()
//	return args.Get(0).(*hook.Hook[*RecordEvent])
//}
//
//func (m *MockBaseApp) DeleteAllAuthOriginsByRecord(record *Record) error {
//	args := m.Called(record)
//	return args.Error(0)
//}
//
//func (m *MockBaseApp) Logger() Logger {
//	args := m.Called()
//	return args.Get(0).(Logger)
//}
//
//func TestRegisterAuthOriginHooks(t *testing.T) {
//	mockBaseApp := new(MockBaseApp)
//	mockHook := &hook.Hook[*RecordEvent]{}
//
//	mockBaseApp.On(" OnRecordUpdate").Return(mockHook)
//
//	app := &BaseApp{}
//	app.registerAuthOriginHooks()
//
//	mockBaseApp.AssertExpectations(t)
//}
//
//type MockLogger struct {
//	mock.Mock
//}
//
//func (m *MockLogger) Warn(message string, args ...any) {
//	m.Called(message, args)
//}
//
//func TestRecordUpdateHook(t *testing.T) {
//	mockBaseApp := new(MockBaseApp)
//	mockLogger := new(MockLogger)
//	mockHook := &hook.Hook[*RecordEvent]{}
//
//	mockBaseApp.On(" OnRecordUpdate").Return(mockHook)
//	mockBaseApp.On(" Logger").Return(mockLogger)
//
//	app := &BaseApp{}
//	app.registerAuthOriginHooks()
//
//	handler := mockHook.handlers[len(mockHook.handlers)-1]
//
//	// Test case: password change
//	event := &RecordEvent{
//		App: mockBaseApp,
//		Record: &Record{
//			collection: &Collection{
//				BaseCollection: BaseCollection{
//					Type: CollectionTypeAuth,
//				},
//			},
//		},
//	}
//	event.Record.Set(FieldNamePassword+" :hash", " newhash")
//	event.Record.(*Record).original = NewRecord(event.Record.Collection())
//	event.Record.(*Record).original.Set(FieldNamePassword+" :hash", " oldhash")
//
//	mockBaseApp.On(" DeleteAllAuthOriginsByRecord", event.Record).Return(nil)
//
//	err := handler.Func(event)
//	assert.NoError(t, err)
//
//	mockBaseApp.AssertExpectations(t)
//}
//
//func TestRecordRefHooks(t *testing.T) {
//	mockApp := new(MockApp)
//	mockHook := &hook.Hook[*RecordEvent]{}
//	mockCollectionHook := &hook.Hook[*CollectionEvent]{}
//
//	mockApp.On(" OnRecordValidate", CollectionNameAuthOrigins).Return(mockHook)
//	mockApp.On(" OnCollectionDeleteExecute").Return(mockCollectionHook)
//	mockApp.On(" OnRecordDeleteExecute").Return(mockHook)
//
//	recordRefHooks[*AuthOrigin](mockApp, CollectionNameAuthOrigins, CollectionTypeAuth)
//
//	mockApp.AssertExpectations(t)
//}
