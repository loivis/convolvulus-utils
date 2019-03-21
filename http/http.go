package http

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
	"github.com/loivis/gs-google-search"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/transform"
)

const userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.96 Safari/537.36"

var gsClient *gs.Client

func init() {
	gsClient = gs.NewClient(gs.WithUserAgent(userAgent))
}

// GetDoc performs http.MethodGet with custom User-Agent and returns goquery document from response body
func GetDoc(url string) (*goquery.Document, error) {
	client := http.Client{}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", userAgent)
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("failed to get %q: %v", url, err)
		return nil, err
	}
	defer resp.Body.Close()

	r, err := decodeHTMLBody(resp.Body)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return doc, nil
}

// decodeHTMLBody returns an decoding reader of the html Body. It tries to guess the encoding from the content.
func decodeHTMLBody(body io.Reader) (io.Reader, error) {
	bnr := bufio.NewReader(body)
	b, err := bnr.Peek(1024)
	if err != nil && err.Error() != "EOF" {
		return nil, err
	}

	e, _, _ := charset.DetermineEncoding(b, "")

	tnr := transform.NewReader(bnr, e.NewDecoder())

	return tnr, nil
}

// Search returns the first link from google search
func Search(q string) (string, error) {
	query := &url.Values{
		"q":   {q},
		"num": {"5"},
		"hl":  {"zh-CN"},
	}

	res, err := gsClient.Search(query)
	if err != nil {
		return "", err
	}

	if len(res) == 0 {
		return "", fmt.Errorf("no search result for %q", q)
	}

	return res[0].Link, nil
}
