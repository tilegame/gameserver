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

 Serving Files
 -------------
  By default, a minimalist message is displayed when visiting the
  homepage (or a non-existent page).  A fileserver can be enabled
  by using the boolean flag
    -serve-files
  which will then serve "index.html" from the working directory
  as the homepage.  This homepage can be changed by using the flag:
    -index <filepath>

 Server Paths
 ------------
      /     routes to files if the file server is enabled.
      /*    routes to any file in the directory and subdirectories.
      /ws   routes to the websocket connection.  Has no files.

 Input and Output
 ----------------
  If the -io flag is used, then stdin and stdout will be enabled.
  Stdin will be scanned line by line (will wait for each line), so
  commands can be sent interactively.

OPTIONS:
`

const (
	HelpAddress    = "Host Address and Port for Standard connections."
	HelpAddressTLS = "Host Address and Port for TLS connections."
	HelpIndex      = "Homepage file. Only matters if file server is enabled."
	HelpIO         = "Enable Stdin input and Stdout output."
	HelpFiles      = "Enables the File Server"
	DefaultAddress = "localhost:8080"
	DefaultIndex   = "index.html"
)

var (
	useStdinStdout bool
	usingTLS       bool
	usingFiles     bool
	addr           string
	index          string
)

var endpoints = map[string]func(http.ResponseWriter, *http.Request){
	"/ws":      serveWebSocket,
	"/ws/echo": serveWebSocketEcho,
}

func init() {
	flag.StringVar(&addr, "address", DefaultAddress, HelpAddress)
	flag.StringVar(&addr, "a", DefaultAddress, HelpAddress)
	flag.StringVar(&index, "index", DefaultIndex, HelpIndex)
	flag.BoolVar(&useStdinStdout, "io", false, HelpIO)
	flag.BoolVar(&usingTLS, "tls", false, HelpAddressTLS)
	flag.BoolVar(&usingFiles, "serve-files", false, HelpFiles)
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, HelpMessage)
		flag.PrintDefaults()
	}

	// waiting until initialization to add this handler allows
	// the user to determine whether or not the program should run
	// a file server, or just a minimialist webpage.
	if usingFiles {
		endpoints["/"] = serveFiles
	} else {
		endpoints["/"] = serveMinimal
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
	for path, handler := range endpoints {
		mux.HandleFunc(path, handler)
	}
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

func serveFiles(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), 405)
		return
	}
	http.ServeFile(w, r, r.URL.Path[1:])
	return
}

func serveMinimal(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	fmt.Fprintln(w, "You've reached The Bach End!")

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
