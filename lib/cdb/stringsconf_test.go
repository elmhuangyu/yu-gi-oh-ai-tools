package cdb

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *DBSuite) TestReadSetName() {
	db := &DB{
		repoPath: s.repoPath, lang: "zh-CN", setName: NewSetCodeAndName(),
	}
	err := db.readSetName()

	// readSetName should succeed
	s.Require().NoError(err, "readSetName should not return error")

	// SetName should be populated
	s.Assert().NotNil(db.setName, "SetName should not be nil")
	s.Assert().True(db.setName.Len() > 0, "SetName should have entries")

	// Test known setname entries from zh-CN strings.conf
	// !setname 0x1 正义盟军	A・O・J
	name, ok := db.setName.GetByUint64(0x1)
	s.Assert().True(ok, "should find setname for code 0x1")
	s.Assert().Equal("正义盟军", name, "setname for 0x1 should be 正义盟军")

	// Test reverse lookup
	code, ok := db.setName.GetByStringFirst("真红眼")
	s.Assert().True(ok, "should find code for setname 真红眼")
	s.Assert().Equal(uint64(0x3b), code, "code for 真红眼 should be 0x3b")
}

func (s *DBSuite) TestReadSetName_InvalidRepoPath() {
	db := &DB{
		repoPath: "/invalid/path", lang: "zh-CN", setName: NewSetCodeAndName(),
	}
	err := db.readSetName()

	// Should return error for invalid path
	s.Assert().Error(err, "readSetName should return error for invalid path")
	s.Assert().True(errors.Is(err, ErrOpenFile), "error should be ErrOpenFile")
}

func (s *DBSuite) TestReadSetName_InvalidLang() {
	db := &DB{
		repoPath: s.repoPath, lang: "InvalidLang", setName: NewSetCodeAndName(),
	}
	err := db.readSetName()

	// Should return error for invalid lang
	s.Assert().Error(err, "readSetName should return error for invalid lang")
	s.Assert().True(errors.Is(err, ErrOpenFile), "error should be ErrOpenFile")
}

func TestParseSetNameLine(t *testing.T) {
	tests := []struct {
		name              string
		line              string
		expectedCode      uint64
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
