package tables

type MetaDataLevel string

var (
	NoMetaData      MetaDataLevel = "nometadata"
	MinimalMetaData MetaDataLevel = "minimalmetadata"
	FullMetaData    MetaDataLevel = "fullmetadata"
)

type GetResultItem struct {
	TableName string `json:"TableName"`

	// Optional, depending on the MetaData Level
	ODataType     string `json:"odata.type,omitempty"`
	ODataID       string `json:"odata.id,omitEmpty"`
	ODataEditLink string `json:"odata.editLink,omitEmpty"`
}

type SignedIdentifier struct {
	Id           string       `xml:"Id"`
	AccessPolicy AccessPolicy `xml:"AccessPolicy"`
}

type AccessPolicy struct {
	Start      string `xml:"Start"`
	Expiry     string `xml:"Expiry"`
	Permission string `xml:"Permission"`
}
