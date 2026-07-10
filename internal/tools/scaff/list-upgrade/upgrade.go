package list_upgrade

import "fmt"

// UpgradeOptions configures which upgrades Upgrade applies to a resource.
type UpgradeOptions struct {
	AddIdentity    bool
	ExtractFlatten bool

	// ReadModel is the SDK read-model type name (e.g. "AzureMonitorWorkspaceResource"),
	// required when extracting a new flatten method.
	ReadModel string

	// ResourceName is the snake-case resource name (e.g. "monitor_workspace") used
	// in the identity-test go:generate directive.
	ResourceName string

	// IdentityProperties is the -properties value for the identity go:generate
	// directive; defaults to "name,resource_group_name" when empty.
	IdentityProperties string
}

// Plan is a read-only assessment of what an upgrade would do to a resource.
type Plan struct {
	Kind          Kind
	Supported     bool
	Reason        string
	HasIdentity   bool
	HasFlatten    bool
	NeedsIdentity bool
	NeedsFlatten  bool
}

// PlanUpgrade assesses the resource and reports what an upgrade would change,
// without mutating anything.
func (r *Resource) PlanUpgrade() Plan {
	p := Plan{
		Kind:          r.Kind,
		HasIdentity:   r.HasIdentity,
		HasFlatten:    r.HasFlatten,
		NeedsIdentity: !r.HasIdentity,
		NeedsFlatten:  !r.HasFlatten,
	}
	switch r.Kind {
	case KindTyped:
		p.Supported = true
	case KindUntyped:
		p.Reason = "untyped resources are not yet supported by `scaff upgrade`"
	default:
		p.Reason = "could not classify the file as a typed or untyped resource"
	}
	return p
}

// Upgrade computes and applies the requested source edits, returning the new
// source. It never writes to disk. changed is false when nothing needed doing.
func (r *Resource) Upgrade(opts UpgradeOptions) (newSrc []byte, changed bool, err error) {
	if r.Kind != KindTyped {
		return nil, false, fmt.Errorf("`scaff upgrade` currently supports typed resources only (got %s)", r.Kind)
	}

	e := newEditor(r.src)
	withIdentity := r.HasIdentity || opts.AddIdentity

	// Extract the flatten method first so a freshly-created flatten already
	// carries the identity write and the identity step below need not touch it.
	if opts.ExtractFlatten && !r.HasFlatten {
		if err := r.extractFlattenEdits(e, opts.ReadModel, withIdentity); err != nil {
			return nil, false, fmt.Errorf("extracting flatten: %w", err)
		}
	}

	if opts.AddIdentity && !r.HasIdentity {
		r.addIdentityAssertionEdit(e)
		r.addIdentityMethodEdit(e)
		if err := r.addCreateIdentityEdit(e); err != nil {
			return nil, false, fmt.Errorf("adding identity to Create: %w", err)
		}
		// Only touch an existing flatten; a newly-extracted one already has it.
		if r.HasFlatten {
			if err := r.injectIdentityIntoExistingFlattenEdit(e); err != nil {
				return nil, false, fmt.Errorf("adding identity to flatten: %w", err)
			}
		}
		props := opts.IdentityProperties
		if props == "" {
			props = "name,resource_group_name"
		}
		if opts.ResourceName != "" {
			r.addGoGenerateEdit(e, opts.ResourceName, props)
		}
	}

	if !e.hasEdits() {
		return r.src, false, nil
	}
	out, err := e.bytes()
	if err != nil {
		return nil, false, err
	}
	return out, true, nil
}
