package ciscoWebex

import (
	"net/http"
	"testing"

	"github.com/amp-labs/connectors"
	"github.com/amp-labs/connectors/common"
	"github.com/amp-labs/connectors/test/utils/mockutils"
	"github.com/amp-labs/connectors/test/utils/mockutils/mockcond"
	"github.com/amp-labs/connectors/test/utils/mockutils/mockserver"
	"github.com/amp-labs/connectors/test/utils/testroutines"
	"github.com/amp-labs/connectors/test/utils/testutils"
)

func TestRead(t *testing.T) {
	t.Parallel()

	responseInvalidPath := testutils.DataFromFile(t, "invalid-path.json")
	responseReadEmpty := testutils.DataFromFile(t, "read-empty.json")
	responseReadPeople := testutils.DataFromFile(t, "read-people.json")
	responseReadPeopleFirstPage := testutils.DataFromFile(t, "read-people-first-page.json")
	responseReadPeopleSecondPage := testutils.DataFromFile(t, "read-people-second-page.json")

	tests := []testroutines.Read{
		{
			Name:         "Read object must be included",
			Server:       mockserver.Dummy(),
			ExpectedErrs: []error{common.ErrMissingObjects},
		},
		{
			Name:         "At least one field is requested",
			Input:        common.ReadParams{ObjectName: "people"},
			Server:       mockserver.Dummy(),
			ExpectedErrs: []error{common.ErrMissingFields},
		},
		{
			Name:         "Unknown objects are not supported",
			Input:        common.ReadParams{ObjectName: "unknown", Fields: connectors.Fields("id")},
			Server:       mockserver.Dummy(),
			ExpectedErrs: []error{common.ErrOperationNotSupportedForObject},
		},
		{
			Name: "Read invalid path",
			Input: common.ReadParams{
				ObjectName: "people",
				Fields:     connectors.Fields("id", "displayName"),
				NextPage:   testroutines.URLTestServer + "/v1/invalidpath",
			},
			Server: mockserver.Conditional{
				Setup: mockserver.ContentJSON(),
				If:    mockcond.Path("/v1/invalidpath"),
				Then:  mockserver.Response(http.StatusNotFound, responseInvalidPath),
			}.Server(),
			ExpectedErrs: []error{common.ErrNotFound},
		},
		{
			Name: "Read empty items",
			Input: common.ReadParams{
				ObjectName: "people",
				Fields:     connectors.Fields("id", "displayName"),
			},
			Server: mockserver.Fixed{
				Setup:  mockserver.ContentJSON(),
				Always: mockserver.Response(http.StatusOK, responseReadEmpty),
			}.Server(),
			Expected: &common.ReadResult{
				Rows: 0,
				Data: []common.ReadResultRow{},
				Done: true,
			},
			ExpectedErrs: nil,
		},
		{
			Name: "Read people",
			Input: common.ReadParams{
				ObjectName: "people",
				Fields:     connectors.Fields("id", "displayName", "emails", "firstName", "lastName"),
			},
			Server: mockserver.Conditional{
				Setup: mockserver.ContentJSON(),
				If:    mockcond.Path("/v1/people"),
				Then:  mockserver.Response(http.StatusOK, responseReadPeople),
			}.Server(),
			Comparator: testroutines.ComparatorSubsetRead,
			Expected: &common.ReadResult{
				Rows: 2,
				Data: []common.ReadResultRow{
					{
						Fields: map[string]any{
							"id":          "Y2lzY29zcGFyazovL3VzL1BFT1BMRS8xNzQ5MmVmNC1iMGFjLTRjMTYtYWRjNS02OTNkMjEyM2Q5MmI",
							"displayname": "admin@example.wbx.ai",
							"emails":      []any{"admin@example.wbx.ai"},
							"firstname":   "admin",
							"lastname":    "admin",
						},
						Raw: map[string]any{
							"id":          "Y2lzY29zcGFyazovL3VzL1BFT1BMRS8xNzQ5MmVmNC1iMGFjLTRjMTYtYWRjNS02OTNkMjEyM2Q5MmI",
							"displayName": "admin@example.wbx.ai",
						},
					},
					{
						Fields: map[string]any{
							"id":          "Y2lzY29zcGFyazovL3VzL1BFT1BMRS9iMTQzNmI1OS02MDExLTQ2OTEtODBjZC0xN2Y0NGRhODVmNTk",
							"displayname": "testuser",
							"emails":      []any{"testuser@example.com"},
							"firstname":   "test",
							"lastname":    "user",
						},
						Raw: map[string]any{
							"id":          "Y2lzY29zcGFyazovL3VzL1BFT1BMRS9iMTQzNmI1OS02MDExLTQ2OTEtODBjZC0xN2Y0NGRhODVmNTk",
							"displayName": "testuser",
						},
					},
				},
				Done: true,
			},
			ExpectedErrs: nil,
		},
		{
			Name: "Read people first page with pagination",
			Input: common.ReadParams{
				ObjectName: "people",
				Fields:     connectors.Fields("id", "displayName", "emails"),
			},
			Server: mockserver.Conditional{
				Setup: mockserver.ContentJSON(),
				If:    mockcond.Path("/v1/people"),
				Then: mockserver.ResponseChainedFuncs(
					mockserver.Header("Link", `<https://webexapis.com/v1/people?max=1&cursor=next_cursor_token>; rel="next"`),
					mockserver.Response(http.StatusOK, responseReadPeopleFirstPage),
				),
			}.Server(),
			Comparator: testroutines.ComparatorSubsetRead,
			Expected: &common.ReadResult{
				Rows: 1,
				Data: []common.ReadResultRow{
					{
						Fields: map[string]any{
							"id":          "Y2lzY29zcGFyazovL3VzL1BFT1BMRS8xNzQ5MmVmNC1iMGFjLTRjMTYtYWRjNS02OTNkMjEyM2Q5MmI",
							"displayname": "admin@example.wbx.ai",
							"emails":      []any{"admin@example.wbx.ai"},
						},
						Raw: map[string]any{
							"id":          "Y2lzY29zcGFyazovL3VzL1BFT1BMRS8xNzQ5MmVmNC1iMGFjLTRjMTYtYWRjNS02OTNkMjEyM2Q5MmI",
							"displayName": "admin@example.wbx.ai",
							"emails":      []any{"admin@example.wbx.ai"},
							"nickName":    "admin",
							"firstName":   "admin",
							"lastName":    "admin",
							"orgId":       "Y2lzY29zcGFyazovL3VzL09SR0FOSVpBVElPTi9iMmYzYzMwNC1kYTI1LTQxZTMtOWY4NC0yOGM1NGEyZmMyZjQ",
							"type":        "person",
						},
					},
				},
				NextPage: "https://webexapis.com/v1/people?max=1&cursor=next_cursor_token",
				Done:     false,
			},
			ExpectedErrs: nil,
		},
		{
			Name: "Read people second page using NextPage token",
			Input: common.ReadParams{
				ObjectName: "people",
				Fields:     connectors.Fields("id", "displayName", "emails"),
				NextPage:   testroutines.URLTestServer + "/v1/people?max=1&cursor=next_cursor_token",
			},
			Server: mockserver.Conditional{
				Setup: mockserver.ContentJSON(),
				If: mockcond.And{
					mockcond.Path("/v1/people"),
					mockcond.QueryParam("cursor", "next_cursor_token"),
				},
				Then: mockserver.Response(http.StatusOK, responseReadPeopleSecondPage),
			}.Server(),
			Comparator: testroutines.ComparatorSubsetRead,
			Expected: &common.ReadResult{
				Rows: 1,
				Data: []common.ReadResultRow{
					{
						Fields: map[string]any{
							"id":          "Y2lzY29zcGFyazovL3VzL1BFT1BMRS9iMTQzNmI1OS02MDExLTQ2OTEtODBjZC0xN2Y0NGRhODVmNTk",
							"displayname": "testuser",
							"emails":      []any{"testuser@example.com"},
						},
						Raw: map[string]any{
							"id":          "Y2lzY29zcGFyazovL3VzL1BFT1BMRS9iMTQzNmI1OS02MDExLTQ2OTEtODBjZC0xN2Y0NGRhODVmNTk",
							"displayName": "testuser",
							"emails":      []any{"testuser@example.com"},
							"nickName":    "testuser",
							"firstName":   "test",
							"lastName":    "user",
							"orgId":       "Y2lzY29zcGFyazovL3VzL09SR0FOSVpBVElPTi9iMmYzYzMwNC1kYTI1LTQxZTMtOWY4NC0yOGM1NGEyZmMyZjQ",
							"type":        "person",
						},
					},
				},
				NextPage: "",
				Done:     true,
			},
			ExpectedErrs: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			t.Parallel()

			tt.Run(t, func() (connectors.ReadConnector, error) {
				return constructTestConnector(tt.Server.URL)
			})
		})
	}
}

func constructTestConnector(serverURL string) (*Connector, error) {
	connector, err := NewConnector(
		common.ConnectorParams{
			AuthenticatedClient: mockutils.NewClient(),
		},
	)
	if err != nil {
		return nil, err
	}

	connector.SetUnitTestBaseURL(mockutils.ReplaceURLOrigin(connector.HTTPClient().Base, serverURL))

	return connector, nil
}
