package main

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type ReverseProxy struct {
	Name string
	RemoteUrl  string
	LocalUrl string
	Port string
	CurrentValue int //1 for remote, -1 for local, 0 for nothing
	IsLocalRunning bool
}

var reverseProxies = []ReverseProxy {
	{
		Name:           "Server 1",
		RemoteUrl:     "https://jsonplaceholder.typicode.com/posts",
		LocalUrl:       "https://jsonplaceholder.typicode.com/comments",
		Port:           ":9091",
		CurrentValue:   1,
		IsLocalRunning: false,

	},
	{
		Name:           "Server 2",
		RemoteUrl:     "https://jsonplaceholder.typicode.com/comments",
		LocalUrl:       "https://jsonplaceholder.typicode.com/posts",
		Port:           ":9092",
		CurrentValue:   1,
		IsLocalRunning: false,
	},
}
var echos [3] *echo.Echo

func main() {
	startReverseProxies()
	http.HandleFunc("/", serveMainPage)
	http.HandleFunc("/switchReverseProxyServer", switchReverseProxyServer)
	http.HandleFunc("/serverList", serverListHandler)

	log.Fatal(http.ListenAndServe(":9090", nil))
}

func serveMainPage(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	p := "." + r.URL.Path
	if p == "./" {
		p = "./static/index.html"
	}
	http.ServeFile(w, r, p)
}

func serverListHandler(w http.ResponseWriter, r *http.Request) {
	for  index , reverseProxy := range reverseProxies {
		_, err := http.Get(reverseProxy.LocalUrl)
		if err != nil {
			reverseProxies[index].IsLocalRunning = false
		} else {
			reverseProxies[index].IsLocalRunning = true
		}
	}

	data, err := json.Marshal(reverseProxies)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Fprintf(w, string(data))
}


func switchReverseProxyServer(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		return
	}

	serverName := r.FormValue("serverName")

	if len(serverName) == 0{
		return
	}

	for  index , reverseProxy := range reverseProxies {
		if reverseProxy.Name == serverName {
			var newServerValue = ""
			reverseProxies[index].CurrentValue = reverseProxies[index].CurrentValue * -1
			if reverseProxies[index].CurrentValue == 1{
				newServerValue = reverseProxy.RemoteUrl
			} else if reverseProxies[index].CurrentValue == -1 {
				newServerValue = reverseProxy.LocalUrl
			}
			killAllReverseProxyServers()
			startReverseProxies()
			fmt.Fprintf(w, "%s", newServerValue)
		}
	}
}

func startReverseProxies(){
	var index = 0
	for  _ , reverseProxy := range reverseProxies {
		if reverseProxy.CurrentValue == 1 {
			echos[index] = runReverseProxyServer(reverseProxy.RemoteUrl, reverseProxy.Port)
			index++
		} else if reverseProxy.CurrentValue == -1{
			echos[index] = runReverseProxyServer(reverseProxy.LocalUrl, reverseProxy.Port)
			index++
		}

	}
}

func killAllReverseProxyServers() {
	for _, echo := range echos {
		if echo != nil {
			echo.Close()
		}
	}
}

func runReverseProxyServer(targetUrl string, targetPort string) *echo.Echo {
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
