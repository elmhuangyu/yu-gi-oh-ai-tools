package cdb

type CardInfoForHuman struct {
	ID        uint64   `json:"id"`
	Name      string   `json:"name"`
	Desc      string   `json:"desc"`
	Atk       int      `json:"atk"`
	Def       int      `json:"def"`
	Type      []string `json:"type"`
	Level     int      `json:"level"`
	Race      string   `json:"race"`
	Attribute string   `json:"attribute"`
	SetNames  []string `json:"setNames"`
}

type CardInfoInDB struct {
	ID        uint64
	Name      string
	Desc      string
	Atk       int
	Def       int
	Level     int
	Type      uint64
	Race      int
	Attribute int
	SetCode   uint64
}

func (s *CardInfoInDB) toCardInfoForHuman(db *DB) *CardInfoForHuman {
	res := &CardInfoForHuman{
		ID:    s.ID,
		Name:  s.Name,
		Desc:  s.Desc,
		Atk:   s.Atk,
		Def:   s.Def,
		Level: s.Level,
	}

	res.Type = GetTypeNames(db.lang, s.Type)
	res.Race = GetRaceName(db.lang, s.Race)
	res.Attribute = GetAttributeName(db.lang, s.Attribute)

	// SetNames is a uint64 where each 16 bits represents a set name code
	// Extract 4 16-bit values and look up each in the setName map
	var setNames []string
	for i := 0; i < 4; i++ {
		// Extract 16 bits at position i*16
		code := int((s.SetCode >> (i * 16)) & 0xFFFF)
		if code > 0 {
			if name, ok := db.setName.GetByInt(code); ok {
				setNames = append(setNames, name)
			}
		}
	}
	res.SetNames = setNames

	return res
}
