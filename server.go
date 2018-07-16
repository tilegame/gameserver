package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/fractalbach/ninjaServer/echoserver"
	"golang.org/x/crypto/acme/autocert"
)

const HelpMessage = `
The Ninja Arena Server!

USAGE:    
  ninjaServer [options]

EXAMPLES:
  ninjaServer -a localhost:8080
  ninjaServer -a=localhost:8080
  ninjaServer -a :http -tls

INFORMATION:
  Files will be served relative to the current directory,
  which will become the root directory for all files.

  Server Paths
      /     routes to index.html
      /*    routes to any file in the directory and subdirectories.
      /ws   routes to the websocket connection.  Has no files.

  If the -io flag is used, then stdin and stdout will be enabled.
  Stdin will be scanned line by line (will wait for each line), so
  commands can be sent interactively.

OPTIONS:
`

const (
	HelpAddress    = "Host Address and Port for Standard connections."
	HelpAddressTLS = "Host Address and Port for TLS connections."
	HelpIndex      = "Homepage file"
	HelpIO         = "Enable Stdin input and Stdout output."
	DefaultAddress = "localhost:8080"
	DefaultIndex   = "index.html"
)

var (
	useStdinStdout bool
	usingTLS       bool
	addr           string
	index          string
)

func init() {
	flag.StringVar(&addr, "address", DefaultAddress, HelpAddress)
	flag.StringVar(&addr, "a", DefaultAddress, HelpAddress)
	flag.StringVar(&index, "index", DefaultIndex, HelpIndex)
	flag.BoolVar(&useStdinStdout, "io", false, HelpIO)
	flag.BoolVar(&usingTLS, "tls", false, HelpAddressTLS)
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, HelpMessage)
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()
	if useStdinStdout {
		go inputLoop()
	}
	runServer()
}

func inputLoop() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		handleStdinCommand(scanner.Text())
	}
}

func handleStdinCommand(line string) {
	log.Println("[Stdin]", line)
	switch line {
	case "hello":
		fmt.Println("well hello to you as well!")
	case "quit", "exit", "goodbye", "stop":
		fmt.Println("Shutting down server...")
		log.Fatal("Shutting down server by request from stdin.")
	}
}

func runServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", serveHome)
	mux.HandleFunc("/ws", serveWebSocket)
	mux.HandleFunc("/ws/echo", serveWebSocketEcho)
	if !usingTLS {
		s := &http.Server{Addr: addr, Handler: mux}
		log.Fatal(s.ListenAndServe())
	}
	m := &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist("thebachend.com"),
		Cache:      autocert.DirCache("certs"),
	}
	go http.ListenAndServe(":http", m.HTTPHandler(nil))
	s := &http.Server{
		Addr:      ":https",
		TLSConfig: &tls.Config{GetCertificate: m.GetCertificate},
		Handler:   mux,
	}
	log.Fatal(s.ListenAndServeTLS("", ""))
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), 405)
		return
	}
	http.ServeFile(w, r, r.URL.Path[1:])
	return
}

func serveWebSocket(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Not yet implemented.")
}

func serveWebSocketEcho(w http.ResponseWriter, r *http.Request) {
	echoserver.HandleWs(w, r)
}

func logRequest(r *http.Request) {
	log.Printf("(%v) %v %v %v", r.RemoteAddr, r.Proto, r.Method, r.URL)
}
