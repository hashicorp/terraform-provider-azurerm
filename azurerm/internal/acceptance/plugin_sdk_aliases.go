package acceptance

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// This file is intended to provide a transition from Plugin SDKv1 to Plugin SDKv2
// without introducing a merge conflict into every PR.

type InstanceState = terraform.InstanceState

type TestStep = resource.TestStep

func ComposeTestCheckFunc(fs ...resource.TestCheckFunc) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(fs...)
}
