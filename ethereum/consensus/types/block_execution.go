package types

import (
	"github.com/ethereum/go-ethereum/common"
	beaconcommon "github.com/protolambda/zrnt/eth2/beacon/common"
)

// BlockExecution is the execution-layer view of a signed beacon block. Which of Payload and Bid is
// set depends on the fork the block belongs to, given by Version.
type BlockExecution struct {
	Version string
	Slot    beaconcommon.Slot

	// Payload is the execution payload embedded in the block. Set before Gloas.
	Payload *ExecutionPayloadRef

	// Bid is the builder bid the block commits to. Set from Gloas on (EIP-7732), where the block no
	// longer embeds its payload.
	Bid *ExecutionPayloadBidRef
}

// ExecutionPayloadRef identifies the execution payload embedded in a pre-Gloas beacon block.
type ExecutionPayloadRef struct {
	BlockNumber uint64
	BlockHash   common.Hash
}

// ExecutionPayloadBidRef holds the execution block hashes a Gloas beacon block commits to.
type ExecutionPayloadBidRef struct {
	// ParentBlockHash is the execution block the beacon block builds on. Once the beacon block is
	// processed it equals state.latest_block_hash, the execution state its post-state reflects.
	ParentBlockHash common.Hash

	// BlockHash is the execution payload the bid commits to. It is revealed separately and may never
	// become canonical.
	BlockHash common.Hash
}
