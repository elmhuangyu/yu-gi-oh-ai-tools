package cdb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetCodeAndName_Add(t *testing.T) {
	tests := []struct {
		name    string
		key     uint64
		value   string
		wantOk  bool
		wantLen int
		wantInt map[uint64]string
		wantStr map[string]uint64
	}{
		{
			name:    "add first entry",
			key:     1,
			value:   "TestCard",
			wantOk:  true,
			wantLen: 1,
			wantInt: map[uint64]string{1: "TestCard"},
			wantStr: map[string]uint64{"TestCard": 1},
		},
		{
			name:    "add duplicate key",
			key:     1,
			value:   "AnotherCard",
			wantOk:  false,
			wantLen: 1,
			wantInt: map[uint64]string{1: "TestCard"},
			wantStr: map[string]uint64{"TestCard": 1},
		},
		{
			name:    "add second entry",
			key:     2,
			value:   "SecondCard",
			wantOk:  true,
			wantLen: 2,
			wantInt: map[uint64]string{1: "TestCard", 2: "SecondCard"},
			wantStr: map[string]uint64{"TestCard": 1, "SecondCard": 2},
		},
		{
			name:    "add duplicate value",
			key:     3,
			value:   "TestCard",
			wantOk:  true,
			wantLen: 3,
			wantInt: map[uint64]string{1: "TestCard", 2: "SecondCard", 3: "TestCard"},
			wantStr: map[string]uint64{"TestCard": 1, "SecondCard": 2},
		},
	}

	scn := NewSetCodeAndName()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := scn.Add(tt.key, tt.value)
			assert.Equal(t, tt.wantOk, got, "Add() = %v, want %v", got, tt.wantOk)
			assert.Equal(t, tt.wantLen, scn.Len(), "Len, want %v() = %v", scn.Len(), tt.wantLen)
		})
	}
}

func TestSetCodeAndName_GetByUint64(t *testing.T) {
	scn := NewSetCodeAndName()
	_ = scn.Add(1, "CardA")
	_ = scn.Add(2, "CardB")

	tests := []struct {
		name    string
		key     uint64
		wantVal string
		wantOk  bool
	}{
		{
			name:    "existing key",
			key:     1,
			wantVal: "CardA",
			wantOk:  true,
		},
		{
			name:    "another existing key",
			key:     2,
			wantVal: "CardB",
			wantOk:  true,
		},
		{
			name:    "non-existing key",
			key:     999,
			wantVal: "",
			wantOk:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, gotOk := scn.GetByUint64(tt.key)
			assert.Equal(t, tt.wantVal, gotVal, "GetByUint64(%d) val = %v, want %v", tt.key, gotVal, tt.wantVal)
			assert.Equal(t, tt.wantOk, gotOk, "GetByUint64(%d) ok = %v, want %v", tt.key, gotOk, tt.wantOk)
		})
	}
}

func TestSetCodeAndName_GetByString(t *testing.T) {
	scn := NewSetCodeAndName()
	_ = scn.Add(1, "CardA")
	_ = scn.Add(2, "CardB")
	_ = scn.Add(3, "CardA") // Add another entry with same name

	tests := []struct {
		name     string
		value    string
		wantKeys []uint64
		wantOk   bool
	}{
		{
			name:     "existing value",
			value:    "CardA",
			wantKeys: []uint64{1, 3},
			wantOk:   true,
		},
		{
			name:     "another existing value",
			value:    "CardB",
			wantKeys: []uint64{2},
			wantOk:   true,
		},
		{
			name:     "non-existing value",
			value:    "NonExistent",
			wantKeys: nil,
			wantOk:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotKeys, gotOk := scn.GetByString(tt.value)
			assert.Equal(t, tt.wantOk, gotOk, "GetByString(%s) ok = %v, want %v", tt.value, gotOk, tt.wantOk)
			assert.Equal(t, tt.wantKeys, gotKeys, "GetByString(%s) keys = %v, want %v", tt.value, gotKeys, tt.wantKeys)
		})
	}
}

func TestSetCodeAndName_GetByStringFirst(t *testing.T) {
	scn := NewSetCodeAndName()
	_ = scn.Add(1, "CardA")
	_ = scn.Add(2, "CardB")
	_ = scn.Add(3, "CardA") // Add another entry with same name

	tests := []struct {
		name    string
		value   string
		wantKey uint64
		wantOk  bool
	}{
		{
			name:    "existing value returns first key",
			value:   "CardA",
			wantKey: 1,
			wantOk:  true,
		},
		{
			name:    "another existing value",
			value:   "CardB",
			wantKey: 2,
			wantOk:  true,
		},
		{
			name:    "non-existing value",
			value:   "NonExistent",
			wantKey: 0,
			wantOk:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotKey, gotOk := scn.GetByStringFirst(tt.value)
			assert.Equal(t, tt.wantKey, gotKey, "GetByStringFirst(%s) key = %v, want %v", tt.value, gotKey, tt.wantKey)
			assert.Equal(t, tt.wantOk, gotOk, "GetByStringFirst(%s) ok = %v, want %v", tt.value, gotOk, tt.wantOk)
		})
	}
}

func TestSetCodeAndName_HasUint64(t *testing.T) {
	scn := NewSetCodeAndName()
	_ = scn.Add(1, "CardA")

	tests := []struct {
		name string
		key  uint64
		want bool
	}{
		{
			name: "existing key",
			key:  1,
			want: true,
		},
		{
			name: "non-existing key",
			key:  2,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := scn.HasUint64(tt.key)
			assert.Equal(t, tt.want, got, "HasUint64(%d) = %v, want %v", tt.key, got, tt.want)
		})
	}
}

func TestSetCodeAndName_HasString(t *testing.T) {
	scn := NewSetCodeAndName()
	_ = scn.Add(1, "CardA")

	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{
			name:  "existing value",
			value: "CardA",
			want:  true,
		},
		{
			name:  "non-existing value",
			value: "CardB",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := scn.HasString(tt.value)
			assert.Equal(t, tt.want, got, "HasString(%s) = %v, want %v", tt.value, got, tt.want)
		})
	}
}

func TestSetCodeAndName_Len(t *testing.T) {
	scn := NewSetCodeAndName()
	assert.Equal(t, 0, scn.Len(), "empty map should have length 0")

	_ = scn.Add(1, "CardA")
	assert.Equal(t, 1, scn.Len(), "map should have length 1")

	_ = scn.Add(2, "CardB")
	assert.Equal(t, 2, scn.Len(), "map should have length 2")

	// Adding duplicate should not change length
	_ = scn.Add(1, "CardC")
	assert.Equal(t, 2, scn.Len(), "map should still have length 2 after failed add")
}

func TestSetCodeAndName_Clear(t *testing.T) {
	scn := NewSetCodeAndName()
	_ = scn.Add(1, "CardA")
	_ = scn.Add(2, "CardB")

	scn.Clear()

	assert.Equal(t, 0, scn.Len(), "map should have length 0 after clear")
	assert.False(t, scn.HasUint64(1), "HasUint64(1) should return false after clear")
	assert.False(t, scn.HasString("CardA"), "HasString(CardA) should return false after clear")
}

func TestSetCodeAndName_Uint64Keys(t *testing.T) {
	scn := NewSetCodeAndName()
	_ = scn.Add(1, "CardA")
	_ = scn.Add(2, "CardB")
	_ = scn.Add(3, "CardC")

	keys := scn.Uint64Keys()
	assert.Len(t, keys, 3, "should have 3 keys")

	// Check all expected keys are present
	keySet := make(map[uint64]bool)
	for _, k := range keys {
		keySet[k] = true
	}
	assert.True(t, keySet[1], "key 1 should be present")
	assert.True(t, keySet[2], "key 2 should be present")
	assert.True(t, keySet[3], "key 3 should be present")
}

func TestSetCodeAndName_StringValues(t *testing.T) {
	scn := NewSetCodeAndName()
	_ = scn.Add(1, "CardA")
	_ = scn.Add(2, "CardB")
	_ = scn.Add(3, "CardC")

	values := scn.StringValues()
	assert.Len(t, values, 3, "should have 3 values")

	// Check all expected values are present
	valueSet := make(map[string]bool)
	for _, v := range values {
		valueSet[v] = true
	}
	assert.True(t, valueSet["CardA"], "value CardA should be present")
	assert.True(t, valueSet["CardB"], "value CardB should be present")
	assert.True(t, valueSet["CardC"], "value CardC should be present")
}
