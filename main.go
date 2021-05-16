package main

import (
	"fmt"
	"github.com/St3ffn/chia-bouncer/bouncer"
	"os"
)

func main() {
	if err := run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

// run will initialise the cli, get all connected full nodes, filter them by given location and
// finally remove filtered full nodes
func run() error {
	ctx, err := bouncer.RunCli()
	if err != nil {
		return err
	}
	if ctx.Done {
		return nil
	}
	chiaCli := bouncer.ChiaCli{ChiaExecutable: ctx.ChiaExecutable}
	nodes, err := chiaCli.ListNodes()
	if err != nil {
		return err
	}
	filter := bouncer.FilterByLocation{
		Nodes:            nodes,
		LocationToFilter: ctx.Location,
	}
	filtered, err := filter.Filter()
	if err != nil {
		return err
	}

	if len(filtered) == 0 {
		_, _ = fmt.Fprintf(os.Stdout, "nothing from %s\n", ctx.Location)
	}

	for _, node := range filtered {
		if err := chiaCli.RemoveNode(node.NodeId); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Can not filter %s:%s from %s...stopping\n",
				node.Ip, node.NodeId, ctx.Location)
			return err
		}
	}
	_, _ = fmt.Fprintf(os.Stdout, "found %d - filtered %d from %s\n", len(nodes), len(filtered), ctx.Location)
	return nil
}
