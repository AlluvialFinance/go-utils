package eth2http

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	beaconcommon "github.com/protolambda/zrnt/eth2/beacon/common"
)

// GetSpec returns Ethreum 2.0 specifications configuration used on the node.
func (c *Client) GetSpec(ctx context.Context) (*beaconcommon.Spec, error) {
	rv, err := c.getSpec(ctx)
	if err != nil {
		c.logger.WithError(err).Errorf("GetSpec failed")
	}

	return rv, err
}

func (c *Client) getSpec(ctx context.Context) (*beaconcommon.Spec, error) {
	req, err := newGetSpecRequest(ctx)
	if err != nil {
		return nil, autorest.NewErrorWithError(err, "eth2http.Client", "GetSpec", nil, "Failure preparing request")
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, autorest.NewErrorWithError(err, "eth2http.Client", "GetSpec", resp, "Failure sending request")
	}

	result, err := inspectGetSpecResponse(resp)
	if err != nil {
		return nil, autorest.NewErrorWithError(err, "eth2http.Client", "GetSpec", resp, "Invalid response")
	}

	return result, nil
}

func newGetSpecRequest(ctx context.Context) (*http.Request, error) {
	return autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithPath("/eth/v1/config/spec"),
	).Prepare(newRequest(ctx))
}

// spec is an intermediary type allowing to properly unmarshal beacon config/spec responses

//nolint:revive,stylecheck // use uppercase as per protolambda/zrnt package
type spec struct {
	beaconcommon.Spec
	BASE_REWARD_FACTOR                         uint64 `json:"BASE_REWARD_FACTOR,string"`
	BYTES_PER_LOGS_BLOOM                       uint64 `json:"BYTES_PER_LOGS_BLOOM,string"`
	CHURN_LIMIT_QUOTIENT                       uint64 `json:"CHURN_LIMIT_QUOTIENT,string"`
	DEPOSIT_CHAIN_ID                           uint64 `json:"DEPOSIT_CHAIN_ID,string"`
	DEPOSIT_NETWORK_ID                         uint64 `json:"DEPOSIT_NETWORK_ID,string"`
	ETH1_FOLLOW_DISTANCE                       uint64 `json:"ETH1_FOLLOW_DISTANCE,string"`
	HISTORICAL_ROOTS_LIMIT                     uint64 `json:"HISTORICAL_ROOTS_LIMIT,string"`
	HYSTERESIS_DOWNWARD_MULTIPLIER             uint64 `json:"HYSTERESIS_DOWNWARD_MULTIPLIER,string"`
	HYSTERESIS_QUOTIENT                        uint64 `json:"HYSTERESIS_QUOTIENT,string"`
	HYSTERESIS_UPWARD_MULTIPLIER               uint64 `json:"HYSTERESIS_UPWARD_MULTIPLIER,string"`
	INACTIVITY_PENALTY_QUOTIENT                uint64 `json:"INACTIVITY_PENALTY_QUOTIENT,string"`
	INACTIVITY_PENALTY_QUOTIENT_ALTAIR         uint64 `json:"INACTIVITY_PENALTY_QUOTIENT_ALTAIR,string"`
	INACTIVITY_PENALTY_QUOTIENT_BELLATRIX      uint64 `json:"INACTIVITY_PENALTY_QUOTIENT_BELLATRIX,string"`
	INACTIVITY_SCORE_BIAS                      uint64 `json:"INACTIVITY_SCORE_BIAS,string"`
	INACTIVITY_SCORE_RECOVERY_RATE             uint64 `json:"INACTIVITY_SCORE_RECOVERY_RATE,string"`
	MAX_ATTESTATIONS                           uint64 `json:"MAX_ATTESTATIONS,string"`
	MAX_ATTESTER_SLASHINGS                     uint64 `json:"MAX_ATTESTER_SLASHINGS,string"`
	MAX_BYTES_PER_TRANSACTION                  uint64 `json:"MAX_BYTES_PER_TRANSACTION,string"`
	MAX_COMMITTEES_PER_SLOT                    uint64 `json:"MAX_COMMITTEES_PER_SLOT,string"`
	MAX_DEPOSITS                               uint64 `json:"MAX_DEPOSITS,string"`
	MAX_EXTRA_DATA_BYTES                       uint64 `json:"MAX_EXTRA_DATA_BYTES,string"`
	MAX_PROPOSER_SLASHINGS                     uint64 `json:"MAX_PROPOSER_SLASHINGS,string"`
	MAX_TRANSACTIONS_PER_PAYLOAD               uint64 `json:"MAX_TRANSACTIONS_PER_PAYLOAD,string"`
	MAX_VALIDATORS_PER_COMMITTEE               uint64 `json:"MAX_VALIDATORS_PER_COMMITTEE,string"`
	MAX_VOLUNTARY_EXITS                        uint64 `json:"MAX_VOLUNTARY_EXITS,string"`
	MIN_GENESIS_ACTIVE_VALIDATOR_COUNT         uint64 `json:"MIN_GENESIS_ACTIVE_VALIDATOR_COUNT,string"`
	MIN_PER_EPOCH_CHURN_LIMIT                  uint64 `json:"MIN_PER_EPOCH_CHURN_LIMIT,string"`
	MIN_SLASHING_PENALTY_QUOTIENT              uint64 `json:"MIN_SLASHING_PENALTY_QUOTIENT,string"`
	MIN_SLASHING_PENALTY_QUOTIENT_ALTAIR       uint64 `json:"MIN_SLASHING_PENALTY_QUOTIENT_ALTAIR,string"`
	MIN_SLASHING_PENALTY_QUOTIENT_BELLATRIX    uint64 `json:"MIN_SLASHING_PENALTY_QUOTIENT_BELLATRIX,string"`
	MIN_SYNC_COMMITTEE_PARTICIPANTS            uint64 `json:"MIN_SYNC_COMMITTEE_PARTICIPANTS,string"`
	PROPORTIONAL_SLASHING_MULTIPLIER           uint64 `json:"PROPORTIONAL_SLASHING_MULTIPLIER,string"`
	PROPORTIONAL_SLASHING_MULTIPLIER_ALTAIR    uint64 `json:"PROPORTIONAL_SLASHING_MULTIPLIER_ALTAIR,string"`
	PROPORTIONAL_SLASHING_MULTIPLIER_BELLATRIX uint64 `json:"PROPORTIONAL_SLASHING_MULTIPLIER_BELLATRIX,string"`
	PROPOSER_REWARD_QUOTIENT                   uint64 `json:"PROPOSER_REWARD_QUOTIENT,string"`
	PROPOSER_SCORE_BOOST                       uint64 `json:"PROPOSER_SCORE_BOOST,string"`
	SAFE_SLOTS_TO_UPDATE_JUSTIFIED             uint64 `json:"SAFE_SLOTS_TO_UPDATE_JUSTIFIED,string"`
	SECONDS_PER_ETH1_BLOCK                     uint64 `json:"SECONDS_PER_ETH1_BLOCK,string"`
	SHUFFLE_ROUND_COUNT                        uint8  `json:"SHUFFLE_ROUND_COUNT,string"`
	SYNC_COMMITTEE_SIZE                        uint64 `json:"SYNC_COMMITTEE_SIZE,string"`
	TARGET_COMMITTEE_SIZE                      uint64 `json:"TARGET_COMMITTEE_SIZE,string"`
	VALIDATOR_REGISTRY_LIMIT                   uint64 `json:"VALIDATOR_REGISTRY_LIMIT,string"`
	WHISTLEBLOWER_REWARD_QUOTIENT              uint64 `json:"WHISTLEBLOWER_REWARD_QUOTIENT,string"`
}

func (s *spec) toSpec() *beaconcommon.Spec {
	rvs := s.Spec
	rvs.BASE_REWARD_FACTOR = s.BASE_REWARD_FACTOR
	rvs.BYTES_PER_LOGS_BLOOM = s.BYTES_PER_LOGS_BLOOM
	rvs.CHURN_LIMIT_QUOTIENT = s.CHURN_LIMIT_QUOTIENT
	rvs.DEPOSIT_CHAIN_ID = s.DEPOSIT_CHAIN_ID
	rvs.DEPOSIT_NETWORK_ID = s.DEPOSIT_NETWORK_ID
	rvs.ETH1_FOLLOW_DISTANCE = s.ETH1_FOLLOW_DISTANCE
	rvs.HISTORICAL_ROOTS_LIMIT = s.HISTORICAL_ROOTS_LIMIT
	rvs.HYSTERESIS_DOWNWARD_MULTIPLIER = s.HYSTERESIS_DOWNWARD_MULTIPLIER
	rvs.HYSTERESIS_QUOTIENT = s.HYSTERESIS_QUOTIENT
	rvs.HYSTERESIS_UPWARD_MULTIPLIER = s.HYSTERESIS_UPWARD_MULTIPLIER
	rvs.INACTIVITY_PENALTY_QUOTIENT = s.INACTIVITY_PENALTY_QUOTIENT
	rvs.INACTIVITY_PENALTY_QUOTIENT_ALTAIR = s.INACTIVITY_PENALTY_QUOTIENT_ALTAIR
	rvs.INACTIVITY_PENALTY_QUOTIENT_BELLATRIX = s.INACTIVITY_PENALTY_QUOTIENT_BELLATRIX
	rvs.INACTIVITY_SCORE_BIAS = s.INACTIVITY_SCORE_BIAS
	rvs.INACTIVITY_SCORE_RECOVERY_RATE = s.INACTIVITY_SCORE_RECOVERY_RATE
	rvs.MAX_ATTESTATIONS = s.MAX_ATTESTATIONS
	rvs.MAX_ATTESTER_SLASHINGS = s.MAX_ATTESTER_SLASHINGS
	rvs.MAX_BYTES_PER_TRANSACTION = s.MAX_BYTES_PER_TRANSACTION
	rvs.MAX_COMMITTEES_PER_SLOT = s.MAX_COMMITTEES_PER_SLOT
	rvs.MAX_DEPOSITS = s.MAX_DEPOSITS
	rvs.MAX_EXTRA_DATA_BYTES = s.MAX_EXTRA_DATA_BYTES
	rvs.MAX_PROPOSER_SLASHINGS = s.MAX_PROPOSER_SLASHINGS
	rvs.MAX_TRANSACTIONS_PER_PAYLOAD = s.MAX_TRANSACTIONS_PER_PAYLOAD
	rvs.MAX_VALIDATORS_PER_COMMITTEE = s.MAX_VALIDATORS_PER_COMMITTEE
	rvs.MAX_VOLUNTARY_EXITS = s.MAX_VOLUNTARY_EXITS
	rvs.MIN_GENESIS_ACTIVE_VALIDATOR_COUNT = s.MIN_GENESIS_ACTIVE_VALIDATOR_COUNT
	rvs.MIN_PER_EPOCH_CHURN_LIMIT = s.MIN_PER_EPOCH_CHURN_LIMIT
	rvs.MIN_SLASHING_PENALTY_QUOTIENT = s.MIN_SLASHING_PENALTY_QUOTIENT
	rvs.MIN_SLASHING_PENALTY_QUOTIENT_ALTAIR = s.MIN_SLASHING_PENALTY_QUOTIENT_ALTAIR
	rvs.MIN_SLASHING_PENALTY_QUOTIENT_BELLATRIX = s.MIN_SLASHING_PENALTY_QUOTIENT_BELLATRIX
	rvs.MIN_SYNC_COMMITTEE_PARTICIPANTS = s.MIN_SYNC_COMMITTEE_PARTICIPANTS
	rvs.PROPORTIONAL_SLASHING_MULTIPLIER = s.PROPORTIONAL_SLASHING_MULTIPLIER
	rvs.PROPORTIONAL_SLASHING_MULTIPLIER_ALTAIR = s.PROPORTIONAL_SLASHING_MULTIPLIER_ALTAIR
	rvs.PROPORTIONAL_SLASHING_MULTIPLIER_BELLATRIX = s.PROPORTIONAL_SLASHING_MULTIPLIER_BELLATRIX
	rvs.PROPOSER_REWARD_QUOTIENT = s.PROPOSER_REWARD_QUOTIENT
	rvs.PROPOSER_SCORE_BOOST = s.PROPOSER_SCORE_BOOST
	rvs.SAFE_SLOTS_TO_UPDATE_JUSTIFIED = s.SAFE_SLOTS_TO_UPDATE_JUSTIFIED
	rvs.SECONDS_PER_ETH1_BLOCK = s.SECONDS_PER_ETH1_BLOCK
	rvs.SHUFFLE_ROUND_COUNT = s.SHUFFLE_ROUND_COUNT
	rvs.SYNC_COMMITTEE_SIZE = s.SYNC_COMMITTEE_SIZE
	rvs.TARGET_COMMITTEE_SIZE = s.TARGET_COMMITTEE_SIZE
	rvs.VALIDATOR_REGISTRY_LIMIT = s.VALIDATOR_REGISTRY_LIMIT
	rvs.WHISTLEBLOWER_REWARD_QUOTIENT = s.WHISTLEBLOWER_REWARD_QUOTIENT

	return &rvs
}

type getSpecResponseMsg struct {
	Data spec `json:"data"`
}

func inspectGetSpecResponse(resp *http.Response) (*beaconcommon.Spec, error) {
	msg := new(getSpecResponseMsg)
	err := inspectResponse(resp, msg)
	if err != nil {
		return nil, err
	}

	return msg.Data.toSpec(), nil
}
