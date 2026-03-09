package cdb

import (
	"os"
	"testing"
	"time"

	"github.com/elmhuangyu/yu-gi-oh-mcp/lib/git"
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
