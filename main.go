package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"

	"github.com/pkg/browser"
)

func printErrorAndExit(context string, err error) {
	fmt.Fprintf(os.Stderr, "[servethis] %s: %v\n", context, err)
	os.Exit(1)
}

func verboseLog(verboseFlag *bool, format string, args ...interface{}) {
	if *verboseFlag {
		fmt.Fprintf(os.Stderr, "[servethis] "+format+"\n", args...)
	}
}

func main() {
	var (
		verbose       *bool = flag.Bool("v", false, "verbose output")
		useFileSystem *bool = flag.Bool("file", false, "open from filesystem instead of http")
		port          *int  = flag.Int("p", 0, "port to serve on (default: a random open port)")
	)

	flag.Parse()

	f, err := os.CreateTemp("", "")
	if err != nil {
		_, err = io.Copy(os.Stdout, os.Stdin)
		printErrorAndExit("creating temp file", err)
	}
	defer os.Remove(f.Name())

	verboseLog(verbose, "using temp file: %s", f.Name())

	cp := io.TeeReader(os.Stdin, os.Stdout)
	_, err = io.Copy(f, cp)
	if err != nil {
		printErrorAndExit("reading input", err)
	}

	if *useFileSystem {
		err = browser.OpenURL(fmt.Sprintf("file://%s", f.Name()))
		if err != nil {
			printErrorAndExit("opening browser", err)
		}
		<-make(chan struct{})
	} else {
		listener, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
		if err != nil {
			printErrorAndExit("starting http server", err)
		}

		verboseLog(verbose, "listening on http://localhost:%d", listener.Addr().(*net.TCPAddr).Port)

		go func() {
			err := browser.OpenURL(fmt.Sprintf("http://localhost:%d", listener.Addr().(*net.TCPAddr).Port))
			if err != nil {
				printErrorAndExit("opening browser", err)
			}
		}()

		http.Serve(listener, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, f.Name())
		}))
	}
}
