package cdb

import (
	"errors"
	"os"
	"testing"

	"github.com/elmhuangyu/yu-gi-oh-mcp/lib/git"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type DBSuite struct {
	suite.Suite
	repoPath string
}

func (s *DBSuite) SetupSuite() {
	repoPath := "/tmp/yugioh-cdb/"
	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		// Clone the ygopro-database repo for testing
		repo := git.NewRepo(repoPath, "https://github.com/mycard/ygopro-database.git")
		err := repo.EnsureRepoUpToDate()
		require.NoError(s.T(), err, "failed to clone ygopro-database repo")
	}
	s.repoPath = repoPath
}

func (s *DBSuite) TestReadSetName() {
	db := New(s.repoPath, "zh-CN")
	err := db.readSetName()

	// readSetName should succeed
	require.NoError(s.T(), err, "readSetName should not return error")

	// SetName should be populated
	assert.NotNil(s.T(), db.SetName, "SetName should not be nil")
	assert.True(s.T(), db.SetName.Len() > 0, "SetName should have entries")

	// Test known setname entries from zh-CN strings.conf
	// !setname 0x1 正义盟军	A・O・J
	name, ok := db.SetName.GetByInt(0x1)
	assert.True(s.T(), ok, "should find setname for code 0x1")
	assert.Equal(s.T(), "正义盟军", name, "setname for 0x1 should be 正义盟军")

	// Test reverse lookup
	code, ok := db.SetName.GetByStringFirst("真红眼")
	assert.True(s.T(), ok, "should find code for setname 真红眼")
	assert.Equal(s.T(), int64(0x3b), int64(code), "code for 真红眼 should be 0x3b")
}

func (s *DBSuite) TestReadSetName_InvalidRepoPath() {
	db := New("/nonexistent/path", "zh-CN")
	err := db.readSetName()

	// Should return error for invalid path
	assert.Error(s.T(), err, "readSetName should return error for invalid path")
	assert.True(s.T(), errors.Is(err, ErrOpenFile), "error should be ErrOpenFile")
}

func (s *DBSuite) TestReadSetName_InvalidLang() {
	db := New(s.repoPath, "invalid-lang")
	err := db.readSetName()

	// Should return error for invalid lang
	assert.Error(s.T(), err, "readSetName should return error for invalid lang")
	assert.True(s.T(), errors.Is(err, ErrOpenFile), "error should be ErrOpenFile")
}

func TestDB(t *testing.T) {
	suite.Run(t, new(DBSuite))
}

func TestParseSetNameLine(t *testing.T) {
	tests := []struct {
		name              string
		line              string
		expectedCode      int
		expectedLocalName string
		expectedDedupKey  string
		expectedErr       error
	}{
		{
			name:              "Japanese name with tab separator",
			line:              "!setname 0x3b 真红眼\tレッドアイズ",
			expectedCode:      0x3b,
			expectedLocalName: "真红眼",
			expectedDedupKey:  "レッドアイズ",
			expectedErr:       nil,
		},
		{
			name:              "English name without Japanese",
			line:              "!setname 0xa008 Masked HERO",
			expectedCode:      0xa008,
			expectedLocalName: "Masked HERO",
			expectedDedupKey:  "",
			expectedErr:       nil,
		},
		{
			name:              "Chinese with Japanese name",
			line:              "!setname 0x2066 磁石战士\tマグネット・ウォリアー",
			expectedCode:      0x2066,
			expectedLocalName: "磁石战士",
			expectedDedupKey:  "マグネット・ウォリアー",
			expectedErr:       nil,
		},
		{
			name:        "invalid line - too few parts",
			line:        "!setname 0x3b",
			expectedErr: ErrInvalidSetNameLine,
		},
		{
			name:        "invalid hex code",
			line:        "!setname xyz test",
			expectedErr: ErrParseCode,
		},
		{
			name:              "empty Japanese name after tab",
			line:              "!setname 0x1 test\t",
			expectedCode:      0x1,
			expectedLocalName: "test",
			expectedDedupKey:  "",
			expectedErr:       nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, localName, dedupKey, err := parseSetNameLine(tt.line)

			if tt.expectedErr != nil {
				assert.Error(t, err, "should return error")
				assert.True(t, errors.Is(err, tt.expectedErr), "error should be %v", tt.expectedErr)
				return
			}

			require.NoError(t, err, "parseSetNameLine should not return error")
			assert.Equal(t, tt.expectedCode, code, "code should match")
			assert.Equal(t, tt.expectedLocalName, localName, "localName should match")
			assert.Equal(t, tt.expectedDedupKey, dedupKey, "dedupKey should match")
		})
	}
}
