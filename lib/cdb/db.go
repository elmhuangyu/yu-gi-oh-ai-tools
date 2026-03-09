package cdb

import (
	"errors"
	"sync"
	"time"

	"github.com/elmhuangyu/yu-gi-oh-mcp/lib/git"
)

var (
	ErrOpenFile           = errors.New("failed to open strings.conf")
	ErrParseCode          = errors.New("failed to parse code")
	ErrParseLine          = errors.New("failed to parse line")
	ErrDuplicate          = errors.New("duplicate key or value")
	ErrInvalidSetNameLine = errors.New("invalid setname line format")
)

const (
	updateInterval = time.Hour * 24
)

type DB struct {
	gitRepo  *git.Repo
	repoPath string
	lang     string
	setName  *DoubleMap
	lock     sync.RWMutex
}

func New(gitRepo *git.Repo, repoPath, lang string) (*DB, error) {
	db := &DB{
		gitRepo:  gitRepo,
		repoPath: repoPath,
		lang:     lang,
		setName:  NewDoubleMap(),
		lock:     sync.RWMutex{},
	}

	err := db.readSetName()
	if err != nil {
		return nil, err
	}

	go db.startUpdateLoop()

	return db, nil
}

func (db *DB) startUpdateLoop() {
	for {
		time.Sleep(updateInterval)
		db.updateRepo()
	}
}

func (db *DB) updateRepo() error {
	db.lock.Lock()
	defer db.lock.Unlock()

	err := db.gitRepo.EnsureRepoUpToDate()
	if err != nil {
		return err
	}

	return db.readSetName()
}
