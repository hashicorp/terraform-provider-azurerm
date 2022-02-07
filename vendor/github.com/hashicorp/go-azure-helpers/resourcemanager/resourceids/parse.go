package resourceids

import (
	"fmt"
	"regexp"
	"strings"
)

type Parser struct {
	// segments is a slice containing the expected (ordered) segments which
	// should be present within this Resource ID
	segments []Segment
}

// NewParser takes a slice of Segments expected for this Resource ID
func NewParser(segments []Segment) Parser {
	return Parser{
		segments: segments,
	}
}

// NewParserFromResourceIdType takes a ResourceId interface and uses its (ordered) Segments
// to create a Parser which can be used to Parse Resource ID's.
func NewParserFromResourceIdType(id ResourceId) Parser {
	segments := id.Segments()
	return NewParser(segments)
}

type ParseResult struct {
	// Parsed is a map of segmentName : segmentValue
	Parsed map[string]string
}

// Parse processes a Resource ID and parses it into a ParseResult containing a map of the
// Known Segments for this Resource ID which callers of this method can then process to
// form a Resource ID struct of those values doing any type conversions as necessary (for
//	example, type-casting/converting Constants).
//
// `input`: the Resource ID to be parsed, which should match the segments for this Resource ID
// `insensitively`: should this Resource ID be parsed case-insensitively and fix up any Constant,
//					Resource Provider and Static Segments to the expected casing.
func (p Parser) Parse(input string, insensitively bool) (*ParseResult, error) {
	if input == "" {
		return nil, fmt.Errorf("cannot parse an empty string")
	}
	if len(p.segments) == 0 {
		return nil, fmt.Errorf("no segments were defined to be able to parse the Resource ID %q", input)
	}

	// if the entire Resource ID is a Scope
	if len(p.segments) == 1 && p.segments[0].Type == ScopeSegmentType {
		return &ParseResult{
			Parsed: map[string]string{
				p.segments[0].Name: input,
			},
		}, nil
	}

	parsed := make(map[string]string)

	hasScopeAtStart := p.segments[0].Type == ScopeSegmentType
	hasScopeAtEnd := p.segments[len(p.segments)-1].Type == ScopeSegmentType

	// go through and build up a regex which will count for the `middle` components of the Resource ID
	nonScopeComponentsRegex := ""
	for i, segment := range p.segments {
		if (i == 0 && hasScopeAtStart) || (i == len(p.segments)-1 && hasScopeAtEnd) {
			continue
		}

		switch segment.Type {
		case ConstantSegmentType:
			{
				if segment.PossibleValues == nil {
					return nil, fmt.Errorf("internal error: constant segment %q had no possible values", segment.Name)
				}

				// e.g. `/(First|Second|Third)`
				nonScopeComponentsRegex += fmt.Sprintf("/(%s)", strings.Join(*segment.PossibleValues, "|"))
				continue
			}

		case ScopeSegmentType:
			{
				return nil, fmt.Errorf("internal error: segment %q is a scope within the middle of a Resource ID which is not supported", segment.Name)
			}

		case ResourceProviderSegmentType, StaticSegmentType:
			{
				if segment.FixedValue == nil {
					return nil, fmt.Errorf("internal error: segment %q is a static/RP without a fixed value", segment.Name)
				}
				nonScopeComponentsRegex += fmt.Sprintf("/%s", *segment.FixedValue)
				continue
			}

		case ResourceGroupSegmentType, SubscriptionIdSegmentType, UserSpecifiedSegmentType:
			{
				nonScopeComponentsRegex += "/(.){1,}"
				continue
			}
		}
	}

	var scopePrefix string
	if hasScopeAtStart {
		prefix, err := p.parseScopePrefix(input, nonScopeComponentsRegex, insensitively)
		if err != nil {
			return nil, fmt.Errorf("parsing scope prefix: %+v", err)
		}

		scopePrefix = *prefix
		parsed[p.segments[0].Name] = *prefix
	}

	// trim off the scopePrefix and the leading `/` to give us the segments we expect plus the final scope string
	// at the end, if present
	uri := input
	if hasScopeAtStart {
		uri = strings.TrimPrefix(uri, scopePrefix)
		uri = strings.TrimPrefix(uri, "/")

		// add a fake start so that the indexes match when we loop around, else we're updating the index below
		uri = fmt.Sprintf("fakestart/%s", uri)
	}

	uri = strings.TrimPrefix(uri, "/")
	split := strings.Split(uri, "/")
	segmentCount := len(split)
	if segmentCount < len(p.segments) {
		return nil, fmt.Errorf("expected %d segments within the Resource ID but got %d for %q", len(p.segments), segmentCount, input)
	}

	if hasScopeAtStart {
		// trim off the fake start since we use any remaining uri as a suffixScope
		uri = strings.TrimPrefix(uri, "fakestart/")
	}

	for i, segment := range p.segments {
		if (i == 0 && hasScopeAtStart) || (i == len(p.segments)-1 && hasScopeAtEnd) {
			continue
		}

		// as we go around each of the segments we're expecting, process the value we should surface
		rawSegment := split[i]
		value, err := p.parseSegment(segment, rawSegment, insensitively)
		if err != nil {
			return nil, fmt.Errorf("parsing segment %q: %+v", segment.Name, err)
		}
		parsed[segment.Name] = *value

		// and then remove rawSegment from `uri` so that any leftovers is the scope
		// since if there's a scope there'll be more segments than we expect
		uri = strings.TrimPrefix(uri, fmt.Sprintf("%s", rawSegment))
		uri = strings.TrimPrefix(uri, "/")
	}

	if uri != "" {
		if !hasScopeAtEnd {
			return nil, fmt.Errorf("unexpected segment %q present at the end of the URI (input %q)", uri, input)
		}

		parsed[p.segments[len(p.segments)-1].Name] = fmt.Sprintf("/%s", uri)
	}

	if len(p.segments) != len(parsed) {
		return nil, fmt.Errorf("expected %d segments but got %d for %q", len(p.segments), len(parsed), input)
	}

	for k, v := range parsed {
		if v == "" {
			return nil, fmt.Errorf("segment %q is required but got an empty value", k)
		}
	}

	return &ParseResult{
		Parsed: parsed,
	}, nil
}

func (p Parser) parseScopePrefix(input, regexForNonScopeSegments string, insensitively bool) (*string, error) {
	regexToUse := fmt.Sprintf("^((.){1,})%s", regexForNonScopeSegments)
	if insensitively {
		regexToUse = fmt.Sprintf("(?i)%s", regexToUse)
	}
	r, err := regexp.Compile(regexToUse)
	if err != nil {
		return nil, fmt.Errorf("internal error: compiling regex %q to find scope prefix: %+v", regexToUse, err)
	}
	// 0 is the entire string, 1 will be the scope prefix, we can ignore the rest
	values := r.FindStringSubmatch(input)
	if len(values) < 2 {
		return nil, fmt.Errorf("unable to find the scope prefix from the value %q with the regex %q", input, regexToUse)
	}
	v := values[1]
	if v == "" {
		return nil, fmt.Errorf("unable to find the scope prefix from the value %q using the regex %q", input, regexToUse)
	}
	return &v, nil
}

func (p Parser) parseSegment(segment Segment, rawValue string, insensitively bool) (*string, error) {
	switch segment.Type {
	case ConstantSegmentType:
		{
			if segment.PossibleValues == nil {
				return nil, fmt.Errorf("internal error: missing PossibleValues for Constant segment %q", segment.Name)
			}
			for _, possibleVal := range *segment.PossibleValues {
				matches := possibleVal == rawValue
				if insensitively {
					matches = strings.EqualFold(possibleVal, rawValue)
				}

				if matches {
					return &possibleVal, nil
				}
			}

			return nil, fmt.Errorf("expected the segment %q to match one of the values %q but got %q", segment.Name, strings.Join(*segment.PossibleValues, ", "), rawValue)
		}

	case ResourceProviderSegmentType, StaticSegmentType:
		{
			if segment.FixedValue == nil {
				return nil, fmt.Errorf("internal error: segment %q is a static/RP segment without a fixed value", segment.Name)
			}

			matches := *segment.FixedValue == rawValue
			if insensitively {
				matches = strings.EqualFold(*segment.FixedValue, rawValue)
			}

			if matches {
				return &*segment.FixedValue, nil
			}

			return nil, fmt.Errorf("expected the segment %q to be %q", rawValue, *segment.FixedValue)
		}

	case ScopeSegmentType:
		{
			return nil, fmt.Errorf("internal error: scope segments aren't supported unless at the start or the end")
		}

	case ResourceGroupSegmentType, SubscriptionIdSegmentType, UserSpecifiedSegmentType:
		{
			return &rawValue, nil
		}
	}

	return nil, fmt.Errorf("internal error: missing parser for segment %q (type %q)", segment.Name, string(segment.Type))
}
