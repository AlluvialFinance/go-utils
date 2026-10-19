//go:build !integration

package eth2http

import (
	"errors"
	"os"
	"testing"

	"github.com/Azure/go-autorest/autorest"
	"github.com/ethereum/go-ethereum/common"
	"github.com/golang/mock/gomock"
	"github.com/kilnfi/go-utils/ethereum/consensus/types"
	httptestutils "github.com/kilnfi/go-utils/net/http/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetBlockExecution(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCli := httptestutils.NewMockSender(ctrl)
	c := NewClientFromClient(mockCli)

	reply := func(status int, body []byte) {
		req := httptestutils.NewGockRequest()
		req.Get("/eth/v2/beacon/blocks/test-block").Reply(status).JSON(body)
		mockCli.EXPECT().Gock(req)
	}

	t.Run("fulu block from mainnet carries its payload", func(t *testing.T) {
		reply(200, readFixture(t, "block_fulu_mainnet.json"))

		block, err := c.GetBlockExecution(t.Context(), "test-block")
		require.NoError(t, err)
		assert.Equal(t, &types.BlockExecution{
			Version: "fulu",
			Slot:    15465088,
			Payload: &types.ExecutionPayloadRef{
				BlockNumber: 26225974,
				BlockHash:   common.HexToHash("0x6f85838e9e8b977d6edc8767501495bfef04199d8e7380fa2ae51e99a6606f7b"),
			},
		}, block)
	})

	t.Run("gloas block from sepolia carries the builder bid", func(t *testing.T) {
		reply(200, readFixture(t, "block_gloas_sepolia.json"))

		block, err := c.GetBlockExecution(t.Context(), "test-block")
		require.NoError(t, err)
		assert.Equal(t, &types.BlockExecution{
			Version: "gloas",
			Slot:    11389280,
			Bid: &types.ExecutionPayloadBidRef{
				ParentBlockHash: common.HexToHash("0xf6480e184d7d4af8219eff2c4915af3decfcec5dee8f84f3c4cb31eb5d680901"),
				BlockHash:       common.HexToHash("0xc77cbaf2f30a58017067a0bd8b947b4411e3e0181358416396c13dcde738100b"),
			},
		}, block)
	})

	t.Run("missing block surfaces the 404", func(t *testing.T) {
		reply(404, []byte(`{"code":404,"message":"NOT_FOUND: beacon block at slot 123"}`))

		_, err := c.GetBlockExecution(t.Context(), "test-block")
		var detailed autorest.DetailedError
		require.ErrorAs(t, err, &detailed)
		assert.Equal(t, 404, detailed.StatusCode)
	})
}

func TestDecodeBlockExecution(t *testing.T) {
	const (
		hashA = "0x1111111111111111111111111111111111111111111111111111111111111111"
		hashB = "0x2222222222222222222222222222222222222222222222222222222222222222"
	)
	payloadBody := `{"message":{"slot":"100","body":{"execution_payload":{"block_number":"42","block_hash":"` + hashA + `"}}}}`
	bidBody := `{"message":{"slot":"100","body":{"signed_execution_payload_bid":{"message":{"parent_block_hash":"` + hashA + `","block_hash":"` + hashB + `"}}}}}`

	for _, version := range []string{"bellatrix", "capella", "deneb", "electra", "fulu"} {
		t.Run(version+" reads the embedded payload", func(t *testing.T) {
			block, err := decodeBlockExecution(version, []byte(payloadBody))
			require.NoError(t, err)
			assert.Equal(t, &types.ExecutionPayloadRef{BlockNumber: 42, BlockHash: common.HexToHash(hashA)}, block.Payload)
			assert.Nil(t, block.Bid)
		})
	}

	t.Run("gloas reads the bid", func(t *testing.T) {
		block, err := decodeBlockExecution("gloas", []byte(bidBody))
		require.NoError(t, err)
		assert.Equal(t, &types.ExecutionPayloadBidRef{ParentBlockHash: common.HexToHash(hashA), BlockHash: common.HexToHash(hashB)}, block.Bid)
		assert.Nil(t, block.Payload)
	})

	errorCases := []struct {
		name    string
		version string
		data    string
		want    error
	}{
		{"unknown future fork", "heze", bidBody, ErrUnsupportedBlockVersion},
		{"pre-merge fork has no payload", "altair", payloadBody, ErrUnsupportedBlockVersion},
		{"missing version", "", payloadBody, ErrUnsupportedBlockVersion},
		// The version decides the shape: a gloas block without a bid is an error, not a zero value.
		{"gloas without bid", "gloas", payloadBody, nil},
		{"fulu without payload", "fulu", bidBody, nil},
		{"gloas bid without parent hash", "gloas", `{"message":{"slot":"1","body":{"signed_execution_payload_bid":{"message":{"block_hash":"` + hashB + `"}}}}}`, nil},
		{"malformed body", "fulu", `{"message":{"slot":"not-a-number"}}`, nil},
	}
	for _, tt := range errorCases {
		t.Run(tt.name, func(t *testing.T) {
			_, err := decodeBlockExecution(tt.version, []byte(tt.data))
			require.Error(t, err)
			if tt.want != nil {
				assert.ErrorIs(t, err, tt.want)
			}
		})
	}
}

func TestGetBlockRejectsGloas(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCli := httptestutils.NewMockSender(ctrl)
	c := NewClientFromClient(mockCli)

	reply := func(body []byte) {
		req := httptestutils.NewGockRequest()
		req.Get("/eth/v2/beacon/blocks/test-block").Reply(200).JSON(body)
		mockCli.EXPECT().Gock(req)
	}

	t.Run("fulu block still decodes", func(t *testing.T) {
		reply(readFixture(t, "block_fulu_mainnet.json"))

		block, err := c.GetBlock(t.Context(), "test-block")
		require.NoError(t, err)
		assert.Equal(t, uint64(26225974), uint64(block.Message.Body.ExecutionPayload.BlockNumber))
	})

	t.Run("gloas block is an error instead of a zero payload", func(t *testing.T) {
		reply(readFixture(t, "block_gloas_sepolia.json"))

		_, err := c.GetBlock(t.Context(), "test-block")
		assert.True(t, errors.Is(err, ErrUnsupportedBlockVersion), "got %v", err)
	})
}

func readFixture(t *testing.T, name string) []byte {
	t.Helper()
	data, err := os.ReadFile("testdata/" + name)
	require.NoError(t, err)
	return data
}
