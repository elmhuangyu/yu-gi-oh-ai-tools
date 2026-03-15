package git

import (
	"os"
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewRepo(t *testing.T) {
	// Test NewRepo function
	repo := NewRepo("/test/path", "https://github.com/test/repo.git")
	assert.NotNil(t, repo)
	assert.Equal(t, "/test/path", repo.BasePath)
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
	repoPath := path.Join(tmpDir, "ygopro-database")
	_, err = os.Stat(repoPath)
	assert.NoError(t, err)

	// Verify .git directory exists
	gitDir := path.Join(repoPath, ".git")
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

func TestNeedsUpdate(t *testing.T) {
	tmpDir := t.TempDir()
	repo := NewRepo(tmpDir, "https://github.com/test/repo.git")
	lastUpdatePath := path.Join(tmpDir, "last_update")

	t.Run("returns true when file does not exist", func(t *testing.T) {
		result := repo.needsUpdate(lastUpdatePath)
		assert.True(t, result)
	})

	t.Run("returns true when file has invalid content", func(t *testing.T) {
		err := os.WriteFile(lastUpdatePath, []byte("invalid-timestamp"), 0644)
		assert.NoError(t, err)

		result := repo.needsUpdate(lastUpdatePath)
		assert.True(t, result)
	})

	t.Run("returns false when timestamp is less than 1 hour old", func(t *testing.T) {
		timestamp := time.Now().Add(-30 * time.Minute)
		err := os.WriteFile(lastUpdatePath, []byte(timestamp.Format(time.RFC3339)), 0644)
		assert.NoError(t, err)

		result := repo.needsUpdate(lastUpdatePath)
		assert.False(t, result)
	})

	t.Run("returns true when timestamp is more than 1 hour old", func(t *testing.T) {
		timestamp := time.Now().Add(-2 * time.Hour)
		err := os.WriteFile(lastUpdatePath, []byte(timestamp.Format(time.RFC3339)), 0644)
		assert.NoError(t, err)

		result := repo.needsUpdate(lastUpdatePath)
		assert.True(t, result)
	})

	t.Run("returns true when timestamp is exactly 1 hour old", func(t *testing.T) {
		timestamp := time.Now().Add(-1 * time.Hour)
		err := os.WriteFile(lastUpdatePath, []byte(timestamp.Format(time.RFC3339)), 0644)
		assert.NoError(t, err)

		result := repo.needsUpdate(lastUpdatePath)
		assert.True(t, result)
	})
}
