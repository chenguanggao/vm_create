/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
func VersionCmd() *cobra.Command {
	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version of the vm tool",
		Long:  `vm tool create .`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("version v1.0 -- HEAD")
		},
	}
	return versionCmd
}
