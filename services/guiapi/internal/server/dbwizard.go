package server

import (
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"log"
)

func (d *dbWiz) auth(username, password string) (int, error) {
	var storedPassword string
	var userID int

	err := d.dbWizard.QueryRow("SELECT user_id, password_hash FROM users WHERE username = $1", username).Scan(&userID, &storedPassword)
	if err != nil {
		log.Printf("No users with this username")
		return 0, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password))
	if err != nil {
		log.Printf("Invalid password: %v\n", err)
		return 0, err
	}

	return userID, nil
}

func (d *dbWiz) reg(username, email string, hashedPassword []byte) error {
	// Вставка пользователя в базу данных
	_, err := d.dbWizard.Exec("INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3)",
		username, email, hashedPassword)
	if err != nil {
		log.Printf("Error inserting user into the database: %v\n", err)
		return err
	}

	return nil
}

func (d *dbWiz) addRecipe(recipeData *Recipe, userID int) (int, error) {
	// Преобразование списка ингредиентов в формат JSON
	ingredientsJSON, err := json.Marshal(recipeData.Ingredients)
	if err != nil {
		log.Fatal(err)
	}

	// Преобразование списка шагов в формат JSON
	stepsJSON, err := json.Marshal(recipeData.Steps)
	if err != nil {
		log.Fatal(err)
	}

	var id int
	// Вставка рецепта в базу данных
	err = d.dbWizard.QueryRow("INSERT INTO recipes (title, description, user_id, ingredients, steps) VALUES ($1, $2, $3, $4, $5) RETURNING recipe_id",
		recipeData.Title, recipeData.Description, userID, ingredientsJSON, stepsJSON).Scan(&id)
	if err != nil {
		log.Printf("Error inserting recipe into the database: %v\n", err)
		return 0, err
	}

	return id, nil
}

func (d *dbWiz) getUserIdFromRecipeId(recipeID int) (int, error) {
	var ownerID int
	err := d.dbWizard.QueryRow("SELECT user_id FROM recipes WHERE recipe_id = $1", recipeID).Scan(&ownerID)
	if err != nil {
		log.Printf("Error checking recipe ownership: %v\n", err)
		return 0, err
	}

	return ownerID, nil
}

func (d *dbWiz) updateRecipe(recipeData *Recipe, recipeID int) error {
	// Преобразование списка ингредиентов в формат JSON
	ingredientsJSON, err := json.Marshal(recipeData.Ingredients)
	if err != nil {
		log.Fatal(err)
	}

	// Преобразование списка шагов в формат JSON
	stepsJSON, err := json.Marshal(recipeData.Steps)
	if err != nil {
		log.Fatal(err)
	}

	_, err = d.dbWizard.Exec("UPDATE recipes SET title = $1, description = $2, ingredients = $3, steps = $4 WHERE recipe_id = $5",
		recipeData.Title, recipeData.Description, ingredientsJSON, stepsJSON, recipeID)
	if err != nil {
		log.Printf("Error updating recipe in the database: %v\n", err)
		return err
	}

	return nil
}

func (d *dbWiz) deleteRecipe(recipeID int) error {
	_, err := d.dbWizard.Exec("DELETE FROM recipes WHERE recipe_id = $1", recipeID)
	if err != nil {
		log.Printf("Error deleting recipe from the database: %v\n", err)
		return err
	}

	return nil
}

func (d *dbWiz) getRecipeById(recipeID int) (Recipe, error) {
	var recipe Recipe
	var ingredientsJSON string
	var stepsJSON string

	err := d.dbWizard.QueryRow("SELECT title, description, ingredients, steps FROM recipes WHERE recipe_id = $1", recipeID).
		Scan(&recipe.Title, &recipe.Description, &ingredientsJSON, &stepsJSON)
	if err != nil {
		log.Printf("Error getting recipe from the database: %v\n", err)
		return recipe, err
	}

	// Декодирование JSON-строки ингредиентов в список структур
	if err := json.Unmarshal([]byte(ingredientsJSON), &recipe.Ingredients); err != nil {
		log.Printf("Error decoding ingredients from JSON: %v\n", err)
		return recipe, err
	}

	// Декодирование JSON-строки шагов в список структур
	if err := json.Unmarshal([]byte(stepsJSON), &recipe.Steps); err != nil {
		log.Printf("Error decoding steps from JSON: %v\n", err)
		return recipe, err
	}

	return recipe, nil
}
