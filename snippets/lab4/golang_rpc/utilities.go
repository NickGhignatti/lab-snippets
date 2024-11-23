package main

import "sync"

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
