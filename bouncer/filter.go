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
