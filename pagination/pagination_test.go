package pagination

import (
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name       string
		pageStr    string
		sizeStr    string
		sortStr    string
		wantedPage int
		wantedSize int
		wantedSort string
		wantError  bool
	}{
		{
			name:       "Default values",
			pageStr:    "",
			sizeStr:    "",
			sortStr:    "",
			wantedPage: 0,
			wantedSize: 25,
			wantedSort: "",
			wantError:  false,
		},
		{
			name:       "Valid custom values",
			pageStr:    "2",
			sizeStr:    "100",
			sortStr:    "created_at:desc",
			wantedPage: 2,
			wantedSize: 100,
			wantedSort: "created_at:desc",
			wantError:  false,
		},
		{
			name:      "Negative page returns error",
			pageStr:   "-1",
			wantError: true,
		},
		{
			name:      "Invalid page format returns error",
			pageStr:   "abc",
			wantError: true,
		},
		{
			name:      "Invalid size string returns error",
			sizeStr:   "abc",
			wantError: true,
		},
		{
			name:      "Unsupported size returns error",
			sizeStr:   "30",
			wantError: true,
		},
		{
			name:       "Max size is accepted",
			pageStr:    "",
			sizeStr:    "1000",
			sortStr:    "",
			wantedPage: 0,
			wantedSize: 1000,
			wantedSort: "",
			wantError:  false,
		},
		{
			name:      "Invalid sort direction returns error",
			sortStr:   "name:up",
			wantError: true,
		},
		{
			name:      "Invalid sort format returns error",
			sortStr:   "name",
			wantError: true,
		},
		{
			name:       "Sort direction ignores case but preserves original string",
			pageStr:    "",
			sizeStr:    "",
			sortStr:    "name:ASC",
			wantedPage: 0,
			wantedSize: 25,
			wantedSort: "name:ASC",
			wantError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Parse(tt.pageStr, tt.sizeStr, tt.sortStr)
			if (err != nil) != tt.wantError {
				t.Errorf("Parse() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if !tt.wantError {
				if result.Page != tt.wantedPage {
					t.Errorf("Page = %d, want %d", result.Page, tt.wantedPage)
				}
				if result.Size != tt.wantedSize {
					t.Errorf("Size = %d, want %d", result.Size, tt.wantedSize)
				}
				if result.Sort != tt.wantedSort {
					t.Errorf("Sort = '%s', want '%s'", result.Sort, tt.wantedSort)
				}
			}
		})
	}
}
