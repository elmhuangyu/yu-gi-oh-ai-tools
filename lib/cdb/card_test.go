package cdb

import "github.com/elmhuangyu/yu-gi-oh-mcp/lib/git"

func (s *DBSuite) Test_toCardInfoForHuman() {
	db, err := New(git.NewRepo(localPath, remoteURL), s.repoPath, "zh-CN")
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
		Atk:       300,
		Def:       200,
		Level:     6,
		Race:      "天使",
		Attribute: "光",
		SetNames:  []string{"羽翼栗子球", "LV", "元素英雄", "至爱"},
		Type:      []string{"怪兽卡", "效果", "特殊召唤"},
	}
	s.Assert().Equal(want, got)
}
