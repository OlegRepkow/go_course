package services

import (
	"context"
	"errors"
	"log"
	"module_6/internal/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	users     *mongo.Collection
	jwtSecret string
}

func NewAuthService(db *mongo.Database, secret string) *AuthService {
	return &AuthService{
		users:     db.Collection("users"),
		jwtSecret: secret,
	}
}

func (s *AuthService) SignUp(username, password string) error {
	log.Printf("[INFO] AuthService.SignUp: creating user=%s\n", username)
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("[ERROR] AuthService.SignUp: bcrypt failed for user=%s: %v\n", username, err)
		return err
	}

	user := models.User{Username: username, Password: string(hash)}
	_, err = s.users.InsertOne(context.Background(), user)
	if err != nil {
		log.Printf("[ERROR] AuthService.SignUp: DB insert failed for user=%s: %v\n", username, err)
		return err
	}

	log.Printf("[INFO] AuthService.SignUp: user=%s created successfully\n", username)
	return nil
}

func (s *AuthService) SignIn(username, password string) (string, error) {
	log.Printf("[INFO] AuthService.SignIn: authenticating user=%s\n", username)
	var user models.User
	err := s.users.FindOne(context.Background(), bson.M{"username": username}).Decode(&user)
	if err != nil {
		log.Printf("[ERROR] AuthService.SignIn: user=%s not found: %v\n", username, err)
		return "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		log.Printf("[WARN] AuthService.SignIn: invalid password for user=%s\n", username)
		return "", errors.New("invalid credentials")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenStr, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		log.Printf("[ERROR] AuthService.SignIn: token generation failed for user=%s: %v\n", username, err)
		return "", err
	}

	log.Printf("[INFO] AuthService.SignIn: user=%s authenticated successfully\n", username)
	return tokenStr, nil
}
