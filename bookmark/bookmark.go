package bookmark

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
	flags "github.com/jessevdk/go-flags"
	"github.com/pkg/errors"
)

type bookmarkCmdOpts struct {
	URL string `short:"u" long:"url" description:"URL address to parse bookmakrs from" required:"true"`
}

var bookmarkCmd = bookmarkCmdOpts{}

func InitBookmark(parser *flags.Parser) error {
	_, err := parser.AddCommand(
		"bookmark",
		"Parse a web page for bookmarks",
		"",
		&bookmarkCmd,
	)
	if err != nil {
		return errors.Wrapf(err, "InitBookmark")
	}

	return nil
}

func (bcmd *bookmarkCmdOpts) Execute(args []string) error {
	// resp, err := http.Get(bcmd.URL)
	// if err != nil {
	// 	return errors.Wrapf(err, "bookmark: get url %s", bcmd.URL)
	// }

	// defer resp.Body.Close()

	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	return errors.Wrapf(err, "bookmarkCmdOpts.Execute read body")
	// }

	// err = bcmd.parseBody(resp.Body)
	// if err != nil {
	// 	return errors.Wrapf(err, "bookmark: parse body")
	// }

	err := bcmd.parse(bcmd.URL)
	if err != nil {
		return errors.Wrapf(err, "bookmark: parse")
	}

	return nil
}

func (bcmd *bookmarkCmdOpts) parse(url string) error {
	doc, err := goquery.NewDocument(bcmd.URL)
	if err != nil {
		return errors.Wrapf(err, "bookmark: new document")
	}

	title := doc.Find("title").Contents().Text()
	fmt.Println(title)

	doc.Find("h2").Each(func(index int, item *goquery.Selection) {
		fmt.Println(item.Text())
		item.Find("h2").Each(func(i int, item *goquery.Selection) {
			href, ok := item.Attr("href")
			if ok {
				fmt.Println("href:", href)
			}
		})
	})
	// fmt.Println(h2.Contents())

	return nil
}

// func (bcmd *bookmarkCmdOpts) parseBody(body io.ReadCloser) error {
// 	z := html.NewTokenizer(body)

// 	depth := 0
// 	for {
// 		tt := z.Next()

// 		switch {
// 		case tt == html.ErrorToken:
// 			// End of the document, we're done
// 			return nil
// 		case tt == html.StartTagToken:
// 			t := z.Token()

// 			err := handleToken(t)
// 			if err != nil {
// 				return errors.Wrapf(err, "")
// 			}
// 		}
// 	}
// }

// func handleToken(token html.Token) error {
// 	switch {
// 	case token.Data == "h2":
// 		fmt.Println(token)
// 	case token.Data == "a":
// 		url, err := getHref(token)
// 		if err != nil {
// 			return errors.Wrapf(err, "bookmark: handle token")
// 		}

// 		if strings.HasPrefix(url, "http") {
// 			fmt.Println("Found href:", url)
// 		}
// 	}

// 	return nil
// }

// func getHref(token html.Token) (string, error) {
// 	for _, a := range token.Attr {
// 		if a.Key == "href" {
// 			return a.Val, nil
// 		}
// 	}

// 	return "", errors.Errorf("bookmark: href not found")
// }
