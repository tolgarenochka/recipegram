package server

import "log"

func (d *dbWiz) SetStatus(recipeID int, reason string) error {
	_, err := d.dbWizard.Exec("UPDATE recipe_nutrition SET reason_why_not = $1 WHERE recipe_id = $2", reason, recipeID)
	if err != nil {
		log.Printf("Error setting status for recipe %d: %v\n", recipeID, err)
		return err
	}

	return nil
}

type Nutrient struct {
	Calories      float64
	Proteins      float64
	Fats          float64
	Carbohydrates float64
}

type Ingredient struct {
	Name     string
	Weight   float64
	Nutrient Nutrient
}

func (d *dbWiz) GetRecipeIngredients(recipeID int) ([]Ingredient, error) {
	var ingredients []Ingredient

	rows, err := d.dbWizard.Query("SELECT ingredients.name, (ingredient->>'weight')::float AS weight, ingredients.calories_per_100g, ingredients.proteins_per_100g, ingredients.fats_per_100g, ingredients.carbohydrates_per_100g FROM recipes, jsonb_array_elements(recipes.ingredients) AS ingredient JOIN ingredients ON ingredient->>'name' = ingredients.name WHERE recipes.recipe_id = $1;", recipeID)
	if err != nil {
		log.Printf("Error querying ingredients for recipe %d: %v\n", recipeID, err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var ingredient Ingredient
		var nutrient Nutrient
		err := rows.Scan(&ingredient.Name, &ingredient.Weight, &nutrient.Calories, &nutrient.Proteins, &nutrient.Fats, &nutrient.Carbohydrates)
		if err != nil {
			log.Printf("Error scanning ingredient rows: %v\n", err)
			return nil, err
		}
		ingredient.Nutrient = nutrient
		ingredients = append(ingredients, ingredient)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating over ingredient rows: %v\n", err)
		return nil, err
	}

	return ingredients, nil
}

func (d *dbWiz) StoreNutrition(recipeID int, weightTotal, caloriesTotal, proteinsTotal, fatsTotal, carbohydratesTotal float64, reasonWhyNot string) error {
	_, err := d.dbWizard.Exec("INSERT INTO recipe_nutrition (recipe_id, weight_total, calories_total, proteins_total, fats_total, carbohydrates_total, reason_why_not) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		recipeID, weightTotal, caloriesTotal, proteinsTotal, fatsTotal, carbohydratesTotal, reasonWhyNot)
	if err != nil {
		log.Printf("Error storing nutrition for recipe %d: %v\n", recipeID, err)
		return err
	}

	return nil
}
