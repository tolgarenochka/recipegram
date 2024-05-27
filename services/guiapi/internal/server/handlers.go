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

// Auth
// @Summary Auth user by login & password
// @Produce json
// @Param user body User true "User data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /auth [post]
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

// Register
// @Summary Register a new user
// @Produce json
// @Param user body UserRegistration true "User data"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /reg [post]
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

type Ingredient struct {
	Name   string  `json:"name"`
	Weight float64 `json:"weight"`
}

type Step struct {
	Step        int    `json:"step"`
	Instruction string `json:"instruction"`
}

// Recipe структура для представления данных рецепта
type Recipe struct {
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Ingredients []Ingredient `json:"ingredients"`
	Steps       []Step       `json:"steps"`
}

// AddRecipe
// @Summary Add a new recipe
// @Produce json
// @Param recipe body Recipe true "Recipe data"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /addRecipe [post]
// @Security ApiKeyAuth
func (s *Server) addRecipe(ctx *fasthttp.RequestCtx) {
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

	// Вставка рецепта в базу данных
	id, err := s.dbWizard.addRecipe(&recipeData, userID)
	if err != nil {
		ctx.Error("Internal Server Error", fasthttp.StatusInternalServerError)
		return
	}

	err = s.sendMessageKafka("count", strconv.Itoa(id))
	if err != nil {
		log.Printf("Error while sending Kafka message: %s", err.Error())
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

// EditRecipe
// @Summary Edit an existing recipe
// @Produce json
// @Param recipeID path int true "Recipe ID"
// @Param recipe body Recipe true "Recipe data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /editRecipe/{recipeID} [post]
// @Security ApiKeyAuth
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

	err = s.sendMessageKafka("count", strconv.Itoa(recipeID))
	if err != nil {
		log.Printf("Error while sending Kafka message: %s", err.Error())
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

// DeleteRecipe
// @Summary Delete an existing recipe
// @Produce json
// @Param recipeID path int true "Recipe ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /deleteRecipe/{recipeID} [delete]
// @Security ApiKeyAuth
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

// GetRecipe
// @Summary Get a recipe by ID
// @Produce json
// @Param recipeID path int true "Recipe ID"
// @Success 200 {object} Recipe
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /getRecipe/{recipeID} [get]
// @Security ApiKeyAuth
func (s *Server) getRecipe(ctx *fasthttp.RequestCtx) {
	_, _, err := validateToken(ctx)
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

	// Запрос к базе данных для получения рецепта
	recipe, err := s.dbWizard.getRecipeById(recipeID)
	if err != nil {
		ctx.Error("Internal Server Error", fasthttp.StatusInternalServerError)
		return
	}

	// Преобразование рецепта в JSON
	jsonResponse, err := json.Marshal(recipe)
	if err != nil {
		log.Printf("Error encoding recipe to JSON: %v\n", err)
		ctx.Error("Internal Server Error", fasthttp.StatusInternalServerError)
		return
	}

	// Успешное получение рецепта
	ctx.Response.SetStatusCode(fasthttp.StatusOK)
	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.Write(jsonResponse)
}

func (s *Server) getRecipesList(ctx *fasthttp.RequestCtx) {

}

// IngrNutritious структура для представления данных об ингредиенте
type IngrNutritious struct {
	Name                 string  `json:"name"`
	CaloriesPer100g      float64 `json:"calories_per_100g"`
	ProteinsPer100g      float64 `json:"proteins_per_100g"`
	FatsPer100g          float64 `json:"fats_per_100g"`
	CarbohydratesPer100g float64 `json:"carbohydrates_per_100g"`
}

// AddIngredient
// @Summary Add a new ingredient
// @Produce json
// @Param ingredient body IngrNutritious true "Ingredient data"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /addIngredient [post]
// @Security ApiKeyAuth
func (s *Server) addIngredient(ctx *fasthttp.RequestCtx) {
	_, _, err := validateToken(ctx)
	if err != nil {
		log.Printf("Error validating token: %v\n", err)
		ctx.Error("Unauthorized", fasthttp.StatusUnauthorized)
		return
	}

	// Получение данных об ингредиенте из тела запроса
	var ingredient IngrNutritious
	if err := json.Unmarshal(ctx.PostBody(), &ingredient); err != nil {
		log.Printf("Error decoding ingredient data: %v\n", err)
		ctx.Error("Bad Request", fasthttp.StatusBadRequest)
		return
	}

	// Добавление ингредиента в базу данных
	if err := s.dbWizard.addIngredient(&ingredient); err != nil {
		log.Printf("Error adding ingredient to the database: %v\n", err)
		ctx.Error("Internal Server Error", fasthttp.StatusInternalServerError)
		return
	}

	// Успешное добавление ингредиента
	ctx.SetStatusCode(fasthttp.StatusCreated)
	ctx.WriteString("Ingredient added successfully")
}
