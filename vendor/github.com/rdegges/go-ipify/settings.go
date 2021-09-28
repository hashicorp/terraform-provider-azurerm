package ipify

import (
	"fmt"
	"runtime"
	"strings"
)

// The version of this library.
const VERSION = "1.0.0"

// The maximum amount of tries to attempt when making API calls.
const MAX_TRIES = 3

// This is the ipify service base URI.  This is where all API requests go.
var API_URI = "https://api.ipify.org"

// The user-agent string is provided so that I can (eventually) keep track of
// what libraries to support over time.  EG: Maybe the service is used
// primarily by Windows developers, and I should invest more time in improving
// those integrations.
var USER_AGENT = fmt.Sprintf(
	"go-ipify/%s go/%s %s",
	VERSION,
	runtime.Version()[2:],
	strings.Title(runtime.GOOS),
)
