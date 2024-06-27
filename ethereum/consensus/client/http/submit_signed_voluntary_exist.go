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
	Epoch          beaconcommon.Epoch `json:"epoch"`
	ValidatorIndex string             `json:"validator_index"`
}

func (c *Client) SignedVoluntaryExits(ctx context.Context, epoch beaconcommon.Epoch, validatorIdx uint64, signature string) (string, error) {
	resp, err := c.signedVoluntaryExits(ctx, epoch, validatorIdx, signature)
	if err != nil {
		c.logger.WithError(err).Errorf("SignedVoluntaryExits failed")
	}

	return resp, err
}

func (c *Client) signedVoluntaryExits(ctx context.Context, epoch beaconcommon.Epoch, validatorIdx uint64, signature string) (string, error) {
	reqBody := newSignedVoluntaryExit(epoch, validatorIdx, signature)
	req, err := newSignedVoluntaryExitsRequest(ctx, reqBody)
	if err != nil {
		return "", autorest.NewErrorWithError(err, "eth2http.Client", "SignedVoluntaryExits", nil, "Failure preparing request")
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return "", autorest.NewErrorWithError(err, "eth2http.Client", "SignedVoluntaryExits", resp, "Failure sending request")
	}

	ve, err := inspectSignedVoluntaryExitsResponse(resp)
	if err != nil {
		return "", autorest.NewErrorWithError(err, "eth2http.Client", "SignedVoluntaryExits", resp, "Invalid response")
	}

	return ve, nil
}

func newSignedVoluntaryExitsRequest(ctx context.Context, signedVoluntaryExit *signedVoluntaryExit) (*http.Request, error) {
	return autorest.CreatePreparer(
		autorest.AsPost(),
		autorest.WithPath("/eth/v1/beacon/pool/voluntary_exits"),
		autorest.AsJSON(),
		autorest.WithJSON(signedVoluntaryExit),
	).Prepare(newRequest(ctx))
}

func newSignedVoluntaryExit(epoch beaconcommon.Epoch, validatorIdx uint64, signature string) *signedVoluntaryExit {
	return &signedVoluntaryExit{
		Message: message{
			Epoch:          epoch,
			ValidatorIndex: strconv.Itoa(int(validatorIdx)),
		},
		Signature: signature,
	}
}

func inspectSignedVoluntaryExitsResponse(resp *http.Response) (string, error) {
	var msg string
	err := inspectResponse(resp, &msg)
	if err != nil {
		return "", err
	}

	return msg, nil
}
