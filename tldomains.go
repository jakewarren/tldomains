package tldomains

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

var tlds = make(map[string]struct{}, 0)

type tldomains struct {
	CacheFile string
}

// New creates a new *tldomains using the specified filepath
func New(cacheFile string) (*tldomains, error) {
	data, err := ioutil.ReadFile(cacheFile)
	if err != nil {
		data, err = download()
		if err != nil {
			return &tldomains{}, err
		}
		if err = ioutil.WriteFile(cacheFile, data, 0644); err != nil {
			return &tldomains{}, err
		}
	}
	list := strings.Split(string(data), "\n")
	for _, item := range list {
		if item == "" || strings.HasPrefix(item, "//") {
			continue
		}
		tlds[item] = struct{}{}
	}

	return &tldomains{CacheFile: cacheFile}, nil
}

// Host contains the parsed info for the domain
type Host struct {
	Subdomain, Root, Suffix string
}

// Parse extracts a domain into it's component parts
func (extract *tldomains) Parse(host string) Host {
	var h Host

	nhost := strings.ToLower(host)
	parts := strings.Split(nhost, ".")

	if len(parts) == 0 {
		h.Root = host
		return h
	}

	var suffix string
	for i := len(parts) - 1; i >= 0; i-- {
		p := parts[i]

		if suffix == "" {
			suffix = p
		} else {
			suffix = fmt.Sprintf("%s.%s", p, suffix)
		}

		if _, ok := tlds[suffix]; ok {
			h.Suffix = suffix
		} else if h.Root == "" {
			h.Root = p
		} else {
			if h.Subdomain == "" {
				h.Subdomain = p
			} else {
				h.Subdomain = p + "." + h.Subdomain
			}
		}
	}

	return h
}

func download() ([]byte, error) {
	u := "https://publicsuffix.org/list/public_suffix_list.dat"
	resp, err := http.Get(u)
	if err != nil {
		return []byte(""), err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	lines := strings.Split(string(body), "\n")
	var buffer bytes.Buffer

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "//") {
			buffer.WriteString(line)
			buffer.WriteString("\n")
		}
	}

	return buffer.Bytes(), nil
}
