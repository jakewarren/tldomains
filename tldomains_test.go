package tldomains

import "testing"

func TestTLDomains(t *testing.T) {
	tests := []struct {
		in  string
		out Host
	}{
		{in: "mMmm.jello.co.uk", out: Host{"mmmm", "jello", "co.uk"}},
		{in: "pressly.com", out: Host{"", "pressly", "com"}},
		{in: "www.pressly.it", out: Host{"www", "pressly", "it"}},
		{in: "sub.www.google.co.uk", out: Host{"sub.www", "google", "co.uk"}},
	}

	extract, _ := New("/tmp/tld.cache")
	for _, tt := range tests {
		h := extract.Parse(tt.in)
		if h.Subdomain != tt.out.Subdomain || h.Root != tt.out.Root || h.Suffix != tt.out.Suffix {
			t.Errorf("expected %v, got %v", tt.out, h)
		}
	}
}
