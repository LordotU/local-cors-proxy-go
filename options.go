package main

import (
	"encoding/json"
	"errors"
	"flag"
	"os"
	"regexp"
	"strconv"
)

type options struct {
	url             string
	port            uint64
	host            string
	urlSection      string
	serverLogging   bool
	headers         string
	origin          string
	reflectOrigin   bool
	parsedHeaders   map[string]string
	cleanURL        string
	cleanURLSection string
	addr            string
}

func getOptions() (*options, error) {
	o := options{}

	envPort, _ := strconv.ParseUint(os.Getenv("LCP_GO_PORT"), 10, 16)
	envServerLogging, _ := strconv.ParseBool(os.Getenv("LCP_GO_SERVER_LOGGING"))
	envOrigin := os.Getenv("LCP_ORIGIN")
	envReflectOrigin, _ := strconv.ParseBool(os.Getenv("LCP_REFLECT_ORIGIN"))

	flag.StringVar(&o.url, "url", os.Getenv("LCP_GO_URL"), "")
	flag.StringVar(&o.url, "u", os.Getenv("LCP_GO_URL"), "")

	flag.Uint64Var(&o.port, "port", envPort, "")
	flag.Uint64Var(&o.port, "p", envPort, "")

	flag.StringVar(&o.host, "host", os.Getenv("LCP_GO_HOST"), "")
	flag.StringVar(&o.host, "h", os.Getenv("LCP_GO_HOST"), "")

	flag.StringVar(&o.urlSection, "urlSection", os.Getenv("LCP_GO_URL_SECTION"), "")
	flag.StringVar(&o.urlSection, "s", os.Getenv("LCP_GO_URL_SECTION"), "")

	flag.BoolVar(&o.serverLogging, "serverLogging", envServerLogging, "")
	flag.BoolVar(&o.serverLogging, "l", envServerLogging, "")

	flag.StringVar(&o.headers, "headers", os.Getenv("LCP_GO_HEADERS"), "")

	flag.StringVar(&o.origin, "origin", envOrigin, "")
	flag.StringVar(&o.origin, "o", envOrigin, "")

	flag.BoolVar(&o.reflectOrigin, "reflectOrigin", envReflectOrigin, "")
	flag.BoolVar(&o.reflectOrigin, "r", envReflectOrigin, "")

	flag.Parse()

	if o.url == "" {
		return nil, errors.New("--url is required")
	}

	if err := json.Unmarshal([]byte(o.headers), &o.parsedHeaders); err != nil {
		return nil, errors.New("--headers is unmarshalled structure")
	}

	o.cleanURL = regexp.MustCompile(`\/+$`).ReplaceAllString(o.url, ``)
	o.cleanURLSection = regexp.MustCompile(`^\/+|\/+$`).ReplaceAllString(o.urlSection, ``)
	o.addr = o.host + ":" + strconv.FormatUint(o.port, 10)

	return &o, nil
}
