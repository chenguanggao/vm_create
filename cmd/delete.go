/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
func DeleteCmd() *cobra.Command {
	var deleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "a deleted vm .",
		Long:  `a deleted vm .`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("delete called")
		},
	}
	return deleteCmd
}
