package parse

import (
	"fmt"
	"strings"
)

const aadDiagnosticSettingIdPrefix = "/providers/Microsoft.AADIAM/diagnosticSettings/"

type MonitorAADDiagnosticSettingId struct {
	Name string
}

func NewMonitorAADDiagnosticSettingID(name string) MonitorAADDiagnosticSettingId {
	return MonitorAADDiagnosticSettingId{Name: name}
}

func (id MonitorAADDiagnosticSettingId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Monitor AAD Diagnostic Setting", segmentsStr)
}

func (id MonitorAADDiagnosticSettingId) ID() string {
	fmtString := aadDiagnosticSettingIdPrefix + "%s"
	return fmt.Sprintf(fmtString, id.Name)
}

// MonitorAADDiagnosticSettingID parses a MonitorAADDiagnosticSetting ID into an MonitorAADDiagnosticSettingId struct
func MonitorAADDiagnosticSettingID(input string) (*MonitorAADDiagnosticSettingId, error) {
	if !strings.HasPrefix(input, aadDiagnosticSettingIdPrefix) {
		return nil, fmt.Errorf("invalid Monitor AAD Diagnostic Setting ID - ID should starts with %s", aadDiagnosticSettingIdPrefix)
	}
	name := strings.TrimPrefix(input, aadDiagnosticSettingIdPrefix)
	if name == "" {
		return nil, fmt.Errorf("ID was missing the 'diagnosticSettings' element")
	}
	return &MonitorAADDiagnosticSettingId{Name: name}, nil
}
