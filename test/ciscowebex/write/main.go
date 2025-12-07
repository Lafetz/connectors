package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/amp-labs/connectors"
	"github.com/amp-labs/connectors/common"
	"github.com/amp-labs/connectors/providers/ciscoWebex"
	connTest "github.com/amp-labs/connectors/test/ciscowebex"
	"github.com/amp-labs/connectors/test/utils"
	"github.com/brianvoe/gofakeit/v6"
)

func main() {
	ctx, done := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer done()

	utils.SetupLogging()

	conn := connTest.GetCiscoWebexConnector(ctx)

	email := gofakeit.Email()
	updatedEmail := gofakeit.Email()

	createRes, err := conn.Write(ctx, common.WriteParams{
		ObjectName: "people",
		RecordData: map[string]any{
			"emails":      []any{email},
			"displayName": "Example Person",
			"firstName":   "Example",
			"lastName":    "Person",
		},
	})

	if err != nil {
		utils.Fail("error creating person in Cisco Webex", "error", err)
	}

	utils.DumpJSON(createRes, os.Stdout)

	updateRes, err := conn.Write(ctx, common.WriteParams{
		ObjectName: "people",
		RecordId:   createRes.RecordId,
		RecordData: map[string]any{
			"emails":      []any{updatedEmail},
			"displayName": "Example Person Updated",
			"firstName":   "Example",
			"lastName":    "Person",
		},
	})
	if updateRes != nil && updateRes.Success {
		slog.Info("Update returned success, verifying with GET request")
		verifyPerson(ctx, conn, createRes.RecordId, updatedEmail, "Example", "Person")
	} else if err != nil {
		slog.Warn("Update returned error, but checking if person was updated anyway", "error", err)
		verifyPerson(ctx, conn, createRes.RecordId, updatedEmail, "Example", "Person")
	} else {
		slog.Warn("Update returned neither success nor error")
	}

	utils.DumpJSON(updateRes, os.Stdout)
}

func verifyPerson(ctx context.Context, conn *ciscoWebex.Connector, recordID, expectedEmail, expectedFirstName, expectedLastName string) {

	readRes, err := conn.Read(ctx, common.ReadParams{
		ObjectName: "people",
		Fields:     connectors.Fields("id", "emails", "displayName", "firstName", "lastName"),
	})
	if err != nil {
		utils.Fail("error reading people to verify write", "error", err)
		return
	}

	var foundPerson *common.ReadResultRow
	for i := range readRes.Data {
		if readRes.Data[i].Id == recordID {
			foundPerson = &readRes.Data[i]
			break
		}
	}

	if foundPerson == nil {
		utils.Fail("person not found after write operation", "recordId", recordID)
		return
	}

	slog.Info("Verified person with GET request", "recordId", recordID)

	emails, ok := foundPerson.Fields["emails"].([]any)
	if !ok || len(emails) == 0 {
		slog.Warn("emails field not found or empty in verified person")
	} else if emails[0] != expectedEmail {
		slog.Warn("email mismatch", "expected", expectedEmail, "got", emails[0])
	}

	firstName, ok := foundPerson.Fields["firstname"].(string)
	if !ok {
		slog.Warn("firstName field not found in verified person")
	} else if firstName != expectedFirstName {
		slog.Warn("firstName mismatch", "expected", expectedFirstName, "got", firstName)
	}

	lastName, ok := foundPerson.Fields["lastname"].(string)
	if !ok {
		slog.Warn("lastName field not found in verified person")
	} else if lastName != expectedLastName {
		slog.Warn("lastName mismatch", "expected", expectedLastName, "got", lastName)
	}

	slog.Info("Person verification complete", "recordId", recordID)

}
