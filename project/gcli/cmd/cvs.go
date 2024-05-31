/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/csv"
	"encoding/json"
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/go-playground/validator/v10"

	"gopkg.in/yaml.v3"
)

type OutputFormat = string

const (
	Json OutputFormat = "json"
	Yaml OutputFormat = "yaml"
)

var inputFile string
var outputFile string
var outputFormat OutputFormat

// cvsCmd represents the cvs command
var cvsCmd = &cobra.Command{
	Use:   "cvs",
	Short: "convert cvs file to json or yaml",
	Long:  `convert cvs file to json or yaml`,
	Run: func(cmd *cobra.Command, args []string) {

		validate := validator.New()
		errs := validate.Var(inputFile, "required")
		if errs != nil {
			cmd.Println(cmd.UsageString())
			cmd.Printf("Error: required flag '%s' not set\n", "input")
			return
		}

		// 检查文件是否存在
		file, err := os.Open(inputFile)
		if err != nil {
			cmd.Println(cmd.UsageString())
			cmd.Printf("Error: file '%s' does not exist\n", inputFile)
			os.Exit(1)
		}
		defer file.Close()
		reader := csv.NewReader(file)

		record, err := reader.ReadAll()
		if err != nil {
			cmd.Println(cmd.UsageString())
			cmd.Printf("Error: file '%s' format is not cvs\n", inputFile)
			os.Exit(1)
		}

		cvs_map, err := csvToMap(record)
		if err != nil {
			cmd.Println(cmd.UsageString())
			cmd.Printf("Error: file '%s' format is not cvs\n", inputFile)
			os.Exit(1)
		}

		if outputFormat == Json {
			jsonData, err := json.Marshal(cvs_map)
			if err != nil {
				panic(err)
			}

			// 打开文件，如果文件不存在则创建
			output, err := os.Create(outputFile)
			if err != nil {
				panic(err)
			}
			defer output.Close() // 确保在函数结束时关闭文件

			// 写入JSON数据到文件
			_, err = output.Write(jsonData)
			if err != nil {
				panic(err)
			}
		} else if outputFormat == Yaml {
			data, err := yaml.Marshal(cvs_map)
			if err != nil {
				log.Fatalf("error: %v", err)
			}
			// 打开文件，如果文件不存在则创建
			output, err := os.Create(outputFile)
			if err != nil {
				panic(err)
			}
			defer output.Close() // 确保在函数结束时关闭文件

			// 写入JSON数据到文件
			_, err = output.Write(data)
			if err != nil {
				panic(err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(cvsCmd)

	cvsCmd.Flags().StringVarP(&inputFile, "input", "i", "", "Input file path")
	cvsCmd.Flags().StringVarP(&outputFile, "output", "o", ".\\output\\output.json", "Output file path")
	cvsCmd.Flags().StringVarP(&outputFormat, "format", "f", "json", "Output format[json, yaml]")

	cvsCmd.MarkFlagRequired("input")
}

func csvToMap(records [][]string) ([]map[string]string, error) {

	// 第一行通常是列标题，用作map的key
	keys := records[0]

	// 创建map的切片
	maps := make([]map[string]string, 0, len(records)-1)

	// 遍历记录，跳过标题行
	for _, record := range records[1:] {
		// 创建map，将列标题和记录值对应起来
		m := make(map[string]string)
		for i, value := range record {
			m[keys[i]] = value
		}
		maps = append(maps, m)
	}

	return maps, nil
}
