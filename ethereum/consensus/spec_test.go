//revive:disable-next-line:package-directory-mismatch
package ethcl

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestGetSpecByChainID(t *testing.T) {
	tests := []struct {
		name            string
		chainID         uint64
		expectError     bool
		expectedName    string
		expectedGenesis int64
	}{
		{
			name:            "Mainnet",
			chainID:         MainnetChainID,
			expectError:     false,
			expectedName:    "Mainnet",
			expectedGenesis: 1606824023,
		},
		{
			name:            "Goerli",
			chainID:         GoerliChainID,
			expectError:     false,
			expectedName:    "Goerli",
			expectedGenesis: 1614588812,
		},
		{
			name:            "Sepolia",
			chainID:         SepoliaChainID,
			expectError:     false,
			expectedName:    "Sepolia",
			expectedGenesis: 1655733600,
		},
		{
			name:            "Holesky",
			chainID:         HoleskyChainID,
			expectError:     false,
			expectedName:    "Holesky",
			expectedGenesis: 1695902400,
		},
		{
			name:            "Hoodi",
			chainID:         HoodiChainID,
			expectError:     false,
			expectedName:    "Hoodi",
			expectedGenesis: 1742213400,
		},
		{
			name:            "Unknown Chain ID",
			chainID:         999999,
			expectError:     true,
			expectedName:    "",
			expectedGenesis: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spec, err := GetSpecByChainID(tt.chainID)

			// Check error expectations
			if tt.expectError && err == nil {
				t.Errorf("Expected error for chain ID %d, but got none", tt.chainID)
				return
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error for chain ID %d: %v", tt.chainID, err)
				return
			}

			// Skip further checks if error was expected
			if tt.expectError {
				return
			}

			// Check network name
			if spec.NetworkName != tt.expectedName {
				t.Errorf("Expected network name %s, got %s", tt.expectedName, spec.NetworkName)
			}

			// Check genesis time
			if spec.GenesisTime != tt.expectedGenesis {
				t.Errorf("Expected genesis time %d, got %d", tt.expectedGenesis, spec.GenesisTime)
			}
		})
	}
}

func TestTimeToSlot(t *testing.T) {
	// Create a test spec
	spec := &Spec{
		ChainID:        MainnetChainID,
		GenesisTime:    1606824023,
		SecondsPerSlot: 12,
		SlotsPerEpoch:  32,
		NetworkName:    "Mainnet",
	}

	tests := []struct {
		name         string
		timestamp    int64
		expectedSlot Slot
	}{
		{
			name:         "Genesis time",
			timestamp:    1606824023,
			expectedSlot: 0,
		},
		{
			name:         "Genesis time plus 1 second",
			timestamp:    1606824024,
			expectedSlot: 0,
		},
		{
			name:         "Genesis time plus 12 seconds (1 slot)",
			timestamp:    1606824035,
			expectedSlot: 1,
		},
		{
			name:         "Genesis time plus 24 seconds (2 slots)",
			timestamp:    1606824047,
			expectedSlot: 2,
		},
		{
			name:         "Genesis time plus 384 seconds (32 slots = 1 epoch)",
			timestamp:    1606824407,
			expectedSlot: 32,
		},
		{
			name:         "Before genesis time",
			timestamp:    1606824022,
			expectedSlot: 0,
		},
		{
			name:         "Far future",
			timestamp:    1706824023, // ~100M seconds after genesis
			expectedSlot: 8333333,    // ~100M / 12
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			slot := spec.TimeToSlot(tt.timestamp)
			if slot != tt.expectedSlot {
				t.Errorf("TimeToSlot(%d) = %d, want %d", tt.timestamp, slot, tt.expectedSlot)
			}
		})
	}
}

func TestSlotToEpoch(t *testing.T) {
	// Create a test spec
	spec := &Spec{
		ChainID:        MainnetChainID,
		GenesisTime:    1606824023,
		SecondsPerSlot: 12,
		SlotsPerEpoch:  32,
		NetworkName:    "Mainnet",
	}

	tests := []struct {
		name          string
		slot          Slot
		expectedEpoch Epoch
	}{
		{
			name:          "Slot 0",
			slot:          0,
			expectedEpoch: 0,
		},
		{
			name:          "Slot 31 (last slot of epoch 0)",
			slot:          31,
			expectedEpoch: 0,
		},
		{
			name:          "Slot 32 (first slot of epoch 1)",
			slot:          32,
			expectedEpoch: 1,
		},
		{
			name:          "Slot 64 (first slot of epoch 2)",
			slot:          64,
			expectedEpoch: 2,
		},
		{
			name:          "Slot 320 (first slot of epoch 10)",
			slot:          320,
			expectedEpoch: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			epoch := spec.SlotToEpoch(tt.slot)
			if epoch != tt.expectedEpoch {
				t.Errorf("SlotToEpoch(%d) = %d, want %d", tt.slot, epoch, tt.expectedEpoch)
			}
		})
	}
}

func TestTimeToEpoch(t *testing.T) {
	// Create a test spec
	spec := &Spec{
		ChainID:        MainnetChainID,
		GenesisTime:    1606824023,
		SecondsPerSlot: 12,
		SlotsPerEpoch:  32,
		NetworkName:    "Mainnet",
	}

	tests := []struct {
		name          string
		timestamp     int64
		expectedEpoch Epoch
	}{
		{
			name:          "Genesis time",
			timestamp:     1606824023,
			expectedEpoch: 0,
		},
		{
			name:          "Genesis time plus 383 seconds (still epoch 0)",
			timestamp:     1606824023 + 383,
			expectedEpoch: 0,
		},
		{
			name:          "Genesis time plus 384 seconds (start of epoch 1)",
			timestamp:     1606824023 + 384,
			expectedEpoch: 1,
		},
		{
			name:          "Genesis time plus 768 seconds (start of epoch 2)",
			timestamp:     1606824023 + 768,
			expectedEpoch: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			epoch := spec.TimeToEpoch(tt.timestamp)
			if epoch != tt.expectedEpoch {
				t.Errorf("TimeToEpoch(%d) = %d, want %d", tt.timestamp, epoch, tt.expectedEpoch)
			}
		})
	}
}

func TestCurrentFunctions(t *testing.T) {
	// This test is somewhat limited as it depends on the current time
	// We'll just check that the functions don't panic and return reasonable values

	// Create a test spec - we'll use mainnet
	spec, err := GetSpecByChainID(HoodiChainID)
	if err != nil {
		t.Fatalf("Failed to get spec: %v", err)
	}

	currentCLEpoch := uint64(10357)
	currentTimeSeconds := time.Date(2025, 5, 2, 12, 59, 0, 0, time.UTC).Unix()
	require.Equal(t, currentCLEpoch, uint64(spec.TimeToEpoch(currentTimeSeconds)))

	// Test CurrentSlot
	currentSlot := spec.CurrentSlot()

	// Calculate expected slot range
	nowSeconds := time.Now().Unix()
	expectedSlot := spec.TimeToSlot(nowSeconds)

	// The slot should be very close to what we calculate
	if currentSlot < expectedSlot-1 || currentSlot > expectedSlot+1 {
		t.Errorf("CurrentSlot() = %d, expected close to %d", currentSlot, expectedSlot)
	}

	// Test CurrentEpoch
	currentEpoch := spec.CurrentEpoch()
	expectedEpoch := spec.SlotToEpoch(expectedSlot)

	// The epoch should match what we calculate
	if currentEpoch != expectedEpoch {
		t.Errorf("CurrentEpoch() = %d, expected %d", currentEpoch, expectedEpoch)
	}
}

// This test table ensures boundary conditions work correctly
func TestEdgeCases(t *testing.T) {
	spec, err := GetSpecByChainID(MainnetChainID)
	require.NoError(t, err)

	// Edge case: exactly at slot boundary
	slotBoundaryTime := spec.GenesisTime + 12000
	expectedSlot := Slot(1000)
	if slot := spec.TimeToSlot(slotBoundaryTime); slot != expectedSlot {
		t.Errorf("At boundary: TimeToSlot(%d) = %d, want %d",
			slotBoundaryTime, slot, expectedSlot)
	}

	// Edge case: exactly at epoch boundary
	epochBoundaryTime := spec.GenesisTime + (32 * 12 * 10) // 10 epochs after genesis
	expectedEpoch := Epoch(10)
	if epoch := spec.TimeToEpoch(epochBoundaryTime); epoch != expectedEpoch {
		t.Errorf("At epoch boundary: TimeToEpoch(%d) = %d, want %d",
			epochBoundaryTime, epoch, expectedEpoch)
	}
}

func TestEpochToSlot(t *testing.T) {
	spec := &Spec{
		ChainID:        MainnetChainID,
		GenesisTime:    1606824023,
		SecondsPerSlot: 12,
		SlotsPerEpoch:  32,
		NetworkName:    "Mainnet",
	}

	tests := []struct {
		name         string
		epoch        Epoch
		expectedSlot Slot
	}{
		{
			name:         "Epoch 0",
			epoch:        0,
			expectedSlot: 0,
		},
		{
			name:         "Epoch 1",
			epoch:        1,
			expectedSlot: 32,
		},
		{
			name:         "Epoch 2",
			epoch:        2,
			expectedSlot: 64,
		},
		{
			name:         "Epoch 10",
			epoch:        10,
			expectedSlot: 320,
		},
		{
			name:         "Epoch 100",
			epoch:        100,
			expectedSlot: 3200,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			slot := spec.EpochToSlot(tt.epoch)
			require.Equal(t, tt.expectedSlot, slot)
		})
	}
}

func TestSlotToTime(t *testing.T) {
	spec := &Spec{
		ChainID:        MainnetChainID,
		GenesisTime:    1606824023,
		SecondsPerSlot: 12,
		SlotsPerEpoch:  32,
		NetworkName:    "Mainnet",
	}

	tests := []struct {
		name              string
		slot              Slot
		expectedTimestamp int64
	}{
		{
			name:              "Slot 0",
			slot:              0,
			expectedTimestamp: 1606824023,
		},
		{
			name:              "Slot 1",
			slot:              1,
			expectedTimestamp: 1606824023 + 12,
		},
		{
			name:              "Slot 32 (epoch 1 start)",
			slot:              32,
			expectedTimestamp: 1606824023 + (32 * 12),
		},
		{
			name:              "Slot 100",
			slot:              100,
			expectedTimestamp: 1606824023 + (100 * 12),
		},
		{
			name:              "Slot 1000",
			slot:              1000,
			expectedTimestamp: 1606824023 + (1000 * 12),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			timestamp := spec.SlotToTime(tt.slot)
			require.Equal(t, tt.expectedTimestamp, timestamp)
		})
	}
}

func TestEpochToTime(t *testing.T) {
	spec := &Spec{
		ChainID:        MainnetChainID,
		GenesisTime:    1606824023,
		SecondsPerSlot: 12,
		SlotsPerEpoch:  32,
		NetworkName:    "Mainnet",
	}

	tests := []struct {
		name              string
		epoch             Epoch
		expectedTimestamp int64
	}{
		{
			name:              "Epoch 0",
			epoch:             0,
			expectedTimestamp: 1606824023,
		},
		{
			name:              "Epoch 1",
			epoch:             1,
			expectedTimestamp: 1606824023 + (32 * 12),
		},
		{
			name:              "Epoch 2",
			epoch:             2,
			expectedTimestamp: 1606824023 + (64 * 12),
		},
		{
			name:              "Epoch 10",
			epoch:             10,
			expectedTimestamp: 1606824023 + (320 * 12),
		},
		{
			name:              "Epoch 100",
			epoch:             100,
			expectedTimestamp: 1606824023 + (3200 * 12),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			timestamp := spec.EpochToTime(tt.epoch)
			require.Equal(t, tt.expectedTimestamp, timestamp)
		})
	}
}

func TestShardCommitteePeriod(t *testing.T) {
	spec, err := GetSpecByChainID(MainnetChainID)
	require.NoError(t, err)

	period := spec.ShardCommitteePeriod()
	require.Equal(t, uint64(ShardCommitteePeriod), period)
	require.Equal(t, uint64(256), period)
}

func TestEpochUint64(t *testing.T) {
	tests := []struct {
		name     string
		epoch    Epoch
		expected uint64
	}{
		{
			name:     "Epoch 0",
			epoch:    0,
			expected: 0,
		},
		{
			name:     "Epoch 1",
			epoch:    1,
			expected: 1,
		},
		{
			name:     "Epoch 100",
			epoch:    100,
			expected: 100,
		},
		{
			name:     "Large epoch",
			epoch:    1000000,
			expected: 1000000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.epoch.Uint64()
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestRoundTripConversions(t *testing.T) {
	spec, err := GetSpecByChainID(MainnetChainID)
	require.NoError(t, err)

	t.Run("Time -> Slot -> Time", func(t *testing.T) {
		originalTime := spec.GenesisTime + 1000*12
		slot := spec.TimeToSlot(originalTime)
		convertedTime := spec.SlotToTime(slot)
		require.Equal(t, originalTime, convertedTime)
	})

	t.Run("Slot -> Epoch -> Slot", func(t *testing.T) {
		originalSlot := Slot(320) // Start of epoch 10
		epoch := spec.SlotToEpoch(originalSlot)
		convertedSlot := spec.EpochToSlot(epoch)
		require.Equal(t, originalSlot, convertedSlot)
	})

	t.Run("Epoch -> Time -> Epoch", func(t *testing.T) {
		originalEpoch := Epoch(10)
		timestamp := spec.EpochToTime(originalEpoch)
		convertedEpoch := spec.TimeToEpoch(timestamp)
		require.Equal(t, originalEpoch, convertedEpoch)
	})

	t.Run("Epoch -> Slot -> Epoch", func(t *testing.T) {
		originalEpoch := Epoch(5)
		slot := spec.EpochToSlot(originalEpoch)
		convertedEpoch := spec.SlotToEpoch(slot)
		require.Equal(t, originalEpoch, convertedEpoch)
	})
}

func TestMultipleNetworks(t *testing.T) {
	networks := []uint64{
		MainnetChainID,
		GoerliChainID,
		SepoliaChainID,
		HoleskyChainID,
		HoodiChainID,
		LocalChainID,
	}

	for _, chainID := range networks {
		t.Run("Chain ID "+strconv.FormatUint(chainID, 10), func(t *testing.T) {
			spec, err := GetSpecByChainID(chainID)
			require.NoError(t, err)
			require.NotNil(t, spec)

			// Test basic conversions work
			slot := spec.TimeToSlot(spec.GenesisTime + 1000)
			require.Greater(t, uint64(slot), uint64(0))

			epoch := spec.SlotToEpoch(100)
			require.GreaterOrEqual(t, uint64(epoch), uint64(0))

			timestamp := spec.SlotToTime(100)
			require.Greater(t, timestamp, spec.GenesisTime)

			// Test constants
			require.Equal(t, SecondsPerSlot, spec.SecondsPerSlot)
			require.Equal(t, SlotsPerEpoch, spec.SlotsPerEpoch)
		})
	}
}

func TestTimeConversionConsistency(t *testing.T) {
	spec, err := GetSpecByChainID(MainnetChainID)
	require.NoError(t, err)

	// Test that TimeToEpoch is consistent with TimeToSlot + SlotToEpoch
	testTime := spec.GenesisTime + 5000

	directEpoch := spec.TimeToEpoch(testTime)

	slot := spec.TimeToSlot(testTime)
	indirectEpoch := spec.SlotToEpoch(slot)

	require.Equal(t, directEpoch, indirectEpoch)
}

func TestSlotBoundaries(t *testing.T) {
	spec, err := GetSpecByChainID(MainnetChainID)
	require.NoError(t, err)

	t.Run("Just before slot boundary", func(t *testing.T) {
		timestamp := spec.GenesisTime + 11 // 1 second before slot 1
		slot := spec.TimeToSlot(timestamp)
		require.Equal(t, Slot(0), slot)
	})

	t.Run("Exactly at slot boundary", func(t *testing.T) {
		timestamp := spec.GenesisTime + 12 // Exactly at slot 1
		slot := spec.TimeToSlot(timestamp)
		require.Equal(t, Slot(1), slot)
	})

	t.Run("Just after slot boundary", func(t *testing.T) {
		timestamp := spec.GenesisTime + 13 // 1 second after slot 1 starts
		slot := spec.TimeToSlot(timestamp)
		require.Equal(t, Slot(1), slot)
	})
}

func TestEpochBoundaries(t *testing.T) {
	spec, err := GetSpecByChainID(MainnetChainID)
	require.NoError(t, err)

	// Calculate epoch duration: 32 slots * 12 seconds = 384 seconds
	// Use const values to avoid uint64->int64 conversion warnings
	const slotsPerEpoch int64 = 32
	const secondsPerSlot int64 = 12
	epochDuration := slotsPerEpoch * secondsPerSlot

	// Verify our constants match the spec
	require.Equal(t, uint64(slotsPerEpoch), spec.SlotsPerEpoch)
	require.Equal(t, uint64(secondsPerSlot), spec.SecondsPerSlot)

	t.Run("Just before epoch boundary", func(t *testing.T) {
		timestamp := spec.GenesisTime + epochDuration - 1
		epoch := spec.TimeToEpoch(timestamp)
		require.Equal(t, Epoch(0), epoch)
	})

	t.Run("Exactly at epoch boundary", func(t *testing.T) {
		timestamp := spec.GenesisTime + epochDuration
		epoch := spec.TimeToEpoch(timestamp)
		require.Equal(t, Epoch(1), epoch)
	})

	t.Run("Just after epoch boundary", func(t *testing.T) {
		timestamp := spec.GenesisTime + epochDuration + 1
		epoch := spec.TimeToEpoch(timestamp)
		require.Equal(t, Epoch(1), epoch)
	})
}

func TestLargeValues(t *testing.T) {
	spec, err := GetSpecByChainID(MainnetChainID)
	require.NoError(t, err)

	t.Run("Very large slot number", func(t *testing.T) {
		largeSlot := Slot(10000000) // ~10 million slots
		timestamp := spec.SlotToTime(largeSlot)
		convertedSlot := spec.TimeToSlot(timestamp)
		require.Equal(t, largeSlot, convertedSlot)
	})

	t.Run("Very large epoch number", func(t *testing.T) {
		largeEpoch := Epoch(100000) // 100k epochs
		timestamp := spec.EpochToTime(largeEpoch)
		convertedEpoch := spec.TimeToEpoch(timestamp)
		require.Equal(t, largeEpoch, convertedEpoch)
	})
}

func TestNetworkSpecificValues(t *testing.T) {
	t.Run("Hoodi network has different average block time", func(t *testing.T) {
		spec, err := GetSpecByChainID(HoodiChainID)
		require.NoError(t, err)
		require.Equal(t, 13.1, spec.AverageBlockTimeSeconds)
		require.Equal(t, uint64(30000), spec.NetworkOffsetBlocks)
	})

	t.Run("Mainnet has standard values", func(t *testing.T) {
		spec, err := GetSpecByChainID(MainnetChainID)
		require.NoError(t, err)
		require.Equal(t, 12.0, spec.AverageBlockTimeSeconds)
		require.Equal(t, uint64(100), spec.NetworkOffsetBlocks)
	})
}

func TestLocalChainID(t *testing.T) {
	spec, err := GetSpecByChainID(LocalChainID)
	require.NoError(t, err)
	require.Equal(t, "Local", spec.NetworkName)
	require.Equal(t, LocalChainID, spec.ChainID)
	require.Equal(t, int64(1748879977), spec.GenesisTime)
}

func TestSpecificNetworkEpochs(t *testing.T) {
	tests := []struct {
		name              string
		chainID           uint64
		epoch             Epoch
		expectedTimestamp int64
		dateDescription   string
	}{
		{
			name:              "Mainnet epoch 435138",
			chainID:           MainnetChainID,
			epoch:             435138,
			expectedTimestamp: 1773917015, // Mar-19-2026 10:43:35 UTC
			dateDescription:   "Mar-19-2026 10:43:35 UTC",
		},
		{
			name:              "Hoodi epoch 82562",
			chainID:           HoodiChainID,
			epoch:             82562,
			expectedTimestamp: 1773917208, // Mar-19-2026 10:46:48 UTC
			dateDescription:   "Mar-19-2026 10:46:48 UTC",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spec, err := GetSpecByChainID(tt.chainID)
			require.NoError(t, err)

			t.Run("Epoch to Time", func(t *testing.T) {
				timestamp := spec.EpochToTime(tt.epoch)
				require.Equal(t, tt.expectedTimestamp, timestamp,
					"Epoch %d should convert to %s", tt.epoch, tt.dateDescription)
			})

			t.Run("Time to Epoch", func(t *testing.T) {
				epoch := spec.TimeToEpoch(tt.expectedTimestamp)
				require.Equal(t, tt.epoch, epoch,
					"Timestamp %s should convert to epoch %d", tt.dateDescription, tt.epoch)
			})

			t.Run("Round trip conversion", func(t *testing.T) {
				// Epoch -> Time -> Epoch
				timestamp := spec.EpochToTime(tt.epoch)
				convertedEpoch := spec.TimeToEpoch(timestamp)
				require.Equal(t, tt.epoch, convertedEpoch)
			})

			t.Run("Slot calculations", func(t *testing.T) {
				// Calculate expected slot (epoch * slots per epoch)
				expectedSlot := Slot(uint64(tt.epoch) * spec.SlotsPerEpoch)
				slot := spec.EpochToSlot(tt.epoch)
				require.Equal(t, expectedSlot, slot)

				// Verify slot to time matches epoch time
				slotTimestamp := spec.SlotToTime(expectedSlot)
				require.Equal(t, tt.expectedTimestamp, slotTimestamp)
			})
		})
	}
}
