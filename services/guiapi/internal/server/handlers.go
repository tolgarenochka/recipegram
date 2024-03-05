package server

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/valyala/fasthttp"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type JWTClaims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

func (s *Server) auth(ctx *fasthttp.RequestCtx) {
	var userData User
	if err := json.Unmarshal(ctx.PostBody(), &userData); err != nil {
		log.Printf("Error decoding JSON: %v\n", err)
		ctx.Error("Invalid JSON", fasthttp.StatusBadRequest)
		return
	}

	// Поиск пользователя в базе данных по имени пользователя
	userID, err := s.dbWizard.auth(userData.Username, userData.Password)
	if err != nil {
		log.Printf("User not found: %v\n", err)
		ctx.Error("Invalid credentials", fasthttp.StatusUnauthorized)
		return
	}

	// Генерация JWT токена
	claims := &JWTClaims{
		UserID:   userID,
		Username: userData.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // Токен действителен в течение 24 часов
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("your-secret-key"))
	if err != nil {
		log.Printf("Error generating token: %v\n", err)
		ctx.Error("Internal Server Error", fasthttp.StatusInternalServerError)
		return
	}

	// Возврат токена в ответе
	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.Response.SetStatusCode(fasthttp.StatusOK)
	response := map[string]string{"token": tokenString}
	jsonResponse, _ := json.Marshal(response)
	ctx.Write(jsonResponse)
}

type UserRegistration struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s *Server) reg(ctx *fasthttp.RequestCtx) {
	var userData UserRegistration
	if err := json.Unmarshal(ctx.PostBody(), &userData); err != nil {
		log.Printf("Error decoding JSON: %v\n", err)
		ctx.Error("Invalid JSON", fasthttp.StatusBadRequest)
		return
	}

	// Хэширование пароля
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userData.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v\n", err)
		ctx.Error("Internal Server Error", fasthttp.StatusInternalServerError)
		return
	}

	err = s.dbWizard.reg(userData.Username, userData.Email, hashedPassword)
	if err != nil {
		ctx.Error("Internal Server Error", fasthttp.StatusInternalServerError)
		return
	}

	// Успешная регистрация
	ctx.Response.SetStatusCode(fasthttp.StatusCreated)
	ctx.Response.Header.Set("Content-Type", "application/json")
	response := map[string]string{"message": "User registered successfully"}
	jsonResponse, _ := json.Marshal(response)
	ctx.Write(jsonResponse)
}

func (s *Server) addRecipe(ctx *fasthttp.RequestCtx) {

}

func (s *Server) editRecipe(ctx *fasthttp.RequestCtx) {

}

func (s *Server) deleteRecipe(ctx *fasthttp.RequestCtx) {

}

func (s *Server) getRecipe(ctx *fasthttp.RequestCtx) {

}

func (s *Server) getRecipesList(ctx *fasthttp.RequestCtx) {

}
