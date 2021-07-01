package pluginsdk

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

type TestCheckFunc = resource.TestCheckFunc

type InstanceState = terraform.InstanceState
