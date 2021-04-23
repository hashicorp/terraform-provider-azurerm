package pluginsdk

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type CustomizeDiffFunc = func(context.Context, *ResourceDiff, interface{}) error
type ValueChangeConditionFunc = func(ctx context.Context, old, new, meta interface{}) bool

// CustomDiffWithAll returns a CustomizeDiffFunc that runs all of the given
// CustomizeDiffFuncs and returns all of the errors produced.
//
// If one function produces an error, functions after it are still run.
// If this is not desirable, use function Sequence instead.
//
// If multiple functions returns errors, the result is a multierror.
func CustomDiffWithAll(funcs ...CustomizeDiffFunc) schema.CustomizeDiffFunc {
	return func(d *schema.ResourceDiff, meta interface{}) error {
		var err error
		for _, f := range funcs {
			ctx := context.TODO()
			thisErr := f(ctx, d, meta)
			if thisErr != nil {
				err = multierror.Append(err, thisErr)
			}
		}
		return err
	}
}

// CustomDiffInSequence returns a CustomizeDiffFunc that runs all of the given
// CustomizeDiffFuncs in sequence, stopping at the first one that returns
// an error and returning that error.
//
// If all functions succeed, the combined function also succeeds.
func CustomDiffInSequence(funcs ...CustomizeDiffFunc) schema.CustomizeDiffFunc {
	return func(d *schema.ResourceDiff, meta interface{}) error {
		for _, f := range funcs {
			err := f(context.TODO(), d, meta)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

// ForceNewIfChange returns a CustomizeDiffFunc that flags the given key as
// requiring a new resource if the given condition function returns true.
//
// The return value of the condition function is ignored if the old and new
// values compare equal, since no attribute diff is generated in that case.
//
// This function is similar to ForceNewIf but provides the condition function
// only the old and new values of the given key, which leads to more compact
// and explicit code in the common case where the decision can be made with
// only the specific field value.
func ForceNewIfChange(key string, f ValueChangeConditionFunc) CustomizeDiffFunc {
	return func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
		old, new := d.GetChange(key)
		if f(ctx, old, new, meta) {
			d.ForceNew(key)
		}
		return nil
	}
}
