package main

import(
        "github.com/vmihailenco/redis"
        "log"
        "net/url"
        "net/http"
        "net/http/httputil"
        "strings"
)

func main() {
        http.HandleFunc("/", handleConn)
        err := http.ListenAndServe(":23457", nil)
        if err != nil {
                panic(err)
        }
        log.Println("listening on port 23457")
}

func handleConn(w http.ResponseWriter, r *http.Request) {
        log.Println(r.URL)

        password := ""  // no password set
        db := int64(-1) // use default DB
        client := redis.NewTCPClient("localhost:6379", password, db)
        defer client.Close()
        
        host  := r.Header.Get("x-host")
        token := strings.Split(host, ".")[0]
        route := client.SRandMember("route::"+token)
        
        log.Println(token)
        log.Println(route.Val())

        remote, err := url.Parse(route.Val())
        if err != nil {
                panic(err)
        }

        proxy := httputil.NewSingleHostReverseProxy(remote)

        r.Header.Add("X-Proxy", "just-in-the-neighborhood-and-thought-i-would-drop-in")
        r.Header.Add("X-Articulated-Route", route.Val())
        r.Header.Add("X-Token", token)

        proxy.ServeHTTP(w, r)
}
