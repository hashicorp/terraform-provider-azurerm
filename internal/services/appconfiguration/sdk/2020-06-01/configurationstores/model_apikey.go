package configurationstores

type ApiKey struct {
	ConnectionString *string `json:"connectionString,omitempty"`
	Id               *string `json:"id,omitempty"`
	LastModified     *string `json:"lastModified,omitempty"`
	Name             *string `json:"name,omitempty"`
	ReadOnly         *bool   `json:"readOnly,omitempty"`
	Value            *string `json:"value,omitempty"`
}
