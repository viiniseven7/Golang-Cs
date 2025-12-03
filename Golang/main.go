package main


import (
	"net/http"
	"log"
	"Golang/router"

 )
func main(){

	r:= router.New();
	const addr =":8081";
	log.Printf( "Starting server on %s\n", addr);
	log.Fatal(http.ListenAndServe(addr, r))
}