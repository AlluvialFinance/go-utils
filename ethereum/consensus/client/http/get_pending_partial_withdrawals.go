package eth2http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/kilnfi/go-utils/ethereum/consensus/types"
)

func (c *Client) GetPendingPartialWithdrawals(ctx context.Context) ([]*types.PendingPartialWithdrawal, error) {
	rv, err := c.getPendingPartialWithdrawals(ctx)
	if err != nil {
		c.logger.WithError(err).Errorf("GetPendingPartialWithdrawals failed")
	}

	return rv, err
}

func (c *Client) getPendingPartialWithdrawals(ctx context.Context) ([]*types.PendingPartialWithdrawal, error) {

	req, err := newGetPendingPartialWithdrawalsRequest(ctx)
	if err != nil {
		return nil, autorest.NewErrorWithError(err, "eth2http.Client", "GetPendingPartialWithdrawals", nil, "Failure preparing request")
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, autorest.NewErrorWithError(err, "eth2http.Client", "GetPendingPartialWithdrawals", resp, "Failure sending request")
	}

	result, err := inspectGetPendingPartialWithdrawalsResponse(resp)
	if err != nil {
		return nil, autorest.NewErrorWithError(err, "eth2http.Client", "GetPendingPartialWithdrawals", resp, "Invalid response")
	}

	return result, nil

}

func newGetPendingPartialWithdrawalsRequest(ctx context.Context) (*http.Request, error) {
	return autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithPath("eth/v1/beacon/states/head/pending_partial_withdrawals"),
	).Prepare(newRequest(ctx))
}

func inspectGetPendingPartialWithdrawalsResponse(resp *http.Response) ([]*types.PendingPartialWithdrawal, error) {
	var pendingPartialWithdrawals types.WithdrawalsResponse
	if err := json.NewDecoder(resp.Body).Decode(&pendingPartialWithdrawals); err != nil {
		return nil, fmt.Errorf("failed to decode pending partial withdrawals: %w", err)
	}

	return pendingPartialWithdrawals.Data, nil
}
