package cdb

import (
	"database/sql"
	"errors"
	"path"
	"path/filepath"
	"time"

	"github.com/elmhuangyu/yu-gi-oh-ai-tools/ygo-db/lib/git"
	"github.com/gofrs/flock"
	_ "modernc.org/sqlite"
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
	gitRepo          *git.Repo
	basePath         string
	lang             string
	setName          *SetCodeAndName
	sqlite           *sql.DB
	lock             *flock.Flock
	enableAutoUpdate bool
}

func New(gitRepo *git.Repo, basePath, lang string, enableAutoUpdate bool) (*DB, error) {
	lockPath := path.Join(basePath, ".lock")
	fl := flock.New(lockPath)

	db := &DB{
		gitRepo:          gitRepo,
		basePath:         basePath,
		lang:             lang,
		setName:          NewSetCodeAndName(),
		lock:             fl,
		enableAutoUpdate: enableAutoUpdate,
	}

	err := fl.RLock()
	if err != nil {
		return nil, err
	}
	defer fl.Unlock()

	err = db.readSetName()
	if err != nil {
		return nil, err
	}
	err = db.connectSQLite()
	if err != nil {
		return nil, err
	}

	if db.enableAutoUpdate {
		go db.startUpdateLoop()
	}

	return db, nil
}

func (db *DB) startUpdateLoop() {
	for {
		time.Sleep(updateInterval)
		db.updateRepo()
	}
}

func (db *DB) updateRepo() error {
	if db.sqlite != nil {
		db.sqlite.Close()
		db.sqlite = nil
	}

	err := db.gitRepo.EnsureRepoUpToDate()
	if err != nil {
		return err
	}

	err = db.readSetName()
	if err != nil {
		return err
	}

	return db.connectSQLite()
}

func (db *DB) connectSQLite() error {
	repoPath := path.Join(db.basePath, "ygopro-database")
	dbPath := filepath.Join(repoPath, "locales", db.lang, "cards.cdb")
	sqlite, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return err
	}
	db.sqlite = sqlite
	return nil
}

func (db *DB) GetCardByID(id uint64) (*CardInfoForHuman, error) {
	query := `
		SELECT d.id, t.name, t.desc, d.atk, d.def, d.level, d.type, d.race, d.attribute, d.setcode
		FROM datas d
		JOIN texts t ON d.id = t.id
		WHERE d.id = ?
	`

	row := db.sqlite.QueryRow(query, id)

	var card CardInfoInDB
	err := row.Scan(
		&card.ID,
		&card.Name,
		&card.Desc,
		&card.Atk,
		&card.Def,
		&card.Level,
		&card.Type,
		&card.Race,
		&card.Attribute,
		&card.SetCode,
	)
	if err != nil {
		return nil, err
	}

	return card.toCardInfoForHuman(db), nil
}

func (db *DB) GetCardsByIDs(ids []uint64) (map[uint64]*CardInfoForHuman, error) {
	if len(ids) == 0 {
		return make(map[uint64]*CardInfoForHuman), nil
	}

	query := `
		SELECT d.id, t.name, t.desc, d.atk, d.def, d.level, d.type, d.race, d.attribute, d.setcode
		FROM datas d
		JOIN texts t ON d.id = t.id
		WHERE d.id IN (` + placeholders(len(ids)) + `)
	`

	args := make([]interface{}, len(ids))
	for i, id := range ids {
		args[i] = id
	}

	rows, err := db.sqlite.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[uint64]*CardInfoForHuman)
	for rows.Next() {
		var card CardInfoInDB
		err := rows.Scan(
			&card.ID,
			&card.Name,
			&card.Desc,
			&card.Atk,
			&card.Def,
			&card.Level,
			&card.Type,
			&card.Race,
			&card.Attribute,
			&card.SetCode,
		)
		if err != nil {
			return nil, err
		}
		result[card.ID] = card.toCardInfoForHuman(db)
	}

	return result, rows.Err()
}

func (db *DB) FindCardByName(name string, page int) (*CardInfoForHuman, []*CardInfoForHuman, int, error) {
	const limitSize = 30
	offset := page * limitSize

	db.lock.RLock()
	defer db.lock.Unlock()

	likePattern := "%" + name + "%"

	var total int
	err := db.sqlite.QueryRow(`
		SELECT COUNT(*)
		FROM texts t
		JOIN datas d ON d.id = t.id
		WHERE t.name LIKE ? AND d.alias = 0
	`, likePattern).Scan(&total)
	if err != nil {
		return nil, nil, 0, err
	}

	var exact *CardInfoForHuman
	if offset == 0 {
		var card CardInfoInDB
		err := db.sqlite.QueryRow(`
			SELECT d.id, t.name, t.desc, d.atk, d.def, d.level, d.type, d.race, d.attribute, d.setcode
			FROM datas d
			JOIN texts t ON d.id = t.id
			WHERE t.name = ? AND d.alias = 0
		`, name).Scan(
			&card.ID,
			&card.Name,
			&card.Desc,
			&card.Atk,
			&card.Def,
			&card.Level,
			&card.Type,
			&card.Race,
			&card.Attribute,
			&card.SetCode,
		)
		if err == nil {
			exact = card.toCardInfoForHuman(db)
		}
	}

	query := `
		SELECT d.id, t.name, t.desc, d.atk, d.def, d.level, d.type, d.race, d.attribute, d.setcode
		FROM datas d
		JOIN texts t ON d.id = t.id
		WHERE t.name LIKE ? AND t.name != ? AND d.alias = 0
		ORDER BY t.name
		LIMIT ? OFFSET ?
	`

	rows, err := db.sqlite.Query(query, likePattern, name, limitSize, offset)
	if err != nil {
		return nil, nil, 0, err
	}
	defer rows.Close()

	var maybe []*CardInfoForHuman

	for rows.Next() {
		var card CardInfoInDB
		err := rows.Scan(
			&card.ID,
			&card.Name,
			&card.Desc,
			&card.Atk,
			&card.Def,
			&card.Level,
			&card.Type,
			&card.Race,
			&card.Attribute,
			&card.SetCode,
		)
		if err != nil {
			return nil, nil, 0, err
		}

		maybe = append(maybe, card.toCardInfoForHuman(db))
	}

	return exact, maybe, total, rows.Err()
}

func (db *DB) FindCardsBySetName(setNames []string, page int) ([]*CardInfoForHuman, int, error) {
	const limitSize = 30
	offset := page * limitSize

	if len(setNames) == 0 || len(setNames) > 4 {
		return nil, 0, nil
	}

	// Collect all set codes from all set names
	var allCodes []uint64
	for _, setName := range setNames {
		codes, ok := db.setName.GetByName(setName)
		if !ok || len(codes) == 0 {
			return nil, 0, nil
		}
		allCodes = append(allCodes, codes...)
	}

	if len(allCodes) == 0 {
		return nil, 0, nil
	}

	// Build the SQL query using AND logic
	// For each set code, we check if it exists in any of the 4 slots
	// Then we AND all conditions together
	// Each slot is 16 bits: slot0 (bits 0-15), slot1 (bits 16-31), slot2 (bits 32-47), slot3 (bits 48-63)

	// Build the WHERE clause for each set code
	whereClause := ""
	var args []interface{}
	for i, code := range allCodes {
		var mask uint64 = 0x0FFF
		if !isRootSetCode(code) {
			mask = 0xFFFF
		}
		targetCode := code & mask

		if i > 0 {
			whereClause += " AND "
		}
		whereClause += "("
		whereClause += "((d.setcode >> 0)  & 0xFFFF & ?) = ? OR "
		whereClause += "((d.setcode >> 16) & 0xFFFF & ?) = ? OR "
		whereClause += "((d.setcode >> 32) & 0xFFFF & ?) = ? OR "
		whereClause += "((d.setcode >> 48) & 0xFFFF & ?) = ?"
		whereClause += ")"

		// Add mask and target for each of the 4 slot checks
		args = append(args, mask, targetCode, mask, targetCode, mask, targetCode, mask, targetCode)
	}

	query := `
		SELECT d.id, t.name, t.desc, d.atk, d.def, d.level, d.type, d.race, d.attribute, d.setcode
		FROM datas d
		JOIN texts t ON d.id = t.id
		WHERE ` + whereClause + ` AND d.alias = 0
	`

	// Get total count
	countQuery := "SELECT COUNT(*) FROM datas d JOIN texts t ON d.id = t.id WHERE " + whereClause + " AND d.alias = 0"
	var total int
	err := db.sqlite.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Add ordering and pagination
	query += " ORDER BY t.name LIMIT ? OFFSET ?"
	args = append(args, limitSize, offset)

	rows, err := db.sqlite.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var results []*CardInfoForHuman
	for rows.Next() {
		var card CardInfoInDB
		err := rows.Scan(
			&card.ID,
			&card.Name,
			&card.Desc,
			&card.Atk,
			&card.Def,
			&card.Level,
			&card.Type,
			&card.Race,
			&card.Attribute,
			&card.SetCode,
		)
		if err != nil {
			return nil, 0, err
		}

		results = append(results, card.toCardInfoForHuman(db))
	}

	return results, total, rows.Err()
}

func placeholders(n int) string {
	if n <= 0 {
		return ""
	}
	result := make([]byte, 0, n*2)
	for i := 0; i < n; i++ {
		if i > 0 {
			result = append(result, ',')
		}
		result = append(result, '?')
	}
	return string(result)
}
