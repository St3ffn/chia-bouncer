package bouncer

import (
	"errors"
	"reflect"
	"testing"
)

type testData struct {
	ip             string
	lookupLocation string // lookup ok
}

func (t testData) node() FullNode {
	return FullNode{
		Ip: t.ip,
	}
}

var marsOne = testData{
	ip:             "1.1.1.1",
	lookupLocation: "Elon on Mars",
}
var marsTwo = testData{
	ip:             "1.1.1.2",
	lookupLocation: "Mars",
}
var korea = testData{
	ip:             "2.1.1.2",
	lookupLocation: "KR, Korea, Republic of",
}
var netherlands = testData{
	ip:             "3.1.1.2",
	lookupLocation: "NL, Netherlands",
}

var testDataSet = []testData{marsOne, marsTwo, korea, netherlands}

func lookupOk(ipV4 string) (_ string, _ error) {
	for _, data := range testDataSet {
		if data.ip == ipV4 {
			return data.lookupLocation, nil
		}
	}
	panic("test is broken")
}

func lookupError(_ string) (_ string, _ error) {
	return "", errors.New("lookup error")
}

func TestFilterByLocation_Filter(t *testing.T) {
	tests := []struct {
		name             string
		getLocation      GetLocation
		nodes            []FullNode
		locationToFilter string
		want             []FullNode
		wantErr          bool
	}{
		{
			name:        "filter two mars",
			getLocation: lookupOk,
			nodes: []FullNode{
				korea.node(),
				marsOne.node(),
				netherlands.node(),
				marsTwo.node(),
			},
			locationToFilter: "mars", // case insensitive
			want: []FullNode{
				marsOne.node(),
				marsTwo.node(),
			},
		},
		{
			name:        "no match",
			getLocation: lookupOk,
			nodes: []FullNode{
				korea.node(),
				marsOne.node(),
				netherlands.node(),
				marsTwo.node(),
			},
			locationToFilter: "pluto", // case insensitive
			want:             nil,
		},
		{
			name:             "filter broken",
			getLocation:      lookupError,
			nodes:            []FullNode{korea.node(), marsOne.node(), netherlands.node()},
			locationToFilter: "mars",
			wantErr:          true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getLocation = tt.getLocation
			f := FilterByLocation{
				Nodes:            tt.nodes,
				LocationToFilter: tt.locationToFilter,
			}
			got, err := f.Filter()
			if (err != nil) != tt.wantErr {
				t.Errorf("Filter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Filter() got = %v, want %v", got, tt.want)
			}
		})
	}
}
