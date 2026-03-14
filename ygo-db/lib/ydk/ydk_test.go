package ydk

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseYDKFile(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantMain  map[int]int
		wantExtra map[int]int
		wantSide  map[int]int
	}{
		{
			name:      "empty input",
			input:     "",
			wantMain:  map[int]int{},
			wantExtra: map[int]int{},
			wantSide:  map[int]int{},
		},
		{
			name: "main deck only",
			input: `#main
16306932
16306932
14105621`,
			wantMain:  map[int]int{16306932: 2, 14105621: 1},
			wantExtra: map[int]int{},
			wantSide:  map[int]int{},
		},
		{
			name: "extra deck only",
			input: `#extra
6218704
13331639`,
			wantMain:  map[int]int{},
			wantExtra: map[int]int{6218704: 1, 13331639: 1},
			wantSide:  map[int]int{},
		},
		{
			name: "side deck only",
			input: `!side
12345678
98765432`,
			wantMain:  map[int]int{},
			wantExtra: map[int]int{},
			wantSide:  map[int]int{12345678: 1, 98765432: 1},
		},
		{
			name: "full deck",
			input: `#main
16306932
16306932
14105621

#extra
6218704

!side
12345678`,
			wantMain:  map[int]int{16306932: 2, 14105621: 1},
			wantExtra: map[int]int{6218704: 1},
			wantSide:  map[int]int{12345678: 1},
		},
		{
			name: "comments are ignored",
			input: `#created by YGO Omega
#main
16306932
#comment here
14105621
#extra
6218704
!side
12345678`,
			wantMain:  map[int]int{16306932: 1, 14105621: 1},
			wantExtra: map[int]int{6218704: 1},
			wantSide:  map[int]int{12345678: 1},
		},
		{
			name: "#main and #extra with same prefix",
			input: `#main
100
#extra
200`,
			wantMain:  map[int]int{100: 1},
			wantExtra: map[int]int{200: 1},
			wantSide:  map[int]int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			main, extra, side := ParseYDKFile(tt.input)

			// Compare main deck
			assert.Equal(t, len(tt.wantMain), len(main))
			for code, count := range tt.wantMain {
				assert.Equal(t, count, main[code])
			}

			// Compare extra deck
			assert.Equal(t, len(tt.wantExtra), len(extra))
			for code, count := range tt.wantExtra {
				assert.Equal(t, count, extra[code])
			}

			// Compare side deck
			assert.Equal(t, len(tt.wantSide), len(side))
			for code, count := range tt.wantSide {
				assert.Equal(t, count, side[code])
			}
		})
	}
}
