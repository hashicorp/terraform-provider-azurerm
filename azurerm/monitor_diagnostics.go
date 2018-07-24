package azurerm

func getAllDiagnosticSettings(targetResourceId string, meta interface{}) (*[]interface{}, *[]interface{}, error) {
	client := meta.(*ArmClient).monitorDiagnosticSettingsCategoryClient
	ctx := meta.(*ArmClient).StopContext
	returnMetricSettings := make([]interface{}, 0)
	returnLogSettings := make([]interface{}, 0)

	categoryList, err := client.List(ctx, targetResourceId)
	if err != nil {
		return nil, nil, err
	}

	for _, item := range *categoryList.Value {
		if item.DiagnosticSettingsCategory.CategoryType == "Metrics" {
			returnMetricSettings = append(returnMetricSettings, *item.Name)
		}
		if item.DiagnosticSettingsCategory.CategoryType == "Logs" {
			returnLogSettings = append(returnLogSettings, *item.Name)
		}
	}

	return &returnMetricSettings, &returnLogSettings, nil
}
