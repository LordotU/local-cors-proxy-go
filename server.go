package main

import (
	"fmt"
	"log"
	"regexp"

	. "github.com/logrusorgru/aurora"
	"github.com/valyala/fasthttp"
)

func getServer(options *Options) *fasthttp.Server {
	proxyPathPattern := regexp.MustCompile(`^\/` + options.cleanUrlSection)

	requestHandler := func(ctx *fasthttp.RequestCtx) {
		proxyPath := string(ctx.Path())

		if proxyPathPattern.MatchString(proxyPath) {
			proxyRequestHandler(
				ctx,
				options,
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
		err := server.ListenAndServe(options.addr)
		if err != nil {
			log.Fatal(err)
		} else {
			fmt.Println(Red("Shutted down!\n").Bold())
		}
	}()

	return server
}

func proxyRequestHandler(ctx *fasthttp.RequestCtx, options *Options, proxyPath string) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	res := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(res)

	proxiedUri := options.cleanUrl + proxyPath

	ctx.Request.CopyTo(req)
	req.SetRequestURI(proxiedUri)
	for headerName, headerValue := range options.parsedHeaders {
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
		Sprintf(
			"%s %s %s %s %s %s %d\n",
			Magenta(ctx.Method()),
			Blue("request proxied:"),
			Green(ctx.RequestURI()),
			Blue("->"),
			Green(proxiedUri),
			Blue("with status code"),
			White(res.StatusCode()),
		),
	)
}
