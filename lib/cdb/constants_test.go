package cdb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTypeByName(t *testing.T) {
	tests := []struct {
		name     string
		lang     string
		input    string
		expected uint64
	}{
		{"Monster in English", LangEnUS, "Monster", 0x1},
		{"Monster in Chinese", LangZhCN, "怪兽卡", 0x1},
		{"Spell in English", LangEnUS, "Spell", 0x2},
		{"Spell in Chinese", LangZhCN, "魔法卡", 0x2},
		{"Trap in English", LangEnUS, "Trap", 0x4},
		{"Trap in Chinese", LangZhCN, "陷阱卡", 0x4},
		{"Normal in English", LangEnUS, "Normal", 0x10},
		{"Normal in Chinese", LangZhCN, "通常怪兽", 0x10},
		{"Effect in English", LangEnUS, "Effect", 0x20},
		{"Effect in Chinese", LangZhCN, "效果", 0x20},
		{"Case insensitive English", LangEnUS, "  monster  ", 0x1},
		{"Case insensitive Chinese", LangZhCN, "  怪兽卡  ", 0x1},
		{"Not found", LangEnUS, "Unknown", 0},
		{"Not found Chinese", LangZhCN, "不存在的类型", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetTypeByName(tt.lang, tt.input)
			if tt.expected == 0 {
				assert.False(t, got.IsSome(), "expected empty optional for %s", tt.input)
			} else {
				assert.True(t, got.IsSome(), "expected present optional for %s", tt.input)
				assert.Equal(t, tt.expected, got.Unwrap())
			}
		})
	}
}

func TestGetTypeNames(t *testing.T) {
	tests := []struct {
		name     string
		lang     string
		input    uint64
		expected []string
	}{
		{
			"Single type English",
			LangEnUS,
			0x1,
			[]string{"Monster"},
		},
		{
			"Single type Chinese",
			LangZhCN,
			0x1,
			[]string{"怪兽卡"},
		},
		{
			"Multiple types English",
			LangEnUS,
			0x1 | 0x10 | 0x20,
			[]string{"Monster", "Normal", "Effect"},
		},
		{
			"Multiple types Chinese",
			LangZhCN,
			0x1 | 0x10 | 0x20,
			[]string{"怪兽卡", "通常怪兽", "效果"},
		},
		{
			"Zero value",
			LangZhCN,
			0,
			[]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetTypeNames(tt.lang, tt.input)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestGetAttributeByName(t *testing.T) {
	tests := []struct {
		name     string
		lang     string
		input    string
		expected int
	}{
		{"Earth in English", LangEnUS, "Earth", 0x01},
		{"Earth in Chinese", LangZhCN, "地", 0x01},
		{"Water in English", LangEnUS, "Water", 0x02},
		{"Water in Chinese", LangZhCN, "水", 0x02},
		{"Fire in English", LangEnUS, "Fire", 0x04},
		{"Fire in Chinese", LangZhCN, "炎", 0x04},
		{"Wind in English", LangEnUS, "Wind", 0x08},
		{"Wind in Chinese", LangZhCN, "风", 0x08},
		{"Light in English", LangEnUS, "Light", 0x10},
		{"Light in Chinese", LangZhCN, "光", 0x10},
		{"Dark in English", LangEnUS, "Dark", 0x20},
		{"Dark in Chinese", LangZhCN, "暗", 0x20},
		{"Divine in English", LangEnUS, "Divine", 0x40},
		{"Divine in Chinese", LangZhCN, "神", 0x40},
		{"Case insensitive English", LangEnUS, "  earth  ", 0x01},
		{"Case insensitive Chinese", LangZhCN, "  地  ", 0x01},
		{"Not found", LangEnUS, "Unknown", 0},
		{"Not found Chinese", LangZhCN, "不存在的属性", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetAttributeByName(tt.lang, tt.input)
			if tt.expected == 0 {
				assert.False(t, got.IsSome(), "expected empty optional for %s", tt.input)
			} else {
				assert.True(t, got.IsSome(), "expected present optional for %s", tt.input)
				assert.Equal(t, tt.expected, got.Unwrap())
			}
		})
	}
}

func TestGetAttributeName(t *testing.T) {
	tests := []struct {
		name     string
		lang     string
		input    int
		expected string
	}{
		{"Earth English", LangEnUS, 0x01, "Earth"},
		{"Earth Chinese", LangZhCN, 0x01, "地"},
		{"Water English", LangEnUS, 0x02, "Water"},
		{"Water Chinese", LangZhCN, 0x02, "水"},
		{"Fire English", LangEnUS, 0x04, "Fire"},
		{"Fire Chinese", LangZhCN, 0x04, "炎"},
		{"Wind English", LangEnUS, 0x08, "Wind"},
		{"Wind Chinese", LangZhCN, 0x08, "风"},
		{"Light English", LangEnUS, 0x10, "Light"},
		{"Light Chinese", LangZhCN, 0x10, "光"},
		{"Dark English", LangEnUS, 0x20, "Dark"},
		{"Dark Chinese", LangZhCN, 0x20, "暗"},
		{"Divine English", LangEnUS, 0x40, "Divine"},
		{"Divine Chinese", LangZhCN, 0x40, "神"},
		{"Not found", LangZhCN, 0x80, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetAttributeName(tt.lang, tt.input)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestGetRaceByName(t *testing.T) {
	tests := []struct {
		name     string
		lang     string
		input    string
		expected int
	}{
		{"Warrior in English", LangEnUS, "Warrior", 0x1},
		{"Warrior in Chinese", LangZhCN, "战士", 0x1},
		{"Spellcaster in English", LangEnUS, "Spellcaster", 0x2},
		{"Spellcaster in Chinese", LangZhCN, "魔法师", 0x2},
		{"Fairy in English", LangEnUS, "Fairy", 0x4},
		{"Fairy in Chinese", LangZhCN, "天使", 0x4},
		{"Dragon in English", LangEnUS, "Dragon", 0x2000},
		{"Dragon in Chinese", LangZhCN, "龙", 0x2000},
		{"Case insensitive English", LangEnUS, "  warrior  ", 0x1},
		{"Case insensitive Chinese", LangZhCN, "  战士  ", 0x1},
		{"Not found", LangEnUS, "Unknown", 0},
		{"Not found Chinese", LangZhCN, "不存在的种族", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetRaceByName(tt.lang, tt.input)
			if tt.expected == 0 {
				assert.False(t, got.IsSome(), "expected empty optional for %s", tt.input)
			} else {
				assert.True(t, got.IsSome(), "expected present optional for %s", tt.input)
				assert.Equal(t, tt.expected, got.Unwrap())
			}
		})
	}
}

func TestGetRaceName(t *testing.T) {
	tests := []struct {
		name     string
		lang     string
		input    int
		expected string
	}{
		{"Warrior English", LangEnUS, 0x1, "Warrior"},
		{"Warrior Chinese", LangZhCN, 0x1, "战士"},
		{"Spellcaster English", LangEnUS, 0x2, "Spellcaster"},
		{"Spellcaster Chinese", LangZhCN, 0x2, "魔法师"},
		{"Fairy English", LangEnUS, 0x4, "Fairy"},
		{"Fairy Chinese", LangZhCN, 0x4, "天使"},
		{"Dragon English", LangEnUS, 0x2000, "Dragon"},
		{"Dragon Chinese", LangZhCN, 0x2000, "龙"},
		{"Not found", LangZhCN, 0x4000000, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetRaceName(tt.lang, tt.input)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestListTypes(t *testing.T) {
	tests := []struct {
		name     string
		lang     string
		expected []string
	}{
		{
			"English",
			LangEnUS,
			func() []string {
				names := make([]string, len(TypeMap))
				for i, e := range TypeMap {
					names[i] = e.EnUSName
				}
				return names
			}(),
		},
		{
			"Chinese",
			LangZhCN,
			func() []string {
				names := make([]string, len(TypeMap))
				for i, e := range TypeMap {
					names[i] = e.ZhCNName
				}
				return names
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ListTypes(tt.lang)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestListAttributes(t *testing.T) {
	tests := []struct {
		name     string
		lang     string
		expected []string
	}{
		{
			"English",
			LangEnUS,
			func() []string {
				names := make([]string, len(AttributeMap))
				for i, e := range AttributeMap {
					names[i] = e.EnUSName
				}
				return names
			}(),
		},
		{
			"Chinese",
			LangZhCN,
			func() []string {
				names := make([]string, len(AttributeMap))
				for i, e := range AttributeMap {
					names[i] = e.ZhCNName
				}
				return names
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ListAttributes(tt.lang)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestListRaces(t *testing.T) {
	tests := []struct {
		name     string
		lang     string
		expected []string
	}{
		{
			"English",
			LangEnUS,
			func() []string {
				names := make([]string, len(RaceMap))
				for i, e := range RaceMap {
					names[i] = e.EnUSName
				}
				return names
			}(),
		},
		{
			"Chinese",
			LangZhCN,
			func() []string {
				names := make([]string, len(RaceMap))
				for i, e := range RaceMap {
					names[i] = e.ZhCNName
				}
				return names
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ListRaces(tt.lang)
			assert.Equal(t, tt.expected, got)
		})
	}
}
