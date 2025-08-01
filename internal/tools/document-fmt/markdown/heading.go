package markdown

import (
	"regexp"
	"strings"
)

type Heading struct {
	Level int
	Text  string
	Raw   string
}

func NewHeading(line string) Heading {
	heading := Heading{}

	r := regexp.MustCompile(`^(#{1,6})(.*)$`)
	matches := r.FindStringSubmatch(line)

	if len(matches) == 3 {
		heading.Raw = matches[0]
		heading.Level = len(matches[1])
		heading.Text = strings.TrimSpace(matches[2])
	}

	return heading
}
