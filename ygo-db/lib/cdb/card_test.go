package cdb

import (
	"testing"

	"github.com/elmhuangyu/yu-gi-oh-ai-tools/ygo-db/lib/git"
	"github.com/moznion/go-optional"
	"github.com/stretchr/testify/assert"
)

func TestCardInfoForHuman_ToCardInfoForAI(t *testing.T) {
	input := &CardInfoForHuman{
		ID:        48486809,
		Name:      "羽翼栗子球 LV6",
		Desc:      "测试描述",
		Atk:       optional.Some(300),
		Def:       optional.Some(200),
		Level:     optional.Some(6),
		Race:      optional.Some("天使"),
		Attribute: optional.Some("光"),
		SetNames:  []string{"羽翼栗子球", "栗子球", "LV", "元素英雄", "英雄", "至爱"},
		Type:      []string{"怪兽卡", "效果", "特殊召唤"},
	}

	got := input.ToCardInfoForAI()

	atk := 300
	def := 200
	level := 6
	race := "天使"
	attribute := "光"

	want := &CardInfoForAI{
		Name:       "羽翼栗子球 LV6",
		Desc:       "测试描述",
		Atk:        &atk,
		Def:        &def,
		Level:      &level,
		Race:       &race,
		Attribute:  &attribute,
		Archetypes: "羽翼栗子球|栗子球|LV|元素英雄|英雄|至爱",
		Type:       "怪兽卡|效果|特殊召唤",
	}

	assert.Equal(t, want, got)
}

func (s *DBSuite) Test_toCardInfoForHuman() {
	db, err := New(git.NewRepo(localPath, remoteURL), s.basePath, "zh-CN", false)
	s.Require().NoError(err, "New should not return error")

	from := &CardInfoInDB{
		ID:        48486809,
		SetCode:   113768701513240740,
		Type:      33554465,
		Atk:       300,
		Def:       200,
		Level:     6,
		Race:      4,
		Attribute: 16,
		Name:      "羽翼栗子球 LV6",
		Desc: `这个卡名在规则上也当作「元素英雄」卡、「至爱」卡使用。这张卡不能通常召唤。「羽翼栗子球 LV6」1回合1次在把自己的手卡·场上（表侧表示）·墓地1只「元素英雄」融合怪兽或「羽翼栗子球」除外的场合才能从手卡·墓地特殊召唤。
①：对方怪兽的攻击宣言时或者对方把场上的怪兽的效果发动时，把这张卡解放才能发动。那1只怪兽破坏，给与对方那个原本攻击力数值的伤害。`,
	}

	got := from.toCardInfoForHuman(db)
	want := &CardInfoForHuman{
		ID:   48486809,
		Name: "羽翼栗子球 LV6",
		Desc: `这个卡名在规则上也当作「元素英雄」卡、「至爱」卡使用。这张卡不能通常召唤。「羽翼栗子球 LV6」1回合1次在把自己的手卡·场上（表侧表示）·墓地1只「元素英雄」融合怪兽或「羽翼栗子球」除外的场合才能从手卡·墓地特殊召唤。
①：对方怪兽的攻击宣言时或者对方把场上的怪兽的效果发动时，把这张卡解放才能发动。那1只怪兽破坏，给与对方那个原本攻击力数值的伤害。`,
		Atk:       optional.Some(300),
		Def:       optional.Some(200),
		Level:     optional.Some(6),
		Race:      optional.Some("天使"),
		Attribute: optional.Some("光"),
		SetNames:  []string{"羽翼栗子球", "栗子球", "LV", "元素英雄", "英雄", "至爱"},
		Type:      []string{"怪兽卡", "效果", "特殊召唤"},
	}
	s.Assert().Equal(want, got)
}

func TestCardInfoForHumanToCSV_Empty(t *testing.T) {
	headers, rows := CardInfoForHumanToCSV(nil)
	assert.Nil(t, headers)
	assert.Nil(t, rows)

	headers, rows = CardInfoForHumanToCSV([]*CardInfoForHuman{})
	assert.Nil(t, headers)
	assert.Nil(t, rows)
}

func TestCardInfoForHumanToCSV_AllFieldsPopulated(t *testing.T) {
	cards := []*CardInfoForHuman{
		{
			ID:        1,
			Name:      "Test Card 1",
			Desc:      "Description 1",
			Atk:       optional.Some(100),
			Def:       optional.Some(200),
			Level:     optional.Some(4),
			Race:      optional.Some("Warrior"),
			Attribute: optional.Some("Earth"),
			Type:      []string{"Monster", "Effect"},
			SetNames:  []string{"Test Set"},
		},
		{
			ID:        2,
			Name:      "Test Card 2",
			Desc:      "Description 2",
			Atk:       optional.Some(300),
			Def:       optional.Some(400),
			Level:     optional.Some(5),
			Race:      optional.Some("Spellcaster"),
			Attribute: optional.Some("Fire"),
			Type:      []string{"Monster", "Union"},
			SetNames:  []string{"Another Set"},
		},
	}

	headers, rows := CardInfoForHumanToCSV(cards)

	// All fields should be present since they have values
	expectedHeaders := []string{"name", "desc", "atk", "def", "level", "race", "attribute", "type", "setNames"}
	assert.Equal(t, expectedHeaders, headers)
	assert.Equal(t, 2, len(rows))

	// First row
	assert.Equal(t, []string{"Test Card 1", "Description 1", "100", "200", "4", "Warrior", "Earth", "Monster|Effect", "Test Set"}, rows[0])
	// Second row
	assert.Equal(t, []string{"Test Card 2", "Description 2", "300", "400", "5", "Spellcaster", "Fire", "Monster|Union", "Another Set"}, rows[1])
}

func TestCardInfoForHumanToCSV_SomeOptionalFieldsMissing(t *testing.T) {
	cards := []*CardInfoForHuman{
		{
			ID:   1,
			Name: "Monster Card",
			Desc: "A monster",
			Atk:  optional.Some(100),
			Def:  optional.Some(200),
			Type: []string{"Monster"},
		},
		{
			ID:        2,
			Name:      "Spell Card",
			Desc:      "A spell",
			Type:      []string{"Spell", "Normal"},
			SetNames:  []string{"Spell Set"},
			Level:     optional.Some(1),
			Race:      optional.Some("Spellcaster"),
			Attribute: optional.Some("Dark"),
		},
	}

	headers, rows := CardInfoForHumanToCSV(cards)

	// level, race, attribute are present because second card has them
	expectedHeaders := []string{"name", "desc", "atk", "def", "level", "race", "attribute", "type", "setNames"}
	assert.Equal(t, expectedHeaders, headers)

	// First row - monster has atk/def but no setnames
	assert.Equal(t, []string{"Monster Card", "A monster", "100", "200", "", "", "", "Monster", ""}, rows[0])
	// Second row - spell has setnames, level, race, attribute but no atk/def
	assert.Equal(t, []string{"Spell Card", "A spell", "", "", "1", "Spellcaster", "Dark", "Spell|Normal", "Spell Set"}, rows[1])
}

func TestCardInfoForHumanToCSV_SkipEmptyColumns(t *testing.T) {
	// Cards where some fields are empty on all cards
	cards := []*CardInfoForHuman{
		{
			ID:   1,
			Name: "Spell Card",
			Desc: "A spell",
			Type: []string{"Spell", "Normal"},
		},
		{
			ID:   2,
			Name: "Trap Card",
			Desc: "A trap",
			Type: []string{"Trap", "Counter"},
		},
	}

	headers, rows := CardInfoForHumanToCSV(cards)

	// Should skip atk, def, level, race, attribute, setNames (all empty on all cards)
	expectedHeaders := []string{"name", "desc", "type"}
	assert.Equal(t, expectedHeaders, headers)

	assert.Equal(t, []string{"Spell Card", "A spell", "Spell|Normal"}, rows[0])
	assert.Equal(t, []string{"Trap Card", "A trap", "Trap|Counter"}, rows[1])
}
