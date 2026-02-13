//go:build !integration

//nolint:revive // package name intentionally reflects domain, not directory name
package eth2http

import (
	"testing"

	eth2client "github.com/kilnfi/go-utils/ethereum/consensus/client"
	"github.com/stretchr/testify/assert"
)

func TestClientImplementsInterface(t *testing.T) {
	iClient := (*eth2client.Client)(nil)
	client := new(Client)
	assert.Implements(t, iClient, client)
}
