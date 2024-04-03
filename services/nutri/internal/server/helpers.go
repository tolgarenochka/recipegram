package server

// Todo: переделать выходхные знаяения в структуру
func countNutriPerWeight(ingredients []Ingredient) (float64, float64, float64, float64, float64) {
	var totalWeight, totalCarbs, totalFats, totalProteins, totalCcal float64

	for _, ingredient := range ingredients {
		totalWeight += ingredient.Weight
		k := ingredient.Weight / 100
		totalFats += ingredient.Nutrient.Fats * k
		totalCarbs += ingredient.Nutrient.Carbohydrates * k
		totalProteins += ingredient.Nutrient.Proteins * k
		totalCcal += ingredient.Nutrient.Calories * k
	}

	return totalWeight, totalCarbs, totalFats, totalProteins, totalCcal
}
