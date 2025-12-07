package main

import (
	"context"
	"os/signal"
	"syscall"

	connTest "github.com/amp-labs/connectors/test/ciscowebex"
	"github.com/amp-labs/connectors/test/utils"
	"github.com/amp-labs/connectors/test/utils/testscenario"
)

func main() {
	ctx, done := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer done()

	utils.SetupLogging()

	conn := connTest.GetCiscoWebexConnector(ctx)

	// Use ContainsRead because some fields are optional (not all users have phoneNumbers, lastActivity, etc.)
	testscenario.ValidateMetadataContainsRead(ctx, conn, "people", nil)
}
