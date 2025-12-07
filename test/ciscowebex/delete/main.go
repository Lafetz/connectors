package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

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

	createRes := createPerson(ctx, conn, email)
	if createRes == nil || createRes.RecordId == "" {
		utils.Fail("failed to create person for delete test", "recordId", createRes)
		return
	}

	slog.Info("Created person for deletion", "recordId", createRes.RecordId)
	utils.DumpJSON(createRes, os.Stdout)

	deleteRes := deletePerson(ctx, conn, createRes.RecordId)
	if deleteRes == nil {
		utils.Fail("delete operation returned nil result")
		return
	}

	slog.Info("Successfully deleted person", "recordId", createRes.RecordId)
	utils.DumpJSON(deleteRes, os.Stdout)

	slog.Info("Attempting to delete already-deleted person (should fail)...")
	_, err := conn.Delete(ctx, common.DeleteParams{
		ObjectName: "people",
		RecordId:   createRes.RecordId,
	})
	if err == nil {
		utils.Fail("expected error when deleting non-existent person", "error", err)
	}
	slog.Info("Failed to delete non-existent person as expected")
}

func createPerson(ctx context.Context, conn *ciscoWebex.Connector, email string) *common.WriteResult {
	res, err := conn.Write(ctx, common.WriteParams{
		ObjectName: "people",
		RecordData: map[string]any{
			"emails":      []any{email},
			"displayName": "Test Person for Deletion",
			"firstName":   "Test",
			"lastName":    "Delete",
		},
	})
	if err != nil {
		utils.Fail("error creating person in Cisco Webex", "error", err)
		return nil
	}

	return res
}

func deletePerson(ctx context.Context, conn *ciscoWebex.Connector, recordID string) *common.DeleteResult {
	deleteRes, err := conn.Delete(ctx, common.DeleteParams{
		ObjectName: "people",
		RecordId:   recordID,
	})
	if err != nil {
		utils.Fail("error deleting person in Cisco Webex", "error", err, "recordId", recordID)
		return nil
	}

	return deleteRes
}
