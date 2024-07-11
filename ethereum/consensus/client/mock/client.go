// Code generated by MockGen. DO NOT EDIT.
// Source: client.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	types "github.com/kilnfi/go-utils/ethereum/consensus/types"
	capella "github.com/protolambda/zrnt/eth2/beacon/capella"
	common "github.com/protolambda/zrnt/eth2/beacon/common"
	phase0 "github.com/protolambda/zrnt/eth2/beacon/phase0"
)

// MockClient is a mock of Client interface.
type MockClient struct {
	ctrl     *gomock.Controller
	recorder *MockClientMockRecorder
}

// MockClientMockRecorder is the mock recorder for MockClient.
type MockClientMockRecorder struct {
	mock *MockClient
}

// NewMockClient creates a new mock instance.
func NewMockClient(ctrl *gomock.Controller) *MockClient {
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &MockClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClient) EXPECT() *MockClientMockRecorder {
	return m.recorder
}

// GetAttestations mocks base method.
func (m *MockClient) GetAttestations(ctx context.Context) (phase0.Attestations, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAttestations", ctx)
	ret0, _ := ret[0].(phase0.Attestations)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAttestations indicates an expected call of GetAttestations.
func (mr *MockClientMockRecorder) GetAttestations(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAttestations", reflect.TypeOf((*MockClient)(nil).GetAttestations), ctx)
}

// GetAttesterSlashings mocks base method.
func (m *MockClient) GetAttesterSlashings(ctx context.Context) (phase0.AttesterSlashings, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAttesterSlashings", ctx)
	ret0, _ := ret[0].(phase0.AttesterSlashings)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAttesterSlashings indicates an expected call of GetAttesterSlashings.
func (mr *MockClientMockRecorder) GetAttesterSlashings(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAttesterSlashings", reflect.TypeOf((*MockClient)(nil).GetAttesterSlashings), ctx)
}

// GetBlock mocks base method.
func (m *MockClient) GetBlock(ctx context.Context, blockID string) (*capella.SignedBeaconBlock, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBlock", ctx, blockID)
	ret0, _ := ret[0].(*capella.SignedBeaconBlock)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBlock indicates an expected call of GetBlock.
func (mr *MockClientMockRecorder) GetBlock(ctx, blockID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBlock", reflect.TypeOf((*MockClient)(nil).GetBlock), ctx, blockID)
}

// GetBlockAttestations mocks base method.
func (m *MockClient) GetBlockAttestations(ctx context.Context, blockID string) (phase0.Attestations, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBlockAttestations", ctx, blockID)
	ret0, _ := ret[0].(phase0.Attestations)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBlockAttestations indicates an expected call of GetBlockAttestations.
func (mr *MockClientMockRecorder) GetBlockAttestations(ctx, blockID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBlockAttestations", reflect.TypeOf((*MockClient)(nil).GetBlockAttestations), ctx, blockID)
}

// GetBlockHeader mocks base method.
func (m *MockClient) GetBlockHeader(ctx context.Context, blockID string) (*types.BeaconBlockHeader, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBlockHeader", ctx, blockID)
	ret0, _ := ret[0].(*types.BeaconBlockHeader)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBlockHeader indicates an expected call of GetBlockHeader.
func (mr *MockClientMockRecorder) GetBlockHeader(ctx, blockID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBlockHeader", reflect.TypeOf((*MockClient)(nil).GetBlockHeader), ctx, blockID)
}

// GetBlockHeaders mocks base method.
func (m *MockClient) GetBlockHeaders(ctx context.Context, slot *common.Slot, parentRoot *common.Root) ([]*types.BeaconBlockHeader, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBlockHeaders", ctx, slot, parentRoot)
	ret0, _ := ret[0].([]*types.BeaconBlockHeader)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBlockHeaders indicates an expected call of GetBlockHeaders.
func (mr *MockClientMockRecorder) GetBlockHeaders(ctx, slot, parentRoot interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBlockHeaders", reflect.TypeOf((*MockClient)(nil).GetBlockHeaders), ctx, slot, parentRoot)
}

// GetBlockRoot mocks base method.
func (m *MockClient) GetBlockRoot(ctx context.Context, blockID string) (*common.Root, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBlockRoot", ctx, blockID)
	ret0, _ := ret[0].(*common.Root)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBlockRoot indicates an expected call of GetBlockRoot.
func (mr *MockClientMockRecorder) GetBlockRoot(ctx, blockID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBlockRoot", reflect.TypeOf((*MockClient)(nil).GetBlockRoot), ctx, blockID)
}

// GetCommittees mocks base method.
func (m *MockClient) GetCommittees(ctx context.Context, stateID string, epoch *common.Epoch, index *common.CommitteeIndex, slot *common.Slot) ([]*types.Committee, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCommittees", ctx, stateID, epoch, index, slot)
	ret0, _ := ret[0].([]*types.Committee)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCommittees indicates an expected call of GetCommittees.
func (mr *MockClientMockRecorder) GetCommittees(ctx, stateID, epoch, index, slot interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCommittees", reflect.TypeOf((*MockClient)(nil).GetCommittees), ctx, stateID, epoch, index, slot)
}

// GetGenesis mocks base method.
func (m *MockClient) GetGenesis(ctx context.Context) (*types.Genesis, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGenesis", ctx)
	ret0, _ := ret[0].(*types.Genesis)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGenesis indicates an expected call of GetGenesis.
func (mr *MockClientMockRecorder) GetGenesis(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGenesis", reflect.TypeOf((*MockClient)(nil).GetGenesis), ctx)
}

// GetNodeVersion mocks base method.
func (m *MockClient) GetNodeVersion(ctx context.Context) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNodeVersion", ctx)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNodeVersion indicates an expected call of GetNodeVersion.
func (mr *MockClientMockRecorder) GetNodeVersion(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNodeVersion", reflect.TypeOf((*MockClient)(nil).GetNodeVersion), ctx)
}

// GetProposerSlashings mocks base method.
func (m *MockClient) GetProposerSlashings(ctx context.Context) (phase0.ProposerSlashings, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProposerSlashings", ctx)
	ret0, _ := ret[0].(phase0.ProposerSlashings)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProposerSlashings indicates an expected call of GetProposerSlashings.
func (mr *MockClientMockRecorder) GetProposerSlashings(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProposerSlashings", reflect.TypeOf((*MockClient)(nil).GetProposerSlashings), ctx)
}

// GetSpec mocks base method.
func (m *MockClient) GetSpec(ctx context.Context) (*common.Spec, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSpec", ctx)
	ret0, _ := ret[0].(*common.Spec)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSpec indicates an expected call of GetSpec.
func (mr *MockClientMockRecorder) GetSpec(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSpec", reflect.TypeOf((*MockClient)(nil).GetSpec), ctx)
}

// GetStateFinalityCheckpoints mocks base method.
func (m *MockClient) GetStateFinalityCheckpoints(ctx context.Context, stateID string) (*types.StateFinalityCheckpoints, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStateFinalityCheckpoints", ctx, stateID)
	ret0, _ := ret[0].(*types.StateFinalityCheckpoints)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetStateFinalityCheckpoints indicates an expected call of GetStateFinalityCheckpoints.
func (mr *MockClientMockRecorder) GetStateFinalityCheckpoints(ctx, stateID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStateFinalityCheckpoints", reflect.TypeOf((*MockClient)(nil).GetStateFinalityCheckpoints), ctx, stateID)
}

// GetStateFork mocks base method.
func (m *MockClient) GetStateFork(ctx context.Context, stateID string) (*common.Fork, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStateFork", ctx, stateID)
	ret0, _ := ret[0].(*common.Fork)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetStateFork indicates an expected call of GetStateFork.
func (mr *MockClientMockRecorder) GetStateFork(ctx, stateID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStateFork", reflect.TypeOf((*MockClient)(nil).GetStateFork), ctx, stateID)
}

// GetStateRoot mocks base method.
func (m *MockClient) GetStateRoot(ctx context.Context, stateID string) (*common.Root, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStateRoot", ctx, stateID)
	ret0, _ := ret[0].(*common.Root)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetStateRoot indicates an expected call of GetStateRoot.
func (mr *MockClientMockRecorder) GetStateRoot(ctx, stateID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStateRoot", reflect.TypeOf((*MockClient)(nil).GetStateRoot), ctx, stateID)
}

// GetSyncCommittees mocks base method.
func (m *MockClient) GetSyncCommittees(ctx context.Context, stateID string, epoch *common.Epoch) (*types.SyncCommittees, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSyncCommittees", ctx, stateID, epoch)
	ret0, _ := ret[0].(*types.SyncCommittees)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSyncCommittees indicates an expected call of GetSyncCommittees.
func (mr *MockClientMockRecorder) GetSyncCommittees(ctx, stateID, epoch interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSyncCommittees", reflect.TypeOf((*MockClient)(nil).GetSyncCommittees), ctx, stateID, epoch)
}

// GetValidator mocks base method.
func (m *MockClient) GetValidator(ctx context.Context, stateID, validatorID string) (*types.Validator, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetValidator", ctx, stateID, validatorID)
	ret0, _ := ret[0].(*types.Validator)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetValidator indicates an expected call of GetValidator.
func (mr *MockClientMockRecorder) GetValidator(ctx, stateID, validatorID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetValidator", reflect.TypeOf((*MockClient)(nil).GetValidator), ctx, stateID, validatorID)
}

// GetValidatorBalances mocks base method.
func (m *MockClient) GetValidatorBalances(ctx context.Context, stateID string, validatorIDs []string) ([]*types.ValidatorBalance, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetValidatorBalances", ctx, stateID, validatorIDs)
	ret0, _ := ret[0].([]*types.ValidatorBalance)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetValidatorBalances indicates an expected call of GetValidatorBalances.
func (mr *MockClientMockRecorder) GetValidatorBalances(ctx, stateID, validatorIDs interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetValidatorBalances", reflect.TypeOf((*MockClient)(nil).GetValidatorBalances), ctx, stateID, validatorIDs)
}

// GetValidators mocks base method.
func (m *MockClient) GetValidators(ctx context.Context, stateID string, validatorIDs, statuses []string) ([]*types.Validator, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetValidators", ctx, stateID, validatorIDs, statuses)
	ret0, _ := ret[0].([]*types.Validator)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetValidators indicates an expected call of GetValidators.
func (mr *MockClientMockRecorder) GetValidators(ctx, stateID, validatorIDs, statuses interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetValidators", reflect.TypeOf((*MockClient)(nil).GetValidators), ctx, stateID, validatorIDs, statuses)
}

// GetVoluntaryExits mocks base method.
func (m *MockClient) GetVoluntaryExits(ctx context.Context) (phase0.VoluntaryExits, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetVoluntaryExits", ctx)
	ret0, _ := ret[0].(phase0.VoluntaryExits)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetVoluntaryExits indicates an expected call of GetVoluntaryExits.
func (mr *MockClientMockRecorder) GetVoluntaryExits(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetVoluntaryExits", reflect.TypeOf((*MockClient)(nil).GetVoluntaryExits), ctx)
}

// MockBeaconClient is a mock of BeaconClient interface.
type MockBeaconClient struct {
	ctrl     *gomock.Controller
	recorder *MockBeaconClientMockRecorder
}

// MockBeaconClientMockRecorder is the mock recorder for MockBeaconClient.
type MockBeaconClientMockRecorder struct {
	mock *MockBeaconClient
}

// NewMockBeaconClient creates a new mock instance.
func NewMockBeaconClient(ctrl *gomock.Controller) *MockBeaconClient {
	mock := &MockBeaconClient{ctrl: ctrl}
	mock.recorder = &MockBeaconClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBeaconClient) EXPECT() *MockBeaconClientMockRecorder {
	return m.recorder
}

// GetAttestations mocks base method.
func (m *MockBeaconClient) GetAttestations(ctx context.Context) (phase0.Attestations, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAttestations", ctx)
	ret0, _ := ret[0].(phase0.Attestations)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAttestations indicates an expected call of GetAttestations.
func (mr *MockBeaconClientMockRecorder) GetAttestations(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAttestations", reflect.TypeOf((*MockBeaconClient)(nil).GetAttestations), ctx)
}

// GetAttesterSlashings mocks base method.
func (m *MockBeaconClient) GetAttesterSlashings(ctx context.Context) (phase0.AttesterSlashings, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAttesterSlashings", ctx)
	ret0, _ := ret[0].(phase0.AttesterSlashings)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAttesterSlashings indicates an expected call of GetAttesterSlashings.
func (mr *MockBeaconClientMockRecorder) GetAttesterSlashings(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAttesterSlashings", reflect.TypeOf((*MockBeaconClient)(nil).GetAttesterSlashings), ctx)
}

// GetBlock mocks base method.
func (m *MockBeaconClient) GetBlock(ctx context.Context, blockID string) (*capella.SignedBeaconBlock, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBlock", ctx, blockID)
	ret0, _ := ret[0].(*capella.SignedBeaconBlock)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBlock indicates an expected call of GetBlock.
func (mr *MockBeaconClientMockRecorder) GetBlock(ctx, blockID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBlock", reflect.TypeOf((*MockBeaconClient)(nil).GetBlock), ctx, blockID)
}

// GetBlockAttestations mocks base method.
func (m *MockBeaconClient) GetBlockAttestations(ctx context.Context, blockID string) (phase0.Attestations, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBlockAttestations", ctx, blockID)
	ret0, _ := ret[0].(phase0.Attestations)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBlockAttestations indicates an expected call of GetBlockAttestations.
func (mr *MockBeaconClientMockRecorder) GetBlockAttestations(ctx, blockID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBlockAttestations", reflect.TypeOf((*MockBeaconClient)(nil).GetBlockAttestations), ctx, blockID)
}

// GetBlockHeader mocks base method.
func (m *MockBeaconClient) GetBlockHeader(ctx context.Context, blockID string) (*types.BeaconBlockHeader, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBlockHeader", ctx, blockID)
	ret0, _ := ret[0].(*types.BeaconBlockHeader)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBlockHeader indicates an expected call of GetBlockHeader.
func (mr *MockBeaconClientMockRecorder) GetBlockHeader(ctx, blockID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBlockHeader", reflect.TypeOf((*MockBeaconClient)(nil).GetBlockHeader), ctx, blockID)
}

// GetBlockHeaders mocks base method.
func (m *MockBeaconClient) GetBlockHeaders(ctx context.Context, slot *common.Slot, parentRoot *common.Root) ([]*types.BeaconBlockHeader, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBlockHeaders", ctx, slot, parentRoot)
	ret0, _ := ret[0].([]*types.BeaconBlockHeader)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBlockHeaders indicates an expected call of GetBlockHeaders.
func (mr *MockBeaconClientMockRecorder) GetBlockHeaders(ctx, slot, parentRoot interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBlockHeaders", reflect.TypeOf((*MockBeaconClient)(nil).GetBlockHeaders), ctx, slot, parentRoot)
}

// GetBlockRoot mocks base method.
func (m *MockBeaconClient) GetBlockRoot(ctx context.Context, blockID string) (*common.Root, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBlockRoot", ctx, blockID)
	ret0, _ := ret[0].(*common.Root)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBlockRoot indicates an expected call of GetBlockRoot.
func (mr *MockBeaconClientMockRecorder) GetBlockRoot(ctx, blockID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBlockRoot", reflect.TypeOf((*MockBeaconClient)(nil).GetBlockRoot), ctx, blockID)
}

// GetCommittees mocks base method.
func (m *MockBeaconClient) GetCommittees(ctx context.Context, stateID string, epoch *common.Epoch, index *common.CommitteeIndex, slot *common.Slot) ([]*types.Committee, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCommittees", ctx, stateID, epoch, index, slot)
	ret0, _ := ret[0].([]*types.Committee)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCommittees indicates an expected call of GetCommittees.
func (mr *MockBeaconClientMockRecorder) GetCommittees(ctx, stateID, epoch, index, slot interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCommittees", reflect.TypeOf((*MockBeaconClient)(nil).GetCommittees), ctx, stateID, epoch, index, slot)
}

// GetGenesis mocks base method.
func (m *MockBeaconClient) GetGenesis(ctx context.Context) (*types.Genesis, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGenesis", ctx)
	ret0, _ := ret[0].(*types.Genesis)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGenesis indicates an expected call of GetGenesis.
func (mr *MockBeaconClientMockRecorder) GetGenesis(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGenesis", reflect.TypeOf((*MockBeaconClient)(nil).GetGenesis), ctx)
}

// GetProposerSlashings mocks base method.
func (m *MockBeaconClient) GetProposerSlashings(ctx context.Context) (phase0.ProposerSlashings, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProposerSlashings", ctx)
	ret0, _ := ret[0].(phase0.ProposerSlashings)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProposerSlashings indicates an expected call of GetProposerSlashings.
func (mr *MockBeaconClientMockRecorder) GetProposerSlashings(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProposerSlashings", reflect.TypeOf((*MockBeaconClient)(nil).GetProposerSlashings), ctx)
}

// GetStateFinalityCheckpoints mocks base method.
func (m *MockBeaconClient) GetStateFinalityCheckpoints(ctx context.Context, stateID string) (*types.StateFinalityCheckpoints, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStateFinalityCheckpoints", ctx, stateID)
	ret0, _ := ret[0].(*types.StateFinalityCheckpoints)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetStateFinalityCheckpoints indicates an expected call of GetStateFinalityCheckpoints.
func (mr *MockBeaconClientMockRecorder) GetStateFinalityCheckpoints(ctx, stateID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStateFinalityCheckpoints", reflect.TypeOf((*MockBeaconClient)(nil).GetStateFinalityCheckpoints), ctx, stateID)
}

// GetStateFork mocks base method.
func (m *MockBeaconClient) GetStateFork(ctx context.Context, stateID string) (*common.Fork, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStateFork", ctx, stateID)
	ret0, _ := ret[0].(*common.Fork)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetStateFork indicates an expected call of GetStateFork.
func (mr *MockBeaconClientMockRecorder) GetStateFork(ctx, stateID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStateFork", reflect.TypeOf((*MockBeaconClient)(nil).GetStateFork), ctx, stateID)
}

// GetStateRoot mocks base method.
func (m *MockBeaconClient) GetStateRoot(ctx context.Context, stateID string) (*common.Root, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStateRoot", ctx, stateID)
	ret0, _ := ret[0].(*common.Root)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetStateRoot indicates an expected call of GetStateRoot.
func (mr *MockBeaconClientMockRecorder) GetStateRoot(ctx, stateID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStateRoot", reflect.TypeOf((*MockBeaconClient)(nil).GetStateRoot), ctx, stateID)
}

// GetSyncCommittees mocks base method.
func (m *MockBeaconClient) GetSyncCommittees(ctx context.Context, stateID string, epoch *common.Epoch) (*types.SyncCommittees, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSyncCommittees", ctx, stateID, epoch)
	ret0, _ := ret[0].(*types.SyncCommittees)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSyncCommittees indicates an expected call of GetSyncCommittees.
func (mr *MockBeaconClientMockRecorder) GetSyncCommittees(ctx, stateID, epoch interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSyncCommittees", reflect.TypeOf((*MockBeaconClient)(nil).GetSyncCommittees), ctx, stateID, epoch)
}

// GetValidator mocks base method.
func (m *MockBeaconClient) GetValidator(ctx context.Context, stateID, validatorID string) (*types.Validator, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetValidator", ctx, stateID, validatorID)
	ret0, _ := ret[0].(*types.Validator)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetValidator indicates an expected call of GetValidator.
func (mr *MockBeaconClientMockRecorder) GetValidator(ctx, stateID, validatorID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetValidator", reflect.TypeOf((*MockBeaconClient)(nil).GetValidator), ctx, stateID, validatorID)
}

// GetValidatorBalances mocks base method.
func (m *MockBeaconClient) GetValidatorBalances(ctx context.Context, stateID string, validatorIDs []string) ([]*types.ValidatorBalance, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetValidatorBalances", ctx, stateID, validatorIDs)
	ret0, _ := ret[0].([]*types.ValidatorBalance)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetValidatorBalances indicates an expected call of GetValidatorBalances.
func (mr *MockBeaconClientMockRecorder) GetValidatorBalances(ctx, stateID, validatorIDs interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetValidatorBalances", reflect.TypeOf((*MockBeaconClient)(nil).GetValidatorBalances), ctx, stateID, validatorIDs)
}

// GetValidators mocks base method.
func (m *MockBeaconClient) GetValidators(ctx context.Context, stateID string, validatorIDs, statuses []string) ([]*types.Validator, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetValidators", ctx, stateID, validatorIDs, statuses)
	ret0, _ := ret[0].([]*types.Validator)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetValidators indicates an expected call of GetValidators.
func (mr *MockBeaconClientMockRecorder) GetValidators(ctx, stateID, validatorIDs, statuses interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetValidators", reflect.TypeOf((*MockBeaconClient)(nil).GetValidators), ctx, stateID, validatorIDs, statuses)
}

// GetVoluntaryExits mocks base method.
func (m *MockBeaconClient) GetVoluntaryExits(ctx context.Context) (phase0.VoluntaryExits, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetVoluntaryExits", ctx)
	ret0, _ := ret[0].(phase0.VoluntaryExits)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetVoluntaryExits indicates an expected call of GetVoluntaryExits.
func (mr *MockBeaconClientMockRecorder) GetVoluntaryExits(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetVoluntaryExits", reflect.TypeOf((*MockBeaconClient)(nil).GetVoluntaryExits), ctx)
}

// MockNodeClient is a mock of NodeClient interface.
type MockNodeClient struct {
	ctrl     *gomock.Controller
	recorder *MockNodeClientMockRecorder
}

// MockNodeClientMockRecorder is the mock recorder for MockNodeClient.
type MockNodeClientMockRecorder struct {
	mock *MockNodeClient
}

// NewMockNodeClient creates a new mock instance.
func NewMockNodeClient(ctrl *gomock.Controller) *MockNodeClient {
	mock := &MockNodeClient{ctrl: ctrl}
	mock.recorder = &MockNodeClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockNodeClient) EXPECT() *MockNodeClientMockRecorder {
	return m.recorder
}

// GetNodeVersion mocks base method.
func (m *MockNodeClient) GetNodeVersion(ctx context.Context) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNodeVersion", ctx)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNodeVersion indicates an expected call of GetNodeVersion.
func (mr *MockNodeClientMockRecorder) GetNodeVersion(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNodeVersion", reflect.TypeOf((*MockNodeClient)(nil).GetNodeVersion), ctx)
}

// MockConfigClient is a mock of ConfigClient interface.
type MockConfigClient struct {
	ctrl     *gomock.Controller
	recorder *MockConfigClientMockRecorder
}

// MockConfigClientMockRecorder is the mock recorder for MockConfigClient.
type MockConfigClientMockRecorder struct {
	mock *MockConfigClient
}

// NewMockConfigClient creates a new mock instance.
func NewMockConfigClient(ctrl *gomock.Controller) *MockConfigClient {
	mock := &MockConfigClient{ctrl: ctrl}
	mock.recorder = &MockConfigClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockConfigClient) EXPECT() *MockConfigClientMockRecorder {
	return m.recorder
}

// GetSpec mocks base method.
func (m *MockConfigClient) GetSpec(ctx context.Context) (*common.Spec, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSpec", ctx)
	ret0, _ := ret[0].(*common.Spec)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSpec indicates an expected call of GetSpec.
func (mr *MockConfigClientMockRecorder) GetSpec(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSpec", reflect.TypeOf((*MockConfigClient)(nil).GetSpec), ctx)
}
