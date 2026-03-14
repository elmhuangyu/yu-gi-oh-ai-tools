# Testing Code Style

All tests should follow these conventions:

1. **Use testify library**: All test files must use [`testify`](https://github.com/stretchr/testify) for assertions and mocking. Import as `github.com/stretchr/testify/assert` and `github.com/stretchr/testify/require`.
2. **Use table-driven tests**: When testing a function with multiple test cases, use table-driven tests (also known as subtests) for better organization and maintainability.

   ```go
   func TestGetTypeByName(t *testing.T) {
       tests := []struct {
           name     string
           lang     Language
           input    string
           expected uint32
       }{
           {"Monster in English", LangEnUS, "Monster", 0x1},
           {"Monster in Chinese", LangZhCN, "怪兽卡", 0x1},
           {"Spell in English", LangEnUS, "Spell", 0x2},
           {"Not found", LangEnUS, "Unknown", 0},
       }
       
       for _, tt := range tests {
           t.Run(tt.name, func(t *testing.T) {
               got := GetTypeByName(tt.lang, tt.input)
               assert.Equal(t, tt.expected, got)
           })
       }
   }
   ```
3. **Use `require` for fatal assertions**: Use `require` instead of `assert` when the test cannot continue if the assertion fails (e.g., setup errors).
