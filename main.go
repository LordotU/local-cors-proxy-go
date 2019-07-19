package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	. "github.com/logrusorgru/aurora"
	"github.com/valyala/fasthttp"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err.Error())
	}

	err = godotenv.Load(cwd + "/.env")
	if err != nil {
		log.Fatal(err.Error())
	}

	options, err := getOptions()
	if err != nil {
		log.Fatal(err.Error())
	}

	server := getServer(options)

	fmt.Print("\n")
	fmt.Println(Black(" Local Cors Proxy Go ").Bold().Underline().BgWhite())
	fmt.Print("\n")
	fmt.Printf(Sprintf(Blue("Proxy url: %v\n"), Green(options.cleanUrl)))
	fmt.Printf(Sprintf(Blue("Proxy url section: %v\n"), Green(options.cleanUrlSection)))
	fmt.Printf(Sprintf(Blue("Proxy port: %v\n"), Green(options.port)))
	fmt.Print("\n")
	fmt.Printf(
		Sprintf(
			Cyan("To start using the proxy simply replace the proxied section of your url with: %s"),
			Sprintf(
				Yellow("http://%v/%v").Bold(),
				options.addr,
				options.cleanUrlSection,
			),
		),
	)
	fmt.Println("\n")

	waitingForGracefulShutdown(server)
}

func waitingForGracefulShutdown(server *fasthttp.Server) {
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGTERM)

	interruptReason := <-interruptChan

	fmt.Println(
		Sprintf(
			Red("Shutting down cause %v..\n"),
			interruptReason,
		),
	)

	_, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	server.Shutdown()
}
