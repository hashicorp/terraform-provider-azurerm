package accesscontrol

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestACEValidate_WithValidACE(t *testing.T) {
	ace := ACE{
		TagType:     TagTypeUser,
		Permissions: "rwx",
	}
	assert.Nil(t, ace.Validate(), "Expected ACE to Validate successfully")
}

func TestACEValidate_WithACEInvalidPermisions(t *testing.T) {
	ace := ACE{
		TagType:     TagTypeUser,
		Permissions: "awx",
	}

	assert.EqualError(t, ace.Validate(), "Permissions must be of the form [r-][w-][x-]")
}

func TestACEValidate_WithACEQualifierSetForMask(t *testing.T) {
	qualifier := uuid.MustParse("ba4662cb-995c-479d-8f7c-a1b3d8ae05a9")
	ace := ACE{
		TagType:      TagTypeMask,
		TagQualifier: &qualifier,
		Permissions:  "rwx",
	}

	assert.EqualError(t, ace.Validate(), "TagQualifier cannot be set for 'mask' or 'other' TagTypes")
}

func TestACEParse_WithValidDefault(t *testing.T) {
	ace, err := ParseACE("default:user:22edd3d8-9253-4463-b8f8-442ebe33b622:rwx")
	assert.NoError(t, err)
	assert.Equal(t, true, ace.IsDefault, "Expected default")
	assert.Equal(t, TagTypeUser, ace.TagType)
	assert.Equal(t, "22edd3d8-9253-4463-b8f8-442ebe33b622", ace.TagQualifier.String())
	assert.Equal(t, "rwx", ace.Permissions)
}
func TestACEParse_WithValidNonDefault(t *testing.T) {
	ace, err := ParseACE("user:885d0d94-9ecb-4e0d-8581-781b56d27b10:rwx")
	assert.NoError(t, err)
	assert.Equal(t, false, ace.IsDefault, "Expected non-default")
	assert.Equal(t, TagTypeUser, ace.TagType)
	assert.Equal(t, "885d0d94-9ecb-4e0d-8581-781b56d27b10", ace.TagQualifier.String())
	assert.Equal(t, "rwx", ace.Permissions)
}

func TestACEParse_WithInvalid4Part(t *testing.T) {
	_, err := ParseACE("test:::")
	assert.EqualError(t, err, "When specifying a 4-part ACE the first part must be 'default'")
}
func TestACEParse_WithInvalidTagType(t *testing.T) {
	_, err := ParseACE("wibble::rwx")
	assert.EqualError(t, err, "Unrecognized TagType: \"wibble\"")
}
func TestACEParse_WithInvalidNumberOfParts(t *testing.T) {
	_, err := ParseACE("user:rwx")
	assert.EqualError(t, err, "ACE string should have either 3 or 4 parts")
}
func TestACEParse_WithInvalidQualifier(t *testing.T) {
	_, err := ParseACE("user:wibble94-9ecb-4e0d-8581-781b56d27b10:rwx")
	assert.EqualError(t, err, "Error parsing qualifer \"wibble94-9ecb-4e0d-8581-781b56d27b10\": invalid UUID format")
}

func TestACEString_WithDefault(t *testing.T) {
	qualifier := uuid.MustParse("ba4662cb-995c-479d-8f7c-a1b3d8ae05a9")
	ace := ACE{
		IsDefault:    true,
		TagType:      TagTypeUser,
		TagQualifier: &qualifier,
		Permissions:  "r-x",
	}
	assert.Equal(t, "default:user:ba4662cb-995c-479d-8f7c-a1b3d8ae05a9:r-x", ace.String())
}
func TestACEString_WithNonDefault(t *testing.T) {
	qualifier := uuid.MustParse("ba4662cb-995c-479d-8f7c-a1b3d8ae05a9")
	ace := ACE{
		IsDefault:    false,
		TagType:      TagTypeGroup,
		TagQualifier: &qualifier,
		Permissions:  "rw-",
	}
	assert.Equal(t, "group:ba4662cb-995c-479d-8f7c-a1b3d8ae05a9:rw-", ace.String())
}
