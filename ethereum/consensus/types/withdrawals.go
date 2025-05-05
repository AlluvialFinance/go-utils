package types

import "github.com/ethereum/go-ethereum/common"

type PendingPartialWithdrawal struct {
	ValidatorIndex    uint64         `json:"validator_index"`
	Address           common.Address `json:"address"`
	WithdrawableEpoch uint64         `json:"withdrawable_epoch"`
}

type WithdrawalsResponse struct {
	Version             string                      `json:"version"`
	ExecutionOptimistic bool                        `json:"execution_optimistic"`
	Finalized           string                      `json:"finalized"`
	Data                []*PendingPartialWithdrawal `json:"data"`
}
