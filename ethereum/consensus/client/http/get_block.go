package eth2http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/protolambda/zrnt/eth2/beacon/bellatrix"
)

// GetBlock returns block details for given block id.
func (c *Client) GetBlock(ctx context.Context, blockID string) (*bellatrix.SignedBeaconBlock, error) {
	return c.getBlock(ctx, blockID)
}

func (c *Client) getBlock(ctx context.Context, blockID string) (*bellatrix.SignedBeaconBlock, error) {
	req, err := newGetBlockRequest(ctx, blockID)
	if err != nil {
		return nil, autorest.NewErrorWithError(err, "eth2http.Client", "GetBlock", nil, "Failure preparing request")
	}

	resp, err := c.client.Do(req)
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, autorest.NewErrorWithError(err, "eth2http.Client", "GetBlock", resp, "Failure sending request")
	}

	result, err := inspectGetBlockResponse(resp)
	if err != nil {
		return nil, autorest.NewErrorWithError(err, "eth2http.Client", "GetBlock", resp, "Invalid response")
	}

	return result, nil
}

func newGetBlockRequest(ctx context.Context, blockID string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"blockID": autorest.Encode("path", blockID),
	}

	return autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithPathParameters("eth/v2/beacon/blocks/{blockID}", pathParameters),
	).Prepare(newRequest(ctx))
}

type getBlockResponseMsg struct {
	Version string                       `json:"version"`
	Data    *bellatrix.SignedBeaconBlock `json:"data"`
}

func inspectGetBlockResponse(resp *http.Response) (*bellatrix.SignedBeaconBlock, error) {
	msg := new(getBlockResponseMsg)
	err := inspectResponse(resp, msg)
	if err != nil {
		return nil, err
	}

	// From Gloas on the body no longer embeds the execution payload, so decoding it as a Bellatrix
	// block would succeed with a zero payload. Use GetBlockExecution for those blocks.
	if !bellatrixCompatibleVersions[msg.Version] {
		return nil, fmt.Errorf("%w %q, use GetBlockExecution", ErrUnsupportedBlockVersion, msg.Version)
	}

	return msg.Data, nil
}

// bellatrixCompatibleVersions are the versions GetBlock can decode into a Bellatrix block. An empty
// version is accepted for nodes that do not report it.
var bellatrixCompatibleVersions = map[string]bool{
	"":          true,
	"phase0":    true,
	"altair":    true,
	"bellatrix": true,
	"capella":   true,
	"deneb":     true,
	"electra":   true,
	"fulu":      true,
}
