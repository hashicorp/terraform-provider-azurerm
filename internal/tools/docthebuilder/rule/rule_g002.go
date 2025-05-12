package rule

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/docthebuilder/data"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/docthebuilder/differror"
)

type G002 struct{}

type note struct {
	level  string
	prefix string
	text   string
}

var (
	_ Rule = G002{}

	// While `>` does not template as a Note block in the Registry docs, it's a common error
	// as `>` is the markdown indicator for a note, hence why we look for this as well so we can fix it
	fullNoteRegex    = regexp.MustCompile(`^[\s\t]*(>|->|~>|!>)[\s|\t]*(\*\*[\w\s:]*\*\*)(.*)`)
	partialNoteRegex = regexp.MustCompile(`^[\s\t]*(>|->|~>|!>)(.*)`)
)

func (r G002) Name() string {
	return "G002"
}

func (r G002) Description() string {
	return fmt.Sprintf("%s - validates notes in documentation are following the expected format", r.Name())
}

func (r G002) Run(data *data.ResourceData, fix bool) []error {
	errs := make([]error, 0)

	for _, section := range data.Document.Sections {
		content := section.GetContent()
		for idx, line := range content {
			if partialNoteRegex.MatchString(line) {
				n := parseNote(line)

				if n == nil {
					errs = append(errs, fmt.Errorf("%s - Unable to parse note: `%s`", r.Name(), line))
					continue
				}

				// If we encounter a markdown note, default to an informational note.
				if n.level == ">" {
					n.level = "->"
				}

				current := line
				expected := note{
					level:  n.level,
					prefix: "Note:",
					text:   n.text,
				}.string()

				if current != expected {
					errs = append(errs, differror.New(fmt.Sprintf("%s - Note not in expected format", r.Name()), current, expected))

					if fix {
						data.Document.HasChange = true
						content[idx] = expected
						section.SetContent(content)
					}
				}
			}
		}
	}

	return errs
}

func parseNote(line string) *note {
	if fullNoteRegex.MatchString(line) {
		groups := fullNoteRegex.FindStringSubmatch(line)

		if len(groups) == 4 {
			return &note{
				level:  groups[1],
				prefix: strings.ReplaceAll(groups[2], "*", ""),
				text:   normalizeNoteText(groups[3]),
			}
		}
	}

	if partialNoteRegex.MatchString(line) {
		// if full regex is not matched, but it starts with `(>|->|~>|!>)` it's likely the `**Note:**` prefix is missing
		// try to parse the text and level, we'll build a correctly formatted note. Best effort, not perfect.
		groups := partialNoteRegex.FindStringSubmatch(line)
		if len(groups) == 3 {
			return &note{
				level: groups[1],
				text:  normalizeNoteText(groups[2]),
			}
		}
	}

	return nil
}

func normalizeNoteText(text string) string {
	// may have to expand on this regex in the future
	regex := regexp.MustCompile(`(?i)(^(note)?(:)?)`)

	return strings.TrimSpace(regex.ReplaceAllString(strings.TrimSpace(text), ""))
}

func (n note) string() string {
	return fmt.Sprintf("%s **%s** %s", n.level, n.prefix, n.text)
}
