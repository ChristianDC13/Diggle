package urltools

import "testing"

func TestExtractDomain(t *testing.T) {
	table := []struct {
		url      string
		expected string
		err      bool
	}{
		{"http://en.wikipedia.org/", "wikipedia.org", false},
		{"https://en.wikipedia.org/", "wikipedia.org", false},
		{"https://en.wikipedia.org/wiki/URL", "wikipedia.org", false},
		{"https://en.wikipedia.org/wiki/URL?query=1", "wikipedia.org", false},
		{"//en.wikipedia.org/", "wikipedia.org", false},
		{"bodyContent", "", true},
		{"/wiki/Special:Random", "", true},
		{"#p-lang-btn", "", true},
		{"#", "", true},
		{"", "", true},
		{"http://", "", true},
		{"http:///", "", true},
	}

	for _, test := range table {
		domain, err := ExtractDomain(test.url)
		if domain != test.expected {
			t.Errorf("Value: ExtractDomain(%s) = got %s; want %s", test.url, domain, test.expected)
		}
		if (err != nil) != test.err {
			t.Errorf("Error: ExtractDomain(%s) = got %v; want %v", test.url, err, test.err)
		}

	}

}
