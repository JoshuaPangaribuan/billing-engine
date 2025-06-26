package pkghttp

import (
	"context"
	"net"
	"net/http"
)

type ipCtxKey struct{}

func GetIPAddressFromContext(ctx context.Context) string {
	ip := ctx.Value(ipCtxKey{})
	if ip == nil {
		return ""
	}
	return ip.(string)
}

func setIPAddressToContext(ctx context.Context, ip string) context.Context {
	return context.WithValue(ctx, ipCtxKey{}, ip)
}

var trueClientIP = http.CanonicalHeaderKey("True-Client-IP")
var xForwardedFor = http.CanonicalHeaderKey("X-Forwarded-For")
var xRealIP = http.CanonicalHeaderKey("X-Real-IP")

func IPAddressExtractorMiddleware(next EndpointHandler) EndpointHandler {

	return func(ctx context.Context, r Request) (response interface{}, err error) {
		// Logic parsing IP Address
		ip := extractRealIP(r)
		ctx = setIPAddressToContext(ctx, ip)
		return next(ctx, r)
	}
}

func extractRealIP(r Request) string {
	var ip string

	if tcip := r.Header().Get(trueClientIP); tcip != "" {
		ip = tcip
	} else if xrip := r.Header().Get(xRealIP); xrip != "" {
		ip = xrip
	} else if xff := r.Header().Get(xForwardedFor); xff != "" {
		ip = xff
	}

	if ip == "" || net.ParseIP(ip) == nil {
		return ""
	}

	return ip
}
