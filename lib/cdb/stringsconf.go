package cdb

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// parseSetNameLine parses a setname line from strings.conf
// Format: `!setname 0x3b 真红眼	レッドアイズ` or `!setname 0xa008 Masked HERO`
// Returns the code, local name, and dedup key (Japanese name if available)
func parseSetNameLine(line string) (code int, localName, dedupKey string, err error) {
	// split by space only on the first two parts: !setname and code
	parts := strings.SplitN(line, " ", 3)
	if len(parts) < 3 {
		return 0, "", "", ErrInvalidSetNameLine
	}

	// parse hex code (e.g., 0x3b)
	codeStr := parts[1]
	code64, err := strconv.ParseInt(codeStr, 0, 32)
	if err != nil {
		return 0, "", "", errors.Join(ErrParseCode, err)
	}
	code = int(code64)

	// get name (may contain spaces, split by tab for Japanese name)
	// format: localName\tJapaneseName
	namePart := parts[2]
	nameParts := strings.Split(namePart, "\t")
	localName = nameParts[0]

	// Check if there's a Japanese name (after tab)
	// For non-Japanese locales, the text after tab might be the Japanese name
	// We use Japanese name as dedup key when available to handle cases like:
	// !setname 0x2066 磁石战士	マグネット・ウォリアー
	// !setname 0xe9 磁石战士	磁石の戦士(じしゃくのせんし)
	// Both have same Chinese name but different Japanese names
	if len(nameParts) > 1 && nameParts[1] != "" {
		// Japanese name exists, use it as dedup key
		dedupKey = nameParts[1]
	} else {
		// No Japanese name, use local name as dedup key
		dedupKey = ""
	}

	return code, localName, dedupKey, nil
}

func (db *DB) readSetName() error {
	setName := NewDoubleMap()

	// read strings.conf file in `$repoPath/locales/$Lang`
	filePath := filepath.Join(db.repoPath, "locales", db.lang, "strings.conf")
	file, err := os.Open(filePath)
	if err != nil {
		return errors.Join(ErrOpenFile, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// only process line begin with `!setname`
		if !strings.HasPrefix(line, "!setname") {
			continue
		}

		// use parseSetNameLine to extract code, localName and dedupKey
		code, localName, dedupKey, err := parseSetNameLine(line)
		if err != nil {
			// skip invalid lines (except parse code errors)
			if errors.Is(err, ErrInvalidSetNameLine) {
				continue
			}
			return err
		}

		if dedupKey == "" {
			if !setName.Add(code, localName) {
				return errors.Join(ErrDuplicate, fmt.Errorf("code=%d, name=%s", code, localName))
			}
		} else {
			// add to double map using Japanese name as dedup key when available
			if !setName.AddWithDedup(code, localName, dedupKey) {
				return errors.Join(ErrDuplicate, fmt.Errorf("code=%d, name=%s, dedupKey=%s", code, localName, dedupKey))
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return errors.Join(ErrParseLine, err)
	}

	db.setName = setName
	return nil
}
