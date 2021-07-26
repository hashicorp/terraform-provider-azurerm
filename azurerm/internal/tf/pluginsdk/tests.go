package pluginsdk

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

type TestCheckFunc = resource.TestCheckFunc

type InstanceState = terraform.InstanceState
