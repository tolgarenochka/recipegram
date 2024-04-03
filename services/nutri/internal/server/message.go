package server

import (
	"log"
	"strconv"
)

func (s *Server) CountNutri(msg []byte) error {
	// Получение ID рецепта из сообщения Kafka
	recipeID, err := strconv.Atoi(string(msg))
	if err != nil {
		log.Printf("Error converting recipe ID: %v", err)
		return err
	}

	// Рассчет КБЖУ для рецепта и запись результатов в БД
	err = s.calculateAndStoreNutrition(recipeID)
	if err != nil {
		log.Printf("Error calculating nutrition for recipe %d: %v", recipeID, err)
		return err
	}

	return nil
}

func (s *Server) calculateAndStoreNutrition(recipeID int) error {
	ingredients, err := s.dbWizard.GetRecipeIngredients(recipeID)
	if err != nil {
		s.dbWizard.SetStatus(recipeID, "Ошибка при вычислении КБЖУ")
		return err
	}

	totalWeight, totalCarbs, totalFats, totalProteins, totalCcal := countNutriPerWeight(ingredients)

	err = s.dbWizard.StoreNutrition(recipeID, totalWeight, totalCcal, totalProteins, totalFats, totalCarbs, "")
	if err != nil {
		s.dbWizard.SetStatus(recipeID, "Ошибка при вычислении КБЖУ")
		return err
	}

	return nil
}
