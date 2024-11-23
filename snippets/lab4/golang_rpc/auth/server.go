package main

import (
	"errors"
	"flag"
	"fmt"
	"net"
	"net/rpc"
	"strconv"
	"sync"
)

type User struct {
	Username string
	Password string
}

type AuthService struct {
	Users map[string]*User
	Mutex sync.Mutex
}

type RegisterRequest struct {
	Username string
	Password string
}

type AuthRequest struct {
	Username string
	Password string
}

type Response struct {
	Message string
	Success bool
}

func newAuthService() *AuthService {
	return &AuthService{
		Users: make(map[string]*User),
	}
}

func (authService *AuthService) RegisterUser(request RegisterRequest, response *Response) error {
	authService.Mutex.Lock()
	defer authService.Mutex.Unlock()

	if _, exists := authService.Users[request.Username]; exists {
		response.Message = "Username already exists"
		return errors.New(response.Message)
	}

	authService.Users[request.Username] = &User{
		Username: request.Username,
		Password: request.Password,
	}
	response.Message = "User registered successfully"
	response.Success = true

	return nil
}

func (authService *AuthService) AuthenticateUser(request AuthRequest, response *Response) error {
	authService.Mutex.Lock()
	defer authService.Mutex.Unlock()

	user, exists := authService.Users[request.Username]

	if !exists || request.Password != user.Password {
		response.Message = "Invalid username or password"
		return errors.New(response.Message)
	}

	response.Message = "Authentication successful"
	response.Success = true

	return nil
}

func main() {
	authService := newAuthService()

	err := rpc.Register(authService)
	if err != nil {
		fmt.Println("Error registering the service : ", err)
		return
	}

	rpcPort := flag.Int("p", 8080, "port for the rpc server")
	flag.Parse()

	listener, err := net.Listen("tcp", "localhost:"+strconv.Itoa(*rpcPort))
	if err != nil {
		fmt.Println("Error starting listener :", err)
		return
	}
	defer listener.Close()

	fmt.Println("Authentication service is running on port : ", *rpcPort)

	for {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting the connection :", err)
			continue
		}
		go rpc.ServeConn(connection)
	}
}
