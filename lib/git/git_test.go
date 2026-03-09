package git

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRepo(t *testing.T) {
	// Test NewRepo function
	repo := NewRepo("/test/path", "https://github.com/test/repo.git")
	assert.NotNil(t, repo)
	assert.Equal(t, "/test/path", repo.Path)
	assert.Equal(t, "https://github.com/test/repo.git", repo.RemoteURL)
}

func TestEnsureRepoUpToDate_WithNonExistentPath(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a test repository
	testRemoteURL := "https://github.com/stretchr/testify.git" // Public repo for testing
	repo := NewRepo(tmpDir, testRemoteURL)

	// Test EnsureRepoUpToDate - should clone the repo
	err := repo.EnsureRepoUpToDate()
	assert.NoError(t, err)

	// Verify the repository was created
	_, err = os.Stat(tmpDir)
	assert.NoError(t, err)

	// Verify .git directory exists
	gitDir := tmpDir + "/.git"
	_, err = os.Stat(gitDir)
	assert.NoError(t, err)
}

func TestEnsureRepoUpToDate_WithExistingPath(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a test repository
	testRemoteURL := "https://github.com/stretchr/testify.git" // Public repo for testing
	repo := NewRepo(tmpDir, testRemoteURL)

	// First call to ensure repo is cloned
	err := repo.EnsureRepoUpToDate()
	assert.NoError(t, err)

	// Second call to ensure repo is up to date
	err = repo.EnsureRepoUpToDate()
	assert.NoError(t, err)
}
