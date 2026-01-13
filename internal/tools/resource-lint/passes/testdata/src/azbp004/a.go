package azbp004

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
)

type Props struct {
	Enabled                  *bool
	Name                     *string
	Count                    *int
	DynamicThrottlingEnabled *bool
	AllowProjectManagement   *bool
}

// Valid cases - should not trigger warnings

func validCase1() string {
	return pointer.From((*string)(nil))
}

func validCase2(props *Props) bool {
	return pointer.From(props.Enabled)
}

func validCase3(props *Props) string {
	return pointer.From(props.Name)
}

func validCase4(props *Props) int {
	return pointer.From(props.Count)
}

func validCase5(props *Props) {
	// Already using pointer.From - correct usage
	enabled := pointer.From(props.Enabled)
	useEnabled(enabled)
	useEnabledAgain(enabled)
}

// Invalid cases - should trigger AZBP004 warnings

func invalidCase1(props *Props) {
	enabled := false // want `AZBP004`
	if props.Enabled != nil {
		enabled = *props.Enabled
	}
	useEnabled(enabled)
}

func invalidCase2(props *Props) {
	name := "" // want `AZBP004`
	if props.Name != nil {
		name = *props.Name
	}
	useName(name)
}

func invalidCase3(props *Props) {
	count := 0 // want `AZBP004`
	if props.Count != nil {
		count = *props.Count
	}
	useCount(count)
}

func invalidCase4(props *Props) {
	dynamicThrottlingEnabled := false // want `AZBP004`
	if props.DynamicThrottlingEnabled != nil {
		dynamicThrottlingEnabled = *props.DynamicThrottlingEnabled
	}
	useDynamicThrottling(dynamicThrottlingEnabled)
}

func invalidCase5(props *Props) {
	enabled := false // want `AZBP004`
	if props.Enabled != nil {
		enabled = *props.Enabled
	}
	useEnabled(enabled)
	useEnabledAgain(enabled)
}

// Edge cases

func edgeCase1(props *Props) {
	// Has else branch - should not trigger
	enabled := false
	if props.Enabled != nil {
		enabled = *props.Enabled
	} else {
		enabled = true
	}
	useEnabled(enabled)
}

func edgeCase2(props *Props) {
	// Not initialized to zero value - should not trigger
	enabled := true
	if props.Enabled != nil {
		enabled = *props.Enabled
	}
	useEnabled(enabled)
}

func edgeCase3(props *Props) {
	// Different variable names - should not trigger
	enabled := false
	if props.Enabled != nil {
		other := *props.Enabled
		useEnabled(other)
	}
	useEnabled(enabled)
}

func edgeCase4(props *Props) {
	// Not adjacent statements - should not trigger
	enabled := false
	doSomething()
	if props.Enabled != nil {
		enabled = *props.Enabled
	}
	useEnabled(enabled)
}

func edgeCase5(props *Props) {
	// Using := in if body - should not trigger
	_ = false
	if props.Enabled != nil {
		enabled := *props.Enabled
		useEnabled(enabled)
	}
}

// Helper functions
func useEnabled(b bool)              {}
func useEnabledAgain(b bool)         {}
func useName(s string)               {}
func useCount(i int)                 {}
func useDynamicThrottling(b bool)    {}
func doSomething()                   {}
