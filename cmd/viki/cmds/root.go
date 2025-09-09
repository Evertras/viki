package cmds

import "github.com/spf13/cobra"

var RootCmd = &cobra.Command{
	Use:   "viki",
	Short: "Viki is a tool for turning Obsidian vaults into websites.",
}
