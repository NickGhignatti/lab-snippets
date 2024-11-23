package main

import (
	"flag"
	"fmt"
	"net/rpc"
	"strconv"
)

func main() {
	rpcPort := flag.Int("p", 8080, "port of the rpc server")
	username := flag.String("username", "DefaultUser", "username of the client")
	password := flag.String("password", "DefaultPassword", "password of the client")
	flag.Parse()

	client, err := rpc.Dial("tcp", "localhost:"+strconv.Itoa(*rpcPort))
	if err != nil {
		fmt.Println("Error connecting to the server :", err)
		return
	}

	registerRequest := RegisterRequest{
		Username: *username,
		Password: *password,
	}

	var response Response
	err = client.Call("AuthService.RegisterUser", registerRequest, &response)
	if err != nil {
		fmt.Println("Error registering user:", err)
	} else {
		fmt.Println("Register respones:", response.Message)
	}

	authRequest := AuthRequest{
		Username: *username,
		Password: *password,
	}

	err = client.Call("AuthService.AuthenticateUser", authRequest, &response)
	if err != nil {
		fmt.Println("Error authenticate user:", err)
	} else {
		fmt.Println("Authentication respones:", response.Message)
	}
}
