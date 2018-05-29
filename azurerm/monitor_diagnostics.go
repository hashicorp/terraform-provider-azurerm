package azurerm

import "strings"

type MonitorDiagnosticId struct {
	ResourceID string
	Name       string
}

func parseMonitorDiagnosticId(monitorId string) MonitorDiagnosticId {
	var returnId MonitorDiagnosticId
	returnId.ResourceID = retrieveResourceId(monitorId)
	returnId.Name = retrieveDiagnosticName(monitorId)
	return returnId
}

func retrieveResourceId(diagnosticSettingId string) string {
	substring := "/providers/microsoft.insights/diagnosticSettings/"
	return diagnosticSettingId[0:strings.Index(diagnosticSettingId, substring)]
}

func retrieveDiagnosticName(diagnosticSettingId string) string {
	substring := "/providers/microsoft.insights/diagnosticSettings/"
	return diagnosticSettingId[len(substring)+strings.Index(diagnosticSettingId, substring):]
}
