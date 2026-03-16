package cdb

import (
	"strings"

	"github.com/moznion/go-optional"
)

// Language constants
const (
	LangZhCN = "zh-CN"
	LangEnUS = "en-US"
)

var (
	SupportLanguages = []string{LangZhCN, LangEnUS}
)

// TypeEntry represents a type with multilingual names
type TypeEntry struct {
	Value    uint64 `json:"value"`
	ZhCNName string `json:"zhCNName"`
	EnUSName string `json:"enUSName"`
}

// AttributeEntry represents an attribute with multilingual names
type AttributeEntry struct {
	Value    int    `json:"value"`
	ZhCNName string `json:"zhCNName"`
	EnUSName string `json:"enUSName"`
}

// RaceEntry represents a race with multilingual names
type RaceEntry struct {
	Value    int    `json:"value"`
	ZhCNName string `json:"zhCNName"`
	EnUSName string `json:"enUSName"`
}

const (
	TypeMonster uint64 = 0x1
)

// TypeMap stores all type entries with embedded constants
var TypeMap = []TypeEntry{
	{0x1, "怪兽卡", "Monster"},
	{0x2, "魔法卡", "Spell"},
	{0x4, "陷阱卡", "Trap"},
	{0x10, "通常怪兽", "Normal"},
	{0x20, "效果", "Effect"},
	{0x40, "融合", "Fusion"},
	{0x80, "仪式", "Ritual"},
	{0x100, "陷阱怪兽", "Trap Monster"},
	{0x200, "灵魂", "Spirit"},
	{0x400, "同盟", "Union"},
	{0x800, "二重", "Dual"},
	{0x1000, "调整", "Tuner"},
	{0x2000, "同调", "Synchro"},
	{0x4000, "衍生物", "Token"},
	{0x10000, "速攻", "Quickplay"},
	{0x20000, "永续", "Continuous"},
	{0x40000, "装备", "Equip"},
	{0x80000, "场地", "Field"},
	{0x100000, "反击", "Counter"},
	{0x200000, "翻转", "Flip"},
	{0x400000, "卡通", "Toon"},
	{0x800000, "超量", "Xyz"},
	{0x1000000, "灵摆", "Pendulum"},
	{0x2000000, "特殊召唤", "Special Summon"},
	{0x4000000, "连接", "Link"},
}

// AttributeMap stores all attribute entries with embedded constants
var AttributeMap = []AttributeEntry{
	{0x01, "地", "Earth"},
	{0x02, "水", "Water"},
	{0x04, "炎", "Fire"},
	{0x08, "风", "Wind"},
	{0x10, "光", "Light"},
	{0x20, "暗", "Dark"},
	{0x40, "神", "Divine"},
}

// RaceMap stores all race entries with embedded constants
var RaceMap = []RaceEntry{
	{0x1, "战士", "Warrior"},
	{0x2, "魔法师", "Spellcaster"},
	{0x4, "天使", "Fairy"},
	{0x8, "恶魔", "Fiend"},
	{0x10, "不死", "Zombie"},
	{0x20, "机械", "Machine"},
	{0x40, "水", "Aqua"},
	{0x80, "炎", "Pyro"},
	{0x100, "岩石", "Rock"},
	{0x200, "鸟兽", "Winged Beast"},
	{0x400, "植物", "Plant"},
	{0x800, "昆虫", "Insect"},
	{0x1000, "雷", "Thunder"},
	{0x2000, "龙", "Dragon"},
	{0x4000, "兽", "Beast"},
	{0x8000, "兽战士", "Beast-Warrior"},
	{0x10000, "恐龙", "Dinosaur"},
	{0x20000, "鱼", "Fish"},
	{0x40000, "海龙", "Sea Serpent"},
	{0x80000, "爬虫", "Reptile"},
	{0x100000, "念动力", "Psychic"},
	{0x200000, "幻神兽", "Divine-Beast"},
	{0x400000, "创造神", "Creator God"},
	{0x800000, "幻龙", "Wyrm"},
	{0x1000000, "电子界", "Cyberse"},
	{0x2000000, "幻想魔", "Illusion"},
}

// GetTypeByName returns the type value by name in the specified language
func GetTypeByName(lang, name string) optional.Option[uint64] {
	name = strings.TrimSpace(strings.ToLower(name))
	for _, entry := range TypeMap {
		var entryName string
		switch lang {
		case LangZhCN:
			entryName = strings.ToLower(entry.ZhCNName)
		case LangEnUS:
			entryName = strings.ToLower(entry.EnUSName)
		default:
			entryName = strings.ToLower(entry.ZhCNName) // default to zh-CN
		}
		if entryName == name {
			return optional.Some(entry.Value)
		}
	}
	return optional.None[uint64]()
}

// GetTypeNames returns all type names for the given type value in the specified language
func GetTypeNames(lang string, t uint64) []string {
	if t == 0 {
		return []string{}
	}
	var names []string
	for _, entry := range TypeMap {
		if t&entry.Value != 0 {
			var name string
			switch lang {
			case LangZhCN:
				name = entry.ZhCNName
			case LangEnUS:
				name = entry.EnUSName
			default:
				name = entry.ZhCNName // default to zh-CN
			}
			names = append(names, name)
		}
	}
	return names
}

// GetAttributeByName returns the attribute value by name in the specified language
func GetAttributeByName(lang, name string) optional.Option[int] {
	name = strings.TrimSpace(strings.ToLower(name))
	for _, entry := range AttributeMap {
		var entryName string
		switch lang {
		case LangZhCN:
			entryName = strings.ToLower(entry.ZhCNName)
		case LangEnUS:
			entryName = strings.ToLower(entry.EnUSName)
		default:
			entryName = strings.ToLower(entry.ZhCNName) // default to zh-CN
		}
		if entryName == name {
			return optional.Some(entry.Value)
		}
	}
	return optional.None[int]()
}

// GetAttributeName returns the attribute name for the given attribute value in the specified language
func GetAttributeName(lang string, attr int) string {
	for _, entry := range AttributeMap {
		if entry.Value == attr {
			switch lang {
			case LangZhCN:
				return entry.ZhCNName
			case LangEnUS:
				return entry.EnUSName
			default:
				return entry.ZhCNName // default to zh-CN
			}
		}
	}
	return ""
}

// GetRaceByName returns the race value by name in the specified language
func GetRaceByName(lang, name string) optional.Option[int] {
	name = strings.TrimSpace(strings.ToLower(name))
	for _, entry := range RaceMap {
		var entryName string
		switch lang {
		case LangZhCN:
			entryName = strings.ToLower(entry.ZhCNName)
		case LangEnUS:
			entryName = strings.ToLower(entry.EnUSName)
		default:
			entryName = strings.ToLower(entry.ZhCNName) // default to zh-CN
		}
		if entryName == name {
			return optional.Some(entry.Value)
		}
	}
	return optional.None[int]()
}

// GetRaceName returns the race name for the given race value in the specified language
func GetRaceName(lang string, race int) string {
	for _, entry := range RaceMap {
		if entry.Value == race {
			switch lang {
			case LangZhCN:
				return entry.ZhCNName
			case LangEnUS:
				return entry.EnUSName
			default:
				return entry.ZhCNName // default to zh-CN
			}
		}
	}
	return ""
}

// ListTypes returns all type names in the specified language
func ListTypes(lang string) []string {
	var names []string
	for _, entry := range TypeMap {
		var name string
		switch lang {
		case LangZhCN:
			name = entry.ZhCNName
		case LangEnUS:
			name = entry.EnUSName
		default:
			name = entry.ZhCNName // default to zh-CN
		}
		names = append(names, name)
	}
	return names
}

// ListAttributes returns all attribute names in the specified language
func ListAttributes(lang string) []string {
	var names []string
	for _, entry := range AttributeMap {
		var name string
		switch lang {
		case LangZhCN:
			name = entry.ZhCNName
		case LangEnUS:
			name = entry.EnUSName
		default:
			name = entry.ZhCNName // default to zh-CN
		}
		names = append(names, name)
	}
	return names
}

// ListRaces returns all race names in the specified language
func ListRaces(lang string) []string {
	var names []string
	for _, entry := range RaceMap {
		var name string
		switch lang {
		case LangZhCN:
			name = entry.ZhCNName
		case LangEnUS:
			name = entry.EnUSName
		default:
			name = entry.ZhCNName // default to zh-CN
		}
		names = append(names, name)
	}
	return names
}
