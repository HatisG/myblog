package main

import (
	"fmt"
	"net/http"
)

func main() {

	http.HandleFunc("/ping", pingHandler)
	http.HandleFunc("/post", postsHandler)

	fmt.Println("Server is running on :8080")
	http.ListenAndServe(":8080", nil)

}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"msg":"success"}`))
}

func postsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`[]`))
}
