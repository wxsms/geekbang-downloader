package cmd

import (
	"fmt"
	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/spf13/cobra"
	"github.com/wxsms/geekbang-downloader/apis"
	"github.com/wxsms/geekbang-downloader/helpers"
	"log"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Download a course",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		//fmt.Println(viper.Get("cookie"))
		cid, _ := cmd.Flags().GetString("course")

		var info apis.ColumnInfoResp
		var list apis.ArticleListResp
		if err := helpers.Request(apis.ColumnInfoApi, fmt.Sprintf(`{"product_id":%s,"with_recommend_article":true}`, cid), &info); err != nil {
			log.Fatal(err)
			return
		}
		helpers.Debug(cmd, "get course info result:", info)

		if err := helpers.Request(apis.ArticleListApi, fmt.Sprintf(`{"cid":"%s","size":500,"prev":0,"order":"earliest","sample":false}`, cid), &list); err != nil {
			log.Fatal(err)
			return
		}
		helpers.Debug(cmd, "get article list result:", list)

		converter := md.NewConverter("", true, nil)
		helpers.Debug(cmd, converter)

	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	getCmd.Flags().StringP("course", "c", "", "course id to download, for example, id is [100093501] in this url: https://time.geekbang.org/column/intro/100093501?tab=catalog")
	getCmd.MarkFlagRequired("course")
}
