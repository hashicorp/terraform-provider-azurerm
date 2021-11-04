package accounts

type HiveMetastore struct {
	Id         *string                  `json:"id,omitempty"`
	Name       *string                  `json:"name,omitempty"`
	Properties *HiveMetastoreProperties `json:"properties,omitempty"`
	Type       *string                  `json:"type,omitempty"`
}
