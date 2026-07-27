package common

import (
	"encoding/binary"
	"hash/fnv"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/zcash/lightwalletd/walletrpc"
	"google.golang.org/protobuf/proto"
)

const (
	// recentBlockThreshold is the number of confirmations below which
	// blocks are considered "recent" and subject to TTL expiration.
	recentBlockThreshold = 100

	// recentBlockTTL is how long cached tree states for recent blocks
	// remain valid before requiring a fresh fetch from zcashd.
	recentBlockTTL = 2 * time.Second
)

// TreeStateCache provides disk-persisted caching for tree states within a configurable
// block window. Tree states for blocks older than the window are not cached.
type TreeStateCache struct {
	dataDir   string
	chainName string
	dataFile  *os.File
	indexFile *os.File

	// In-memory index: height -> entry
	heightIndex map[uint64]*treeStateCacheEntry

	// Configuration
	windowSize int // max blocks to cache (default 525600 = ~1 year)
	minHeight  int // don't cache below this height

	mutex sync.RWMutex
}

// treeStateCacheEntry represents a cached tree state entry
type treeStateCacheEntry struct {
	offset   int64     // offset in data file
	length   int       // length of serialized data (not including checksum)
	hash     string    // block hash for hash-based lookups
	cachedAt time.Time // when this entry was cached (for TTL)
}

// Index file entry format:
// - height: 8 bytes (uint64, little endian)
// - offset: 8 bytes (int64, little endian)
// - length: 4 bytes (uint32, little endian)
// - hashLen: 2 bytes (uint16, little endian)
// - hash: variable bytes
const indexEntryBaseSize = 8 + 8 + 4 + 2 // 22 bytes + hash

// NewTreeStateCache creates a new tree state cache.
// windowSize is the number of blocks to cache (0 = disabled).
// currentHeight is the current chain tip height.
func NewTreeStateCache(dbPath string, chainName string, windowSize int, currentHeight int) *TreeStateCache {
	if windowSize <= 0 {
		return nil
	}

	c := &TreeStateCache{
		dataDir:     dbPath,
		chainName:   chainName,
		heightIndex: make(map[uint64]*treeStateCacheEntry),
		windowSize:  windowSize,
	}

	// Calculate minimum height to cache
	c.minHeight = currentHeight - windowSize
	if c.minHeight < 0 {
		c.minHeight = 0
	}

	// Create directory if needed
	dir := filepath.Join(dbPath, chainName)
	if err := os.MkdirAll(dir, 0755); err != nil {
		Log.Fatal("mkdir ", dir, " failed: ", err)
	}

	// Open data file
	dataPath := filepath.Join(dir, "treestates.dat")
	var err error
	c.dataFile, err = os.OpenFile(dataPath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		Log.Fatal("open ", dataPath, " failed: ", err)
	}

	// Open index file
	indexPath := filepath.Join(dir, "treestates.idx")
	c.indexFile, err = os.OpenFile(indexPath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		Log.Fatal("open ", indexPath, " failed: ", err)
	}

	// Load existing index
	c.loadIndex()

	return c
}

// loadIndex reads the index file and populates the in-memory index
func (c *TreeStateCache) loadIndex() {
	indexData, err := os.ReadFile(c.indexFile.Name())
	if err != nil {
		Log.Warning("failed to read tree state index: ", err)
		return
	}

	offset := 0
	entriesLoaded := 0
	entriesPruned := 0
	corrupted := false

	for offset < len(indexData) {
		if offset+indexEntryBaseSize > len(indexData) {
			Log.Warning("tree state index truncated at offset ", offset)
			corrupted = true
			break
		}

		height := binary.LittleEndian.Uint64(indexData[offset:])
		fileOffset := int64(binary.LittleEndian.Uint64(indexData[offset+8:]))
		length := int(binary.LittleEndian.Uint32(indexData[offset+16:]))
		hashLen := int(binary.LittleEndian.Uint16(indexData[offset+20:]))

		// Sanity checks
		if length < 0 || length > 1000000 {
			Log.Warning("tree state index has invalid length ", length, " at offset ", offset)
			corrupted = true
			break
		}
		if hashLen < 0 || hashLen > 1000 {
			Log.Warning("tree state index has invalid hashLen ", hashLen, " at offset ", offset)
			corrupted = true
			break
		}
		if fileOffset < 0 {
			Log.Warning("tree state index has invalid file offset ", fileOffset, " at offset ", offset)
			corrupted = true
			break
		}

		if offset+indexEntryBaseSize+hashLen > len(indexData) {
			Log.Warning("tree state index hash truncated at offset ", offset)
			corrupted = true
			break
		}

		hash := string(indexData[offset+22 : offset+22+hashLen])
		offset += indexEntryBaseSize + hashLen

		// Only load entries within the current window
		if int(height) >= c.minHeight {
			c.heightIndex[height] = &treeStateCacheEntry{
				offset:   fileOffset,
				length:   length,
				hash:     hash,
				cachedAt: time.Now(), // Treat as freshly cached on load
			}
			entriesLoaded++
		} else {
			entriesPruned++
		}
	}

	if corrupted {
		Log.Warning("TreeStateCache: index corruption detected during load, clearing cache")
		c.heightIndex = make(map[uint64]*treeStateCacheEntry)
		c.dataFile.Truncate(0)
		c.indexFile.Truncate(0)
		c.dataFile.Sync()
		c.indexFile.Sync()
		return
	}

	if entriesLoaded > 0 || entriesPruned > 0 {
		Log.Info("TreeStateCache: loaded ", entriesLoaded, " entries, pruned ", entriesPruned, " entries outside window")
	}
}

// treeStateChecksum calculates the checksum for a tree state entry
func treeStateChecksum(height uint64, data []byte) []byte {
	h := make([]byte, 8)
	binary.LittleEndian.PutUint64(h, height)
	cs := fnv.New64a()
	cs.Write(h)
	cs.Write(data)
	return cs.Sum(nil)
}

// Get retrieves a tree state by height from the cache.
// Returns nil if not found, if height is outside the cache window,
// or if it's a recent block and the TTL has expired.
// currentTip is the current chain tip height, used for TTL checks on recent blocks.
func (c *TreeStateCache) Get(height uint64, currentTip int) *walletrpc.TreeState {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Check if height is within window
	if int(height) < c.minHeight {
		return nil
	}

	entry, ok := c.heightIndex[height]
	if !ok {
		return nil
	}

	// For recent blocks (< 100 confirmations), check TTL
	if currentTip-int(height) < recentBlockThreshold {
		if time.Since(entry.cachedAt) > recentBlockTTL {
			// TTL expired - delete entry and return nil to force fresh fetch
			delete(c.heightIndex, height)
			return nil
		}
	}

	return c.readEntry(height, entry)
}

// GetByHash retrieves a tree state by block hash from the cache.
// Returns nil if not found or if TTL expired for recent blocks.
// currentTip is the current chain tip height, used for TTL checks on recent blocks.
func (c *TreeStateCache) GetByHash(hash string, currentTip int) *walletrpc.TreeState {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Linear scan through index (hash lookups are less common)
	for height, entry := range c.heightIndex {
		if entry.hash == hash {
			// For recent blocks (< 100 confirmations), check TTL
			if currentTip-int(height) < recentBlockThreshold {
				if time.Since(entry.cachedAt) > recentBlockTTL {
					// TTL expired - delete entry and return nil
					delete(c.heightIndex, height)
					return nil
				}
			}
			return c.readEntry(height, entry)
		}
	}
	return nil
}

// readEntry reads a tree state from the data file
// Caller must hold c.mutex.Lock() (not RLock) because corruption triggers reset
func (c *TreeStateCache) readEntry(height uint64, entry *treeStateCacheEntry) *walletrpc.TreeState {
	// Read data from file
	buf := make([]byte, entry.length+8) // +8 for checksum
	n, err := c.dataFile.ReadAt(buf, entry.offset)
	if err != nil || n != len(buf) {
		Log.Warning("TreeStateCache: read failed at offset ", entry.offset, ": ", err)
		c.handleCorruption("read failure")
		return nil
	}

	// Verify checksum
	diskChecksum := buf[:8]
	data := buf[8:]
	if string(treeStateChecksum(height, data)) != string(diskChecksum) {
		Log.Warning("TreeStateCache: checksum mismatch at height ", height)
		c.handleCorruption("checksum mismatch")
		return nil
	}

	// Unmarshal
	ts := &walletrpc.TreeState{}
	if err := proto.Unmarshal(data, ts); err != nil {
		Log.Warning("TreeStateCache: unmarshal failed at height ", height, ": ", err)
		c.handleCorruption("unmarshal failure")
		return nil
	}

	return ts
}

// handleCorruption deletes the cache files and resets to empty state.
// Caller must hold c.mutex.Lock()
func (c *TreeStateCache) handleCorruption(reason string) {
	Log.Warning("TreeStateCache: corruption detected (", reason, "), deleting cache and restarting fresh")

	// Clear in-memory index
	c.heightIndex = make(map[uint64]*treeStateCacheEntry)

	// Truncate both files to zero
	if c.dataFile != nil {
		c.dataFile.Truncate(0)
		c.dataFile.Sync()
	}
	if c.indexFile != nil {
		c.indexFile.Truncate(0)
		c.indexFile.Sync()
	}

	Log.Info("TreeStateCache: cache cleared, will rebuild from RPC requests")
}

// Put stores a tree state in the cache.
// Only caches if the height is within the current window.
func (c *TreeStateCache) Put(ts *walletrpc.TreeState) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	height := ts.Height

	// Don't cache if outside window
	if int(height) < c.minHeight {
		return nil
	}

	// Don't cache if already present
	if _, ok := c.heightIndex[height]; ok {
		return nil
	}

	// Serialize the tree state
	data, err := proto.Marshal(ts)
	if err != nil {
		return err
	}

	// Get current file position (end of file)
	fileInfo, err := c.dataFile.Stat()
	if err != nil {
		return err
	}
	offset := fileInfo.Size()

	// Write checksum + data to data file
	checksum := treeStateChecksum(height, data)
	if _, err := c.dataFile.WriteAt(checksum, offset); err != nil {
		return err
	}
	if _, err := c.dataFile.WriteAt(data, offset+8); err != nil {
		return err
	}

	// Write index entry
	indexEntry := make([]byte, indexEntryBaseSize+len(ts.Hash))
	binary.LittleEndian.PutUint64(indexEntry[0:], height)
	binary.LittleEndian.PutUint64(indexEntry[8:], uint64(offset))
	binary.LittleEndian.PutUint32(indexEntry[16:], uint32(len(data)))
	binary.LittleEndian.PutUint16(indexEntry[20:], uint16(len(ts.Hash)))
	copy(indexEntry[22:], ts.Hash)

	indexInfo, err := c.indexFile.Stat()
	if err != nil {
		return err
	}
	if _, err := c.indexFile.WriteAt(indexEntry, indexInfo.Size()); err != nil {
		return err
	}

	// Update in-memory index
	c.heightIndex[height] = &treeStateCacheEntry{
		offset:   offset,
		length:   len(data),
		hash:     ts.Hash,
		cachedAt: time.Now(),
	}

	return nil
}

// Reorg invalidates all cached entries at or above the given height.
// This should be called when a blockchain reorganization is detected.
func (c *TreeStateCache) Reorg(height int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Remove entries at or above the reorg height
	removed := 0
	for h := range c.heightIndex {
		if int(h) >= height {
			delete(c.heightIndex, h)
			removed++
		}
	}

	if removed > 0 {
		Log.Info("TreeStateCache: invalidated ", removed, " entries due to reorg at height ", height)
		// Rebuild the index file with remaining entries
		c.rebuildIndex()
	}
}

// UpdateWindow updates the minimum cached height based on the current chain tip.
// This should be called periodically as the chain grows.
func (c *TreeStateCache) UpdateWindow(currentHeight int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	newMinHeight := currentHeight - c.windowSize
	if newMinHeight < 0 {
		newMinHeight = 0
	}

	if newMinHeight <= c.minHeight {
		return // Window hasn't moved
	}

	// Remove entries that are now outside the window
	removed := 0
	for h := range c.heightIndex {
		if int(h) < newMinHeight {
			delete(c.heightIndex, h)
			removed++
		}
	}

	c.minHeight = newMinHeight

	if removed > 0 {
		Log.Info("TreeStateCache: pruned ", removed, " entries outside new window (minHeight=", newMinHeight, ")")
		c.rebuildIndex()
	}
}

// rebuildIndex rewrites the index file with only the current in-memory entries.
// Caller must hold c.mutex.Lock()
func (c *TreeStateCache) rebuildIndex() {
	// Truncate and rewrite index file
	if err := c.indexFile.Truncate(0); err != nil {
		Log.Warning("TreeStateCache: failed to truncate index: ", err)
		return
	}

	offset := int64(0)
	for height, entry := range c.heightIndex {
		indexEntry := make([]byte, indexEntryBaseSize+len(entry.hash))
		binary.LittleEndian.PutUint64(indexEntry[0:], height)
		binary.LittleEndian.PutUint64(indexEntry[8:], uint64(entry.offset))
		binary.LittleEndian.PutUint32(indexEntry[16:], uint32(entry.length))
		binary.LittleEndian.PutUint16(indexEntry[20:], uint16(len(entry.hash)))
		copy(indexEntry[22:], entry.hash)

		if _, err := c.indexFile.WriteAt(indexEntry, offset); err != nil {
			Log.Warning("TreeStateCache: failed to write index entry: ", err)
			return
		}
		offset += int64(len(indexEntry))
	}

	c.indexFile.Sync()
}

// GetMinHeight returns the minimum height that will be cached.
func (c *TreeStateCache) GetMinHeight() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.minHeight
}

// Sync flushes the cache files to disk.
func (c *TreeStateCache) Sync() {
	c.dataFile.Sync()
	c.indexFile.Sync()
}

// Close closes the cache files.
func (c *TreeStateCache) Close() {
	if c.dataFile != nil {
		c.dataFile.Close()
		c.dataFile = nil
	}
	if c.indexFile != nil {
		c.indexFile.Close()
		c.indexFile = nil
	}
}
