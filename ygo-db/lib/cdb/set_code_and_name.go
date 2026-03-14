package cdb

// SetCodeAndName is a bidirectional map that maintains two-way mappings between uint64 and string.
// This is used for cdb setname code <-> name lookups.
// It ensures no duplicate keys on either side, but allows multiple codes to have the same name.
type SetCodeAndName struct {
	// codeToName: code -> name
	codeToName map[uint64]string
	// nameToCodes: name -> []code (supports multiple codes per name)
	nameToCodes map[string][]uint64
	// dedup: dedupKey -> code (deduplication key for AddWithDedup)
	dedup map[string]uint64
}

// NewSetCodeAndName creates a new SetCodeAndName instance.
func NewSetCodeAndName() *SetCodeAndName {
	return &SetCodeAndName{
		codeToName:  make(map[uint64]string),
		nameToCodes: make(map[string][]uint64),
		dedup:       make(map[string]uint64),
	}
}

// Add inserts a new key-value pair into the SetCodeAndName.
// It returns false if the key already exists in intToString.
// Multiple keys can map to the same string value (many-to-one).
func (scn *SetCodeAndName) Add(key uint64, value string) bool {
	// Check if key already exists in intToString
	if _, exists := scn.codeToName[key]; exists {
		return false
	}
	scn.codeToName[key] = value
	// Append to the slice for this value (supports many-to-one)
	scn.nameToCodes[value] = append(scn.nameToCodes[value], key)
	return true
}

// AddWithDedup inserts a new key-value pair with a custom deduplication key.
// The dedupKey is used for deduplication in the dedup map instead of value.
// This is useful when the setname has same local names but different Japanese name.
// It returns false if the key already exists in intToString OR if dedupKey already exists in dedup map.
func (scn *SetCodeAndName) AddWithDedup(key uint64, value string, dedupKey string) bool {
	// Check if key already exists in intToString
	if _, exists := scn.codeToName[key]; exists {
		return false
	}

	// Check if dedupKey already exists in dedup map
	if _, exists := scn.dedup[dedupKey]; exists {
		return false
	}

	scn.codeToName[key] = value
	// Add to stringToInt (supports one-to-many by value)
	scn.nameToCodes[value] = append(scn.nameToCodes[value], key)
	// Add to dedup map (unique dedupKey -> key)
	scn.dedup[dedupKey] = key
	return true
}

// GetByCode retrieves the setName by setCode.
// Returns empty string and false if not found.
func (scn *SetCodeAndName) GetByCode(key uint64) (string, bool) {
	val, ok := scn.codeToName[key]
	return val, ok
}

// GetByName retrieves all setCodes by setName.
// Returns empty slice and false if not found.
func (scn *SetCodeAndName) GetByName(value string) ([]uint64, bool) {
	keys, ok := scn.nameToCodes[value]
	return keys, ok
}

func (scn *SetCodeAndName) HasSetCode(key uint64) bool {
	_, ok := scn.codeToName[key]
	return ok
}

func (scn *SetCodeAndName) HasSetName(value string) bool {
	_, ok := scn.nameToCodes[value]
	return ok
}

// Len returns the number of entries in the SetCodeAndName.
func (scn *SetCodeAndName) Len() int {
	return len(scn.codeToName)
}

// Clear removes all entries from the SetCodeAndName.
func (scn *SetCodeAndName) Clear() {
	scn.codeToName = make(map[uint64]string)
	scn.nameToCodes = make(map[string][]uint64)
	scn.dedup = make(map[string]uint64)
}

// SetCodes returns all setCode.
func (scn *SetCodeAndName) SetCodes() []uint64 {
	keys := make([]uint64, 0, len(scn.codeToName))
	for k := range scn.codeToName {
		keys = append(keys, k)
	}
	return keys
}

// SetNames returns all unique setName.
func (scn *SetCodeAndName) SetNames() []string {
	values := make([]string, 0, len(scn.nameToCodes))
	for v := range scn.nameToCodes {
		values = append(values, v)
	}
	return values
}
