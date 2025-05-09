package types

type PendingPartialWithdrawal struct {
	ValidatorIndex    string `json:"validator_index"`
	Amount            string `json:"amount"`
	WithdrawableEpoch string `json:"withdrawable_epoch"`
}
