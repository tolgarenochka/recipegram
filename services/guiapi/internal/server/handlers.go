package server

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/valyala/fasthttp"
	"golang.org/x/crypto/bcrypt"
	"log"
	"strconv"
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

// Функция для проверки валидности JWT-токена
func validateToken(ctx *fasthttp.RequestCtx) (int, string, error) {
	tokenString := string(ctx.Request.Header.Peek("Authorization"))
	if tokenString == "" {
		return 0, "", fmt.Errorf("token not found")
	}

	claims := &JWTClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("your-secret-key"), nil
	})

	if err != nil {
		return 0, "", fmt.Errorf("error parsing token: %v", err)
	}

	if !token.Valid {
		return 0, "", fmt.Errorf("invalid token")
	}

	return claims.UserID, claims.Username, nil
}

// Recipe структура для представления данных рецепта
type Recipe struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Ingredients []string `json:"ingredients"`
	Steps       []struct {
		Step        int    `json:"step"`
		Instruction string `json:"instruction"`
	} `json:"steps"`
}

func (s *Server) addRecipe(ctx *fasthttp.RequestCtx) { // Проверка валидности токена
	userID, _, err := validateToken(ctx)
	if err != nil {
		log.Printf("Error validating token: %v\n", err)
		ctx.Error("Unauthorized", fasthttp.StatusUnauthorized)
		return
	}

	// Получение данных из тела запроса
	var recipeData Recipe
	if err := json.Unmarshal(ctx.PostBody(), &recipeData); err != nil {
		log.Printf("Error decoding JSON: %v\n", err)
		ctx.Error("Invalid JSON", fasthttp.StatusBadRequest)
		return
	}

	// Вставка рецепта в базу данных
	err = s.dbWizard.addRecipe(&recipeData, userID)
	if err != nil {
		ctx.Error("Internal Server Error", fasthttp.StatusInternalServerError)
		return
	}

	// Успешное добавление рецепта
	ctx.Response.SetStatusCode(fasthttp.StatusCreated)
	ctx.Response.Header.Set("Content-Type", "application/json")
	response := map[string]string{"message": "Recipe added successfully"}
	jsonResponse, _ := json.Marshal(response)
	ctx.Write(jsonResponse)
}

func (s *Server) editRecipe(ctx *fasthttp.RequestCtx) {
	// Проверка валидности токена
	userID, _, err := validateToken(ctx)
	if err != nil {
		log.Printf("Error validating token: %v\n", err)
		ctx.Error("Unauthorized", fasthttp.StatusUnauthorized)
		return
	}

	// Получение данных из тела запроса
	var recipeData Recipe
	if err := json.Unmarshal(ctx.PostBody(), &recipeData); err != nil {
		log.Printf("Error decoding JSON: %v\n", err)
		ctx.Error("Invalid JSON", fasthttp.StatusBadRequest)
		return
	}

	// Получение ID рецепта из пути запроса (/editRecipe/{recipeID})
	recipeID, err := strconv.Atoi(ctx.UserValue("recipeID").(string))
	if err != nil {
		log.Printf("Error parsing recipeID: %v\n", err)
		ctx.Error("Bad Request", fasthttp.StatusBadRequest)
		return
	}

	// Проверка, принадлежит ли рецепт пользователю
	ownerID, err := s.dbWizard.getUserIdFromRecipeId(recipeID)
	if err != nil {
		ctx.Error("Internal Server Error", fasthttp.StatusInternalServerError)
		return
	}

	if userID != ownerID {
		log.Println("Recipe does not belong to the user")
		ctx.Error("Forbidden", fasthttp.StatusForbidden)
		return
	}

	// Обновление данных рецепта в базе данных
	err = s.dbWizard.updateRecipe(&recipeData, recipeID)
	if err != nil {
		ctx.Error("Internal Server Error", fasthttp.StatusInternalServerError)
		return
	}

	// Успешное редактирование рецепта
	ctx.Response.SetStatusCode(fasthttp.StatusOK)
	ctx.Response.Header.Set("Content-Type", "application/json")
	response := map[string]string{"message": "Recipe updated successfully"}
	jsonResponse, _ := json.Marshal(response)
	ctx.Write(jsonResponse)
}

func (s *Server) deleteRecipe(ctx *fasthttp.RequestCtx) { // Проверка валидности токена
	userID, _, err := validateToken(ctx)
	if err != nil {
		log.Printf("Error validating token: %v\n", err)
		ctx.Error("Unauthorized", fasthttp.StatusUnauthorized)
		return
	}

	// Получение ID рецепта из пути запроса (/deleteRecipe/{recipeID})
	recipeID, err := strconv.Atoi(ctx.UserValue("recipeID").(string))
	if err != nil {
		log.Printf("Error parsing recipeID: %v\n", err)
		ctx.Error("Bad Request", fasthttp.StatusBadRequest)
		return
	}

	// Проверка, принадлежит ли рецепт пользователю
	ownerID, err := s.dbWizard.getUserIdFromRecipeId(recipeID)
	if err != nil {
		ctx.Error("Internal Server Error", fasthttp.StatusInternalServerError)
		return
	}

	if userID != ownerID {
		log.Println("Recipe does not belong to the user")
		ctx.Error("Forbidden", fasthttp.StatusForbidden)
		return
	}

	// Удаление рецепта из базы данных
	err = s.dbWizard.deleteRecipe(recipeID)
	if err != nil {
		ctx.Error("Internal Server Error", fasthttp.StatusInternalServerError)
		return
	}

	// Успешное удаление рецепта
	ctx.Response.SetStatusCode(fasthttp.StatusOK)
	ctx.Response.Header.Set("Content-Type", "application/json")
	response := map[string]string{"message": "Recipe deleted successfully"}
	jsonResponse, _ := json.Marshal(response)
	ctx.Write(jsonResponse)
}

func (s *Server) getRecipe(ctx *fasthttp.RequestCtx) {

}

func (s *Server) getRecipesList(ctx *fasthttp.RequestCtx) {

}
