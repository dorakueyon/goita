package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/codegangsta/cli"
	"github.com/fatih/color"
	"github.com/mattn/go-runewidth"
	"log"
	"os"
	"strings"
	"unicode/utf8"
)

type entry struct {
	Title     string
	URL       string
	LikeCount string
}

type QueryResult struct {
	Title   string
	Entries []*entry
}

func main() {
	app := cli.NewApp()
	app.Name = "goita"
	app.Usage = "Command Line Client for Qiita ranking"
	app.HideHelp = true
	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "number, n",
			Value: 10,
			Usage: "number of output lines",
		},
	}
	app.Action = func(c *cli.Context) {
		number := c.Int("number")
		url := buildUrl()
		result, err := crawl(url, number)
		if err != nil {
			log.Fatal(err)
		}
		showResult(result, url)
	}
	app.Run(os.Args)
}

func buildUrl() string {
	return fmt.Sprintf("http://qrank.wbsrv.net/entries?days=1&orderby=like_count")
}

func crawl(url string, number int) (QueryResult, error) {
	entries := []*entry{}
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return QueryResult{
			Title:   "",
			Entries: entries,
		}, err
	}
	doc.Find("tr").Each(func(_ int, line *goquery.Selection) {
		cells := [3]interface{}{"", "", 0}
		cells[0] = line.Find("a").Text()
		cells[1], _ = line.Find("a").Attr("href")
		cells[2] = line.Find("td").Eq(2).Text()
		fmt.Println(err)
		entry := entry{
			Title:     strings.TrimSpace(cells[0].(string)),
			URL:       cells[1].(string),
			LikeCount: cells[2].(string),
		}
		if entry.URL != "" {
			entries = append(entries, &entry)
		}
	})
	if number > len(entries) {
		number = len(entries)
	}
	return QueryResult{
		Title:   doc.Find(".subtitle1").Text(),
		Entries: entries[:number],
	}, nil
}

func maxTitleWidth(entries []*entry) int {
	width := 0
	for _, e := range entries {
		count := runewidth.StringWidth(e.Title)
		if count > width {
			width = count
		}
	}
	return width
}

func maxURLWidth(entries []*entry) int {
	width := 0
	for _, e := range entries {
		count := utf8.RuneCountInString(e.URL)
		if count > width {
			width = count
		}
	}
	return width
}

func showResult(result QueryResult, url string) {
	entries := result.Entries
	if len(entries) == 0 {
		fmt.Println("Rankingｹﾞﾄできなかたyo!")
		fmt.Printf("  url: %s \n\n", url)
		return
	}
	fmt.Printf("%s : %d 件\n",
		result.Title,
		len(entries),
	)
	titleWidth := maxTitleWidth(entries)
	titleFmt := fmt.Sprintf("%%-%ds", titleWidth)

	urlWidth := maxURLWidth(entries)
	urlFmt := fmt.Sprintf("%%-%ds", urlWidth)

	fmt.Fprintf(color.Output, " %s | %s | %s \n",
		color.BlueString(titleFmt, "Title"),
		fmt.Sprintf(urlFmt, "Url"),
		color.CyanString("Like"),
	)
	fmt.Println(strings.Repeat("-", titleWidth+urlWidth+16))
	for _, e := range entries {
		fmt.Fprintf(color.Output, " %s | %s | %s \n",
			color.BlueString(runewidth.FillRight(e.Title, titleWidth)),
			fmt.Sprintf(urlFmt, e.URL),
			color.CyanString(e.LikeCount),
		)
	}
}
