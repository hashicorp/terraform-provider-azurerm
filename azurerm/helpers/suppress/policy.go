package suppress

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/structure"
	"strings"
)

// Policy Definition Metadata adds additional fields into the metadata string. We'll remove those before checking for a diff in the JSON
func SuppressPolicyDefinitionMetadata(_, old, new string, _ *schema.ResourceData) bool {
	splitMetaData := make([]string, 0)

	oldMetaDataKeys := strings.Split(old, ",")
	for _, key := range oldMetaDataKeys {
		if !strings.Contains(key, "createdBy") && !strings.Contains(key, "createdOn") && !strings.Contains(key, "updatedBy") && !strings.Contains(key, "updatedOn") {
			splitMetaData = append(splitMetaData, key)
		}
	}

	return structure.SuppressJsonDiff("", strings.Join(splitMetaData, ","), new, nil)
}
