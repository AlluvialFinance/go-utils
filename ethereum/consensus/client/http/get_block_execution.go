package eth2http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/ethereum/go-ethereum/common"
	"github.com/kilnfi/go-utils/ethereum/consensus/types"
	beaconcommon "github.com/protolambda/zrnt/eth2/beacon/common"
)

// ErrUnsupportedBlockVersion is returned for a beacon block whose fork version the decoder does not
// know, rather than decoding it into the wrong shape.
var ErrUnsupportedBlockVersion = errors.New("unsupported beacon block version")

// Fork versions whose block body embeds the execution payload (Bellatrix to Fulu), and the first one
// whose body carries a builder bid instead (EIP-7732).
var payloadBlockVersions = map[string]bool{
	"bellatrix": true,
	"capella":   true,
	"deneb":     true,
	"electra":   true,
	"fulu":      true,
}

const bidBlockVersion = "gloas"

// GetBlockExecution returns the execution-layer view of the block for given block id. The fork is
// taken from the response's version, never inferred from the block's shape.
func (c *Client) GetBlockExecution(ctx context.Context, blockID string) (*types.BlockExecution, error) {
	req, err := newGetBlockRequest(ctx, blockID)
	if err != nil {
		return nil, autorest.NewErrorWithError(err, "eth2http.Client", "GetBlockExecution", nil, "Failure preparing request")
	}

	resp, err := c.client.Do(req)
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, autorest.NewErrorWithError(err, "eth2http.Client", "GetBlockExecution", resp, "Failure sending request")
	}

	result, err := inspectGetBlockExecutionResponse(resp)
	if err != nil {
		return nil, autorest.NewErrorWithError(err, "eth2http.Client", "GetBlockExecution", resp, "Invalid response")
	}

	return result, nil
}

type getBlockExecutionResponseMsg struct {
	Version string          `json:"version"`
	Data    json.RawMessage `json:"data"`
}

// blockExecutionMsg decodes only the fields of a signed beacon block that BlockExecution exposes.
type blockExecutionMsg struct {
	Message blockExecutionMessageMsg `json:"message"`
}

type blockExecutionMessageMsg struct {
	Slot beaconcommon.Slot     `json:"slot"`
	Body blockExecutionBodyMsg `json:"body"`
}

type blockExecutionBodyMsg struct {
	ExecutionPayload          *executionPayloadRefMsg `json:"execution_payload"`
	SignedExecutionPayloadBid *signedPayloadBidRefMsg `json:"signed_execution_payload_bid"`
}

type executionPayloadRefMsg struct {
	BlockNumber uint64      `json:"block_number,string"`
	BlockHash   common.Hash `json:"block_hash"`
}

type signedPayloadBidRefMsg struct {
	Message payloadBidRefMsg `json:"message"`
}

type payloadBidRefMsg struct {
	ParentBlockHash common.Hash `json:"parent_block_hash"`
	BlockHash       common.Hash `json:"block_hash"`
}

func inspectGetBlockExecutionResponse(resp *http.Response) (*types.BlockExecution, error) {
	msg := new(getBlockExecutionResponseMsg)
	if err := inspectResponse(resp, msg); err != nil {
		return nil, err
	}

	return decodeBlockExecution(msg.Version, msg.Data)
}

func decodeBlockExecution(version string, data json.RawMessage) (*types.BlockExecution, error) {
	if !payloadBlockVersions[version] && version != bidBlockVersion {
		return nil, fmt.Errorf("%w %q", ErrUnsupportedBlockVersion, version)
	}

	block := new(blockExecutionMsg)
	if err := json.Unmarshal(data, block); err != nil {
		return nil, fmt.Errorf("failed to decode %s block: %w", version, err)
	}

	// A field missing for the version is an error, never a zero value. Only the genesis block has slot
	// 0, and its body is empty, so slot 0 means the field is missing.
	if block.Message.Slot == 0 {
		return nil, fmt.Errorf("%s block has no slot", version)
	}

	body := block.Message.Body
	result := &types.BlockExecution{Version: version, Slot: block.Message.Slot}

	if version == bidBlockVersion {
		bid := body.SignedExecutionPayloadBid
		if bid == nil || bid.Message.ParentBlockHash == (common.Hash{}) || bid.Message.BlockHash == (common.Hash{}) {
			return nil, fmt.Errorf("%s block at slot %d has no execution payload bid", version, block.Message.Slot)
		}
		result.Bid = &types.ExecutionPayloadBidRef{
			ParentBlockHash: bid.Message.ParentBlockHash,
			BlockHash:       bid.Message.BlockHash,
		}
		return result, nil
	}

	// No beacon block body carries the EL genesis block, so block number 0 means the field is missing.
	payload := body.ExecutionPayload
	if payload == nil || payload.BlockHash == (common.Hash{}) || payload.BlockNumber == 0 {
		return nil, fmt.Errorf("%s block at slot %d has no execution payload", version, block.Message.Slot)
	}
	result.Payload = &types.ExecutionPayloadRef{
		BlockNumber: payload.BlockNumber,
		BlockHash:   payload.BlockHash,
	}
	return result, nil
}
