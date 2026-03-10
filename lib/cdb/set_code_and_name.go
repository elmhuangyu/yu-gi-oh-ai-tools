package cdb

// SetCodeAndName is a bidirectional map that maintains two-way mappings between uint64 and string.
// This is used for cdb setname code <-> name lookups.
// It ensures no duplicate keys on either side, but allows multiple codes to have the same name.
type SetCodeAndName struct {
	// intToString: code -> name
	intToString map[uint64]string
	// stringToInt: name -> []code (supports multiple codes per name)
	stringToInt map[string][]uint64
	// dedup: dedupKey -> code (deduplication key for AddWithDedup)
	dedup map[string]uint64
}

// NewSetCodeAndName creates a new SetCodeAndName instance.
func NewSetCodeAndName() *SetCodeAndName {
	return &SetCodeAndName{
		intToString: make(map[uint64]string),
		stringToInt: make(map[string][]uint64),
		dedup:       make(map[string]uint64),
	}
}

// Add inserts a new key-value pair into the SetCodeAndName.
// It returns false if the key already exists in intToString.
// Multiple keys can map to the same string value (many-to-one).
func (scn *SetCodeAndName) Add(key uint64, value string) bool {
	// Check if key already exists in intToString
	if _, exists := scn.intToString[key]; exists {
		return false
	}
	scn.intToString[key] = value
	// Append to the slice for this value (supports many-to-one)
	scn.stringToInt[value] = append(scn.stringToInt[value], key)
	return true
}

// AddWithDedup inserts a new key-value pair with a custom deduplication key.
// The dedupKey is used for deduplication in the dedup map instead of value.
// This is useful when the setname has same local names but different Japanese name.
// It returns false if the key already exists in intToString OR if dedupKey already exists in dedup map.
func (scn *SetCodeAndName) AddWithDedup(key uint64, value string, dedupKey string) bool {
	// Check if key already exists in intToString
	if _, exists := scn.intToString[key]; exists {
		return false
	}

	// Check if dedupKey already exists in dedup map
	if _, exists := scn.dedup[dedupKey]; exists {
		return false
	}

	scn.intToString[key] = value
	// Add to stringToInt (supports one-to-many by value)
	scn.stringToInt[value] = append(scn.stringToInt[value], key)
	// Add to dedup map (unique dedupKey -> key)
	scn.dedup[dedupKey] = key
	return true
}

// GetByUint64 retrieves the string value by uint64 key.
// Returns empty string and false if not found.
func (scn *SetCodeAndName) GetByUint64(key uint64) (string, bool) {
	val, ok := scn.intToString[key]
	return val, ok
}

// GetByString retrieves all uint64 keys by string value.
// Returns empty slice and false if not found.
func (scn *SetCodeAndName) GetByString(value string) ([]uint64, bool) {
	keys, ok := scn.stringToInt[value]
	return keys, ok
}

// GetByStringFirst retrieves the first uint64 key by string value.
// Returns 0 and false if not found.
// This is useful when you expect only one value per name.
func (scn *SetCodeAndName) GetByStringFirst(value string) (uint64, bool) {
	keys, ok := scn.stringToInt[value]
	if !ok || len(keys) == 0 {
		return 0, false
	}
	return keys[0], true
}

// HasUint64 checks if a uint64 key exists.
func (scn *SetCodeAndName) HasUint64(key uint64) bool {
	_, ok := scn.intToString[key]
	return ok
}

// HasString checks if a string value exists.
func (scn *SetCodeAndName) HasString(value string) bool {
	_, ok := scn.stringToInt[value]
	return ok
}

// Len returns the number of entries in the SetCodeAndName.
func (scn *SetCodeAndName) Len() int {
	return len(scn.intToString)
}

// Clear removes all entries from the SetCodeAndName.
func (scn *SetCodeAndName) Clear() {
	scn.intToString = make(map[uint64]string)
	scn.stringToInt = make(map[string][]uint64)
	scn.dedup = make(map[string]uint64)
}

// Uint64Keys returns all uint64 keys.
func (scn *SetCodeAndName) Uint64Keys() []uint64 {
	keys := make([]uint64, 0, len(scn.intToString))
	for k := range scn.intToString {
		keys = append(keys, k)
	}
	return keys
}

// StringValues returns all unique string values.
func (scn *SetCodeAndName) StringValues() []string {
	values := make([]string, 0, len(scn.stringToInt))
	for v := range scn.stringToInt {
		values = append(values, v)
	}
	return values
}
