package diagutils

import (
	"errors"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

type ErrorDiags diag.Diagnostics

func (diags ErrorDiags) Errors() []error {
	var es []error
	for i := range diags {
		if diags[i].Severity == diag.Error {
			s := fmt.Sprintf("Error: %s", diags[i].Summary)
			if diags[i].Detail != "" {
				s = fmt.Sprintf("%s: %s", s, diags[i].Detail)
			}
			es = append(es, errors.New(s))
		}
	}
	return es
}

func (diags ErrorDiags) Error() string {
	return multierror.ListFormatFunc(diags.Errors())
}

type WarningDiags diag.Diagnostics

func (diags WarningDiags) Warnings() []string {
	var ws []string
	for i := range diags {
		if diags[i].Severity == diag.Warning {
			s := fmt.Sprintf("Warning: %s", diags[i].Summary)
			if diags[i].Detail != "" {
				s = fmt.Sprintf("%s: %s", s, diags[i].Detail)
			}
			ws = append(ws, s)
		}
	}
	return ws
}
