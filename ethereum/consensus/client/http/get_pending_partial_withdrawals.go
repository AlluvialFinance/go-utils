//nolint:revive // package name intentionally reflects domain, not directory name
package eth2http

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/kilnfi/go-utils/ethereum/consensus/types"
)

func (c *Client) GetPendingPartialWithdrawals(ctx context.Context, stateID string) ([]*types.PendingPartialWithdrawal, error) {
	rv, err := c.getPendingPartialWithdrawals(ctx, stateID)
	if err != nil {
		c.logger.WithError(err).Errorf("GetPendingPartialWithdrawals failed")
	}

	return rv, err
}

func (c *Client) getPendingPartialWithdrawals(ctx context.Context, stateID string) ([]*types.PendingPartialWithdrawal, error) {
	req, err := newGetPendingPartialWithdrawalsRequest(ctx, stateID)
	if err != nil {
		return nil, autorest.NewErrorWithError(err, "eth2http.Client", "GetPendingPartialWithdrawals", nil, "Failure preparing request")
	}

	resp, err := c.client.Do(req) //nolint:bodyclose // response body is closed by inspect*Response via autorest.ByClosing
	if err != nil {
		return nil, autorest.NewErrorWithError(err, "eth2http.Client", "GetPendingPartialWithdrawals", resp, "Failure sending request")
	}

	result, err := inspectGetPendingPartialWithdrawalsResponse(resp)
	if err != nil {
		return nil, autorest.NewErrorWithError(err, "eth2http.Client", "GetPendingPartialWithdrawals", resp, "Invalid response")
	}

	return result, nil
}

func newGetPendingPartialWithdrawalsRequest(ctx context.Context, stateID string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"stateID": autorest.Encode("path", stateID),
	}

	return autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithPathParameters("eth/v1/beacon/states/{stateID}/pending_partial_withdrawals", pathParameters),
	).Prepare(newRequest(ctx))
}

type getPendingPartialWithdrawalsResponseMsg struct {
	Data []*types.PendingPartialWithdrawal `json:"data"`
}

func inspectGetPendingPartialWithdrawalsResponse(resp *http.Response) ([]*types.PendingPartialWithdrawal, error) {
	msg := new(getPendingPartialWithdrawalsResponseMsg)
	err := inspectResponse(resp, msg)
	if err != nil {
		return nil, err
	}

	return msg.Data, nil
}
