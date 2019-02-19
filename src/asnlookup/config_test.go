package asnlookup

import (
	"os"
	"testing"
)

func TestGetConfig(t *testing.T) {
	testCases := []struct {
		name             string
		setUpFunc        func()
		ipToFindStr      string
		ipAddressListStr []string
		asnList          []int
		err              error
	}{
		{
			name: "Read Configuration From URL",
			setUpFunc: func() {
				os.Args = append(os.Args, "8.8.8.8")

			},
			ipToFindStr: "8.8.8.8",
			err:         nil,
		},
		{
			name: "Read Configuration From Local File",
			setUpFunc: func() {
				os.Setenv("CONFIG_FILE_PATH", "./config_file_test.txt")
			},
			ipToFindStr: "8.8.8.8",
			ipAddressListStr: []string{
				"8.8.8.0",
				"8.0.0.0",
				"8.0.0.0",
				"192.121.43.0",
			},
			asnList: []int{350, 352, 351, 156},
			err:     nil,
		},
	}

	_, cfgFileDefined := os.LookupEnv("CONFIG_FILE_PATH")
	if cfgFileDefined {
		t.Fatal("Please unset CONFIG_FILE_PATH environment variable before running tests")
	}

	for _, testCase := range testCases {
		testCase.setUpFunc()
		cfg, err := GetConfig(testCase.ipToFindStr)
		if err != testCase.err {
			t.Fatalf("%s: received error does not match: got %v, want %v", testCase.name, err, testCase.err)
		}

		if cfg.IPToFind.GetString() != testCase.ipToFindStr {
			t.Fatalf("%s: received IPToFind does not match: got %v, want %v", testCase.name, cfg.IPToFind.GetString(), testCase.ipToFindStr)
		}
		if len(testCase.ipAddressListStr) > 0 &&
			len(testCase.asnList) > 0 {
			for i := 0; i < len(testCase.ipAddressListStr) &&
				i < len(testCase.asnList); i++ {
				if testCase.ipAddressListStr[i] != cfg.IPAddressList[i].GetString() {
					t.Fatalf("%s: received IP address does not match: got %v, want %v", testCase.name, cfg.IPAddressList[i].GetString(), testCase.ipAddressListStr[i])
				}

				if testCase.asnList[i] != cfg.IPAddressList[i].GetAsn() {
					t.Fatalf("%s: received IP address does not match: got %v, want %v", testCase.name, cfg.IPAddressList[i].GetAsn(), testCase.asnList[i])
				}
			}
		}
	}
}
