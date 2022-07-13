// Copyright (C) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package main

import (
	"fmt"
	"os"

	"github.com/axiacoin/axia-network-v2-runner/cmd/runner/control"
	"github.com/axiacoin/axia-network-v2-runner/cmd/runner/ping"
	"github.com/axiacoin/axia-network-v2-runner/cmd/runner/server"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:        "axia-network-runner",
	Short:      "axia-network-runner commands",
	SuggestFor: []string{"network-runner"},
}

func init() {
	cobra.EnablePrefixMatching = true
}

func init() {
	rootCmd.AddCommand(
		server.NewCommand(),
		ping.NewCommand(),
		control.NewCommand(),
	)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "axia-network-runner failed %v\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}
