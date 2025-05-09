package types

type PendingPartialWithdrawal struct {
	ValidatorIndex    string `json:"validator_index"`
	Amount            string `json:"amount"` // Amount represents the amount to be withdrawn, specified in Gwei.
	WithdrawableEpoch string `json:"withdrawable_epoch"`
}
