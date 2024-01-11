package cmd

import (
	"fmt"
	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/flytam/filenamify"
	"github.com/spf13/cobra"
	"github.com/wxsms/geekbang-downloader/apis"
	"github.com/wxsms/geekbang-downloader/helpers"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func downloadCourse(dest string, cid string, localImage bool) {

	var info apis.ColumnInfoResp
	var list apis.ArticleListResp
	if err := helpers.Request(apis.ColumnInfoApi, fmt.Sprintf(`{"product_id":%s,"with_recommend_article":true}`, cid), &info); err != nil {
		log.Fatal(err)
		return
	}

	log.Printf("----------------------")
	log.Printf("ID:\t%v", cid)
	log.Printf("Title:\t%v", info.Data.Title)
	log.Printf("Images:\t%v", localImage)
	log.Printf("----------------------")

	if err := helpers.Request(apis.ArticleListApi, fmt.Sprintf(`{"cid":"%s","size":500,"prev":0,"order":"earliest","sample":false}`, cid), &list); err != nil {
		log.Fatal(err)
		return
	}

	courseTitle, _ := filenamify.FilenamifyV2(info.Data.Title)
	courseTitle = strings.TrimSpace(courseTitle)
	path := filepath.Join(dest, courseTitle)
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		log.Fatal(err)
		return
	}

	converter := md.NewConverter("", true, nil)
	for i, a := range list.Data.List {
		index := i + 1
		log.Printf("[%d/%d] %s", index, len(list.Data.List), a.Title)
		var article apis.ArticleDetailResp
		if err := helpers.Request(apis.ArticleDetailApi, fmt.Sprintf(`{"id":"%d","include_neighbors":true,"is_freelyread":true}`, a.Id), &article); err != nil {
			log.Fatal(err)
			return
		}
		markdown, err := converter.ConvertString(article.Data.ArticleContent)
		if err != nil {
			log.Fatal(err)
			return
		}
		markdown = fmt.Sprintf("# %s\n\n%s", article.Data.ArticleTitle, markdown)
		title, _ := filenamify.FilenamifyV2(article.Data.ArticleTitle)
		markdownFile := fmt.Sprintf("%d__%s.md", index, title)
		if err := os.WriteFile(filepath.Join(path, markdownFile), []byte(markdown), 0o644); err != nil {
			log.Fatal(err)
			return
		}
		if localImage {
			helpers.ReplaceRemoteImagesWithLocal(path, markdownFile)
		}
		// api rate limit
		time.Sleep(5 * time.Second)
	}

	log.Printf("Course %s successfully downloaded.", cid)
}

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Download a course by a given id into current folder",
	Run: func(cmd *cobra.Command, args []string) {
		//fmt.Println(viper.Get("cookie"))
		cid, _ := cmd.Flags().GetString("course")
		dest, _ := cmd.Flags().GetString("dest")
		localImage, _ := cmd.Flags().GetBool("local-image")
		downloadCourse(dest, cid, localImage)
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
	getCmd.Flags().StringP("dest", "d", ".", "download dest, current dir by default")
	getCmd.Flags().Bool("local-image", false, "download image to local")
	err := getCmd.MarkFlagRequired("course")
	if err != nil {
		panic(err)
	}
}
