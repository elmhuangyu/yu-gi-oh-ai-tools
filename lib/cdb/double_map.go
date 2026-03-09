package cdb

// DoubleMap is a bidirectional map that maintains two-way mappings between int and string.
// This is used for cdb setname code <-> name lookups.
// It ensures no duplicate keys on either side, but allows multiple codes to have the same name.
type DoubleMap struct {
	// intToString: code -> name
	intToString map[int]string
	// stringToInt: name -> []code (supports multiple codes per name)
	stringToInt map[string][]int
	// dedup: dedupKey -> code (deduplication key for AddWithDedup)
	dedup map[string]int
}

// NewDoubleMap creates a new DoubleMap instance.
func NewDoubleMap() *DoubleMap {
	return &DoubleMap{
		intToString: make(map[int]string),
		stringToInt: make(map[string][]int),
		dedup:       make(map[string]int),
	}
}

// Add inserts a new key-value pair into the double map.
// It returns false if the key already exists in intToString.
// Multiple keys can map to the same string value (many-to-one).
func (dm *DoubleMap) Add(key int, value string) bool {
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
func (dm *DoubleMap) AddWithDedup(key int, value string, dedupKey string) bool {
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

// GetByInt retrieves the string value by int key.
// Returns empty string and false if not found.
func (dm *DoubleMap) GetByInt(key int) (string, bool) {
	val, ok := dm.intToString[key]
	return val, ok
}

// GetByString retrieves all int keys by string value.
// Returns empty slice and false if not found.
func (dm *DoubleMap) GetByString(value string) ([]int, bool) {
	keys, ok := dm.stringToInt[value]
	return keys, ok
}

// GetByStringFirst retrieves the first int key by string value.
// Returns 0 and false if not found.
// This is useful when you expect only one value per name.
func (dm *DoubleMap) GetByStringFirst(value string) (int, bool) {
	keys, ok := dm.stringToInt[value]
	if !ok || len(keys) == 0 {
		return 0, false
	}
	return keys[0], true
}

// HasInt checks if an int key exists.
func (dm *DoubleMap) HasInt(key int) bool {
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
	dm.intToString = make(map[int]string)
	dm.stringToInt = make(map[string][]int)
	dm.dedup = make(map[string]int)
}

// IntKeys returns all int keys.
func (dm *DoubleMap) IntKeys() []int {
	keys := make([]int, 0, len(dm.intToString))
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
