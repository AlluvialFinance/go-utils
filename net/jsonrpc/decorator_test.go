//go:build !integration

package jsonrpc_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kilnfi/go-utils/net/jsonrpc"
	jsonrpctestutils "github.com/kilnfi/go-utils/net/jsonrpc/testutils"
	"github.com/stretchr/testify/require"
)

func TestWithVersion(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCli := jsonrpctestutils.NewMockClient(ctrl)
	c := jsonrpc.WithVersion("2.0")(mockCli)

	mockCli.EXPECT().Call(
		gomock.Any(),
		jsonrpctestutils.HasVersion("2.0"),
		gomock.Any())
	err := c.Call(t.Context(), &jsonrpc.Request{}, nil)
	require.NoError(t, err)
}

func TestWithIncrementalID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCli := jsonrpctestutils.NewMockClient(ctrl)
	c := jsonrpc.WithIncrementalID()(mockCli)

	mockCli.EXPECT().Call(
		gomock.Any(),
		jsonrpctestutils.HasID(uint32(0)),
		gomock.Any())
	err := c.Call(t.Context(), &jsonrpc.Request{}, nil)
	require.NoError(t, err)

	mockCli.EXPECT().Call(
		gomock.Any(),
		jsonrpctestutils.HasID(uint32(1)),
		gomock.Any())
	err = c.Call(t.Context(), &jsonrpc.Request{}, nil)
	require.NoError(t, err)
}
