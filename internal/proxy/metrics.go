package proxy

import "github.com/valyala/fasthttp"

// Metrics is a snapshot of the cumulative proxy traffic counters.
type Metrics struct {
	BytesIn  int64 `json:"bytesIn"`
	BytesOut int64 `json:"bytesOut"`
	Requests int64 `json:"requests"`
}

// countTraffic adds one request's byte totals to the cumulative counters. The
// sizes are approximate — request/response header bytes plus any in-memory
// body — which is accurate enough for a live traffic meter.
func (p *ProxyServer) countTraffic(ctx *fasthttp.RequestCtx) {
	in := int64(len(ctx.Request.Header.Header())) + int64(len(ctx.Request.Body()))
	out := int64(len(ctx.Response.Header.Header())) + int64(len(ctx.Response.Body()))
	p.bytesIn.Add(in)
	p.bytesOut.Add(out)
	p.requests.Add(1)
}

// Metrics returns a snapshot of the cumulative traffic counters. It is safe to
// call concurrently and survives in-process reloads (the counters live on the
// server, not the router).
func (p *ProxyServer) Metrics() Metrics {
	return Metrics{
		BytesIn:  p.bytesIn.Load(),
		BytesOut: p.bytesOut.Load(),
		Requests: p.requests.Load(),
	}
}
