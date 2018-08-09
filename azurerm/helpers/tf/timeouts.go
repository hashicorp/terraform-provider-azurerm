package tf

import "github.com/hashicorp/terraform/helper/schema"

func TimeoutForCreateUpdate(d *schema.ResourceData) string {
	if d.IsNewResource() {
		return schema.TimeoutCreate
	}

	return schema.TimeoutUpdate
}
