package configs_test

import (
	"net/url"
	"testing"
)

func TestUrl(t *testing.T) {
	dns := "cloudflare:?CF_API_TOKEN=${CF_API_TOKEN}"

	uri, err := url.Parse(dns)
	if err != nil {
		t.Fatal(err)
	}
	println(uri.Scheme)

	if uri.Scheme != "cloudflare" {
		t.Fatalf("expected scheme to be cloudflare, got %q", uri.Scheme)
	}

	res := uri.Query().Get("CF_API_TOKEN")
	println(res)
	if res != "${CF_API_TOKEN}" {

		t.Fatalf("expected CF_API_TOKEN to be ${CF_API_TOKEN}, got %q", res)
	}
}
