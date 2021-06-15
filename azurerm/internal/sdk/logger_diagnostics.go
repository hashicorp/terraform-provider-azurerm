package sdk

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

var _ Logger = &DiagnosticsLogger{}

type DiagnosticsLogger struct {
	diagnostics diag.Diagnostics
}

func (d *DiagnosticsLogger) Info(message string) {
	log.Printf("[INFO] %s", message)
}

func (d *DiagnosticsLogger) Infof(format string, args ...interface{}) {
	log.Printf("[INFO] "+format, args...)
}

func (d *DiagnosticsLogger) Warn(message string) {
	d.diagnostics = append(d.diagnostics, diag.Diagnostic{
		Severity:      diag.Warning,
		Summary:       message,
		Detail:        message,
		AttributePath: nil,
	})
}

func (d *DiagnosticsLogger) Warnf(format string, args ...interface{}) {
	d.diagnostics = append(d.diagnostics, diag.Diagnostic{
		Severity:      diag.Warning,
		Summary:       fmt.Sprintf(format, args...),
		Detail:        fmt.Sprintf(format, args...),
		AttributePath: nil,
	})
}
