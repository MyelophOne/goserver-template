package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/myelophone/goserver/web/runtime"
)

type cacheRule struct {
	pattern string
	rule    runtime.RouteRule
}

func main() {
	output := flag.String("out", "dist/nginx/route-cache.conf", "generated Nginx route cache configuration")
	flag.Parse()

	config, err := runtime.UseRuntimeConfig()
	if err != nil {
		fatal(err)
	}
	if err := os.MkdirAll(filepath.Dir(*output), 0o755); err != nil {
		fatal(err)
	}
	if err := os.WriteFile(*output, render(config.RouteRules), 0o644); err != nil {
		fatal(err)
	}
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func render(rules map[string]runtime.RouteRule) []byte {
	items := make([]cacheRule, 0, len(rules))
	for pattern, rule := range rules {
		items = append(items, cacheRule{pattern: pattern, rule: rule})
	}
	sort.Slice(items, func(i, j int) bool {
		if len(items[i].pattern) == len(items[j].pattern) {
			return items[i].pattern < items[j].pattern
		}
		return len(items[i].pattern) > len(items[j].pattern)
	})

	var out bytes.Buffer
	out.WriteString("# Generated from resolved production websettings routeRules. Do not edit.\n")
	out.WriteString("map $uri $goserver_skip_public_cache {\n    default 1;\n")
	for _, item := range items {
		value := "1"
		if item.rule.PublicStatic && cacheTTL(item.rule) > 0 {
			value = "0"
		}
		fmt.Fprintf(&out, "    ~%s %s; # %s\n", nginxPattern(item.pattern, item.rule.Exclude), value, item.pattern)
	}
	out.WriteString("}\n")
	return out.Bytes()
}

func cacheTTL(rule runtime.RouteRule) int {
	if rule.Cache != nil && rule.Cache.MaxAge > 0 {
		return rule.Cache.MaxAge
	}
	return rule.SWR
}

func nginxPattern(pattern string, exclude []string) string {
	base := routePattern(pattern)
	if len(exclude) == 0 {
		return base
	}
	parts := make([]string, 0, len(exclude))
	for _, item := range exclude {
		if item != "" {
			parts = append(parts, strings.TrimPrefix(routePattern(item), "^"))
		}
	}
	if len(parts) == 0 {
		return base
	}
	return "^(?!(?:" + strings.Join(parts, "|") + "))" + strings.TrimPrefix(base, "^")
}

func routePattern(pattern string) string {
	if strings.HasPrefix(pattern, "*") && strings.HasSuffix(pattern, "*") && len(pattern) > 2 {
		return "^.*" + regexp.QuoteMeta(strings.Trim(pattern, "*")) + ".*$"
	}
	if strings.HasSuffix(pattern, "/**") {
		return "^" + regexp.QuoteMeta(strings.TrimSuffix(strings.TrimSuffix(pattern, "**"), "/"))
	}
	if !strings.Contains(pattern, "*") {
		return "^" + regexp.QuoteMeta(pattern) + "$"
	}

	segments := strings.Split(strings.Trim(pattern, "/"), "/")
	for i, segment := range segments {
		if segment == "*" {
			segments[i] = "[^/]+"
		} else {
			segments[i] = regexp.QuoteMeta(segment)
		}
	}
	return "^/" + strings.Join(segments, "/") + "$"
}
