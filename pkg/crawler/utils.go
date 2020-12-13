package crawler

import (
	"net/url"
	"path"
	"strings"

	"github.com/drsigned/gos"
)

func fixURL(URL string, site *gos.URL) string {
	var fixedURL string

	if strings.HasPrefix(URL, "http") {
		// `http://google.com` OR `https://google.com`
		fixedURL = URL
	} else if strings.HasPrefix(URL, "//") {
		// `//google.com/example.php`
		fixedURL = site.Scheme + ":" + URL
	} else if !strings.HasPrefix(URL, "//") {
		if strings.HasPrefix(URL, "/") {
			// `/?thread=10`
			fixedURL = site.Scheme + "://" + site.Host + URL
		} else {
			if strings.HasPrefix(URL, ".") {
				if strings.HasPrefix(URL, "..") {
					fixedURL = site.Scheme + "://" + site.Host + URL[2:]
				} else {
					fixedURL = site.Scheme + "://" + site.Host + URL[1:]
				}
			} else {
				// `console/test.php`
				fixedURL = site.Scheme + "://" + site.Host + "/" + URL
			}
		}
	}

	return fixedURL
}

func decodeChars(s string) string {
	source, err := url.QueryUnescape(s)
	if err == nil {
		s = source
	}

	// In case json encoded chars
	replacer := strings.NewReplacer(
		`\u002f`, "/",
		`\u0026`, "&",
	)
	s = replacer.Replace(strings.ToLower(s))
	return s
}

func getExtType(URL string) string {
	u, err := gos.ParseURL(URL)
	if err != nil {
		return ""
	}

	return path.Ext(u.Path)
}

func strMatch(s1 string, s2 string) bool {
	return strings.ToLower(s1) == strings.ToLower(s2)
}
