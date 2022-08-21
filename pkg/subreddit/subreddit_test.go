package subreddit

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadSubreddits(t *testing.T) {
	type args struct {
		file string
	}

	subRedditList = make(map[string]bool)
	subRedditList["/r/100yearsago"] = true
	subRedditList["/r/1200isplenty"] = true

	tests := []struct {
		name         string
		args         args
		wantErr      bool
		expectedErr  error
		expectedList map[string]bool
	}{
		{name: "success",
			args:         args{file: "./testdata/subreddits.txt"},
			wantErr:      false,
			expectedErr:  nil,
			expectedList: subRedditList},
		{name: "no-list",
			args:         args{file: "./testdata/subredditX.txt"},
			wantErr:      true,
			expectedErr:  fmt.Errorf("error opening subreddits file: open ./testdata/subredditX.txt: no such file or directory"),
			expectedList: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := LoadSubreddits(tt.args.file)
			if tt.wantErr {
				assert.Equal(t, tt.expectedErr, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedList, subRedditList)
			}
		})
	}
}
