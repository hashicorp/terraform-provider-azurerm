package acceptance

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
)

// This file is intended to provide a transition from Plugin SDKv1 to Plugin SDKv2
// without introducing a merge conflict into every PR.

type InstanceState = terraform.InstanceState

type State = terraform.State

type TestStep = resource.TestStep

type StateChangeConf = resource.StateChangeConf

type TestCheckFunc = resource.TestCheckFunc

func ComposeTestCheckFunc(fs ...resource.TestCheckFunc) pluginsdk.TestCheckFunc {
	return resource.ComposeTestCheckFunc(fs...)
}

func ComposeAggregateTestCheckFunc(fs ...resource.TestCheckFunc) pluginsdk.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(fs...)
}

// @tombuildsstuff:
// Below this point are convenience methods which exist to allow existing code
// to compile. Whilst these are not initially deprecated, they will be in the
// future - this is done to allow a gradual transition path for existing PR's.

// TestCheckResourceAttr is a wrapper to enable builds to continue
func TestCheckResourceAttr(name, key, value string) pluginsdk.TestCheckFunc {
	// TODO: move this comment up a level in the future
	// Deprecated: use `check.That(name).Key(key).HasValue(value)` instead
	return resource.TestCheckResourceAttr(name, key, value)
}

func TestCheckResourceAttrSet(name, key string) pluginsdk.TestCheckFunc {
	// TODO: move this comment up a level in the future
	// Deprecated: use `check.That(name).Key(key).HasValue(value)` instead
	return resource.TestCheckResourceAttrSet(name, key)
}

// TestCheckOutput is a wrapper to enable builds to continue
func TestCheckOutput(name, value string) pluginsdk.TestCheckFunc {
	// TODO: move this comment up a level in the future
	// Deprecated: use `check.That(name).Key(key).HasValue(value)` instead
	return resource.TestCheckOutput(name, value)
}

// TestMatchOutput is a wrapper to enable builds to continue
func TestMatchOutput(name string, r *regexp.Regexp) pluginsdk.TestCheckFunc {
	// TODO: move this comment up a level in the future
	return resource.TestMatchOutput(name, r)
}

// TestMatchResourceAttr is a TestCheckFunc which checks that the value
// in state for the given name/key combination matches the given regex.
func TestMatchResourceAttr(name, key string, r *regexp.Regexp) pluginsdk.TestCheckFunc {
	// TODO: move this comment up a level in the future
	// Deprecated: use `check.That(name).Key(key).MatchesRegex(r)` instead
	return resource.TestMatchResourceAttr(name, key, r)
}

// TestCheckResourceAttrPair is a TestCheckFunc which validates that the values
// in state for a pair of name/key combinations are equal.
func TestCheckResourceAttrPair(nameFirst, keyFirst, nameSecond, keySecond string) resource.TestCheckFunc {
	// TODO: move this comment up a level in the future
	// Deprecated: use this instead:
	//  check.That(nameFirst).Key(keyFirst).MatchesOtherKey(
	//    check.That(nameSecond).Key(keySecond),
	//  ),
	return resource.TestCheckResourceAttrPair(nameFirst, keyFirst, nameSecond, keySecond)
}

// TestCheckNoResourceAttr is a TestCheckFunc which ensures that
// NO value exists in state for the given name/key combination.
func TestCheckNoResourceAttr(name, key string) resource.TestCheckFunc {
	// TODO: move this comment up a level in the future
	// Deprecated: use `check.That(name).Key(key).DoesNotExist()` instead
	return resource.TestCheckNoResourceAttr(name, key)
}
