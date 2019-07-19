package main

import (
	"encoding/json"
	"errors"
	"flag"
	"os"
	"regexp"
	"strconv"
)

type Options struct {
	url             string
	cleanUrl        string
	port            uint64
	host            string
	urlSection      string
	cleanUrlSection string
	serverLogging   bool
	headers         string
	parsedHeaders   map[string]string
	addr            string
}

func getOptions() (*Options, error) {
	options := Options{}

	envPort, _ := strconv.ParseUint(os.Getenv("LCP_GO_PORT"), 10, 16)
	envServerLogging, _ := strconv.ParseBool(os.Getenv("LCP_GO_SERVER_LOGGING"))

	flag.StringVar(&options.url, "url", os.Getenv("LCP_GO_URL"), "")
	flag.StringVar(&options.url, "u", os.Getenv("LCP_GO_URL"), "")

	flag.Uint64Var(&options.port, "port", envPort, "")
	flag.Uint64Var(&options.port, "p", envPort, "")

	flag.StringVar(&options.host, "host", os.Getenv("LCP_GO_HOST"), "")
	flag.StringVar(&options.host, "h", os.Getenv("LCP_GO_HOST"), "")

	flag.StringVar(&options.urlSection, "urlSection", os.Getenv("LCP_GO_URL_SECTION"), "")
	flag.StringVar(&options.urlSection, "s", os.Getenv("LCP_GO_URL_SECTION"), "")

	flag.BoolVar(&options.serverLogging, "serverLogging", envServerLogging, "")
	flag.BoolVar(&options.serverLogging, "l", envServerLogging, "")

	flag.StringVar(&options.headers, "headers", os.Getenv("LCP_GO_HEADERS"), "")

	flag.Parse()

	if options.url == "" {
		return nil, errors.New("--url is required!")
	}

	if err := json.Unmarshal([]byte(options.headers), &options.parsedHeaders); err != nil {
		return nil, errors.New("--headers is unmarshalled structure!")
	}

	options.cleanUrl = regexp.MustCompile(`\/+$`).ReplaceAllString(options.url, ``)
	options.cleanUrlSection = regexp.MustCompile(`^\/+|\/+$`).ReplaceAllString(options.urlSection, ``)
	options.addr = options.host + ":" + strconv.FormatUint(options.port, 10)

	return &options, nil
}
