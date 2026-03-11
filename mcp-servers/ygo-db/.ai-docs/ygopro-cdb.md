# YGOCore CDB (SQLite3) Database Schema

The `cards.cdb` file is a SQLite3 database used by YGOPRO/YGOCore to store card information. It typically contains two main tables: `datas` and `texts`.

## Tables

### 1. `datas` Table
Stores numerical card data and properties.

| Column | Type | Description |
| :--- | :--- | :--- |
| **id** | INTEGER | Card's unique passcode (8-digit number). Primary Key. |
| **ot** | INTEGER | Bitmask for Origin/Legality (1: OCG, 2: TCG, 3: OCG/TCG, 8: Anime, etc.). |
| **alias** | INTEGER | ID of the card this card is treated as. |
| **setcode** | INTEGER | Bitmask for Archetypes/Sets (e.g., "Elemental HERO"). |
| **type** | INTEGER | Bitmask for Card Type (Monster, Spell, Trap, etc.). |
| **atk** | INTEGER | Attack points. |
| **def** | INTEGER | Defense points. |
| **level** | INTEGER | Level/Rank (lower 12 bits). Also encodes Pendulum Scales in higher bits. |
| **race** | INTEGER | Bitmask for Monster Race (Warrior, Spellcaster, etc.). |
| **attribute** | INTEGER | Bitmask for Monster Attribute (Light, Dark, etc.). |
| **category** | INTEGER | Bitmask for card effect categories. |

### 2. `texts` Table
Stores localized strings for the card.

| Column | Type | Description |
| :--- | :--- | :--- |
| **id** | INTEGER | Card's unique passcode. Primary Key (links to `datas.id`). |
| **name** | TEXT | Name of the card. |
| **desc** | TEXT | Card's effect or flavor text. |
| **str1** to **str16** | TEXT | Strings for specific effect prompts or menu options. |

## SQL Schema

```sql
CREATE TABLE datas (
    id INTEGER PRIMARY KEY,
    ot INTEGER,
    alias INTEGER,
    setcode INTEGER,
    type INTEGER,
    atk INTEGER,
    def INTEGER,
    level INTEGER,
    race INTEGER,
    attribute INTEGER,
    category INTEGER
);

CREATE TABLE texts (
    id INTEGER PRIMARY KEY,
    name TEXT,
    desc TEXT,
    str1 TEXT, str2 TEXT, str3 TEXT, str4 TEXT,
    str5 TEXT, str6 TEXT, str7 TEXT, str8 TEXT,
    str9 TEXT, str10 TEXT, str11 TEXT, str12 TEXT,
    str13 TEXT, str14 TEXT, str15 TEXT, str16 TEXT
);
```

## Go Implementation Context

In this project, `CardInfoInDB` (defined in `lib/cdb/card.go`) represents a joined view of these tables:

```go
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
```

### Bitmask Constants

Refer to `lib/cdb/constants.go` for the full list of bitmask values for `Type`, `Race`, and `Attribute`.

- **Type Examples**: `0x1` (Monster), `0x2` (Spell), `0x4` (Trap), `0x10` (Normal), `0x20` (Effect).
- **Attribute Examples**: `0x1` (Earth), `0x2` (Water), `0x4` (Fire), `0x8` (Wind), `0x10` (Light), `0x20` (Dark).
- **Race Examples**: `0x1` (Warrior), `0x2` (Spellcaster), `0x4` (Fairy), `0x8` (Fiend).

### SetCode Parsing

The `SetCode` is a `uint64` where each 16 bits represents a set name code. Up to 4 set names can be associated with a single card.
Logic to extract set names (from `lib/cdb/card.go`):

```go
for i := 0; i < 4; i++ {
    code := int((s.SetCode >> (i * 16)) & 0xFFFF)
    // ... look up code in setname map ...
}
```
