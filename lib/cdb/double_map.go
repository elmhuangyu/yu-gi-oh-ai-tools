package cdb

// DoubleMap is a bidirectional map that maintains two-way mappings between uint64 and string.
// This is used for cdb setname code <-> name lookups.
// It ensures no duplicate keys on either side, but allows multiple codes to have the same name.
type DoubleMap struct {
	// intToString: code -> name
	intToString map[uint64]string
	// stringToInt: name -> []code (supports multiple codes per name)
	stringToInt map[string][]uint64
	// dedup: dedupKey -> code (deduplication key for AddWithDedup)
	dedup map[string]uint64
}

// NewDoubleMap creates a new DoubleMap instance.
func NewDoubleMap() *DoubleMap {
	return &DoubleMap{
		intToString: make(map[uint64]string),
		stringToInt: make(map[string][]uint64),
		dedup:       make(map[string]uint64),
	}
}

// Add inserts a new key-value pair into the double map.
// It returns false if the key already exists in intToString.
// Multiple keys can map to the same string value (many-to-one).
func (dm *DoubleMap) Add(key uint64, value string) bool {
	// Check if key already exists in intToString
	if _, exists := dm.intToString[key]; exists {
		return false
	}
	dm.intToString[key] = value
	// Append to the slice for this value (supports many-to-one)
	dm.stringToInt[value] = append(dm.stringToInt[value], key)
	return true
}

// AddWithDedup inserts a new key-value pair with a custom deduplication key.
// The dedupKey is used for deduplication in the dedup map instead of value.
// This is useful when the setname has same local names but different Japanese name.
// It returns false if the key already exists in intToString OR if dedupKey already exists in dedup map.
func (dm *DoubleMap) AddWithDedup(key uint64, value string, dedupKey string) bool {
	// Check if key already exists in intToString
	if _, exists := dm.intToString[key]; exists {
		return false
	}

	// Check if dedupKey already exists in dedup map
	if _, exists := dm.dedup[dedupKey]; exists {
		return false
	}

	dm.intToString[key] = value
	// Add to stringToInt (supports one-to-many by value)
	dm.stringToInt[value] = append(dm.stringToInt[value], key)
	// Add to dedup map (unique dedupKey -> key)
	dm.dedup[dedupKey] = key
	return true
}

// GetByUint64 retrieves the string value by uint64 key.
// Returns empty string and false if not found.
func (dm *DoubleMap) GetByUint64(key uint64) (string, bool) {
	val, ok := dm.intToString[key]
	return val, ok
}

// GetByString retrieves all uint64 keys by string value.
// Returns empty slice and false if not found.
func (dm *DoubleMap) GetByString(value string) ([]uint64, bool) {
	keys, ok := dm.stringToInt[value]
	return keys, ok
}

// GetByStringFirst retrieves the first uint64 key by string value.
// Returns 0 and false if not found.
// This is useful when you expect only one value per name.
func (dm *DoubleMap) GetByStringFirst(value string) (uint64, bool) {
	keys, ok := dm.stringToInt[value]
	if !ok || len(keys) == 0 {
		return 0, false
	}
	return keys[0], true
}

// HasUint64 checks if a uint64 key exists.
func (dm *DoubleMap) HasUint64(key uint64) bool {
	_, ok := dm.intToString[key]
	return ok
}

// HasString checks if a string value exists.
func (dm *DoubleMap) HasString(value string) bool {
	_, ok := dm.stringToInt[value]
	return ok
}

// Len returns the number of entries in the double map.
func (dm *DoubleMap) Len() int {
	return len(dm.intToString)
}

// Clear removes all entries from the double map.
func (dm *DoubleMap) Clear() {
	dm.intToString = make(map[uint64]string)
	dm.stringToInt = make(map[string][]uint64)
	dm.dedup = make(map[string]uint64)
}

// Uint64Keys returns all uint64 keys.
func (dm *DoubleMap) Uint64Keys() []uint64 {
	keys := make([]uint64, 0, len(dm.intToString))
	for k := range dm.intToString {
		keys = append(keys, k)
	}
	return keys
}

// StringValues returns all unique string values.
func (dm *DoubleMap) StringValues() []string {
	values := make([]string, 0, len(dm.stringToInt))
	for v := range dm.stringToInt {
		values = append(values, v)
	}
	return values
}
