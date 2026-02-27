//revive:disable-next-line:package-directory-mismatch
package ethcl

import (
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
	currentTime := time.Date(2025, 05, 02, 12, 59, 0, 0, time.UTC).Unix()
	require.Equal(t, currentCLEpoch, uint64(spec.TimeToEpoch(currentTime)))

	// Test CurrentSlot
	currentSlot := spec.CurrentSlot()

	// Calculate expected slot range
	now := time.Now().Unix()
	expectedSlot := spec.TimeToSlot(now)

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
