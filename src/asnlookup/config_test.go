package asnlookup

import (
	"os"
	"testing"
)

func TestGetConfig(t *testing.T) {
	testCases := []struct {
		name      string
		setUpFunc func()
		err       error
	}{
		{
			name: "Read Configuration From Local File",
			setUpFunc: func() {
				os.Setenv("CONFIG_FILE_PATH", "./cfgfile.txt")
			},
			err: nil,
		},
		{
			name: "Read Configuration From URL",
			setUpFunc: func() {
			},
			err: nil,
		},
	}

	for _, testCase := range testCases {
		testCase.setUpFunc()
		_, err := GetConfig()
		if err != testCase.err {
			t.Fatalf("%s: received error does not match: got %v, want %v", testCase.name, err, testCase.err)
		}
	}
}
