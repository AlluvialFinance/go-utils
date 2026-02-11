package eth2http

import (
	"context"
	"net/http"
	"strconv"

	"github.com/Azure/go-autorest/autorest"
	beaconcommon "github.com/protolambda/zrnt/eth2/beacon/common"
)

type signedVoluntaryExit struct {
	Message   message `json:"message"`
	Signature string  `json:"signature"`
}

type message struct {
	Epoch          string `json:"epoch"`
	ValidatorIndex string `json:"validator_index"`
}

// SubmitSignedVoluntaryExit submits a signed voluntary exit to the beacon node.
func (c *Client) SubmitSignedVoluntaryExit(ctx context.Context, epoch beaconcommon.Epoch, validatorIdx uint64, signature string) (string, error) {
	resp, err := c.submitSignedVoluntaryExit(ctx, epoch, validatorIdx, signature)
	if err != nil {
		return "", err
	}

	return resp.Message, nil
}

func (c *Client) submitSignedVoluntaryExit(ctx context.Context, epoch beaconcommon.Epoch, validatorIdx uint64, signature string) (*SubmitSignedVoluntaryExitResponse, error) {
	reqBody := newSignedVoluntaryExit(epoch, validatorIdx, signature)
	req, err := newSignedVoluntaryExitsRequest(ctx, reqBody)
	if err != nil {
		return nil, autorest.NewErrorWithError(err, "eth2http.Client", "SubmitSignedVoluntaryExit", nil, "Failure preparing request")
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, autorest.NewErrorWithError(err, "eth2http.Client", "SubmitSignedVoluntaryExit", resp, "Failure sending request")
	}

	ve, err := inspectSubmitSignedVoluntaryExitResponse(resp)
	if err != nil {
		return nil, autorest.NewErrorWithError(err, "eth2http.Client", "SubmitSignedVoluntaryExit", resp, "Invalid response")
	}

	return ve, nil
}

func newSignedVoluntaryExitsRequest(ctx context.Context, signedVoluntaryExit *signedVoluntaryExit) (*http.Request, error) {
	return autorest.CreatePreparer(
		autorest.AsPost(),
		autorest.AsJSON(),
		autorest.WithJSON(signedVoluntaryExit),
		autorest.WithPath("/eth/v1/beacon/pool/voluntary_exits"),
	).Prepare(newRequest(ctx))
}

func newSignedVoluntaryExit(epoch beaconcommon.Epoch, validatorIdx uint64, signature string) *signedVoluntaryExit {
	return &signedVoluntaryExit{
		Message: message{
			Epoch:          strconv.Itoa(int(epoch)),
			ValidatorIndex: strconv.Itoa(int(validatorIdx)),
		},
		Signature: signature,
	}
}

type SubmitSignedVoluntaryExitResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func inspectSubmitSignedVoluntaryExitResponse(resp *http.Response) (*SubmitSignedVoluntaryExitResponse, error) {
	msg := new(SubmitSignedVoluntaryExitResponse)
	err := inspectResponse(resp, msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}
