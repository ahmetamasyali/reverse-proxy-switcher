package main

import (
	"fmt"
	"github.com/labstack/echo"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

var echos [1] *echo.Echo

var firstUrlValueMap= map[string]string{
	"postComment": "https://jsonplaceholder.typicode.com/posts",
}

var secondUrlValueMap = map[string]string{
	"postComment": "https://jsonplaceholder.typicode.com/comments",
}

var serverPortMap = map[string]string{
	"postComment": ":9091",
}

//1 for first value, 2 for second value, 0 for nothing
var currentServerValues = map[string]int{
	"postComment": 1,
}

func main() {


	startServers()

	http.HandleFunc("/switchServer", switchServer)
	http.HandleFunc("/", serveMainPage)

	log.Fatal(http.ListenAndServe(":9090", nil))
}

func startServersHandler(w http.ResponseWriter, r *http.Request) {
	startServers()
	fmt.Fprintf(w, "servers started", r.URL.Path[1:])
}

func serveMainPage(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	p := "." + r.URL.Path
	if p == "./" {
		p = "./static/index.html"
	}
	http.ServeFile(w, r, p)
}

func startServers(){
	var index = 0
	for  serverName, value := range currentServerValues {
		if value == 1 {
			echos[index] = runProxyServer(firstUrlValueMap[serverName], serverPortMap[serverName])
			index++
		} else if value == 2{
			echos[index] = runProxyServer(secondUrlValueMap[serverName], serverPortMap[serverName])
			index++
		}

	}
}

func switchServer(w http.ResponseWriter, r *http.Request) {
	serverName := r.FormValue("serverName")

	if len(serverName) == 0{
		return
	}

	if _, ok := currentServerValues[serverName]; ok {
		var newServerValue = ""
		if currentServerValues[serverName] == 1{
			currentServerValues[serverName] = 2
			newServerValue = secondUrlValueMap[serverName]
		} else if currentServerValues[serverName] == 2{
			currentServerValues[serverName] = 1
			newServerValue = firstUrlValueMap[serverName]
		}
		killAllProxyServers()
		startServers()

		fmt.Fprintf(w, "%s", newServerValue)
	}
}

func killServers(w http.ResponseWriter, r *http.Request) {
	killAllProxyServers()
	fmt.Fprintf(w, "Servers are killed", r.URL.Path[1:])
}

func killAllProxyServers() {
	for _, echo := range echos {
		 if echo != nil {
			 echo.Close()
		 }
	}
}

func runProxyServer(targetUrl string, targetPort string) *echo.Echo {
	e := echo.New()

	// create the reverse proxy
	url, _ := url.Parse(targetUrl)
	proxy := httputil.NewSingleHostReverseProxy(url)

		reverseProxyRoutePrefix := ""
	routerGroup := e.Group(reverseProxyRoutePrefix)
	routerGroup.Use(func(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
		return func(context echo.Context) error {

			req := context.Request()
			res := context.Response().Writer

			// Update the headers to allow for SSL redirection
			req.Host = url.Host
			req.URL.Host = url.Host
			req.URL.Scheme = url.Scheme

			//trim reverseProxyRoutePrefix
			path := req.URL.Path
			req.URL.Path = strings.TrimLeft(path, reverseProxyRoutePrefix)

			// ServeHttp is non blocking and uses a go routine under the hood
			proxy.ServeHTTP(res, req)
			return nil
		}
	})

	go e.Start(targetPort)

	return e
}
