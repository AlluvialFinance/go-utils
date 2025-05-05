//go:build !integration
// +build !integration

package eth2http

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	httptestutils "github.com/kilnfi/go-utils/net/http/testutils"
)

func TestGetPendingPartialWithdrawals(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCli := httptestutils.NewMockSender(ctrl)
	c := NewClientFromClient(mockCli)

	t.Run("StatusOK", func(t *testing.T) { testGetPendingPartialWithdrawalsStatusOK(t, c, mockCli) })
	t.Run("Status400", func(t *testing.T) { testGetPendingPartialWithdrawalsStatus400(t, c, mockCli) })
}

func testGetPendingPartialWithdrawalsStatusOK(t *testing.T, c *Client, mockCli *httptestutils.MockSender) {
	req := httptestutils.NewGockRequest()
	req.Get("/eth/v1/beacon/states/head/pending_partial_withdrawals").
		Reply(200).
		JSON([]byte(`{
			"data": [
				{
					"validator_index": 123,
					"address": "0x1234567890123456789012345678901234567890",
					"withdrawable_epoch": 123456
				}
			]
		}`))

	mockCli.EXPECT().Gock(req)

	withdrawals, err := c.GetPendingPartialWithdrawals(context.Background())

	require.NoError(t, err)
	assert.Equal(t, 1, len(withdrawals))
	assert.Equal(t, uint64(123), withdrawals[0].ValidatorIndex)
	assert.Equal(t, common.HexToAddress("0x1234567890123456789012345678901234567890"), withdrawals[0].Address)
	assert.Equal(t, uint64(123456), withdrawals[0].WithdrawableEpoch)
}

func testGetPendingPartialWithdrawalsStatus400(t *testing.T, c *Client, mockCli *httptestutils.MockSender) {
	req := httptestutils.NewGockRequest()
	req.Get("/eth/v1/beacon/states/head/pending_partial_withdrawals").
		Reply(400)

	mockCli.EXPECT().Gock(req)

	_, err := c.GetPendingPartialWithdrawals(context.Background())

	require.Error(t, err)
}
