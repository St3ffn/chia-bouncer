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

func TestFilterByDown_Filter(t *testing.T) {
	nodes := []FullNode{
		{Ip: "1.1.1.1", Ports: "61075/8444", NodeId: "bbbbbbbb", LastConnect: "May 14 18:35:01", Up: 12.8, Down: 5.8},
		{Ip: "2.2.2.2", Ports: "8444/8444", NodeId: "cccccccc", LastConnect: "May 14 18:35:01", Up: 3.2, Down: 12.1},
		{Ip: "3.3.3.3", Ports: "8444/8444", NodeId: "dddddddd", LastConnect: "May 14 11:55:26", Up: 0.0, Down: 0.0},
		{Ip: "4.4.4.4", Ports: "8444/8444", NodeId: "wewewewe", LastConnect: "May 14 18:35:01", Up: 25.8, Down: 2.9},
		{Ip: "5.5.5.5", Ports: "2323/8444", NodeId: "ffffffff", LastConnect: "May 14 14:25:26", Up: 123.0, Down: 0.0},
	}

	tests := []struct {
		name      string
		Nodes     []FullNode
		threshold float64
		wantIps   []string
	}{
		{name: "ok threshold 0", Nodes: nodes, threshold: 0, wantIps: []string{"3.3.3.3", "5.5.5.5"}},
		{name: "ok threshold 6", Nodes: nodes, threshold: 6, wantIps: []string{"1.1.1.1", "3.3.3.3", "4.4.4.4", "5.5.5.5"}},
		{name: "ok threshold negative", Nodes: nodes, threshold: -2, wantIps: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := FilterByDown{
				Nodes:     tt.Nodes,
				Threshold: tt.threshold,
			}
			got, err := f.Filter()
			if err != nil {
				t.Errorf("Filter() error = %v, no error wanted", err)
				return
			}
			if len(got) != len(tt.wantIps) {
				t.Errorf("Filter() got = %v, want %v", got, tt.wantIps)
				return
			}

			for i, node := range got {
				if node.Ip != tt.wantIps[i] {
					t.Errorf("Filter() got = %v, want %v", node.Ip, tt.wantIps[i])
					return
				}
			}
		})
	}
}
