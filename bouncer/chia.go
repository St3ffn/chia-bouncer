package bouncer

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
)

// getLocation get the location for the given ip
var getLocation GetLocation = geoIpLookup

// FullNode represents a connected full node
type FullNode struct {
	// Ip V4 e.g. 127.0.0.1, 123.11.12.13
	Ip string
	// Ports in format 61075/8444
	Ports string
	// NodeId represents the node id (e.g. aaaabbcc)
	NodeId string
	// LastConnect in format May 14 11:15:20
	LastConnect string
	// UpDown represents the down- and uplink in MiB (e.g. 0.2|0.0)
	UpDown string
}

// IpLocation to get the corresponding location of the ipV4 of the FullNode
func (f FullNode) IpLocation() (string, error) {
	return getLocation(f.Ip)
}

// Chia interface to provide basic actions for Chia
type Chia interface {
	// ListNodes will list all connected FullNodes
	ListNodes() (nodes []FullNode, err error)
	// RemoveNode will remove a connected Full Node by it's nodeId
	RemoveNode(nodeId string) error
}

// ChiaCli represents the chia command line interface
type ChiaCli struct {
	// ChiaExecutable represents the chia executable e.g. /home/steffen/chia-blockchain/venv/bin/chia
	ChiaExecutable string
}

// ListNodes will list all connected FullNodes.
// It basically acts as wrapper around "chia show -c"
func (c ChiaCli) ListNodes() (nodes []FullNode, err error) {
	out, err := execCmd(c.ChiaExecutable, "show", "-c")
	if err != nil {
		return nil, err
	}
	return convertToNodes(out)
}

// convertToNodes filters the stdout of "chia show -c" command for FullNodes and
// converts them to list of FullNode
func convertToNodes(input []byte) (nodes []FullNode, err error) {
	scanner := bufio.NewScanner(bytes.NewReader(input))
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "FULL_NODE") {
			continue
		}

		node, err := convertToNode(scanner.Text())
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, *node)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return nodes, nil
}

// convertToNode will convert a single FULL_NODE line of the stdout of the "chia show -c" command.
func convertToNode(line string) (node *FullNode, err error) {
	// FULL_NODE 1.2.3.4                          61075/8444  aaaaaaaa... May 14 17:03:40      9.4|5.0
	fields := strings.Fields(line)
	if len(fields) != 8 {
		return nil, fmt.Errorf("line can not be converted: %s", line)
	}

	return &FullNode{
		Ip:          fields[1],
		Ports:       fields[2],
		NodeId:      strings.TrimSuffix(fields[3], "..."),                     // remove ... of string "62b29c64..."
		LastConnect: fmt.Sprintf("%s %s %s", fields[4], fields[5], fields[6]), // like "May 14 17:03:40"
		UpDown:      fields[7],
	}, nil
}

// RemoveNode will remove the connected full node via "chia show -r nodeId"
func (c ChiaCli) RemoveNode(nodeId string) error {
	if _, err := execCmd(c.ChiaExecutable, "show", "-r", nodeId); err != nil {
		return err
	}
	return nil
}
