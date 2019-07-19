# Local Cors Proxy Go

[![License](https://img.shields.io/badge/License-MIT-000000.svg)](https://opensource.org/licenses/MIT)

## Description

Simple and fast (built on top of [fasthttp](https://github.com/valyala/fasthttp) proxy to bypass [CORS](https://developer.mozilla.org/ru/docs/Web/HTTP/CORS) issues during prototyping against existing APIs without having to worry about CORS

It was built to solve the issue of getting errors like this:

```text
... has been blocked by CORS policy: No 'Access-Control-Allow-Origin' header is present on the requested resource. If an opaque response serves your needs, set the request's mode to 'no-cors' to fetch the resource with CORS disabled.
```

### Limitations

Because [`func Do(req *Request, resp *Response) error`](https://godoc.org/github.com/valyala/fasthttp#Do) is using for all type of queries it doesn't support redirects!

## Usage

Let's imagine API endpoint that we want to request that has CORS issues:

```text
https://licenseapi.herokuapp.com/licenses/mit
```

Pull Docker image and run a container:

```bash
docker pull lordotu/lcp-go

docker run -dti \
  -e LCP_GO_URL=https://licenseapi.herokuapp.com \
  -e LCP_GO_HOST=0.0.0.0 \
  -p 8118:8118 \
  --name lcp-go \
  lordotu/lcp-go
```

Then in your client code, new API endpoint:

```text
http://localhost:8118/proxy/licenses/mit
```

End result will be a request to `https://licenseapi.herokuapp.com/licenses/mit` without the CORS issues!

Alternatively you can build binary (for Linux) from sources with command `CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o lcp-go` and run it like:

```bash
lcp-go --url https://licenseapi.herokuapp.com
```

Or:

```bash
LCP_GO_URL=https://licenseapi.herokuapp.com lcp-go
```

## Configuring

You may set params via command line args or via env variables. All defaults are stored in `.env` file in the working directory.

Only one argument is **required**: `--url` (or `LCP_GO_URL` if you prefer env variables).

### Options

| Option          | Shorthand | Example                           | Default   |
| --------------- | --------- | --------------------------------- | --------: |
| --url           | -u        | https://licenseapi.herokuapp.com  |           |
| --port          | -p        | 8119                              |      8118 |
| --host          | -h        | 0.0.0.0                           | localhost |
| --urlSection    | -s        | through                           |     proxy |
| --serverLogging | -l        | true                              |     false |
| --headers       |           | {"X-Requested-With": "Corsyusha"} |        {} |

### Environment variables

| Option                   | Example                           | Default   |
| ------------------------ | --------------------------------- | --------: |
| LCP_GO_URL            | https://licenseapi.herokuapp.com     |           |
| LCP_GO_PORT           | 8119                                 |      8118 |
| LCP_GO_HOST           | 0.0.0.0                              | localhost |
| LCP_GO_URL_SECTION    | through                              |     proxy |
| LCP_GO_SERVER_LOGGING | true                                 |     false |
| LCP_GO_HEADERS        | {"X-Requested-With": "Corsyusha"}    |        {} |

---

###### Inspired by: https://github.com/LordotU/corsyusha
