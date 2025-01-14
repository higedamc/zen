package rule

import (
	"net/http"
	"strings"
)

// https://adguard.com/kb/general/ad-filtering/create-own-filters/#third-party-modifier
type thirdPartyModifier struct {
	inverted bool
}

func (m *thirdPartyModifier) Parse(modifier string) error {
	if modifier[0] == '~' {
		m.inverted = true
	}
	return nil
}

func (m *thirdPartyModifier) ShouldMatch(req *http.Request) bool {
	if req.Header.Get("Sec-Fetch-Site") == "cross-site" {
		return !m.inverted
	}
	if referer := req.Header.Get("Referer"); referer != "" {
		host := req.Host
		if host == "" {
			host = req.URL.Hostname()
		}
		refererURL, err := req.URL.Parse(referer)
		if err != nil {
			return false
		}
		refererHost := refererURL.Hostname()
		if strings.HasSuffix(refererHost, host) {
			return m.inverted
		}
		return !m.inverted
	}
	return false
}
