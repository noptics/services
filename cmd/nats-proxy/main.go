package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
)

type proxy struct {
}

func (p *proxy) wrapRoute(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log.Println(r.Method, " ", r.URL)

	dest := r.URL.Query().Get("url")
	if len(dest) == 0 {
		log.Println("empty reuqest url")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message": "error", "error": "must provide endpoint in query as 'url'`))
		return
	}

	fmt.Println(dest)

	resp, err := http.Get(dest)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()

	log.Println(resp.Status)
	setHeaders(w, r.Header)
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func setHeaders(w http.ResponseWriter, rh http.Header) {
	allowMethod := "GET"
	allowHeaders := "Content-Type"
	w.Header().Set("Cache-Control", "must-revalidate")
	w.Header().Set("Allow", allowMethod)
	w.Header().Set("Access-Control-Allow-Methods", allowMethod)
	w.Header().Set("Access-Control-Allow-Headers", allowHeaders)

	o := rh.Get("Origin")
	if o == "" {
		o = "*"
	}
	w.Header().Set("Access-Control-Allow-Origin", o)
	w.Header().Set("Access-Control-Allow-Credentials", "true")
}

func main() {
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8282"
	}

	handler := &proxy{}

	log.Println("Starting proxy server on", fmt.Sprintf("%s:%s", host, port))
	r := httprouter.New()
	r.HandleOPTIONS = true
	r.GlobalOPTIONS = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("options")
		if r.Header.Get("Access-Control-Request-Method") != "" {
			// Set CORS headers
			setHeaders(w, r.Header)
		}

		// Adjust status code to 204
		w.WriteHeader(http.StatusNoContent)
	})

	r.GET("/proxy", handler.wrapRoute)
	if err := http.ListenAndServe(fmt.Sprintf("%s:%s", host, port), r); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
