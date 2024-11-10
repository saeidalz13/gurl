package conninfo

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

var testCir = ConnInfoResolver{}

func TestIsDomainLocalhost(t *testing.T) {
	tests := []struct {
		name        string
		domain      string
		expectedRes bool
	}{
		{name: "should_give_true_with_localhost_ip", domain: "127.0.0.1", expectedRes: true},
		{name: "should_give_true_with_localhost_ip", domain: "127.0.0.1:1111", expectedRes: true},
		{name: "should_give_true_with_localhost_string", domain: "localhost", expectedRes: true},
		{name: "should_give_true_with_localhost_string", domain: "localhost:9999", expectedRes: true},
		{name: "should_give_false_with_invalid_string", domain: "google.com", expectedRes: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testCir.domain = test.domain

			resp := testCir.isDomainLocalHost()
			if resp != test.expectedRes {
				t.Fatalf("expected: %v\tgot:%v", test.expectedRes, resp)
			}
		})
	}
}

func TestExtractPort(t *testing.T) {
	tests := []struct {
		name        string
		domain      string
		expectedRes int
		expectedErr string
	}{
		{name: "should_give_correct_port", domain: "127.0.0.1:1111", expectedRes: 1111},
		{name: "should_give_correct_port", domain: "127.0.0.1:9999", expectedRes: 9999},
		{name: "should_give_correct_port", domain: "127.0.0.1", expectedRes: 0, expectedErr: "domain must be in format of ip:port"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testCir.domain = test.domain

			res, err := testCir.extractPort()

			if test.expectedRes != res {
				t.Fatalf("expected result: %d\tgot: %d", test.expectedRes, res)
			}

			if err != nil && test.expectedErr != err.Error() {
				t.Fatalf("expected err: %s\tgot: %s", test.expectedErr, err.Error())
			}
		})
	}
}

func TestCacheIP(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatal(err)
	}

	ipCacheDir := filepath.Join(homeDir, ".gurl", "test")
	err = os.MkdirAll(ipCacheDir, 0o700)
	if err != nil {
		t.Fatal(err)
	}

	ipToSave := "0.0.0.0"
	testCir.ipCacheDir = ipCacheDir
	testCir.domain = "myunittestgurl.com"
	domainFile := filepath.Join(testCir.ipCacheDir, testCir.domain)

	err = testCir.cacheDomainIp(ipToSave)
	if err != nil {
		t.Fatal(err)
	}

	savedIp, err := os.ReadFile(domainFile)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(savedIp, []byte(ipToSave)) {
		t.Fatalf("expected saved IP: %s\tgot: %s", ipToSave, string(savedIp))
	}

	err = os.Remove(domainFile)
	if err != nil {
		t.Fatal(err)
	}

	err = os.Remove(ipCacheDir)
	if err != nil {
		t.Fatal(err)
	}
}
