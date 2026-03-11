package cdb

import (
	"os"
	"testing"
	"time"

	"github.com/elmhuangyu/yu-gi-oh-ai-tools/mcp-servers/ygo-db/lib/git"
	"github.com/stretchr/testify/suite"
)

type DBSuite struct {
	suite.Suite
	repoPath string
}

const (
	localPath = "/tmp/yugioh-cdb/"
	remoteURL = "https://github.com/mycard/ygopro-database.git"
)

func (s *DBSuite) SetupSuite() {
	if _, err := os.Stat(localPath); os.IsNotExist(err) {
		// Clone the ygopro-database repo for testing
		repo := git.NewRepo(localPath, remoteURL)
		err := repo.EnsureRepoUpToDate()
		s.Require().NoError(err, "failed to clone ygopro-database repo")
	}
	s.repoPath = localPath
}

func TestDB(t *testing.T) {
	suite.Run(t, new(DBSuite))
}

func (s *DBSuite) Test_New() {
	db, err := New(git.NewRepo(localPath, remoteURL), s.repoPath, "zh-CN")
	s.Require().NoError(err, "New should not return error")
	s.Require().NotNil(db, "db should not be nil")
	s.Require().NotNil(db.setName, "SetName should not be nil")
	s.Assert().Less(0, db.setName.Len(), "SetName should have entries")
}

func (s *DBSuite) Test_updateRepo() {
	db, err := New(git.NewRepo(localPath, remoteURL), s.repoPath, "zh-CN")
	s.Require().NoError(err, "New should not return error")

	err = db.updateRepo()
	s.Require().NoError(err, "updateRepo should not return error")
	s.Assert().Less(0, db.setName.Len(), "SetName should have entries")
}

func (s *DBSuite) Test_updateRepo_BlockedByReadLock() {
	db, err := New(git.NewRepo(localPath, remoteURL), s.repoPath, "zh-CN")
	s.Require().NoError(err, "New should not return error")

	// Acquire a read lock
	db.lock.RLock()
	// Start a goroutine that holds a read lock for 1 second
	startTime := time.Now()
	go func() {
		// Hold the read lock for 1 second
		time.Sleep(1 * time.Second)
		db.lock.RUnlock()
	}()

	// Call updateRepo which requires a write lock
	// It should be blocked by the read lock and take at least 1 second
	err = db.updateRepo()
	elapsed := time.Since(startTime)

	s.Require().NoError(err, "updateRepo should not return error")
	s.Assert().GreaterOrEqual(elapsed, 900*time.Millisecond,
		"updateRepo should have been blocked by the read lock for at least ~1 second")
}

func (s *DBSuite) Test_GetCardByID() {
	db, err := New(git.NewRepo(localPath, remoteURL), s.repoPath, "zh-CN")
	s.Require().NoError(err, "New should not return error")

	card, err := db.GetCardByID(48486809)
	s.Require().NoError(err, "GetCardByID should not return error")
	s.Require().NotNil(card, "card should not be nil")

	s.Assert().Equal(uint64(48486809), card.ID)
	s.Assert().Equal("羽翼栗子球 LV6", card.Name)
	s.Assert().NotEmpty(card.Desc, "card desc should not be empty")
	s.Assert().Equal(300, card.Atk)
	s.Assert().Equal(200, card.Def)
	s.Assert().Equal(6, card.Level)
	s.Assert().Equal("天使", card.Race)
	s.Assert().Equal("光", card.Attribute)
	s.Assert().Equal([]string{"羽翼栗子球", "栗子球", "LV", "元素英雄", "英雄", "至爱"}, card.SetNames)
	s.Assert().Equal([]string{"怪兽卡", "效果", "特殊召唤"}, card.Type)
}

func (s *DBSuite) Test_GetCardByID_NotFound() {
	db, err := New(git.NewRepo(localPath, remoteURL), s.repoPath, "zh-CN")
	s.Require().NoError(err, "New should not return error")

	_, err = db.GetCardByID(99999999)
	s.Assert().Error(err, "GetCardByID should return error for non-existent card")
}

func (s *DBSuite) Test_FindCardByName() {
	db, err := New(git.NewRepo(localPath, remoteURL), s.repoPath, "zh-CN")
	s.Require().NoError(err, "New should not return error")

	exact, maybe, total, err := db.FindCardByName("青眼白龙", 0)
	s.Require().NoError(err, "FindCardByName should not return error")

	s.Require().NotNil(exact, "exact match should not be nil")
	s.Assert().Equal("青眼白龙", exact.Name)

	s.Assert().Equal(2, total)
	s.Assert().Equal(1, len(maybe), "should return 1 partial match (罪 青眼白龙)")
	s.Assert().Equal("罪 青眼白龙", maybe[0].Name)
}

func (s *DBSuite) Test_FindCardByName_Pagination() {
	db, err := New(git.NewRepo(localPath, remoteURL), s.repoPath, "zh-CN")
	s.Require().NoError(err, "New should not return error")

	exact, maybe1, total, err := db.FindCardByName("龙", 0)
	s.Require().NoError(err, "FindCardByName should not return error")
	s.Require().Nil(exact, "exact match should be nil")
	s.Assert().Equal(30, len(maybe1), "first page should have 30 results")

	_, maybe2, total2, err := db.FindCardByName("龙", 30)
	s.Require().NoError(err, "FindCardByName with offset should not return error")
	s.Assert().Equal(total, total2)
	s.Assert().Equal(30, len(maybe2), "second page should have 30 results")

	s.Assert().NotEqual(maybe1[0].Name, maybe2[0].Name, "pages should have different results")
}

func (s *DBSuite) Test_FindCardsBySetName() {
	db, err := New(git.NewRepo(localPath, remoteURL), s.repoPath, "zh-CN")
	s.Require().NoError(err, "New should not return error")

	results, total, err := db.FindCardsBySetName([]string{"英雄"}, 0)
	s.Require().NoError(err, "FindCardsBySetName should not return error")
	s.Assert().Greater(total, 0, "should have results for 英雄 set")
	s.Assert().LessOrEqual(len(results), 30, "results should be limited to 30 for first page")

	// Verify that each result contains the 英雄 set name
	for _, card := range results {
		hasHeroSet := false
		for _, setName := range card.SetNames {
			if setName == "英雄" {
				hasHeroSet = true
				break
			}
		}
		s.Assert().True(hasHeroSet, "card %s should have 英雄 set", card.Name)
	}
}

func (s *DBSuite) Test_FindCardsBySetName_Pagination() {
	db, err := New(git.NewRepo(localPath, remoteURL), s.repoPath, "zh-CN")
	s.Require().NoError(err, "New should not return error")

	maybe1, total, err := db.FindCardsBySetName([]string{"英雄"}, 0)
	s.Require().NoError(err, "FindCardsBySetName should not return error")
	s.Assert().LessOrEqual(len(maybe1), 30, "first page may have less than 30 results")

	if total > 30 {
		maybe2, total2, err := db.FindCardsBySetName([]string{"英雄"}, 30)
		s.Require().NoError(err, "FindCardsBySetName with offset should not return error")
		s.Assert().Equal(total, total2)
		s.Assert().Greater(len(maybe2), 0, "second page should have results")
	}
}

func (s *DBSuite) Test_FindCardsBySetName_MultipleSetNames() {
	db, err := New(git.NewRepo(localPath, remoteURL), s.repoPath, "zh-CN")
	s.Require().NoError(err, "New should not return error")

	// Test searching with 2 set names: "栗子球" and "英雄"
	// Should return cards that have BOTH set names (AND logic)
	results, total, err := db.FindCardsBySetName([]string{"栗子球", "英雄"}, 0)
	s.Require().NoError(err, "FindCardsBySetName with multiple set names should not return error")
	s.Assert().Greater(total, 0, "should have results for 栗子球+英雄 set combination")

	// Verify that each result contains both set names (栗子球 AND 英雄)
	foundWingedMagnet := false
	for _, card := range results {
		hasChestnut := false
		hasHero := false
		for _, setName := range card.SetNames {
			if setName == "栗子球" {
				hasChestnut = true
			}
			if setName == "英雄" {
				hasHero = true
			}
		}
		s.Assert().True(hasChestnut && hasHero, "card %s should have both 栗子球 and 英雄 sets", card.Name)

		if card.Name == "羽翼栗子球 LV6" {
			foundWingedMagnet = true
		}
	}
	s.Assert().True(foundWingedMagnet, "should find 羽翼栗子球 LV6 in results")
}
