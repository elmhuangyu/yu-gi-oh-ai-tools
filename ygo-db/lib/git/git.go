package git

import (
	"os"
	"os/exec"
	"path"
	"time"

	"github.com/gofrs/flock"
)

type Repo struct {
	BasePath  string
	RemoteURL string
}

// NewRepo creates a new Repo instance. This is not for general purpose - it only handles
// cloning and pulling the shadow repo.
func NewRepo(basePath, remoteURL string) *Repo {
	return &Repo{
		BasePath:  basePath,
		RemoteURL: remoteURL,
	}
}

func (r *Repo) EnsureRepoUpToDate() (err error) {
	if _, err := os.Stat(r.BasePath); os.IsNotExist(err) {
		if err := os.MkdirAll(r.BasePath, 0777); err != nil {
			return err
		}
	}

	lastUpdatePath := path.Join(r.BasePath, "last_update")

	// First check: read timestamp file without lock
	if !r.needsUpdate(lastUpdatePath) {
		return nil
	}

	lockPath := path.Join(r.BasePath, ".lock")
	fl := flock.New(lockPath)
	if err := fl.Lock(); err != nil {
		return err
	}
	defer func() { err = fl.Unlock() }()

	// Double-check: re-read timestamp after acquiring lock
	if !r.needsUpdate(lastUpdatePath) {
		return nil
	}

	repoPath := path.Join(r.BasePath, "ygopro-database")

	// Check if repository path exists
	if _, err := os.Stat(path.Join(repoPath, ".git")); os.IsNotExist(err) {
		// Repository doesn't exist, clone it
		cmd := exec.Command("git", "clone", "--depth", "1", r.RemoteURL, repoPath)
		if _, err := cmd.CombinedOutput(); err != nil {
			return err
		}
	} else if err != nil {
		// Error checking if path exists
		return err
	} else {
		// Repository exists, pull latest changes
		cmd := exec.Command("git", "pull", "--depth", "1")
		cmd.Dir = repoPath
		if _, err := cmd.CombinedOutput(); err != nil {
			return err
		}
	}

	// Update the timestamp file after successful update
	if err := os.WriteFile(lastUpdatePath, []byte(time.Now().Format(time.RFC3339)), 0644); err != nil {
		return err
	}

	return nil
}

// needsUpdate checks if the last update was more than 1 hour ago.
// Returns true if update is needed, false if within 1 hour.
func (r *Repo) needsUpdate(lastUpdatePath string) bool {
	data, err := os.ReadFile(lastUpdatePath)
	if err != nil {
		return true
	}

	timestamp, err := time.Parse(time.RFC3339, string(data))
	if err != nil {
		return true
	}

	return time.Since(timestamp) >= time.Hour
}
