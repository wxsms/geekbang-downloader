package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/spf13/cobra"
	"log"
	"net/url"
	"os"
)

func printPage(url string, index int, title string) chromedp.Tasks {
	htmlStr := ""

	return chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.WaitVisible(".catalog .article-title .title"),
		chromedp.Click(fmt.Sprintf(".catalog .article-title .title[title=\"%s\"", title)),
		chromedp.WaitVisible("#article-content-container"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			if data, err := page.CaptureSnapshot().Do(ctx); err != nil {
				return err
			} else {
				htmlStr = data
				return nil
			}
		}),
		chromedp.ActionFunc(func(ctx context.Context) error {
			return os.WriteFile(fmt.Sprintf("%d-%s.mhtml", index, title), []byte(htmlStr), 0o644)
		}),
	}
}

func getAllLessons(url string, nodes *[]*cdp.Node) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.WaitVisible(".catalog .article-title .title"),
		chromedp.Nodes(`.catalog .article-title .title`, nodes),
	}
}

func listenForNetworkEvent(ctx context.Context) {
	var id network.RequestID

	chromedp.ListenTarget(ctx, func(ev interface{}) {
		switch ev := ev.(type) {
		case *network.EventRequestWillBeSent:
			uri, _ := url.Parse(ev.Request.URL)
			if uri.Path != "/serv/v1/column/articles" {
				break
			}
			id = ev.RequestID
		case *network.EventLoadingFinished:
			if ev.RequestID != id {
				return
			}
			go func() {
				_ = chromedp.Run(ctx, chromedp.ActionFunc(func(ctx context.Context) error {
					body, err := network.GetResponseBody(ev.RequestID).Do(ctx)
					if err != nil {
						fmt.Println(err)
					} else {

						var res struct {
							Data struct {
								List []struct {
									Title string `json:"article_title"`
									Id    int    `json:"id"`
								} `json:"list"`
							} `json:"data"`
						}
						err := json.Unmarshal(body, &res)
						if err != nil {
							fmt.Println(err.Error())
						} else {
							fmt.Println("ok", res)
						}
					}
					return nil
				}))
			}()

		}
	})
}

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
		u := cmd.Flag("url").Value.String()

		fmt.Println(u, args)
		fmt.Println("get called")

		opts := append(chromedp.DefaultExecAllocatorOptions[:],
			// if user-data-dir is set, chrome won't load the default profile,
			// even if it's set to the directory where the default profile is stored.
			// set it to empty to prevent chromedp from setting it to a temp directory.
			chromedp.UserDataDir(""),
			// in headless mode, chrome won't load the default profile.
			chromedp.Flag("headless", false),
			chromedp.Flag("disable-extensions", false),
		)

		// create context
		ctx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
		defer cancel()
		ctx, cancel = chromedp.NewContext(ctx)
		defer cancel()
		listenForNetworkEvent(ctx)

		lessons := make([]*cdp.Node, 0)
		if err := chromedp.Run(ctx, getAllLessons(u, &lessons)); err != nil {
			log.Fatal(err)
		}
		//for i, n := range lessons {
		//	title, _ := n.Attribute("title")
		//	//fmt.Printf("working on %d,%s...\n", i, title)
		//	//if err := chromedp.Run(ctx, printPage(u, i + 1, title)); err != nil {
		//	//	log.Fatal(err)
		//	//}
		//}

		log.Println("done")
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
	getCmd.Flags().StringP("url", "u", "", "course url to download, for example: https://time.geekbang.org/column/intro/100093501?tab=catalog")
}
