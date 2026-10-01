package main

import (
	"strings"
	"testing"

	"github.com/myelophone/goserver/web/runtime"
)

func TestRenderUsesOnlyCacheablePublicStaticRules(t *testing.T) {
	rules := map[string]runtime.RouteRule{
		"/**":         {PublicStatic: true, Cache: &runtime.CachePolicy{MaxAge: 60}},
		"/account/**": {Cache: &runtime.CachePolicy{MaxAge: 60}},
		"/pricing":    {PublicStatic: true, SWR: 300},
	}
	config := string(render(rules))

	for _, want := range []string{
		"~^/account 1; # /account/**",
		"~^/pricing$ 0; # /pricing",
		"~^ 0; # /**",
	} {
		if !strings.Contains(config, want) {
			t.Errorf("generated config is missing %q:\n%s", want, config)
		}
	}
}

func TestNginxPatternMatchesRouteRuleSyntax(t *testing.T) {
	if got, want := routePattern("/products/*"), "^/products/[^/]+$"; got != want {
		t.Errorf("segment wildcard = %q, want %q", got, want)
	}
	if got, want := routePattern("*preview*"), "^.*preview.*$"; got != want {
		t.Errorf("fragment wildcard = %q, want %q", got, want)
	}
	if got, want := routePattern("/catalog/**"), "^/catalog"; got != want {
		t.Errorf("prefix wildcard = %q, want %q", got, want)
	}
}
