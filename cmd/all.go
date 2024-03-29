package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/wxsms/geekbang-downloader/apis"
	"github.com/wxsms/geekbang-downloader/helpers"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var allCmd = &cobra.Command{
	Use:   "all",
	Short: "Download all courses",
	Run: func(cmd *cobra.Command, args []string) {
		localImage, _ := cmd.Flags().GetBool("local-image")
		skipExist, _ := cmd.Flags().GetBool("skip-exist")
		dest, _ := cmd.Flags().GetString("dest")
		page := 1
		for {
			var res apis.ProductListResp
			if err := helpers.Request(apis.ProductListApi, fmt.Sprintf(`{
    "tag_ids": [],
    "product_type": 0,
    "product_form": 1,
    "pvip": 0,
    "prev": %d,
    "size": 20,
    "sort": 8,
    "with_articles": true
}`, page), &res); err != nil {
				log.Fatal(err)
				return
			}

			if res.Code != 0 {
				log.Fatalf("bad response code: %d", res.Code)
			}

			for _, product := range res.Data.Products {
				if skipExist {
					courseTitle := helpers.ToFilename(product.Title)
					path := filepath.Join(dest, courseTitle)
					if _, err := os.Stat(path); err == nil {
						log.Printf("skip: %d %s\n", product.Id, product.Title)
						continue
					}
				}

				downloadCourse(dest, strconv.Itoa(product.Id), localImage)
				time.Sleep(2 * time.Second)
			}

			if !res.Data.Page.More {
				break
			}
			page++
		}
	},
}

func init() {
	rootCmd.AddCommand(allCmd)
	allCmd.Flags().StringP("dest", "d", ".", "download dest, current dir by default")
	allCmd.Flags().Bool("local-image", false, "download image to local")
	allCmd.Flags().Bool("skip-exist", false, "skip exist course (if the folder present in dest)")
}
