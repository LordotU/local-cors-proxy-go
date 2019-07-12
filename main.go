package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	. "github.com/logrusorgru/aurora"
	"github.com/valyala/fasthttp"
)

var (
	addr          string
	url           string
	port          uint64
	host          string
	urlSection    string
	serverLogging bool
	headers       string
	parsedHeaders map[string]string
)

var proxyHostClient *fasthttp.Client

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err.Error())
	}

	if err := processFlags(); err != nil {
		log.Fatal(err.Error())
	}

	cleanUrl := regexp.MustCompile(`\/+$`).ReplaceAllString(url, ``)
	cleanUrlSection := regexp.MustCompile(`^\/+|\/+$`).ReplaceAllString(urlSection, ``)
	addr = host + ":" + strconv.FormatUint(port, 10)

	proxyHostClient = &fasthttp.Client{}
	proxyPathPattern := regexp.MustCompile(`^\/` + cleanUrlSection)

	requestHandler := func(ctx *fasthttp.RequestCtx) {
		proxyPath := string(ctx.Path())

		if proxyPathPattern.MatchString(proxyPath) {
			proxyRequestHandler(
				ctx,
				cleanUrl+proxyPathPattern.ReplaceAllString(proxyPath, ``),
			)
		} else {
			ctx.Error("Not found", fasthttp.StatusNotFound)
		}
	}

	server := &fasthttp.Server{
		Handler: requestHandler,
	}

	go func() {
		err := server.ListenAndServe(addr)
		if err != nil {
			log.Fatal(err)
		} else {
			fmt.Println(Red("Shutted down!\n").Bold())
		}
	}()

	fmt.Print("\n")
	fmt.Println(Black(" Local Cors Proxy Go ").Bold().Underline().BgWhite())
	fmt.Print("\n")
	fmt.Printf(Sprintf(Blue("Proxy url: %v\n"), Green(cleanUrl)))
	fmt.Printf(Sprintf(Blue("Proxy url section: %v\n"), Green(cleanUrlSection)))
	fmt.Printf(Sprintf(Blue("Proxy port: %v\n"), Green(port)))
	fmt.Print("\n")
	fmt.Printf(
		Sprintf(
			Cyan("To start using the proxy simply replace the proxied section of your url with: %s"),
			Sprintf(
				Yellow("http://%v/%v").Bold(),
				addr,
				cleanUrlSection,
			),
		),
	)
	fmt.Println("\n")

	waitingForGracefulShutdown(server)
}

func processFlags() error {
	envPort, _ := strconv.ParseUint(os.Getenv("LCP_GO_PORT"), 10, 16)
	envServerLogging, _ := strconv.ParseBool(os.Getenv("LCP_GO_SERVER_LOGGING"))

	flag.StringVar(&url, "url", os.Getenv("LCP_GO_URL"), "")
	flag.StringVar(&url, "u", os.Getenv("LCP_GO_URL"), "")

	flag.Uint64Var(&port, "port", envPort, "")
	flag.Uint64Var(&port, "p", envPort, "")

	flag.StringVar(&host, "host", os.Getenv("LCP_GO_HOST"), "")
	flag.StringVar(&host, "h", os.Getenv("LCP_GO_HOST"), "")

	flag.StringVar(&urlSection, "urlSection", os.Getenv("LCP_GO_URL_SECTION"), "")
	flag.StringVar(&urlSection, "s", os.Getenv("LCP_GO_URL_SECTION"), "")

	flag.BoolVar(&serverLogging, "serverLogging", envServerLogging, "")
	flag.BoolVar(&serverLogging, "l", envServerLogging, "")

	flag.StringVar(&headers, "headers", os.Getenv("LCP_GO_HEADERS"), "")

	flag.Parse()

	if url == "" {
		return errors.New("--url is required!")
	}

	if err := json.Unmarshal([]byte(headers), &parsedHeaders); err != nil {
		return errors.New("--headers is unmarshalled structure!")
	}

	return nil
}

func proxyRequestHandler(ctx *fasthttp.RequestCtx, proxiedUri string) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	res := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(res)

	ctx.Request.CopyTo(req)

	req.SetRequestURI(proxiedUri)
	for headerName, headerValue := range parsedHeaders {
		req.Header.Set(headerName, headerValue)
	}

	if err := proxyHostClient.Do(req, res); err != nil {
		ctx.Error(err.Error(), 500)
	}

	res.Header.Set("Access-Control-Allow-Origin", "*")
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

func waitingForGracefulShutdown(srv *fasthttp.Server) {
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGTERM)

	interruptReason := <-interruptChan

	fmt.Printf(
		Sprintf(
			Red("Shutting down cause %v..\n\n"),
			interruptReason,
		),
	)

	_, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	srv.Shutdown()
}
