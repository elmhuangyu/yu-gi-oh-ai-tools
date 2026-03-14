package cdb

import "strings"

type CardInfoForAI struct {
	Name       string `json:"name"`
	Desc       string `json:"desc"`
	Atk        int    `json:"atk,omitempty"`
	Def        int    `json:"def,omitempty"`
	Level      int    `json:"level,omitempty"`
	Race       string `json:"race,omitempty"`
	Attribute  string `json:"attribute,omitempty"`
	Type       string `json:"type,omitempty"`
	Archetypes string `json:"archetypes,omitempty"`
	Count      int    `json:"count,omitempty"`
}

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

func (s *CardInfoForHuman) ToCardInfoForAI() *CardInfoForAI {
	return &CardInfoForAI{
		Name:       s.Name,
		Desc:       s.Desc,
		Atk:        s.Atk,
		Def:        s.Def,
		Level:      s.Level,
		Race:       s.Race,
		Attribute:  s.Attribute,
		Type:       strings.Join(s.Type, "|"),
		Archetypes: strings.Join(s.SetNames, "|"),
	}
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
	// Also append the root set code -> name for each set code
	var setNames []string
	for i := 0; i < 4; i++ {
		// Extract 16 bits at position i*16
		code := (s.SetCode >> (i * 16)) & 0xFFFF
		if code > 0 {
			// Look up the full set code
			if name, ok := db.setName.GetByCode(code); ok {
				setNames = append(setNames, name)
			}
			// Also look up the root set code (lower 12 bits)
			rootCode := getRootSetCode(uint64(code))
			if rootCode > 0 {
				if name, ok := db.setName.GetByCode(rootCode); ok {
					// Avoid duplicates: only add if different from full set code name
					if rootCode != code {
						setNames = append(setNames, name)
					}
				}
			}
		}
	}
	res.SetNames = setNames

	return res
}
