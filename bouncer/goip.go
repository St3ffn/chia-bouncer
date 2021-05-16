package bouncer

import "strings"

// Lookup interface to perform basic lookups
type Lookup interface {
	// IpLocation to perform a lookup based on the Geo IP  Location
	IpLocation() (string, error)
}

// GetLocation to get the Geo Ip Location for a ipV4 Address
type GetLocation func(ipV4 string) (string, error)

// geoIpLookup performs the Geo Ip Location lookup for the ipV4 Address via the geoiplookup cmd tool
// geoiplookup is expected to be in the current PATH for execution
func geoIpLookup(ipV4 string) (location string, err error) {
	out, err := execCmd("geoiplookup", ipV4)
	if err != nil {
		return "", err
	}
	// like "GeoIP Country Edition: KR, Korea, Republic of\n   "
	withoutSuffix := strings.TrimPrefix(string(out[:]), "GeoIP Country Edition: ")
	return strings.TrimSuffix(withoutSuffix, "\n"), nil
}
