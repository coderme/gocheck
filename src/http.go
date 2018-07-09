package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

const (
	src          = "https://github.com/codermeorg/gocheck"
	connectError = "xxx"
)

// Link embrace a URL and its parent page
// where its found
type Link struct {
	Parent,
	URL string
}

var (
	userAgent = fmt.Sprintf(`"Mozilla/5.0 %s/%s (%s)`,
		binName, version,
		src,
	)
	client = &http.Client{
		CheckRedirect: noRedirect,
	}
	rePatterns = map[string]*regexp.Regexp{
		"src":  regexp.MustCompile(`(?i) src=["']?([^<>"']+)`),
		"href": regexp.MustCompile(`(?i) href=["']?([^<>"']+)`),
	}
)

func fetch(link string) *Result {
	defer func() {
		<-limit
	}()

	if link == "" {
		return nil
	}

	r := &Result{}
	r.URL = link
	r.Time = time.Now()
	req, err := http.NewRequest(`GET`, link, nil)
	r.Elapsed = time.Since(r.Time)

	if err != nil {
		return patchedResult(connectError, r)
	}

	req.Header.Set(`User-Agent`, userAgent)

	resp, err := client.Do(req)
	if err != nil {
		return patchedResult(connectError, r)
	}

	defer resp.Body.Close()
	if !isHTML(resp.Header.Get("Content-Type")) {
		return nil
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return patchedResult(connectError, r)
	}

	go discoverURLs(link, string(b))

	return patchedResult(resp.StatusCode, r)
}

func patchedResult(code interface{}, r *Result) *Result {
	var c string

	switch code.(type) {
	case int:
		c = fmt.Sprintf("%s", code)
	case string:
		c = code.(string)
	}
	r.Code = c

	switch {
	case strings.HasPrefix(c, "x"):
		r.ErrorConnect = true
	case *check5xx && strings.HasPrefix(c, "5"):
		r.ErrorServer = true
	case *check4xx && strings.HasPrefix(c, "4"):
		r.ErrorConnect = true
	case *check3xx && strings.HasPrefix(c, "3"):
		r.ErrorRedirect = true
	default:
		return nil
	}

	return r
}

// resolveURL resolve URL relative to Parent page
func resolveURL(parent, u string) (string, error) {
	uP, err := url.Parse(u)
	if err != nil {
		return "", err
	}

	pP, err := url.Parse(parent)
	if err != nil {
		return "", err
	}
	if uP.Host == "" {
		uP.Host = pP.Host
		uP.Scheme = pP.Scheme
	}
	return uP.String(), nil

}

// isHTML checks if the Content-Type of the resp is indeed
// of type HTML or XHTML
func isHTML(c string) bool {
	c = strings.ToLower(c)
	if strings.Contains(c, `text/html`) ||
		strings.Contains(c, `text/xhtml`) {
		return true
	}
	return false
}

func discoverURLs(pageURL, content string) {
	ps := []*regexp.Regexp{}

	if *watchHREF {
		ps = append(ps, rePatterns["href"])
	}

	if *watchSRC {
		ps = append(ps, rePatterns["src"])
	}

	for _, p := range ps {
		matches := p.FindAllStringSubmatch(content, -1)
		if matches == nil {
			continue
		}

		for _, m := range matches {
			// m[0] is the whole matched string
			// m[1] is URL
			u, err := resolveURL(pageURL, m[1])
			if err != nil {
				continue
			}
			if !isSameHost(hostName, u) &&
				!*spanHosts {
				continue
			}

			urlQueue <- u

		}

	}

}

func isSameHost(host, u string) bool {
	p, err := url.Parse(u)
	if err != nil {
		return false
	}
	h := strings.Trim(host, "/. ")
	if h == p.Host {
		return true
	}
	return false
}

func noRedirect(req *http.Request, via []*http.Request) error {
	return http.ErrUseLastResponse
}
