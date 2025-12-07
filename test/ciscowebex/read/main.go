package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/amp-labs/connectors"
	"github.com/amp-labs/connectors/common"
	connTest "github.com/amp-labs/connectors/test/ciscowebex"
	"github.com/amp-labs/connectors/test/utils"
)

func main() {
	ctx, done := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer done()

	utils.SetupLogging()

	conn := connTest.GetCiscoWebexConnector(ctx)

	res, err := conn.Read(ctx, common.ReadParams{
		ObjectName: "people",
		Fields:     connectors.Fields("id"),
	})
	if err != nil {
		utils.Fail("error reading from Cisco Webex", "error", err)
	}

	utils.DumpJSON(res, os.Stdout)
}
