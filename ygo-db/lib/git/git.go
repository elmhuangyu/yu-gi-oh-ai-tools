package git

import (
	"os"
	"os/exec"
	"path"
)

type Repo struct {
	Path      string
	RemoteURL string
}

// NewRepo creates a new Repo instance. This is not for general purpose - it only handles
// cloning and pulling the shadow repo.
func NewRepo(path, remoteURL string) *Repo {
	return &Repo{
		Path:      path,
		RemoteURL: remoteURL,
	}
}

func (r *Repo) EnsureRepoUpToDate() error {
	// Check if repository path exists
	if _, err := os.Stat(path.Join(r.Path, ".git")); os.IsNotExist(err) {
		// Repository doesn't exist, clone it
		cmd := exec.Command("git", "clone", "--depth", "1", r.RemoteURL, r.Path)
		if _, err := cmd.CombinedOutput(); err != nil {
			return err
		}
	} else if err != nil {
		// Error checking if path exists
		return err
	} else {
		// Repository exists, pull latest changes
		cmd := exec.Command("git", "pull", "--depth", "1")
		cmd.Dir = r.Path
		if _, err := cmd.CombinedOutput(); err != nil {
			return err
		}
	}

	return nil
}
