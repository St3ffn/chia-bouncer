package bouncer

import (
	"testing"
)

func Test_convertToNode(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		wantNode FullNode
		wantErr  bool
	}{
		{
			name: "line is ok",
			line: "FULL_NODE 12.13.14.15                          12550/8444  e49d1edd... May 14 18:34:56    199.9|0.0",
			wantNode: FullNode{
				Ip:          "12.13.14.15",
				Ports:       "12550/8444",
				NodeId:      "eeee1edd",
				LastConnect: "May 14 18:34:56",
				UpDown:      "199.9|0.0",
			},
		},
		{
			name:    "line is nok",
			line:    "FULL_NODE 12.13.14.15                          12550/8444  eeee1edd... May 14 18:34:00",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNode, err := convertToNode(tt.line)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("convertToNode() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			testNode(*gotNode, tt.wantNode, t)
		})
	}
}

func testNode(got, want FullNode, t *testing.T) {
	if got.NodeId != want.NodeId {
		t.Errorf("convertToNode() differ in NodeId gotNode = %v, want %v", got, want)
	}
	if got.Ip != want.Ip {
		t.Errorf("convertToNode() differ in ip GotNode = %v, want %v", got, want)
	}
	if got.Ports != want.Ports {
		t.Errorf("convertToNode() differ in Ports gotNode = %v, want %v", got, want)
	}
	if got.LastConnect != want.LastConnect {
		t.Errorf("convertToNode() differ in LastConnect gotNode = %v, want %v", got, want)
	}
	if got.UpDown != want.UpDown {
		t.Errorf("convertToNode() differ in UpDown gotNode = %v, want %v", got, want)
	}
}

const conncetions = `Connections:
Type      IP                                     Ports       NodeID      Last Connect      MiB Up|Dwn
FARMER    127.0.0.1                              47198/8447  aaaaaaaa... May 14 11:15:20      0.4|0.0
FULL_NODE 1.1.1.1                          61075/8444  bbbbbbbb... May 14 18:35:01     12.8|5.8
                                                 -SB Height:   280545    -Hash: aaa5d7b0...
FULL_NODE 2.2.2.2                          8444/8444  cccccccc... May 14 18:35:01      3.2|12.1
                                                 -SB Height:   280545    -Hash: abb5d7b0...
FULL_NODE 3.3.3.3                           8444/8444  dddddddd... May 14 11:55:26      0.0|0.0
                                                 -SB Height:        0    -Hash:  Info...
WALLET    127.0.0.1                              51578/8449  asasasas... May 14 18:34:46      5.6|0.0
FULL_NODE 4.4.4.4                          8444/8444  wewewewe... May 14 18:35:01     25.8|2.9
                                                 -SB Height:   280544    -Hash: aayr0731...

`

const connectionsNoFullNode = `Connections:
Type      IP                                     Ports       NodeID      Last Connect      MiB Up|Dwn
FARMER    127.0.0.1                              47198/8447  aaaaaaaa... May 14 11:15:20      0.4|0.0
WALLET    127.0.0.1                              51578/8449  asasasas... May 14 18:34:46      5.6|0.0

`

func Test_convertToNodes(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantNodes []FullNode
		wantErr   bool
	}{
		{name: "no full nodes", input: connectionsNoFullNode, wantNodes: nil},
		{name: "several full nodes", input: conncetions, wantNodes: []FullNode{
			{Ip: "1.1.1.1", Ports: "61075/8444", NodeId: "bbbbbbbb", LastConnect: "May 14 18:35:01", UpDown: "12.8|5.8"},
			{Ip: "2.2.2.2", Ports: "8444/8444", NodeId: "cccccccc", LastConnect: "May 14 18:35:01", UpDown: "3.2|12.1"},
			{Ip: "3.3.3.3", Ports: "8444/8444", NodeId: "dddddddd", LastConnect: "May 14 11:55:26", UpDown: "0.0|0.0"},
			{Ip: "4.4.4.4", Ports: "8444/8444", NodeId: "wewewewe", LastConnect: "May 14 18:35:01", UpDown: "25.8|2.9"},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNodes, err := convertToNodes([]byte(tt.input))
			if err != nil {
				if !tt.wantErr {
					t.Errorf("convertToNodes() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if len(gotNodes) != len(tt.wantNodes) {

				t.Errorf("convertToNodes() got and wanted nodes differ in size gotNodes = %v, wantNodes %v",
					gotNodes, tt.wantNodes)
			}
			for i, gotNode := range gotNodes {
				testNode(gotNode, tt.wantNodes[i], t)
			}
		})
	}
}
