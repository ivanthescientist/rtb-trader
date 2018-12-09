package bidrequest

import (
	"errors"
	"github.com/bsm/openrtb"
	"github.com/mssola/user_agent"
	"net/url"
	"strings"
)

var deviceTypeMapping = map[int]string{
	openrtb.DeviceTypeMobile: DeviceTypeMobile,
	openrtb.DeviceTypeTablet: DeviceTypeMobile,
	openrtb.DeviceTypePhone:  DeviceTypeMobile,

	openrtb.DeviceTypePC: DeviceTypeDesktop,
	openrtb.DeviceTypeTV: DeviceTypeDesktop,
}

var (
	ErrDeviceUnknown = errors.New("device field is not present")
)

func ExtractFields(request *openrtb.BidRequest) (FieldMap, error) {
	var fields = make(FieldMap)
	if request.Device == nil {
		return nil, ErrDeviceUnknown
	}

	ua := user_agent.New(request.Device.UA)

	extractOSName(fields, request, ua)
	extractOSVersion(fields, request, ua)
	extractDeviceType(fields, request, ua)
	extractBrowserName(fields, ua)
	extractBrowserType(fields, request, ua)
	extractIP(fields, request)
	extractCountry(fields, request)
	extractSiteDomain(fields, request)

	return fields, nil
}

func extractOSName(fields FieldMap, request *openrtb.BidRequest, ua *user_agent.UserAgent) {
	if request.Device.OS != "" {
		fields[FieldOSName] = request.Device.OS
		return
	}

	if ua != nil && ua.OSInfo().Name != "" {
		fields[FieldOSName] = ua.OSInfo().Name
	}
}

func extractOSVersion(fields FieldMap, request *openrtb.BidRequest, ua *user_agent.UserAgent) {
	if request.Device.OSVer != "" {
		fields[FieldOSVersion] = request.Device.OSVer
	}

	if ua != nil && ua.OSInfo().Version != "" {
		fields[FieldOSVersion] = ua.OSInfo().Version
	}
}

func extractDeviceType(fields FieldMap, request *openrtb.BidRequest, ua *user_agent.UserAgent) {
	if deviceType, isPresent := deviceTypeMapping[request.Device.DeviceType]; isPresent {
		fields[FieldDeviceType] = deviceType
		return
	}

	if ua == nil {
		return
	}

	fields[FieldDeviceType] = DeviceTypeDesktop

	if ua.Mobile() {
		fields[FieldDeviceType] = DeviceTypeMobile
	}
}

func extractBrowserName(fields FieldMap, ua *user_agent.UserAgent) {
	if ua == nil {
		return
	}

	fields[FieldBrowserName], _ = ua.Browser()
}

func extractBrowserType(fields FieldMap, request *openrtb.BidRequest, ua *user_agent.UserAgent) {
	var mobileOSSet = map[string]struct{}{
		"android": {},
		"ios":     {},
		"tizen":   {},
	}

	var desktopOSSet = map[string]struct{}{
		"windows": {},
		"macosx":  {},
		"linux":   {},
	}

	// App is only present in mobile app' web views.
	if request.App != nil {
		fields[FieldBrowserType] = BrowserTypeInApp
		return
	}

	_, isMobileOS := mobileOSSet[strings.ToLower(request.Device.OS)]
	if isMobileOS {
		fields[FieldBrowserType] = BrowserTypeMobile
		return
	}

	_, isDesktopOS := desktopOSSet[strings.ToLower(request.Device.OS)]
	if isDesktopOS {
		fields[FieldBrowserType] = BrowserTypeDesktop
		return
	}

	// If UA is not present, then we can't do anything about it
	if ua == nil {
		return
	}

	if !ua.Mobile() {
		fields[FieldBrowserType] = BrowserTypeDesktop
		return
	}

	fields[FieldBrowserType] = BrowserTypeMobile
	// "wv" is a special user-agent marker for Android to signal that its a WebView
	if strings.ToLower(ua.OSInfo().Name) == "android" && strings.Contains(ua.UA(), "wv") {
		fields[FieldBrowserType] = BrowserTypeInApp
	}
}

func extractIP(fields FieldMap, request *openrtb.BidRequest) {
	if request.Device.IP == "" {
		return
	}

	fields[FieldIP] = request.Device.IP
}

func extractCountry(fields FieldMap, request *openrtb.BidRequest) {
	if request.Device.Geo != nil && request.Device.Geo.Country != "" {
		fields[FieldCountry] = request.Device.Geo.Country
	}

	if request.Device.IP == "" {
		return
	}

	if country := ipToCountry(request.Device.IP); country != "" {
		fields[FieldCountry] = country
	}
}

// Please imagine that this searches an in-memory IP-to-country prefix tree of sorts
func ipToCountry(_ string) string {
	return "World"
}

func extractSiteDomain(fields map[string]string, request *openrtb.BidRequest) {
	if request.Site == nil {
		return
	}

	if request.Site.Domain != "" {
		fields[FieldSiteDomain] = request.Site.Domain
		return
	}

	if request.Site.Page == "" {
		return
	}

	pageUrl, err := url.Parse(request.Site.Page)
	if err != nil {
		return
	}

	fields[FieldSiteDomain] = pageUrl.Hostname()
}
