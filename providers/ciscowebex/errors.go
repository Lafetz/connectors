package ciscowebex

import (
	"fmt"

	"github.com/amp-labs/connectors/common/interpreter"
)

// nolint:gochecknoglobals
var errorFormats = interpreter.NewFormatSwitch(
	[]interpreter.FormatTemplate{
		{
			MustKeys: []string{"message"},
			Template: func() interpreter.ErrorDescriptor { return &ResponseMessageError{} },
		},
	}...,
)

type ResponseMessageError struct {
	Message string `json:"message"`
}

func (r ResponseMessageError) CombineErr(base error) error {
	if len(r.Message) != 0 {
		return fmt.Errorf("%w: %s", base, r.Message)
	}

	return base
}
