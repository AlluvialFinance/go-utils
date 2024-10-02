//go:build !integration
// +build !integration

package eth2http

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	httptestutils "github.com/kilnfi/go-utils/net/http/testutils"
	"github.com/protolambda/zrnt/eth2/configs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetSpec(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCli := httptestutils.NewMockSender(ctrl)
	c := NewClientFromClient(mockCli)

	t.Run("StatusOK", func(t *testing.T) { testGetSpecStatusOK(t, c, mockCli) })
}

func testGetSpecStatusOK(t *testing.T, c *Client, mockCli *httptestutils.MockSender) {
	req := httptestutils.NewGockRequest()
	req.Get("/eth/v1/config/spec").
		Reply(200).
		JSON([]byte(`{"data":{"ALTAIR_FORK_EPOCH":"74240","ALTAIR_FORK_VERSION":"0x01000000","BASE_REWARD_FACTOR":"64","BELLATRIX_FORK_EPOCH":"18446744073709551615","BELLATRIX_FORK_VERSION":"0x02000000","BLS_WITHDRAWAL_PREFIX":"0x00","BYTES_PER_LOGS_BLOOM":"256","CHURN_LIMIT_QUOTIENT":"65536","CONFIG_NAME":"mainnet","DEPOSIT_CHAIN_ID":"1","DEPOSIT_CONTRACT_ADDRESS":"0x00000000219ab540356cBB839Cbe05303d7705Fa","DEPOSIT_NETWORK_ID":"1","DOMAIN_AGGREGATE_AND_PROOF":"0x06000000","DOMAIN_BEACON_ATTESTER":"0x01000000","DOMAIN_BEACON_PROPOSER":"0x00000000","DOMAIN_CONTRIBUTION_AND_PROOF":"0x09000000","DOMAIN_DEPOSIT":"0x03000000","DOMAIN_RANDAO":"0x02000000","DOMAIN_SELECTION_PROOF":"0x05000000","DOMAIN_SYNC_COMMITTEE":"0x07000000","DOMAIN_SYNC_COMMITTEE_SELECTION_PROOF":"0x08000000","DOMAIN_VOLUNTARY_EXIT":"0x04000000","EFFECTIVE_BALANCE_INCREMENT":"1000000000","EJECTION_BALANCE":"16000000000","EPOCHS_PER_ETH1_VOTING_PERIOD":"64","EPOCHS_PER_HISTORICAL_VECTOR":"65536","EPOCHS_PER_RANDOM_SUBNET_SUBSCRIPTION":"256","EPOCHS_PER_SLASHINGS_VECTOR":"8192","EPOCHS_PER_SYNC_COMMITTEE_PERIOD":"256","ETH1_FOLLOW_DISTANCE":"2048","GENESIS_DELAY":"604800","GENESIS_FORK_VERSION":"0x00000000","HISTORICAL_ROOTS_LIMIT":"16777216","HYSTERESIS_DOWNWARD_MULTIPLIER":"1","HYSTERESIS_QUOTIENT":"4","HYSTERESIS_UPWARD_MULTIPLIER":"5","INACTIVITY_PENALTY_QUOTIENT":"67108864","INACTIVITY_PENALTY_QUOTIENT_ALTAIR":"50331648","INACTIVITY_PENALTY_QUOTIENT_BELLATRIX":"16777216","INACTIVITY_PENALTY_QUOTIENT_MERGE":"16777216","INACTIVITY_SCORE_BIAS":"4","INACTIVITY_SCORE_RECOVERY_RATE":"16","MAX_ATTESTATIONS":"128","MAX_ATTESTER_SLASHINGS":"2","MAX_BYTES_PER_TRANSACTION":"1073741824","MAX_COMMITTEES_PER_SLOT":"64","MAX_DEPOSITS":"16","MAX_EFFECTIVE_BALANCE":"32000000000","MAX_EXTRA_DATA_BYTES":"32","MAX_PROPOSER_SLASHINGS":"16","MAX_SEED_LOOKAHEAD":"4","MAX_TRANSACTIONS_PER_PAYLOAD":"1048576","MAX_VALIDATORS_PER_COMMITTEE":"2048","MAX_VOLUNTARY_EXITS":"16","MERGE_FORK_EPOCH":"18446744073709551615","MERGE_FORK_VERSION":"0x02000000","MIN_ATTESTATION_INCLUSION_DELAY":"1","MIN_DEPOSIT_AMOUNT":"1000000000","MIN_EPOCHS_TO_INACTIVITY_PENALTY":"4","MIN_GENESIS_ACTIVE_VALIDATOR_COUNT":"16384","MIN_GENESIS_TIME":"1606824000","MIN_PER_EPOCH_CHURN_LIMIT":"4","MIN_SEED_LOOKAHEAD":"1","MIN_SLASHING_PENALTY_QUOTIENT":"128","MIN_SLASHING_PENALTY_QUOTIENT_ALTAIR":"64","MIN_SLASHING_PENALTY_QUOTIENT_BELLATRIX":"32","MIN_SLASHING_PENALTY_QUOTIENT_MERGE":"32","MIN_SYNC_COMMITTEE_PARTICIPANTS":"1","MIN_VALIDATOR_WITHDRAWABILITY_DELAY":"256","PRESET_BASE":"mainnet","PROPORTIONAL_SLASHING_MULTIPLIER":"1","PROPORTIONAL_SLASHING_MULTIPLIER_ALTAIR":"2","PROPORTIONAL_SLASHING_MULTIPLIER_BELLATRIX":"3","PROPORTIONAL_SLASHING_MULTIPLIER_MERGE":"3","PROPOSER_REWARD_QUOTIENT":"8","PROPOSER_SCORE_BOOST":"70","RANDOM_SUBNETS_PER_VALIDATOR":"1","SAFE_SLOTS_TO_UPDATE_JUSTIFIED":"8","SECONDS_PER_ETH1_BLOCK":"14","SECONDS_PER_SLOT":"12","SHARD_COMMITTEE_PERIOD":"256","SHARDING_FORK_EPOCH":"18446744073709551615","SHARDING_FORK_VERSION":"0x03000000","SHUFFLE_ROUND_COUNT":"90","SLOTS_PER_EPOCH":"32","SLOTS_PER_HISTORICAL_ROOT":"8192","SYNC_COMMITTEE_SIZE":"512","SYNC_COMMITTEE_SUBNET_COUNT":"4","TARGET_AGGREGATORS_PER_COMMITTEE":"16","TARGET_AGGREGATORS_PER_SYNC_SUBCOMMITTEE":"16","TARGET_COMMITTEE_SIZE":"128","TERMINAL_BLOCK_HASH":"0x0000000000000000000000000000000000000000000000000000000000000000","TERMINAL_BLOCK_HASH_ACTIVATION_EPOCH":"18446744073709551615","TERMINAL_TOTAL_DIFFICULTY":"115792089237316195423570985008687907853269984665640564039457584007913129638912","VALIDATOR_REGISTRY_LIMIT":"1099511627776","WHISTLEBLOWER_REWARD_QUOTIENT":"512"}}`))

	mockCli.EXPECT().Gock(req)

	spec, err := c.GetSpec(context.Background())

	require.NoError(t, err)
	assert.Equal(
		t,
		configs.Mainnet.Phase0Preset,
		spec.Phase0Preset,
	)
	assert.Equal(
		t,
		configs.Mainnet.AltairPreset,
		spec.AltairPreset,
	)
	assert.Equal(
		t,
		configs.Mainnet.BellatrixPreset,
		spec.BellatrixPreset,
	)

	// TODO: add assert.Equal for spec.ShardPreset
	assert.Equal(
		t,
		configs.Mainnet.Config,
		spec.Config,
	)
}

// Phase0Preset    `json:",inline" yaml:",inline"`
// 	AltairPreset    `json:",inline" yaml:",inline"`
// 	BellatrixPreset `json:",inline" yaml:",inline"`
// 	ShardingPreset  `json:",inline" yaml:",inline"`
// 	Config          `json:",inline" yaml:",inline"`
// 	Setup           `json:",inline" yaml:",inline"`
