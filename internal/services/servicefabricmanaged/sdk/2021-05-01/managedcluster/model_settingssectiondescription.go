package managedcluster

type SettingsSectionDescription struct {
	Name       string                         `json:"name"`
	Parameters []SettingsParameterDescription `json:"parameters"`
}
