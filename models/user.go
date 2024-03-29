package models

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"github.com/go-redis/redis"
)

var (
	ErrUserNotFound = errors.New("User not found")
	ErrInvalidLogin = errors.New("Invalid Login") 
)

func RegisterUser(username, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return client.Set("user:" + username, hash, 0).Err()
}

func AuthenticateUser(username, password string) error {
	hash, err := client.Get("user:" + username).Bytes()
	if err == redis.Nil {
		return ErrUserNotFound
	} else if err != nil {
		return err
	}	
	err = bcrypt.CompareHashAndPassword(hash, []byte(password))
	if err != nil {
		return ErrInvalidLogin
	}
	return nil
}