package accesscontrol

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

type TagType string

const (
	TagTypeUser  TagType = "user"
	TagTypeGroup TagType = "group"
	TagTypeMask  TagType = "mask"
	TagTypeOther TagType = "other"
)

// ACE is roughly modelled on https://linux.die.net/man/5/acl
// Have collapsed ACL_USER_OBJ to be TagType of user with a nil Qualifier

type ACE struct {
	IsDefault    bool
	TagType      TagType
	TagQualifier *uuid.UUID
	Permissions  string // TODO break the rwx into permission flags?
}

var permissionsRegex *regexp.Regexp

func init() {
	permissionsRegex = regexp.MustCompile("[r-][w-][x-]")
}

// ValidateACEPermissions checks the format of the ACE permission string. Returns nil on success.
func ValidateACEPermissions(permissions string) error {
	if !permissionsRegex.MatchString(permissions) {
		return fmt.Errorf("Permissions must be of the form [r-][w-][x-]")
	}
	return nil
}

// Validate checks the formatting and combination of values in the ACE. Returns nil on success
func (ace *ACE) Validate() error {
	switch ace.TagType {
	case TagTypeMask, TagTypeOther:
		if ace.TagQualifier != nil {
			return fmt.Errorf("TagQualifier cannot be set for 'mask' or 'other' TagTypes")
		}
	}

	if err := ValidateACEPermissions(ace.Permissions); err != nil {
		return err
	}

	if err := validateTagType(ace.TagType); err != nil {
		return err
	}

	return nil
}

// ParseACE parses an ACE string and returns the ACE
func ParseACE(input string) (ACE, error) {
	ace := ACE{}

	parts := strings.Split(input, ":")
	if len(parts) == 4 {
		if parts[0] == "default" {
			ace.IsDefault = true
			parts = parts[1:]
		} else {
			return ACE{}, fmt.Errorf("When specifying a 4-part ACE the first part must be 'default'")
		}
	}

	if len(parts) != 3 {
		return ACE{}, fmt.Errorf("ACE string should have either 3 or 4 parts")
	}

	ace.TagType = TagType(parts[0])

	qualiferString := parts[1]
	if qualiferString != "" {
		qualifier, err := uuid.Parse(qualiferString)
		if err != nil {
			return ACE{}, fmt.Errorf("Error parsing qualifer %q: %s", qualiferString, err)
		}
		ace.TagQualifier = &qualifier
	}

	ace.Permissions = parts[2]

	if err := ace.Validate(); err != nil {
		return ACE{}, err
	}
	return ace, nil
}

// String returns the string form of the ACE - this does not check that it is valid
func (ace *ACE) String() string {
	prefix := ""
	if ace.IsDefault {
		prefix = "default:"
	}
	qualifierString := ""
	if ace.TagQualifier != nil {
		qualifierString = ace.TagQualifier.String()
	}
	return fmt.Sprintf("%s%s:%s:%s", prefix, ace.TagType, qualifierString, ace.Permissions)
}

func validateTagType(tagType TagType) error {
	switch tagType {
	case TagTypeUser,
		TagTypeGroup,
		TagTypeMask,
		TagTypeOther:
		return nil
	}
	return fmt.Errorf("Unrecognized TagType: %q", tagType)
}
