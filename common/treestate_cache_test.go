package common

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/zcash/lightwalletd/walletrpc"
)

func TestTreeStateCacheBasic(t *testing.T) {
	// Setup logging
	logger := logrus.New()
	Log = logger.WithFields(logrus.Fields{"app": "test"})

	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "treestate_cache_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create cache with window of 100 blocks, current height 1000
	cache := NewTreeStateCache(tmpDir, "testnet", 100, 1000)
	if cache == nil {
		t.Fatal("NewTreeStateCache returned nil")
	}
	defer cache.Close()

	// minHeight should be 900 (1000 - 100)
	if cache.GetMinHeight() != 900 {
		t.Errorf("Expected minHeight 900, got %d", cache.GetMinHeight())
	}

	// Test Put and Get for height within window (deep block, > 100 confirmations)
	ts := &walletrpc.TreeState{
		Network:     "testnet",
		Height:      950,
		Hash:        "0000000000000000000000000000000000000000000000000000000000000950",
		Time:        1234567890,
		SaplingTree: "sapling_tree_data",
		OrchardTree: "orchard_tree_data",
	}

	err = cache.Put(ts)
	if err != nil {
		t.Fatalf("Put failed: %v", err)
	}

	// Retrieve by height (currentTip=1100, so height 950 is deep with 150 confirmations)
	retrieved := cache.Get(950, 1100)
	if retrieved == nil {
		t.Fatal("Get returned nil for cached entry")
	}
	if retrieved.Height != 950 {
		t.Errorf("Expected height 950, got %d", retrieved.Height)
	}
	if retrieved.Hash != ts.Hash {
		t.Errorf("Hash mismatch")
	}
	if retrieved.SaplingTree != ts.SaplingTree {
		t.Errorf("SaplingTree mismatch")
	}

	// Retrieve by hash
	retrievedByHash := cache.GetByHash(ts.Hash, 1100)
	if retrievedByHash == nil {
		t.Fatal("GetByHash returned nil for cached entry")
	}
	if retrievedByHash.Height != 950 {
		t.Errorf("GetByHash: Expected height 950, got %d", retrievedByHash.Height)
	}
}

func TestTreeStateCacheOutsideWindow(t *testing.T) {
	logger := logrus.New()
	Log = logger.WithFields(logrus.Fields{"app": "test"})

	tmpDir, err := os.MkdirTemp("", "treestate_cache_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create cache with window of 100 blocks, current height 1000
	cache := NewTreeStateCache(tmpDir, "testnet", 100, 1000)
	if cache == nil {
		t.Fatal("NewTreeStateCache returned nil")
	}
	defer cache.Close()

	// Try to put a tree state outside the window (height 800 < minHeight 900)
	ts := &walletrpc.TreeState{
		Network:     "testnet",
		Height:      800,
		Hash:        "0000000000000000000000000000000000000000000000000000000000000800",
		Time:        1234567890,
		SaplingTree: "sapling_tree_data",
	}

	err = cache.Put(ts)
	if err != nil {
		t.Fatalf("Put should not fail (just skip): %v", err)
	}

	// Should not be cached
	retrieved := cache.Get(800, 1000)
	if retrieved != nil {
		t.Error("Get should return nil for height outside window")
	}
}

func TestTreeStateCacheReorg(t *testing.T) {
	logger := logrus.New()
	Log = logger.WithFields(logrus.Fields{"app": "test"})

	tmpDir, err := os.MkdirTemp("", "treestate_cache_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	cache := NewTreeStateCache(tmpDir, "testnet", 100, 1000)
	if cache == nil {
		t.Fatal("NewTreeStateCache returned nil")
	}
	defer cache.Close()

	// Put several tree states
	for h := uint64(950); h <= 960; h++ {
		ts := &walletrpc.TreeState{
			Network: "testnet",
			Height:  h,
			Hash:    "hash_" + string(rune(h)),
		}
		cache.Put(ts)
	}

	// Verify they're cached (use tip=1100 so they're deep blocks)
	if cache.Get(955, 1100) == nil {
		t.Error("Height 955 should be cached")
	}
	if cache.Get(960, 1100) == nil {
		t.Error("Height 960 should be cached")
	}

	// Simulate reorg at height 955
	cache.Reorg(955)

	// Heights >= 955 should be invalidated
	if cache.Get(955, 1100) != nil {
		t.Error("Height 955 should be invalidated after reorg")
	}
	if cache.Get(960, 1100) != nil {
		t.Error("Height 960 should be invalidated after reorg")
	}

	// Heights < 955 should still be cached
	if cache.Get(954, 1100) == nil {
		t.Error("Height 954 should still be cached after reorg")
	}
}

func TestTreeStateCacheUpdateWindow(t *testing.T) {
	logger := logrus.New()
	Log = logger.WithFields(logrus.Fields{"app": "test"})

	tmpDir, err := os.MkdirTemp("", "treestate_cache_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Window of 50, current height 1000 -> minHeight = 950
	cache := NewTreeStateCache(tmpDir, "testnet", 50, 1000)
	if cache == nil {
		t.Fatal("NewTreeStateCache returned nil")
	}
	defer cache.Close()

	// Put entries at heights 950-960
	for h := uint64(950); h <= 960; h++ {
		ts := &walletrpc.TreeState{
			Network: "testnet",
			Height:  h,
			Hash:    "hash_" + string(rune(h)),
		}
		cache.Put(ts)
	}

	// Verify 950 is cached (use tip=1100 so it's a deep block)
	if cache.Get(950, 1100) == nil {
		t.Error("Height 950 should be cached")
	}

	// Update window: current height now 1020 -> minHeight = 970
	cache.UpdateWindow(1020)

	// Height 950 should now be outside window and removed
	if cache.Get(950, 1100) != nil {
		t.Error("Height 950 should be pruned after window update")
	}

	// Height 960 should also be removed (960 < 970)
	if cache.Get(960, 1100) != nil {
		t.Error("Height 960 should be pruned after window update (960 < 970)")
	}
}

func TestTreeStateCachePersistence(t *testing.T) {
	logger := logrus.New()
	Log = logger.WithFields(logrus.Fields{"app": "test"})

	tmpDir, err := os.MkdirTemp("", "treestate_cache_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create cache and add an entry
	cache := NewTreeStateCache(tmpDir, "testnet", 100, 1000)
	if cache == nil {
		t.Fatal("NewTreeStateCache returned nil")
	}

	ts := &walletrpc.TreeState{
		Network:     "testnet",
		Height:      950,
		Hash:        "0000000000000000000000000000000000000000000000000000000000000950",
		Time:        1234567890,
		SaplingTree: "sapling_tree_data",
		OrchardTree: "orchard_tree_data",
	}
	cache.Put(ts)
	cache.Sync()
	cache.Close()

	// Reopen the cache
	cache2 := NewTreeStateCache(tmpDir, "testnet", 100, 1000)
	if cache2 == nil {
		t.Fatal("NewTreeStateCache (reopen) returned nil")
	}
	defer cache2.Close()

	// Verify the entry persisted (use tip=1100 so it's a deep block)
	retrieved := cache2.Get(950, 1100)
	if retrieved == nil {
		t.Fatal("Persisted entry not found after reopen")
	}
	if retrieved.Hash != ts.Hash {
		t.Errorf("Hash mismatch after reopen")
	}
	if retrieved.SaplingTree != ts.SaplingTree {
		t.Errorf("SaplingTree mismatch after reopen")
	}
}

func TestTreeStateCacheDisabled(t *testing.T) {
	// When windowSize is 0, NewTreeStateCache should return nil
	cache := NewTreeStateCache("/tmp", "testnet", 0, 1000)
	if cache != nil {
		t.Error("Expected nil cache when windowSize is 0")
	}
}

func TestTreeStateCacheFiles(t *testing.T) {
	logger := logrus.New()
	Log = logger.WithFields(logrus.Fields{"app": "test"})

	tmpDir, err := os.MkdirTemp("", "treestate_cache_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	cache := NewTreeStateCache(tmpDir, "testnet", 100, 1000)
	if cache == nil {
		t.Fatal("NewTreeStateCache returned nil")
	}

	// Check that cache files were created
	dataPath := filepath.Join(tmpDir, "testnet", "treestates.dat")
	indexPath := filepath.Join(tmpDir, "testnet", "treestates.idx")

	if _, err := os.Stat(dataPath); os.IsNotExist(err) {
		t.Error("Data file not created")
	}
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		t.Error("Index file not created")
	}

	cache.Close()
}

func TestTreeStateCacheTTLForRecentBlocks(t *testing.T) {
	logger := logrus.New()
	Log = logger.WithFields(logrus.Fields{"app": "test"})

	tmpDir, err := os.MkdirTemp("", "treestate_cache_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	cache := NewTreeStateCache(tmpDir, "testnet", 200, 1000)
	if cache == nil {
		t.Fatal("NewTreeStateCache returned nil")
	}
	defer cache.Close()

	// Cache a tree state at height 950
	ts := &walletrpc.TreeState{
		Network: "testnet",
		Height:  950,
		Hash:    "hash_950",
	}
	cache.Put(ts)

	// With currentTip=1000, height 950 has 50 confirmations (< 100, so it's "recent")
	// Should be cached initially (within 2-second TTL)
	if cache.Get(950, 1000) == nil {
		t.Error("Recent block should be cached immediately after Put")
	}

	// With currentTip=1100, height 950 has 150 confirmations (>= 100, so it's "deep")
	// Should always be cached (no TTL for deep blocks)
	if cache.Get(950, 1100) == nil {
		t.Error("Deep block should always be cached")
	}

	// Wait for TTL to expire
	time.Sleep(3 * time.Second)

	// With currentTip=1000, height 950 is still "recent" and TTL should have expired
	if cache.Get(950, 1000) != nil {
		t.Error("Recent block should expire after TTL")
	}

	// Re-cache and test that deep blocks don't expire
	cache.Put(ts)
	time.Sleep(3 * time.Second)

	// With currentTip=1100, height 950 is "deep" - should still be cached despite time passing
	if cache.Get(950, 1100) == nil {
		t.Error("Deep block should not expire based on TTL")
	}
}

func TestTreeStateCacheTTLBurstRequests(t *testing.T) {
	logger := logrus.New()
	Log = logger.WithFields(logrus.Fields{"app": "test"})

	tmpDir, err := os.MkdirTemp("", "treestate_cache_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	cache := NewTreeStateCache(tmpDir, "testnet", 200, 1000)
	if cache == nil {
		t.Fatal("NewTreeStateCache returned nil")
	}
	defer cache.Close()

	// Cache a recent block
	ts := &walletrpc.TreeState{
		Network: "testnet",
		Height:  990,
		Hash:    "hash_990",
	}
	cache.Put(ts)

	// Simulate burst requests within TTL window
	// With currentTip=1000, height 990 has 10 confirmations (recent)
	for i := 0; i < 5; i++ {
		if cache.Get(990, 1000) == nil {
			t.Errorf("Burst request %d should hit cache", i)
		}
	}
}
