package main

import (
	"errors"
	"fmt"
)

func CalculerStatistiquesBase(nombres ...int) (int, int, float64) {
	somme := 0
	count := len(nombres)

	for _, nombre := range nombres {
		somme += nombre
	}

	if count == 0 {
		return 0, 0, 0
	}

	moyenne := float64(somme) / float64(count)
	return somme, count, moyenne
}

func CalculerStatistiquesCompletes(nombres ...float64) (float64, float64, float64, float64, int, error) {
	count := len(nombres)
	if count == 0 {
		return 0, 0, 0, 0, 0, errors.New("aucun argument fourni")
	}

	minimum := nombres[0]
	maximum := nombres[0]
	somme := 0.0

	for _, nombre := range nombres {
		if nombre < minimum {
			minimum = nombre
		}

		if nombre > maximum {
			maximum = nombre
		}

		somme += nombre
	}

	moyenne := somme / float64(count)
	return minimum, maximum, somme, moyenne, count, nil
}

func AnalyserDonneesCapteur(releves ...float64) (float64, float64, float64, int, int, error) {
	relevesValides := make([]float64, 0, len(releves))
	invalides := 0

	for _, releve := range releves {
		if releve > 0.0 && releve <= 100.0 {
			relevesValides = append(relevesValides, releve)
		} else {
			invalides++
		}
	}

	minimum, maximum, _, moyenne, count, err := CalculerStatistiquesCompletes(relevesValides...)
	if err != nil {
		return 0, 0, 0, 0, invalides, errors.New("aucun releve valide trouve")
	}

	return minimum, maximum, moyenne, count, invalides, nil
}

func afficherStatistiquesBase(nombres ...int) {
	somme, count, moyenne := CalculerStatistiquesBase(nombres...)
	fmt.Printf("Nombres: %v | Somme: %d, Count: %d, Moyenne: %.2f\n", nombres, somme, count, moyenne)
}

func afficherStatistiquesCompletes(nombres ...float64) {
	minimum, maximum, somme, moyenne, count, err := CalculerStatistiquesCompletes(nombres...)
	if err != nil {
		fmt.Printf("Nombres: %v | Erreur: %v\n", nombres, err)
		return
	}

	fmt.Printf("Nombres: %v | Min: %.2f, Max: %.2f, Somme: %.2f, Moyenne: %.2f, Count: %d\n", nombres, minimum, maximum, somme, moyenne, count)
}

func afficherAnalyseCapteur(releves ...float64) {
	minimum, maximum, moyenne, valides, invalides, err := AnalyserDonneesCapteur(releves...)
	if err != nil {
		fmt.Printf("Releves: %v | Erreur: %v, Valides: %d, Invalides: %d\n", releves, err, valides, invalides)
		return
	}

	fmt.Printf("Releves: %v | Temp Min: %.2f, Max: %.2f, Moyenne: %.2f, Valides: %d, Invalides: %d\n", releves, minimum, maximum, moyenne, valides, invalides)
}

func main() {
	fmt.Println("Exercice 1 : statistiques de base")
	afficherStatistiquesBase()
	afficherStatistiquesBase(10)
	afficherStatistiquesBase(10, 20, 30, 40)

	fmt.Println()
	fmt.Println("Exercice 2 : statistiques completes")
	afficherStatistiquesCompletes()
	afficherStatistiquesCompletes(7.5)
	afficherStatistiquesCompletes(1.5, 2.8, 0.7, 3.1)

	fmt.Println()
	fmt.Println("Exercice 3 : analyse de donnees de capteur")
	afficherAnalyseCapteur(22.5, 23.1, -5.0, 101.0, 21.9, 0.0, 24.0)
	afficherAnalyseCapteur(18.2, 19.7, 20.4)
	afficherAnalyseCapteur(-10.0, 105.0, 0.0)
}
