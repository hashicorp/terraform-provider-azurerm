package entities

type MetaDataLevel string

var (
	NoMetaData      MetaDataLevel = "nometadata"
	MinimalMetaData MetaDataLevel = "minimalmetadata"
	FullMetaData    MetaDataLevel = "fullmetadata"
)
