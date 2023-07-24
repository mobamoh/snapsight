package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", handleFun)
	fmt.Println("server listening at :1313...")
	if err := http.ListenAndServe(":1313", nil); err != nil {
		log.Fatal(err)
	}

}

func handleFun(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Hello Server</h1>")
}
