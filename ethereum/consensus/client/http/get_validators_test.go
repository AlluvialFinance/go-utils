//go:build !integration

//nolint:revive // package name intentionally reflects domain, not directory name
package eth2http

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kilnfi/go-utils/ethereum/consensus/types"
	httptestutils "github.com/kilnfi/go-utils/net/http/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetValidators(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCli := httptestutils.NewMockSender(ctrl)
	c := NewClientFromClient(mockCli)

	t.Run("StatusOK", func(t *testing.T) { testGetValidatorsStatusOK(t, c, mockCli) })
}

func testGetValidatorsStatusOK(t *testing.T, c *Client, mockCli *httptestutils.MockSender) {
	t.Helper()
	req := httptestutils.NewGockRequest()
	req.Get("/eth/v1/beacon/states/test-state/validators").
		MatchParams(map[string]string{
			"status": "sA,sB",
			"id":     "vA,vB,vC",
		}).
		Reply(200).
		JSON([]byte(`{"data":[]}`))

	mockCli.EXPECT().Gock(req)

	vals, err := c.GetValidators(t.Context(), "test-state", []string{"vA", "vB", "vC"}, []string{"sA", "sB"})
	require.NoError(t, err)
	assert.Equal(
		t,
		[]*types.Validator{},
		vals,
	)
}
