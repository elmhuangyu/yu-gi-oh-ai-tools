package cdb

import (
	"fmt"
	"strings"

	"github.com/moznion/go-optional"
)

type CardInfoForAI struct {
	Name       string  `json:"name" yaml:"name"`
	Desc       string  `json:"desc" yaml:"desc"`
	Atk        *int    `json:"atk,omitempty" yaml:"atk,omitempty"`
	Def        *int    `json:"def,omitempty" yaml:"def,omitempty"`
	Level      *int    `json:"level,omitempty" yaml:"level,omitempty"`
	Race       *string `json:"race,omitempty" yaml:"race,omitempty"`
	Attribute  *string `json:"attribute,omitempty" yaml:"attribute,omitempty"`
	Type       string  `json:"type,omitempty" yaml:"type,omitempty"`
	Archetypes string  `json:"archetypes,omitempty" yaml:"archetypes,omitempty"`
	Count      int     `json:"count,omitempty" yaml:"count,omitempty"`
}

type CardInfoForHuman struct {
	ID        uint64                  `json:"id"`
	Name      string                  `json:"name"`
	Desc      string                  `json:"desc"`
	Atk       optional.Option[int]    `json:"atk"`
	Def       optional.Option[int]    `json:"def"`
	Level     optional.Option[int]    `json:"level"`
	Race      optional.Option[string] `json:"race"`
	Attribute optional.Option[string] `json:"attribute"`
	Type      []string                `json:"type"`
	SetNames  []string                `json:"setNames"`
	Count     optional.Option[int]    `json:"count"`
	Deck      optional.Option[string] `json:"deck"`
}

// CardInfoForHumanToCSV if a field is empty on all cards, skip the column.
// Returns headers & data.
func CardInfoForHumanToCSV(cards []*CardInfoForHuman) ([]string, [][]string) {
	if len(cards) == 0 {
		return nil, nil
	}

	// First pass: determine which fields have at least one non-empty value
	hasName := false
	hasDesc := false
	hasAtk := false
	hasDef := false
	hasLevel := false
	hasRace := false
	hasAttribute := false
	hasType := false
	hasSetNames := false
	hasCount := false
	hasDeck := false
	for _, card := range cards {
		if card.Name != "" {
			hasName = true
		}
		if card.Desc != "" {
			hasDesc = true
		}
		if card.Atk.IsSome() {
			hasAtk = true
		}
		if card.Def.IsSome() {
			hasDef = true
		}
		if card.Level.IsSome() {
			hasLevel = true
		}
		if card.Race.IsSome() {
			hasRace = true
		}
		if card.Attribute.IsSome() {
			hasAttribute = true
		}
		if len(card.Type) > 0 {
			hasType = true
		}
		if len(card.SetNames) > 0 {
			hasSetNames = true
		}
		if card.Count.IsSome() {
			hasCount = true
		}
		if card.Deck.IsSome() {
			hasDeck = true
		}
	}

	// Build headers
	var headers []string
	if hasName {
		headers = append(headers, "name")
	}
	if hasDesc {
		headers = append(headers, "desc")
	}
	if hasAtk {
		headers = append(headers, "atk")
	}
	if hasDef {
		headers = append(headers, "def")
	}
	if hasLevel {
		headers = append(headers, "level")
	}
	if hasRace {
		headers = append(headers, "race")
	}
	if hasAttribute {
		headers = append(headers, "attribute")
	}
	if hasType {
		headers = append(headers, "type")
	}
	if hasSetNames {
		headers = append(headers, "setNames")
	}
	if hasCount {
		headers = append(headers, "count")
	}
	if hasDeck {
		headers = append(headers, "deck")
	}

	// Second pass: build data rows
	var rows [][]string
	for _, card := range cards {
		var row []string
		if hasName {
			row = append(row, card.Name)
		}
		if hasDesc {
			row = append(row, card.Desc)
		}
		if hasAtk {
			if card.Atk.IsSome() {
				row = append(row, formatInt(card.Atk.Unwrap()))
			} else {
				row = append(row, "")
			}
		}
		if hasDef {
			if card.Def.IsSome() {
				row = append(row, formatInt(card.Def.Unwrap()))
			} else {
				row = append(row, "")
			}
		}
		if hasLevel {
			if card.Level.IsSome() {
				row = append(row, formatInt(card.Level.Unwrap()))
			} else {
				row = append(row, "")
			}
		}
		if hasRace {
			if card.Race.IsSome() {
				row = append(row, card.Race.Unwrap())
			} else {
				row = append(row, "")
			}
		}
		if hasAttribute {
			if card.Attribute.IsSome() {
				row = append(row, card.Attribute.Unwrap())
			} else {
				row = append(row, "")
			}
		}
		if hasType {
			row = append(row, strings.Join(card.Type, "|"))
		}
		if hasSetNames {
			row = append(row, strings.Join(card.SetNames, "|"))
		}
		if hasCount {
			if card.Count.IsSome() {
				row = append(row, formatInt(card.Count.Unwrap()))
			} else {
				row = append(row, "")
			}
		}
		if hasDeck {
			if card.Deck.IsSome() {
				row = append(row, card.Deck.Unwrap())
			} else {
				row = append(row, "")
			}
		}
		rows = append(rows, row)
	}

	return headers, rows
}

func formatInt(v int) string {
	return fmt.Sprintf("%d", v)
}

func (s *CardInfoForHuman) ToCardInfoForAI() *CardInfoForAI {
	res := &CardInfoForAI{
		Name:       s.Name,
		Desc:       s.Desc,
		Type:       strings.Join(s.Type, "|"),
		Archetypes: strings.Join(s.SetNames, "|"),
	}

	if s.Atk.IsSome() {
		v := s.Atk.Unwrap()
		res.Atk = &v
	}
	if s.Def.IsSome() {
		v := s.Def.Unwrap()
		res.Def = &v
	}
	if s.Level.IsSome() {
		v := s.Level.Unwrap()
		res.Level = &v
	}
	if s.Race.IsSome() {
		v := s.Race.Unwrap()
		res.Race = &v
	}
	if s.Attribute.IsSome() {
		v := s.Attribute.Unwrap()
		res.Attribute = &v
	}
	if s.Count.IsSome() {
		res.Count = s.Count.Unwrap()
	}

	return res
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
	isMonster := s.Type&TypeMonster != 0

	res := &CardInfoForHuman{
		ID:   s.ID,
		Name: s.Name,
		Desc: s.Desc,
	}

	res.Type = GetTypeNames(db.lang, s.Type)

	if isMonster {
		res.Atk = optional.Some(s.Atk)
		res.Def = optional.Some(s.Def)
		res.Level = optional.Some(s.Level)
		res.Race = optional.Some(GetRaceName(db.lang, s.Race))
		res.Attribute = optional.Some(GetAttributeName(db.lang, s.Attribute))
	}
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
