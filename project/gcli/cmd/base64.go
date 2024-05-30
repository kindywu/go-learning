/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// base64Cmd represents the base64 command
var base64Cmd = &cobra.Command{
	Use:   "base64",
	Short: "Base64 command family",
	Long:  `base64 is a command family that provides encoding and decoding functionalities`,
}

var base64EncodeCmd = &cobra.Command{
	Use:   "encode",
	Short: "Encode data to base64",
	Long:  `encode takes input data and encodes it to a base64 string representation`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("base64 encode called")
	},
}

var base64DecodeCmd = &cobra.Command{
	Use:   "decode",
	Short: "Decode base64 data",
	Long:  `decode takes a base64 string and decodes it back to the original data`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("base64 decode called")
	},
}

func init() {
	base64Cmd.AddCommand(base64EncodeCmd)
	base64Cmd.AddCommand(base64DecodeCmd)
	rootCmd.AddCommand(base64Cmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// base64Cmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// base64Cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
