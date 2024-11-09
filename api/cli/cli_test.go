package cli

import (
	"testing"

	"github.com/saeidalz13/gurl/internal/httpconstants"
)

func TestMustDetermineDataInfo(t *testing.T) {
	tests := []struct {
		expectedDataType uint8
		name             string
		jsonPtr          string
		textPtr          string
		expectedData     string
	}{
		{name: "both_ptr_empty", jsonPtr: "", textPtr: "", expectedDataType: 0, expectedData: ""},
		{name: "json_ptr_value", jsonPtr: "s", textPtr: "", expectedDataType: httpconstants.DataTypeJson, expectedData: "s"},
		{name: "text_ptr_value", jsonPtr: "", textPtr: "s", expectedDataType: httpconstants.DataTypeText, expectedData: "s"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dataType, data := mustDetermineDataInfo(&test.jsonPtr, &test.textPtr)

			if dataType != test.expectedDataType {
				t.Fatalf("expected data type:%d\tgot:%d\t", test.expectedDataType, dataType)
			}

			if data != test.expectedData {
				t.Fatalf("expected data:%s\tgot:%s\t", test.expectedData, data)
			}
		})
	}
}
