package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	message   string
	signalled bool
)

func init() {
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	message = fmt.Sprintf("Hello from %s", hostname)
}

func main() {
	sigchan := make(chan os.Signal, 1)
	defer func() {
		close(sigchan)
	}()
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	http.HandleFunc("/", handler)

	svr := &http.Server{
		Addr:    ":9090",
		Handler: http.DefaultServeMux,
	}

	go func() {
		sig := <-sigchan

		message = "REQUEST RECEIVED AFTER SHUTDOWN SIGNAL"
		signalled = true

		fmt.Println("Received signal:", sig)
		fmt.Println("Waiting 10s before stopping..")

		time.Sleep(10 * time.Second)

		svr.Close()
	}()

	fmt.Println("Listening on port 9090")

	err := svr.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	if signalled {
		fmt.Println("REQUEST RECEIVED AFTER SHUTDOWN SIGNAL")
	}

	fmt.Fprintln(w, message)
}
