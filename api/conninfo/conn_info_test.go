package conninfo

import "testing"

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
