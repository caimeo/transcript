package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"

	"github.com/caimeo/console"
	"github.com/caimeo/iniflags"
)

var port = flag.Int("port", 8080, "Port to listen on")
var showbody = flag.Bool("body", true, "Output body?")

func main() {
	iniflags.SetConfigFile(".settings")
	iniflags.SetAllowMissingConfigFile(true)
	iniflags.Parse()

	portchannel := make(chan int, 1)

	go ServeConent(*port, portchannel)

	finalport := <-portchannel

	console.Always("Listening on port:", finalport)

	for {

	}
}

func provideData(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "OK")

	if r.URL.String() == "/favicon.ico" {
		return
	}
	out, _ := httputil.DumpRequest(r, *showbody)
	console.Always(string(out))
}

func ServeConent(bindport int, p chan int) {
	http.HandleFunc("/", provideData)

	var server *http.Server
	var listener net.Listener
	var finalport = 0
	for finalport == 0 {
		console.Debug("bindport ", bindport)
		addr := fmt.Sprintf(":%d", bindport)
		server = &http.Server{Addr: addr, Handler: nil}
		ln, err := net.Listen("tcp", addr)
		if err != nil {
			if strings.Contains(err.Error(), "bind: permission denied") || strings.Contains(err.Error(), "bind: address already in use") {
				bindport++
			} else {
				check(err)
			}
		} else {
			finalport = bindport
		}
		listener = ln
	}

	p <- finalport
	check(server.Serve(tcpKeepAliveListener{listener.(*net.TCPListener)}))

}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// Unfortunately because tcpKeepAliveListener in net/html is private I have to reimplement it
// tcpKeepAliveListener sets TCP keep-alive timeouts on accepted
// connections. It's used by ListenAndServe and ListenAndServeTLS so
// dead TCP connections (e.g. closing laptop mid-download) eventually
// go away.
type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}
