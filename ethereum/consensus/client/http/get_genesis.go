package eth2http

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"

	"github.com/kilnfi/go-utils/ethereum/consensus/types"
)

// GetGenesis returns genesis block
func (c *Client) GetGenesis(ctx context.Context) (*types.Genesis, error) {
	return c.getGenesis(ctx)
}

func (c *Client) getGenesis(ctx context.Context) (*types.Genesis, error) {
	req, err := newGetGenesisRequest(ctx)
	if err != nil {
		return nil, autorest.NewErrorWithError(err, "eth2http.Client", "GetGenesis", nil, "Failure preparing request")
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, autorest.NewErrorWithError(err, "eth2http.Client", "GetGenesis", resp, "Failure sending request")
	}

	result, err := inspectGetGenesisResponse(resp)
	if err != nil {
		return nil, autorest.NewErrorWithError(err, "eth2http.Client", "GetGenesis", resp, "Invalid response")
	}

	return result, nil
}

func newGetGenesisRequest(ctx context.Context) (*http.Request, error) {
	return autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithPath("eth/v1/beacon/genesis"),
	).Prepare(newRequest(ctx))
}

type getGenesisResponseMsg struct {
	Data *types.Genesis `json:"data"`
}

func inspectGetGenesisResponse(resp *http.Response) (*types.Genesis, error) {
	msg := new(getGenesisResponseMsg)
	err := inspectResponse(resp, msg)
	if err != nil {
		return nil, err
	}

	return msg.Data, nil
}
