//revive:disable-next-line:package-directory-mismatch
package ethcl

import (
	"fmt"
	"time"
)

// Common constants for Ethereum networks
const (
	// ShardCommitteePeriod defaults to 256 epochs
	ShardCommitteePeriod = 256

	// SecondsPerSlot Seconds per slot
	SecondsPerSlot uint64 = 12

	// SlotsPerEpoch Slots per epoch
	SlotsPerEpoch uint64 = 32

	// Chain IDs
	MainnetChainID uint64 = 1
	GoerliChainID  uint64 = 5
	SepoliaChainID uint64 = 11155111
	HoleskyChainID uint64 = 17000
	HoodiChainID   uint64 = 560048
	LocalChainID   uint64 = 3151908
)

// Timestamp represents a UNIX timestamp in seconds

// Slot represents a consensus layer slot number
type Slot uint64

// Epoch represents a consensus layer epoch number
type Epoch uint64

// Spec contains the parameters for a specific Ethereum network
type Spec struct {
	ChainID                 uint64
	GenesisTime             int64
	SecondsPerSlot          uint64
	SlotsPerEpoch           uint64
	NetworkName             string
	AverageBlockTimeSeconds float64 // data retrieved from https://eth.blockscout.com
	NetworkOffsetBlocks     uint64  // Offset to align local/test networks with mainnet time
}

// GetSpecByChainID returns the appropriate Spec based on the chain ID
func GetSpecByChainID(chainID uint64) (*Spec, error) {
	switch chainID {
	case MainnetChainID:
		return &Spec{
			ChainID:                 MainnetChainID,
			GenesisTime:             1606824023, // Dec 1, 2020, 12:00:23 PM UTC
			SecondsPerSlot:          SecondsPerSlot,
			SlotsPerEpoch:           SlotsPerEpoch,
			NetworkName:             "Mainnet",
			AverageBlockTimeSeconds: 12.0,
			NetworkOffsetBlocks:     100,
		}, nil
	case GoerliChainID:
		return &Spec{
			ChainID:                 GoerliChainID,
			GenesisTime:             1614588812, // Mar 1, 2021, 13:20:12 PM UTC
			SecondsPerSlot:          SecondsPerSlot,
			SlotsPerEpoch:           SlotsPerEpoch,
			NetworkName:             "Goerli",
			AverageBlockTimeSeconds: 12.0,
			NetworkOffsetBlocks:     100,
		}, nil
	case SepoliaChainID:
		return &Spec{
			ChainID:                 SepoliaChainID,
			GenesisTime:             1655733600, // Jun 20, 2022, 00:00:00 UTC
			SecondsPerSlot:          SecondsPerSlot,
			SlotsPerEpoch:           SlotsPerEpoch,
			NetworkName:             "Sepolia",
			AverageBlockTimeSeconds: 12.0,
			NetworkOffsetBlocks:     100,
		}, nil
	case HoleskyChainID:
		return &Spec{
			ChainID:                 HoleskyChainID,
			GenesisTime:             1695902400, // Sep 28, 2023, 11:00:00 UTC
			SecondsPerSlot:          SecondsPerSlot,
			SlotsPerEpoch:           SlotsPerEpoch,
			NetworkName:             "Holesky",
			AverageBlockTimeSeconds: 12.0,
			NetworkOffsetBlocks:     100,
		}, nil
	case HoodiChainID:
		return &Spec{
			ChainID:                 HoodiChainID,
			GenesisTime:             1742213400,
			SecondsPerSlot:          SecondsPerSlot,
			SlotsPerEpoch:           SlotsPerEpoch,
			NetworkName:             "Hoodi",
			AverageBlockTimeSeconds: 13.1,
			NetworkOffsetBlocks:     30000,
		}, nil
	case LocalChainID:
		return &Spec{
			ChainID:                 LocalChainID,
			GenesisTime:             1748879977,
			SecondsPerSlot:          SecondsPerSlot,
			SlotsPerEpoch:           SlotsPerEpoch,
			NetworkName:             "Local",
			AverageBlockTimeSeconds: 12.0,
			NetworkOffsetBlocks:     100,
		}, nil
	default:
		return nil, fmt.Errorf("unknown chain ID: %d", chainID)
	}
}

// TimeToSlot converts a timestamp to a slot number based on the network genesis time
func (spec *Spec) TimeToSlot(t int64) Slot {
	if t < spec.GenesisTime {
		return 0
	}
	return Slot(uint64(t-spec.GenesisTime) / spec.SecondsPerSlot)
}

// SlotToEpoch converts a slot number to its corresponding epoch
func (spec *Spec) SlotToEpoch(s Slot) Epoch {
	return Epoch(uint64(s) / spec.SlotsPerEpoch)
}

// TimeToEpoch converts a timestamp directly to an epoch number
func (spec *Spec) TimeToEpoch(t int64) Epoch {
	slot := spec.TimeToSlot(t)
	return spec.SlotToEpoch(slot)
}

// CurrentSlot returns the current slot based on the current time
func (spec *Spec) CurrentSlot() Slot {
	currentTime := time.Now().Unix()
	return spec.TimeToSlot(currentTime)
}

func (spec *Spec) ShardCommitteePeriod() uint64 {
	return ShardCommitteePeriod
}

// CurrentEpoch returns the current epoch based on the current time
func (spec *Spec) CurrentEpoch() Epoch {
	return spec.SlotToEpoch(spec.CurrentSlot())
}

func (e Epoch) Uint64() uint64 {
	return uint64(e)
}
