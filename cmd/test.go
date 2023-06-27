/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/wxsms/geekbang-downloader/apis"
	"github.com/wxsms/geekbang-downloader/helpers"
	"log"

	"github.com/spf13/cobra"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "test if cookie in setting file is valid",
	Run: func(cmd *cobra.Command, args []string) {
		var article apis.ArticleDetailResp
		if err := helpers.Request(apis.ArticleDetailApi, fmt.Sprintf(`{"id":"%d","include_neighbors":true,"is_freelyread":true}`, 426265), &article); err != nil {
			log.Fatal(err)
			return
		} else {
			log.Println("cookie is valid!")
		}
	},
}

func init() {
	rootCmd.AddCommand(testCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// testCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// testCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
