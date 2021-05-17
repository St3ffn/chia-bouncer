package bouncer

import (
	"strings"
)

// Filter acts as base for filtering certain FullNodes
type Filter interface {
	// Filter certain Fullnodes
	Filter() ([]FullNode, error)
}

// FilterByLocation can filter FullNodes by it's location
type FilterByLocation struct {
	Nodes            []FullNode
	LocationToFilter string
}

// Filter performs filtering of FullNodes by LocationToFilter.
// The Filtering is be performed in a case insensitive manner.
func (f FilterByLocation) Filter() ([]FullNode, error) {
	var filtered []FullNode
	for _, node := range f.Nodes {
		location, err := node.IpLocation()

		if err != nil {
			return nil, err
		}
		// ignore case
		if strings.Contains(strings.ToLower(location), strings.ToLower(f.LocationToFilter)) {
			filtered = append(filtered, node)
		}

	}
	return filtered, nil
}

// FilterByDown can filter FullNodes by it's down speed.
// All Nodes which have a lower or equal down speed than specified by threshold will be present
type FilterByDown struct {
	Nodes []FullNode
	// threshold the minimum (inclusive down speed)
	Threshold float64
}

// Filter performs filtering of FullNodes by down threshold.
// All Nodes which have a lower or equal down speed than specified by threshold will be present
func (f FilterByDown) Filter() ([]FullNode, error) {
	var filtered []FullNode
	for _, node := range f.Nodes {
		if node.Down <= f.Threshold {
			filtered = append(filtered, node)
		}
	}
	return filtered, nil
}
