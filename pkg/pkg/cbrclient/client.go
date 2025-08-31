package cbrclient

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"

	"golang.org/x/net/html/charset"
)

type ValCurs struct {
	Date   string   `xml:"Date,attr"`
	Name   string   `xml:"name,attr"`
	Valute []Valute `xml:"Valute"`
}

type Valute struct {
	ID       string `xml:"ID,attr"`
	NumCode  string `xml:"NumCode"`
	CharCode string `xml:"CharCode"`
	Nominal  int    `xml:"Nominal"`
	Name     string `xml:"Name"`
	Value    string `xml:"Value"`
}

type Client struct {
	http         *http.Client
	base         string
	fallbackBase string
	ua           string
}

func New(timeout time.Duration) *Client {
	tr := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	return &Client{
		http:         &http.Client{Timeout: timeout, Transport: tr},
		base:         "https://www.cbr.ru/scripts/XML_daily_eng.asp",
		fallbackBase: "http://www.cbr.ru/scripts/XML_daily_eng.asp",
		ua:           "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) CBRStats/1.0 Chrome/120 Safari/537.36",
	}
}

func (c *Client) FetchByDate(ctx context.Context, ddmmyyyy string) (*ValCurs, error) {

	v, code, err := c.fetchOnce(ctx, c.base, ddmmyyyy)
	if err == nil {
		return v, nil
	}
 
	if code == http.StatusForbidden {
		if v2, _, err2 := c.fetchOnce(ctx, c.fallbackBase, ddmmyyyy); err2 == nil {
			return v2, nil
		}
	}
	return nil, err
}

func (c *Client) fetchOnce(ctx context.Context, base, ddmmyyyy string) (*ValCurs, int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		fmt.Sprintf("%s?date_req=%s", base, ddmmyyyy), nil)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("User-Agent", c.ua)
	req.Header.Set("Accept", "application/xml,text/xml,*/*;q=0.1")
	req.Header.Set("Connection", "keep-alive")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode/100 != 2 {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return nil, resp.StatusCode, fmt.Errorf("cbr http %d: %s", resp.StatusCode, string(b))
	}

	dec := xml.NewDecoder(resp.Body)
	dec.CharsetReader = charset.NewReaderLabel
	var v ValCurs
	if err := dec.Decode(&v); err != nil {
		return nil, resp.StatusCode, err
	}
	return &v, resp.StatusCode, nil
}
