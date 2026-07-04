package parser

import (
	_ "embed"
	"encoding/hex"
	"strings"
	"testing"
)

//go:embed cb4134000.hex
var realCoinbaseHex string

// TestRealNU6_3CoinbaseParses parses the actual V6 coinbase from testnet block
// 4,134,000 (the NU6.3 activation block, produced by zfnd/zebra:6.0.0-rc.0).
// This is the regression that the placeholder version-group/branch IDs failed.
func TestRealNU6_3CoinbaseParses(t *testing.T) {
	raw, err := hex.DecodeString(strings.TrimSpace(realCoinbaseHex))
	if err != nil {
		t.Fatal(err)
	}
	tx := NewTransaction()
	rest, err := tx.ParseFromSlice(raw)
	if err != nil {
		t.Fatalf("real NU6.3 coinbase failed to parse: %v", err)
	}
	if len(rest) != 0 {
		t.Fatalf("did not consume entire coinbase, %d bytes remain", len(rest))
	}
	if tx.version != IRONWOOD_NU6_3_TX_VERSION {
		t.Fatalf("version = %d, want 6", tx.version)
	}
	if tx.nVersionGroupID != IRONWOOD_NU6_3_VERSION_GROUP_ID {
		t.Fatalf("nVersionGroupID = 0x%08X, want 0x%08X", tx.nVersionGroupID, IRONWOOD_NU6_3_VERSION_GROUP_ID)
	}
	if tx.consensusBranchID != IRONWOOD_NU6_3_CONSENSUS_BRANCH_ID {
		t.Fatalf("consensusBranchID = 0x%08X, want 0x%08X", tx.consensusBranchID, IRONWOOD_NU6_3_CONSENSUS_BRANCH_ID)
	}
}
