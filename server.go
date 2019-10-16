package main

import (
	"fmt"
	"log"
	"regexp"

	"github.com/logrusorgru/aurora"
	"github.com/valyala/fasthttp"
)

func getServer(o *options) *fasthttp.Server {
	proxyPathPattern := regexp.MustCompile(`^\/` + o.cleanURLSection)

	requestHandler := func(ctx *fasthttp.RequestCtx) {
		proxyPath := string(ctx.Path())

		if proxyPathPattern.MatchString(proxyPath) {
			proxyRequestHandler(
				ctx,
				o,
				proxyPathPattern.ReplaceAllString(proxyPath, ``),
			)
		} else {
			ctx.Error("Not found", fasthttp.StatusNotFound)
		}
	}

	server := &fasthttp.Server{
		Handler: requestHandler,
	}

	go func() {
		err := server.ListenAndServe(o.addr)
		if err != nil {
			log.Fatal(err)
		} else {
			fmt.Println(aurora.Red("Shutted down!\n").Bold())
		}
	}()

	return server
}

func proxyRequestHandler(ctx *fasthttp.RequestCtx, o *options, proxyPath string) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	res := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(res)

	proxiedURI := o.cleanURL + proxyPath

	ctx.Request.CopyTo(req)
	req.SetRequestURI(proxiedURI)
	for headerName, headerValue := range o.parsedHeaders {
		req.Header.Set(headerName, headerValue)
	}

	if err := fasthttp.Do(req, res); err != nil {
		ctx.Error(err.Error(), 500)
	}

	res.Header.Set("Vary", "*")
	res.Header.Set("Access-Control-Allow-Credentials", "true")
	res.Header.Set("Access-Control-Allow-Methods", "GET,HEAD,PUT,PATCH,POST,DELETE")
	res.Header.Set("Access-Control-Allow-Origin", "*")
	if ctx.IsOptions() {
		accessControlRequestHeaders := string(ctx.Request.Header.Peek("Access-Control-Request-Headers"))
		if accessControlRequestHeaders != "" {
			res.Header.Set("Access-Control-Allow-Headers", accessControlRequestHeaders)
		}
		res.Header.Set("Content-Length", "0")
		res.SetStatusCode(204)
	}

	res.WriteTo(ctx.Conn())

	defer fmt.Printf(
		aurora.Sprintf(
			"%s %s %s %s %s %s %d\n",
			aurora.Magenta(ctx.Method()),
			aurora.Blue("request proxied:"),
			aurora.Green(ctx.RequestURI()),
			aurora.Blue("->"),
			aurora.Green(proxiedURI),
			aurora.Blue("with status code"),
			aurora.White(res.StatusCode()),
		),
	)
}
