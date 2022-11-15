package pluginsdk

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"

var (
	ExpandJsonFromString = structure.ExpandJsonFromString
	FlattenJsonToString  = structure.FlattenJsonToString
	SuppressJsonDiff     = structure.SuppressJsonDiff
)
