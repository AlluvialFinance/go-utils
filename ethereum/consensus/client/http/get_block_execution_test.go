//go:build !integration

package eth2http

import (
	"encoding/json"
	"errors"
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/Azure/go-autorest/autorest"
	"github.com/ethereum/go-ethereum/common"
	"github.com/golang/mock/gomock"
	"github.com/kilnfi/go-utils/ethereum/consensus/types"
	httptestutils "github.com/kilnfi/go-utils/net/http/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	hashA    = "0x1111111111111111111111111111111111111111111111111111111111111111"
	hashB    = "0x2222222222222222222222222222222222222222222222222222222222222222"
	zeroHash = "0x0000000000000000000000000000000000000000000000000000000000000000"
)

// payloadVersions are the forks whose block body embeds the execution payload.
var payloadVersions = []string{"bellatrix", "capella", "deneb", "electra", "fulu"}

// Builders for the data of a /eth/v2/beacon/blocks response, trimmed to the fields the decoder reads.
func blockData(slot, body string) string {
	return `{"message":{"slot":"` + slot + `","body":` + body + `}}`
}

func payloadBody(fields string) string { return `{"execution_payload":{` + fields + `}}` }

func bidBody(fields string) string {
	return `{"signed_execution_payload_bid":{"message":{` + fields + `}}}`
}

var (
	validPayloadBody = payloadBody(`"block_number":"42","block_hash":"` + hashA + `"`)
	validBidBody     = bidBody(`"parent_block_hash":"` + hashA + `","block_hash":"` + hashB + `"`)
	// bothBody carries both shapes, so only the version can pick one.
	bothBody = `{"execution_payload":{"block_number":"42","block_hash":"` + hashA + `"},` +
		`"signed_execution_payload_bid":{"message":{"parent_block_hash":"` + hashA + `","block_hash":"` + hashB + `"}}}`
)

func TestGetBlockExecution(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCli := httptestutils.NewMockSender(ctrl)
	c := NewClientFromClient(mockCli)

	reply := func(blockID string, status int, body []byte) {
		req := httptestutils.NewGockRequest()
		req.Get("/eth/v2/beacon/blocks/" + blockID).Reply(status).JSON(body)
		mockCli.EXPECT().Gock(req)
	}

	t.Run("fulu block from mainnet carries its payload", func(t *testing.T) {
		reply("test-block", 200, readFixture(t, "block_fulu_mainnet.json"))

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
		reply("test-block", 200, readFixture(t, "block_gloas_sepolia.json"))

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

	for _, blockID := range []string{"finalized", "head", "11389280", hashA} {
		t.Run("requests block id "+blockID, func(t *testing.T) {
			reply(blockID, 200, readFixture(t, "block_gloas_sepolia.json"))

			_, err := c.GetBlockExecution(t.Context(), blockID)
			require.NoError(t, err)
		})
	}

	t.Run("missing block surfaces the 404", func(t *testing.T) {
		reply("test-block", 404, []byte(`{"code":404,"message":"NOT_FOUND: beacon block at slot 123"}`))

		_, err := c.GetBlockExecution(t.Context(), "test-block")
		var detailed autorest.DetailedError
		require.ErrorAs(t, err, &detailed)
		assert.Equal(t, 404, detailed.StatusCode)
	})

	t.Run("server error surfaces its status", func(t *testing.T) {
		reply("test-block", 500, []byte(`{"code":500,"message":"internal error"}`))

		_, err := c.GetBlockExecution(t.Context(), "test-block")
		var detailed autorest.DetailedError
		require.ErrorAs(t, err, &detailed)
		assert.Equal(t, 500, detailed.StatusCode)
	})

	t.Run("unsupported version is ErrUnsupportedBlockVersion", func(t *testing.T) {
		reply("test-block", 200, []byte(`{"version":"heze","data":`+blockData("1", validBidBody)+`}`))

		_, err := c.GetBlockExecution(t.Context(), "test-block")
		assert.ErrorIs(t, err, ErrUnsupportedBlockVersion)
	})

	t.Run("block missing a required field is an error", func(t *testing.T) {
		reply("test-block", 200, []byte(`{"version":"gloas","data":`+blockData("1", `{}`)+`}`))

		_, err := c.GetBlockExecution(t.Context(), "test-block")
		assert.ErrorContains(t, err, "has no execution payload bid")
	})

	t.Run("invalid JSON is an error", func(t *testing.T) {
		reply("test-block", 200, []byte(`{"version":"gloas","data":`))

		_, err := c.GetBlockExecution(t.Context(), "test-block")
		require.Error(t, err)
	})

	t.Run("transport failure is an error", func(t *testing.T) {
		errConn := errors.New("connection refused")
		mockCli.EXPECT().Do(gomock.Any()).Return(nil, errConn)

		_, err := c.GetBlockExecution(t.Context(), "test-block")
		assert.ErrorIs(t, err, errConn)
	})
}

func TestDecodeBlockExecution(t *testing.T) {
	// The version decides which shape is read, even when the body carries both.
	for _, version := range payloadVersions {
		for name, body := range map[string]string{"payload": validPayloadBody, "payload and bid": bothBody} {
			t.Run(version+" reads the embedded payload from a body with "+name, func(t *testing.T) {
				block, err := decodeBlockExecution(version, []byte(blockData("100", body)))
				require.NoError(t, err)
				assert.Equal(t, &types.BlockExecution{
					Version: version,
					Slot:    100,
					Payload: &types.ExecutionPayloadRef{BlockNumber: 42, BlockHash: common.HexToHash(hashA)},
				}, block)
			})
		}
	}

	for name, body := range map[string]string{"bid": validBidBody, "payload and bid": bothBody} {
		t.Run("gloas reads the bid from a body with "+name, func(t *testing.T) {
			block, err := decodeBlockExecution("gloas", []byte(blockData("100", body)))
			require.NoError(t, err)
			assert.Equal(t, &types.BlockExecution{
				Version: "gloas",
				Slot:    100,
				Bid:     &types.ExecutionPayloadBidRef{ParentBlockHash: common.HexToHash(hashA), BlockHash: common.HexToHash(hashB)},
			}, block)
		})
	}
}

// The spec only requires the bid's block_hash to differ from parent_block_hash, so a valid block may
// carry a zero one. Only the parent_block_hash is required to be non-zero.
func TestDecodeBlockExecutionGloasBidZeroBlockHash(t *testing.T) {
	t.Run("zero block hash decodes with the parent hash", func(t *testing.T) {
		data := blockData("100", bidBody(`"parent_block_hash":"`+hashA+`","block_hash":"`+zeroHash+`"`))

		block, err := decodeBlockExecution("gloas", []byte(data))
		require.NoError(t, err)
		assert.Equal(t, &types.BlockExecution{
			Version: "gloas",
			Slot:    100,
			Bid:     &types.ExecutionPayloadBidRef{ParentBlockHash: common.HexToHash(hashA), BlockHash: common.Hash{}},
		}, block)
	})

	t.Run("zero parent hash is an error", func(t *testing.T) {
		data := blockData("100", bidBody(`"parent_block_hash":"`+zeroHash+`","block_hash":"`+zeroHash+`"`))

		block, err := decodeBlockExecution("gloas", []byte(data))
		assert.Nil(t, block)
		assert.ErrorContains(t, err, "has no execution payload bid")
	})

	t.Run("missing bid is an error", func(t *testing.T) {
		block, err := decodeBlockExecution("gloas", []byte(blockData("100", `{}`)))
		assert.Nil(t, block)
		assert.ErrorContains(t, err, "has no execution payload bid")
	})
}

func TestDecodeBlockExecutionRejectsUnsupportedVersions(t *testing.T) {
	// Checked before decoding: an unknown fork may not have either shape.
	for _, version := range []string{"", "phase0", "altair", "heze", "Fulu", "GLOAS", " fulu"} {
		t.Run("version "+version, func(t *testing.T) {
			for _, data := range []string{blockData("100", validPayloadBody), blockData("100", validBidBody), `{`} {
				_, err := decodeBlockExecution(version, []byte(data))
				assert.ErrorIs(t, err, ErrUnsupportedBlockVersion)
			}
		})
	}
}

// Every required field, missing or zero, is an error for every version that requires it.
func TestDecodeBlockExecutionRejectsIncompleteBlocks(t *testing.T) {
	type errorCase struct {
		name    string
		data    string
		wantErr string
	}
	const (
		noPayload = "has no execution payload"
		noBid     = "has no execution payload bid"
		noSlot    = "has no slot"
		badDecode = "failed to decode"
	)

	// Malformed for any version.
	commonCases := []errorCase{
		{"missing data", ``, badDecode},
		{"data is not JSON", `{`, badDecode},
		{"data is not an object", `[]`, badDecode},
		{"slot is not a number", blockData("not-a-number", validPayloadBody), badDecode},
		{"null data", `null`, noSlot},
		{"missing message", `{}`, noSlot},
		{"missing slot", `{"message":{"body":` + validPayloadBody + `}}`, noSlot},
		{"slot 0", blockData("0", validPayloadBody), noSlot},
	}

	payloadCases := []errorCase{
		{"missing body", `{"message":{"slot":"100"}}`, noPayload},
		{"missing payload", blockData("100", `{}`), noPayload},
		{"null payload", blockData("100", `{"execution_payload":null}`), noPayload},
		{"bid instead of payload", blockData("100", validBidBody), noPayload},
		{"missing block hash", blockData("100", payloadBody(`"block_number":"42"`)), noPayload},
		{"zero block hash", blockData("100", payloadBody(`"block_number":"42","block_hash":"`+zeroHash+`"`)), noPayload},
		{"missing block number", blockData("100", payloadBody(`"block_hash":"`+hashA+`"`)), noPayload},
		{"block number 0", blockData("100", payloadBody(`"block_number":"0","block_hash":"`+hashA+`"`)), noPayload},
		{"unquoted block number", blockData("100", payloadBody(`"block_number":42,"block_hash":"`+hashA+`"`)), badDecode},
		{"block number is not a number", blockData("100", payloadBody(`"block_number":"forty-two","block_hash":"`+hashA+`"`)), badDecode},
		{"negative block number", blockData("100", payloadBody(`"block_number":"-1","block_hash":"`+hashA+`"`)), badDecode},
		{"short block hash", blockData("100", payloadBody(`"block_number":"42","block_hash":"0x1234"`)), badDecode},
		{"block hash is not hex", blockData("100", payloadBody(`"block_number":"42","block_hash":"0x`+strings.Repeat("z", 64)+`"`)), badDecode},
	}

	bidCases := []errorCase{
		{"missing body", `{"message":{"slot":"100"}}`, noBid},
		{"missing bid", blockData("100", `{}`), noBid},
		{"null bid", blockData("100", `{"signed_execution_payload_bid":null}`), noBid},
		{"bid without message", blockData("100", `{"signed_execution_payload_bid":{}}`), noBid},
		{"payload instead of bid", blockData("100", validPayloadBody), noBid},
		{"missing parent hash", blockData("100", bidBody(`"block_hash":"`+hashB+`"`)), noBid},
		{"zero parent hash", blockData("100", bidBody(`"parent_block_hash":"`+zeroHash+`","block_hash":"`+hashB+`"`)), noBid},
		{"missing block hash", blockData("100", bidBody(`"parent_block_hash":"`+hashA+`"`)), noBid},
		{"null block hash", blockData("100", bidBody(`"parent_block_hash":"`+hashA+`","block_hash":null`)), noBid},
		{"short parent hash", blockData("100", bidBody(`"parent_block_hash":"0x1234","block_hash":"`+hashB+`"`)), badDecode},
		{"short block hash", blockData("100", bidBody(`"parent_block_hash":"`+hashA+`","block_hash":"0x1234"`)), badDecode},
	}

	run := func(t *testing.T, version string, cases []errorCase) {
		t.Helper()
		for _, tt := range cases {
			t.Run(version+" "+tt.name, func(t *testing.T) {
				block, err := decodeBlockExecution(version, []byte(tt.data))
				assert.Nil(t, block)
				assert.ErrorContains(t, err, tt.wantErr)
				if tt.wantErr == noPayload {
					assert.NotContains(t, err.Error(), "bid", "a payload error must not come from the bid branch")
				}
			})
		}
	}

	for _, version := range payloadVersions {
		run(t, version, slices.Concat(commonCases, payloadCases))
	}
	run(t, "gloas", slices.Concat(commonCases, bidCases))
}

func TestGetBlock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCli := httptestutils.NewMockSender(ctrl)
	c := NewClientFromClient(mockCli)

	// withVersion returns the mainnet fulu fixture relabelled as version.
	withVersion := func(t *testing.T, version string) []byte {
		t.Helper()
		var resp map[string]json.RawMessage
		require.NoError(t, json.Unmarshal(readFixture(t, "block_fulu_mainnet.json"), &resp))
		v, err := json.Marshal(version)
		require.NoError(t, err)
		resp["version"] = v
		body, err := json.Marshal(resp)
		require.NoError(t, err)
		return body
	}

	reply := func(body []byte) {
		req := httptestutils.NewGockRequest()
		req.Get("/eth/v2/beacon/blocks/test-block").Reply(200).JSON(body)
		mockCli.EXPECT().Gock(req)
	}

	// An empty version is accepted for nodes that do not report it.
	for _, version := range []string{"", "phase0", "altair", "bellatrix", "capella", "deneb", "electra", "fulu"} {
		t.Run("version "+version+" decodes", func(t *testing.T) {
			reply(withVersion(t, version))

			block, err := c.GetBlock(t.Context(), "test-block")
			require.NoError(t, err)
			assert.Equal(t, uint64(26225974), uint64(block.Message.Body.ExecutionPayload.BlockNumber))
		})
	}

	t.Run("gloas block is an error instead of a zero payload", func(t *testing.T) {
		reply(readFixture(t, "block_gloas_sepolia.json"))

		_, err := c.GetBlock(t.Context(), "test-block")
		assert.True(t, errors.Is(err, ErrUnsupportedBlockVersion), "got %v", err)
	})

	for _, version := range []string{"heze", "Fulu"} {
		t.Run("version "+version+" is an error", func(t *testing.T) {
			reply(withVersion(t, version))

			_, err := c.GetBlock(t.Context(), "test-block")
			assert.ErrorIs(t, err, ErrUnsupportedBlockVersion)
		})
	}

	t.Run("missing block surfaces the 404", func(t *testing.T) {
		req := httptestutils.NewGockRequest()
		req.Get("/eth/v2/beacon/blocks/test-block").Reply(404).JSON([]byte(`{"code":404,"message":"NOT_FOUND: beacon block at slot 123"}`))
		mockCli.EXPECT().Gock(req)

		_, err := c.GetBlock(t.Context(), "test-block")
		var detailed autorest.DetailedError
		require.ErrorAs(t, err, &detailed)
		assert.Equal(t, 404, detailed.StatusCode)
	})

	t.Run("transport failure is an error", func(t *testing.T) {
		errConn := errors.New("connection refused")
		mockCli.EXPECT().Do(gomock.Any()).Return(nil, errConn)

		_, err := c.GetBlock(t.Context(), "test-block")
		assert.ErrorIs(t, err, errConn)
	})
}

func readFixture(t *testing.T, name string) []byte {
	t.Helper()
	data, err := os.ReadFile("testdata/" + name)
	require.NoError(t, err)
	return data
}
