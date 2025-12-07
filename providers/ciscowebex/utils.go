package ciscowebex

import (
	"github.com/amp-labs/connectors/common"
	"github.com/amp-labs/connectors/internal/components"
	"github.com/amp-labs/connectors/internal/httpkit"
	"github.com/spyzhov/ajson"
)

const objectNamePeople = "people"

func supportedOperations() components.EndpointRegistryInput {
	return components.EndpointRegistryInput{
		common.ModuleRoot: {
			{
				Endpoint: objectNamePeople,
				Support:  components.ReadSupport,
			},
			{
				Endpoint: objectNamePeople,
				Support:  components.WriteSupport,
			},
			{
				Endpoint: objectNamePeople,
				Support:  components.DeleteSupport,
			},
		},
	}
}

func getNextRecordsURL(resp *common.JSONHTTPResponse) common.NextPageFunc {
	return func(_ *ajson.Node) (string, error) {
		return httpkit.HeaderLink(resp, "next"), nil
	}
}
