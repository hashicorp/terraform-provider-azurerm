package client

import (
	"errors"
	"net/http"
	"net/url"
	"strings"
)

const (
	// VCRReplayHeader is set by VCR-aware transports on replayed HTTP responses.
	VCRReplayHeader = "X-Go-Azure-SDK-VCR-Replay"

	SkipPollingDelayHeader = "X-Go-Azure-SDK-Skip-Polling-Delay"

	VCRInteractionNotFoundErrMsg = "requested interaction not found"
)

// The VCR recorder sets the X-Go-Azure-SDK-VCR-Replay header to "true" when a response is returned from a cassette.
// In ModeReplayWithNewEpisodes, VCR may still make live HTTP requests, so we check for the presence of this header
// instead of relying only on the VCR mode.
func IsVCRRecordedResponse(resp *http.Response) bool {
	return resp != nil && resp.Header != nil &&
		strings.EqualFold(resp.Header.Get(VCRReplayHeader), "true")
}

// IsVCRReplayMissError returns true when the error indicates that a replayed cassette has no matching interaction.
func IsVCRReplayMissError(err error) bool {
	var urlErr *url.Error
	if err == nil || !errors.As(err, &urlErr) {
		return false
	}
	return strings.Contains(strings.ToLower(urlErr.Error()), VCRInteractionNotFoundErrMsg)
}

// Deprecated: Use IsVCRReplayMissError instead. This is a temporary fallback, as pollers wrap errors using %+v,
// which loses the original error type information. Once pollers are updated to use %w (which preserves type information),
// this function will be removed.
func IsVCRReplayMissErrorDeprecated(err error) bool {
	return err != nil && strings.Contains(strings.ToLower(err.Error()), VCRInteractionNotFoundErrMsg)
}

func IsVCRReplaying(c *Client) bool {
	return c != nil && c.TransportMode == TransportModeVCRReplay
}
