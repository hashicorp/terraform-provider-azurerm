package runbookdraft

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HTTPStatusCode string

const (
	HTTPStatusCodeAccepted                     HTTPStatusCode = "Accepted"
	HTTPStatusCodeAmbiguous                    HTTPStatusCode = "Ambiguous"
	HTTPStatusCodeBadGateway                   HTTPStatusCode = "BadGateway"
	HTTPStatusCodeBadRequest                   HTTPStatusCode = "BadRequest"
	HTTPStatusCodeConflict                     HTTPStatusCode = "Conflict"
	HTTPStatusCodeContinue                     HTTPStatusCode = "Continue"
	HTTPStatusCodeCreated                      HTTPStatusCode = "Created"
	HTTPStatusCodeExpectationFailed            HTTPStatusCode = "ExpectationFailed"
	HTTPStatusCodeForbidden                    HTTPStatusCode = "Forbidden"
	HTTPStatusCodeFound                        HTTPStatusCode = "Found"
	HTTPStatusCodeGatewayTimeout               HTTPStatusCode = "GatewayTimeout"
	HTTPStatusCodeGone                         HTTPStatusCode = "Gone"
	HTTPStatusCodeHTTPVersionNotSupported      HTTPStatusCode = "HttpVersionNotSupported"
	HTTPStatusCodeInternalServerError          HTTPStatusCode = "InternalServerError"
	HTTPStatusCodeLengthRequired               HTTPStatusCode = "LengthRequired"
	HTTPStatusCodeMethodNotAllowed             HTTPStatusCode = "MethodNotAllowed"
	HTTPStatusCodeMoved                        HTTPStatusCode = "Moved"
	HTTPStatusCodeMovedPermanently             HTTPStatusCode = "MovedPermanently"
	HTTPStatusCodeMultipleChoices              HTTPStatusCode = "MultipleChoices"
	HTTPStatusCodeNoContent                    HTTPStatusCode = "NoContent"
	HTTPStatusCodeNonAuthoritativeInformation  HTTPStatusCode = "NonAuthoritativeInformation"
	HTTPStatusCodeNotAcceptable                HTTPStatusCode = "NotAcceptable"
	HTTPStatusCodeNotFound                     HTTPStatusCode = "NotFound"
	HTTPStatusCodeNotImplemented               HTTPStatusCode = "NotImplemented"
	HTTPStatusCodeNotModified                  HTTPStatusCode = "NotModified"
	HTTPStatusCodeOK                           HTTPStatusCode = "OK"
	HTTPStatusCodePartialContent               HTTPStatusCode = "PartialContent"
	HTTPStatusCodePaymentRequired              HTTPStatusCode = "PaymentRequired"
	HTTPStatusCodePreconditionFailed           HTTPStatusCode = "PreconditionFailed"
	HTTPStatusCodeProxyAuthenticationRequired  HTTPStatusCode = "ProxyAuthenticationRequired"
	HTTPStatusCodeRedirect                     HTTPStatusCode = "Redirect"
	HTTPStatusCodeRedirectKeepVerb             HTTPStatusCode = "RedirectKeepVerb"
	HTTPStatusCodeRedirectMethod               HTTPStatusCode = "RedirectMethod"
	HTTPStatusCodeRequestEntityTooLarge        HTTPStatusCode = "RequestEntityTooLarge"
	HTTPStatusCodeRequestTimeout               HTTPStatusCode = "RequestTimeout"
	HTTPStatusCodeRequestUriTooLong            HTTPStatusCode = "RequestUriTooLong"
	HTTPStatusCodeRequestedRangeNotSatisfiable HTTPStatusCode = "RequestedRangeNotSatisfiable"
	HTTPStatusCodeResetContent                 HTTPStatusCode = "ResetContent"
	HTTPStatusCodeSeeOther                     HTTPStatusCode = "SeeOther"
	HTTPStatusCodeServiceUnavailable           HTTPStatusCode = "ServiceUnavailable"
	HTTPStatusCodeSwitchingProtocols           HTTPStatusCode = "SwitchingProtocols"
	HTTPStatusCodeTemporaryRedirect            HTTPStatusCode = "TemporaryRedirect"
	HTTPStatusCodeUnauthorized                 HTTPStatusCode = "Unauthorized"
	HTTPStatusCodeUnsupportedMediaType         HTTPStatusCode = "UnsupportedMediaType"
	HTTPStatusCodeUnused                       HTTPStatusCode = "Unused"
	HTTPStatusCodeUpgradeRequired              HTTPStatusCode = "UpgradeRequired"
	HTTPStatusCodeUseProxy                     HTTPStatusCode = "UseProxy"
)

func PossibleValuesForHTTPStatusCode() []string {
	return []string{
		string(HTTPStatusCodeAccepted),
		string(HTTPStatusCodeAmbiguous),
		string(HTTPStatusCodeBadGateway),
		string(HTTPStatusCodeBadRequest),
		string(HTTPStatusCodeConflict),
		string(HTTPStatusCodeContinue),
		string(HTTPStatusCodeCreated),
		string(HTTPStatusCodeExpectationFailed),
		string(HTTPStatusCodeForbidden),
		string(HTTPStatusCodeFound),
		string(HTTPStatusCodeGatewayTimeout),
		string(HTTPStatusCodeGone),
		string(HTTPStatusCodeHTTPVersionNotSupported),
		string(HTTPStatusCodeInternalServerError),
		string(HTTPStatusCodeLengthRequired),
		string(HTTPStatusCodeMethodNotAllowed),
		string(HTTPStatusCodeMoved),
		string(HTTPStatusCodeMovedPermanently),
		string(HTTPStatusCodeMultipleChoices),
		string(HTTPStatusCodeNoContent),
		string(HTTPStatusCodeNonAuthoritativeInformation),
		string(HTTPStatusCodeNotAcceptable),
		string(HTTPStatusCodeNotFound),
		string(HTTPStatusCodeNotImplemented),
		string(HTTPStatusCodeNotModified),
		string(HTTPStatusCodeOK),
		string(HTTPStatusCodePartialContent),
		string(HTTPStatusCodePaymentRequired),
		string(HTTPStatusCodePreconditionFailed),
		string(HTTPStatusCodeProxyAuthenticationRequired),
		string(HTTPStatusCodeRedirect),
		string(HTTPStatusCodeRedirectKeepVerb),
		string(HTTPStatusCodeRedirectMethod),
		string(HTTPStatusCodeRequestEntityTooLarge),
		string(HTTPStatusCodeRequestTimeout),
		string(HTTPStatusCodeRequestUriTooLong),
		string(HTTPStatusCodeRequestedRangeNotSatisfiable),
		string(HTTPStatusCodeResetContent),
		string(HTTPStatusCodeSeeOther),
		string(HTTPStatusCodeServiceUnavailable),
		string(HTTPStatusCodeSwitchingProtocols),
		string(HTTPStatusCodeTemporaryRedirect),
		string(HTTPStatusCodeUnauthorized),
		string(HTTPStatusCodeUnsupportedMediaType),
		string(HTTPStatusCodeUnused),
		string(HTTPStatusCodeUpgradeRequired),
		string(HTTPStatusCodeUseProxy),
	}
}

func (s *HTTPStatusCode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseHTTPStatusCode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseHTTPStatusCode(input string) (*HTTPStatusCode, error) {
	vals := map[string]HTTPStatusCode{
		"accepted":                     HTTPStatusCodeAccepted,
		"ambiguous":                    HTTPStatusCodeAmbiguous,
		"badgateway":                   HTTPStatusCodeBadGateway,
		"badrequest":                   HTTPStatusCodeBadRequest,
		"conflict":                     HTTPStatusCodeConflict,
		"continue":                     HTTPStatusCodeContinue,
		"created":                      HTTPStatusCodeCreated,
		"expectationfailed":            HTTPStatusCodeExpectationFailed,
		"forbidden":                    HTTPStatusCodeForbidden,
		"found":                        HTTPStatusCodeFound,
		"gatewaytimeout":               HTTPStatusCodeGatewayTimeout,
		"gone":                         HTTPStatusCodeGone,
		"httpversionnotsupported":      HTTPStatusCodeHTTPVersionNotSupported,
		"internalservererror":          HTTPStatusCodeInternalServerError,
		"lengthrequired":               HTTPStatusCodeLengthRequired,
		"methodnotallowed":             HTTPStatusCodeMethodNotAllowed,
		"moved":                        HTTPStatusCodeMoved,
		"movedpermanently":             HTTPStatusCodeMovedPermanently,
		"multiplechoices":              HTTPStatusCodeMultipleChoices,
		"nocontent":                    HTTPStatusCodeNoContent,
		"nonauthoritativeinformation":  HTTPStatusCodeNonAuthoritativeInformation,
		"notacceptable":                HTTPStatusCodeNotAcceptable,
		"notfound":                     HTTPStatusCodeNotFound,
		"notimplemented":               HTTPStatusCodeNotImplemented,
		"notmodified":                  HTTPStatusCodeNotModified,
		"ok":                           HTTPStatusCodeOK,
		"partialcontent":               HTTPStatusCodePartialContent,
		"paymentrequired":              HTTPStatusCodePaymentRequired,
		"preconditionfailed":           HTTPStatusCodePreconditionFailed,
		"proxyauthenticationrequired":  HTTPStatusCodeProxyAuthenticationRequired,
		"redirect":                     HTTPStatusCodeRedirect,
		"redirectkeepverb":             HTTPStatusCodeRedirectKeepVerb,
		"redirectmethod":               HTTPStatusCodeRedirectMethod,
		"requestentitytoolarge":        HTTPStatusCodeRequestEntityTooLarge,
		"requesttimeout":               HTTPStatusCodeRequestTimeout,
		"requesturitoolong":            HTTPStatusCodeRequestUriTooLong,
		"requestedrangenotsatisfiable": HTTPStatusCodeRequestedRangeNotSatisfiable,
		"resetcontent":                 HTTPStatusCodeResetContent,
		"seeother":                     HTTPStatusCodeSeeOther,
		"serviceunavailable":           HTTPStatusCodeServiceUnavailable,
		"switchingprotocols":           HTTPStatusCodeSwitchingProtocols,
		"temporaryredirect":            HTTPStatusCodeTemporaryRedirect,
		"unauthorized":                 HTTPStatusCodeUnauthorized,
		"unsupportedmediatype":         HTTPStatusCodeUnsupportedMediaType,
		"unused":                       HTTPStatusCodeUnused,
		"upgraderequired":              HTTPStatusCodeUpgradeRequired,
		"useproxy":                     HTTPStatusCodeUseProxy,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := HTTPStatusCode(input)
	return &out, nil
}
