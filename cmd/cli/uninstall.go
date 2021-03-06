package main

import (
	"io"

	"github.com/spf13/cobra"
	"helm.sh/helm/v3/pkg/action"
)

const uninstallDescription = `
This command consists of multiple subcommands related to uninstalling the osm
control plane.
`

func newUninstallCmd(config *action.Configuration, in io.Reader, out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "uninstall",
		Short:   "uninstall osm",
		Aliases: []string{"delete", "del"},
		Long:    uninstallDescription,
		Args:    cobra.NoArgs,
	}

	if !settings.IsManaged() {
		cmd.AddCommand(newUninstallMeshCmd(config, in, out))
	}

	// The command to uninstall osm cluster-wide resources should be available in both managed and unmanaged environments.
	cmd.AddCommand(newUninstallClusterWideResourcesCmd(in, out))

	return cmd
}
