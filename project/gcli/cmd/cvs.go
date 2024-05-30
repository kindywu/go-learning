/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// cvsCmd represents the cvs command
var cvsCmd = &cobra.Command{
	Use:   "cvs",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("cvs called")
	},
}

func init() {
	rootCmd.AddCommand(cvsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cvsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cvsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
