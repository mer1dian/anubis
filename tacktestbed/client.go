package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/tack/tackgo/tack"
	"github.com/tack/tackgo/tacktls"
)

type TestClient struct { 
	store *tack.DefaultStore
}


func NewTestClient() *TestClient {
	return &TestClient{tack.NewDefaultStore(true)}
}

func (client *TestClient) run() {

	config := tls.Config{}
	config.Tack = true
	config.InsecureSkipVerify = true

	for {
		// Make TLS connection
		conn, err := tls.Dial("tcp", "localhost:8443", &config)
		if err != nil {log.Fatal(err)}
		
		err = conn.Handshake()
		if err != nil {log.Fatal(err)}
		
		connState := conn.ConnectionState()
		if connState.TackExtension != nil {
			te := connState.TackExtension
			fmt.Println(te)
		} else {
			fmt.Println("No Tack Extension")
		}

		status, err := tack.ProcessStore(client.store, connState.TackExtension, 
			"localhost", 100)
		fmt.Printf("Status = %s, err = %v\n", status.String(), err)
		fmt.Printf("Pinlist =\n%v\n", client.store.String())

		// Make HTTP request
		_, err = http.Get("http://localhost:8080/next")
	}
}
