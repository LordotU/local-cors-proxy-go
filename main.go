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
	"github.com/logrusorgru/aurora"
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
	fmt.Println(aurora.Black(" Local Cors Proxy Go ").Bold().Underline().BgWhite())
	fmt.Print("\n")
	fmt.Printf(aurora.Sprintf(aurora.Blue("Proxy url: %v\n"), aurora.Green(options.cleanURL)))
	fmt.Printf(aurora.Sprintf(aurora.Blue("Proxy url section: %v\n"), aurora.Green(options.cleanURLSection)))
	fmt.Printf(aurora.Sprintf(aurora.Blue("Proxy port: %v\n"), aurora.Green(options.port)))
	fmt.Print("\n")
	fmt.Printf(
		aurora.Sprintf(
			aurora.Cyan("To start using the proxy simply replace the proxied section of your url with: %s"),
			aurora.Sprintf(
				aurora.Yellow("http://%v/%v").Bold(),
				options.addr,
				options.cleanURLSection,
			),
		),
	)
	fmt.Print("\n")

	waitingForGracefulShutdown(server)
}

func waitingForGracefulShutdown(server *fasthttp.Server) {
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGTERM)

	interruptReason := <-interruptChan

	fmt.Print("\n")
	fmt.Println(
		aurora.Sprintf(
			aurora.Red("Shutting down cause %v..\n"),
			interruptReason,
		),
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)

	go func() {
		server.Shutdown()
		defer cancel()
	}()

	<-ctx.Done()
}
