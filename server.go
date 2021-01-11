/*command gameserver starts up and controls the endpoints for the backend.
 */
package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/tilegame/gameserver/cookiez"
	"github.com/tilegame/gameserver/echoserver"
	"github.com/tilegame/gameserver/wshandle"
	"golang.org/x/crypto/acme/autocert"
)

// HelpMessage is the extra information given when using flags -h or --help.
// it is prefixed to the automatically generated info by flag.PrintDefaults()
const HelpMessage = `
The Ninja Arena Server!

USAGE:
  gameserver [options]

EXAMPLES:
  gameserver -a localhost:8080
  gameserver -a=localhost:8080
  gameserver -a :http -tls

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
	HelpAddress = "Host Address and Port to listen on."
	HelpTLS     = "Enable TLS and AutoCert for LetsEncrypt."
	HelpIndex   = "Homepage file. Only matters if file server is enabled."
	HelpIO      = "Enable Stdin input and Stdout output."
	HelpFiles   = "Enables the File Server"
)

const (
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

var cookieServer = cookiez.NewCookieServer()

var endpoints = map[string]func(http.ResponseWriter, *http.Request){
	"/ws":       serveWebSocket,
	"/ws/echo":  serveWebSocketEcho,
	"/cookie":   cookieServer.ServeCookies,
	"/sessions": cookieServer.HandleInfo,
}

var endpointDescriptions = map[string]string{
	"/ws":       "Main websocket connection for game (not implemented yet)",
	"/ws/echo":  "echo server used for testing connection speeds",
	"/cookie":   "generates and/or validates new cookies for clients",
	"/sessions": "generates a list of active sessions",
}

var (
	htmlpagepath   string
	BasicPagePlate *template.Template
)

func init() {
	flag.StringVar(&addr, "address", DefaultAddress, HelpAddress)
	flag.StringVar(&addr, "a", DefaultAddress, HelpAddress)
	flag.StringVar(&index, "index", DefaultIndex, HelpIndex)
	flag.BoolVar(&useStdinStdout, "io", false, HelpIO)
	flag.BoolVar(&usingTLS, "tls", false, HelpTLS)
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

	// Retrieve the directory that contians the source code.  This
	// allows us to locate the resources folder, which contains
	// lots of useful things.
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("No caller information")
	}
	thefolder := filepath.Dir(filename)

	// Compile the homepage html template.
	a := "resources/templates/endpoints_info.html"
	htmlpagepath = filepath.Join(thefolder, a)
	log.Println("homepage template located at:", htmlpagepath)
	BasicPagePlate = template.Must(template.ParseFiles(htmlpagepath))
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
		cookieServer.SetCookieSecurity(false)
		log.Println("started server on:", s.Addr)
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
	log.Println("started server on:", s.Addr)
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
	err := BasicPagePlate.Execute(w, endpointDescriptions)
	if err != nil {
		fmt.Fprintln(w, "You've reached The Bach End!")
		log.Println(err)
	}
}

// TODO: move this somewhere other than server.go
var clientroom = wshandle.NewClientRoom()

func serveWebSocket(w http.ResponseWriter, r *http.Request) {
	clientroom.Handle(w, r)
}

// TODO: use this /ws/echo endpoint for literally just an echo testing
// server, or change it to /ws/speed and use it as a connection speed
// testing server.
func serveWebSocketEcho(w http.ResponseWriter, r *http.Request) {
	echoserver.HandleWs(w, r)
}

func logRequest(r *http.Request) {
	log.Printf("(%v) %v %v %v", r.RemoteAddr, r.Proto, r.Method, r.URL)
}
