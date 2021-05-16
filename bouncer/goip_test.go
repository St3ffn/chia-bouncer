package bouncer

import (
	"errors"
	"testing"
)

func cmdSuccessKorea(_ string, _ ...string) (out []byte, err error) {
	return []byte("GeoIP Country Edition: KR, Korea, Republic of"), nil
}

func cmdFailed(_ string, _ ...string) (out []byte, err error) {
	return nil, errors.New("something went wrong")
}

func Test_geoIpLookup(t *testing.T) {
	tests := []struct {
		name         string
		ip           string
		execCmd      ExecCmd
		wantLocation string
		wantErr      error
	}{
		{name: "ok", ip: "1.2.3.4", execCmd: cmdSuccessKorea, wantLocation: "KR, Korea, Republic of"},
		{name: "nok", ip: "1.2.3.4", execCmd: cmdFailed, wantErr: errors.New("something went wrong")},
	}
	for _, tt := range tests {
		execCmd = tt.execCmd
		t.Run(tt.name, func(t *testing.T) {
			gotLocation, err := geoIpLookup(tt.ip)
			if err != nil {
				if tt.wantErr == nil {
					t.Errorf("geoIpLookup() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				if err.Error() != tt.wantErr.Error() {
					t.Errorf("geoIpLookup() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}
			if gotLocation != tt.wantLocation {
				t.Errorf("geoIpLookup() gotLocation = %v, want %v", gotLocation, tt.wantLocation)
			}
		})
	}
}
