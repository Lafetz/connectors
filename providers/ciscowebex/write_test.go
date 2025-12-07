package ciscowebex

import (
	"net/http"
	"testing"

	"github.com/amp-labs/connectors"
	"github.com/amp-labs/connectors/common"
	"github.com/amp-labs/connectors/test/utils/mockutils/mockcond"
	"github.com/amp-labs/connectors/test/utils/mockutils/mockserver"
	"github.com/amp-labs/connectors/test/utils/testroutines"
	"github.com/amp-labs/connectors/test/utils/testutils"
)

// nolint:funlen,gocognit,cyclop
func TestWrite(t *testing.T) {
	t.Parallel()

	responseCreatePerson := testutils.DataFromFile(t, "write-create-person.json")
	responseUpdatePerson := testutils.DataFromFile(t, "write-update-person.json")

	tests := []testroutines.Write{
		{
			Name:         "Write object must be included",
			Server:       mockserver.Dummy(),
			ExpectedErrs: []error{common.ErrMissingObjects},
		},
		{
			Name:         "Write needs data payload",
			Input:        common.WriteParams{ObjectName: "people"},
			Server:       mockserver.Dummy(),
			ExpectedErrs: []error{common.ErrMissingRecordData},
		},
		{
			Name:     "Unknown object name is not supported",
			Input:    common.WriteParams{ObjectName: "invalid", RecordData: "dummy"},
			Server:   mockserver.Dummy(),
			Expected: nil,
			ExpectedErrs: []error{
				common.ErrOperationNotSupportedForObject,
			},
		},
		{
			Name: "Write must act as a Create (POST)",
			Input: common.WriteParams{
				ObjectName: "people",
				RecordData: map[string]any{
					"emails":      []any{"exmple@example.com"},
					"displayName": "Example Person",
					"firstName":   "Example",
					"lastName":    "Person",
				},
			},
			Server: mockserver.Conditional{
				Setup: mockserver.ContentJSON(),
				If:    mockcond.MethodPOST(),
				Then:  mockserver.Response(http.StatusOK, responseCreatePerson),
			}.Server(),
			Comparator: testroutines.ComparatorSubsetWrite,
			Expected: &common.WriteResult{
				Success:  true,
				RecordId: "Y2lzY29zcGFyazovL3VzL1BFT1BMRS9hZWRiZmVlYy1kZGI1LTRhNTItYjZhZS1lYzExMTVlZmZjZTg",
				Errors:   nil,
				Data: map[string]any{
					"id":          "Y2lzY29zcGFyazovL3VzL1BFT1BMRS9hZWRiZmVlYy1kZGI1LTRhNTItYjZhZS1lYzExMTVlZmZjZTg",
					"displayName": "Example Person",
					"firstName":   "Example",
					"lastName":    "Person",
					"emails":      []any{"exmple@example.com"},
					"type":        "person",
				},
			},
			ExpectedErrs: nil,
		},
		{
			Name: "Write must act as an Update (PUT)",
			Input: common.WriteParams{
				ObjectName: "people",
				RecordId:   "Y2lzY29zcGFyazovL3VzL1BFT1BMRS9hZWRiZmVlYy1kZGI1LTRhNTItYjZhZS1lYzExMTVlZmZjZTg",
				RecordData: map[string]any{
					"emails":      []any{"exmple.updated@example.com"},
					"displayName": "Example Person Updated",
					"firstName":   "Example",
					"lastName":    "Person",
				},
			},
			Server: mockserver.Conditional{
				Setup: mockserver.ContentJSON(),
				If: mockcond.And{
					mockcond.MethodPUT(),
					mockcond.Path("/v1/people/Y2lzY29zcGFyazovL3VzL1BFT1BMRS9hZWRiZmVlYy1kZGI1LTRhNTItYjZhZS1lYzExMTVlZmZjZTg"),
				},
				Then: mockserver.Response(http.StatusOK, responseUpdatePerson),
			}.Server(),
			Comparator: testroutines.ComparatorSubsetWrite,
			Expected: &common.WriteResult{
				Success:  true,
				RecordId: "Y2lzY29zcGFyazovL3VzL1BFT1BMRS9hZWRiZmVlYy1kZGI1LTRhNTItYjZhZS1lYzExMTVlZmZjZTg",
				Errors:   nil,
				Data: map[string]any{
					"id":          "Y2lzY29zcGFyazovL3VzL1BFT1BMRS9hZWRiZmVlYy1kZGI1LTRhNTItYjZhZS1lYzExMTVlZmZjZTg",
					"displayName": "Example Person Updated",
					"firstName":   "Example",
					"lastName":    "Person",
					"emails":      []any{"exmple.updated@example.com"},
					"type":        "person",
				},
			},
			ExpectedErrs: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			t.Parallel()

			tt.Run(t, func() (connectors.WriteConnector, error) {
				return constructTestConnector(tt.Server.URL)
			})
		})
	}
}
