package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type PersonneBasique struct {
	Nom   string
	Age   int
	Email string
	Actif bool
}

type Personne struct {
	Nom        string `json:"full_name"`
	Age        int    `json:"age_in_years"`
	Email      string `json:"contact_email,omitempty"`
	Actif      bool   `json:"is_active"`
	MotDePasse string `json:"-"`
}

type Produit struct {
	ID      int     `json:"product_id"`
	Nom     string  `json:"item_name"`
	Prix    float64 `json:"unit_price"`
	EnStock bool    `json:"in_stock"`
}

type Livre struct {
	ID               int      `json:"book_id"`
	Titre            string   `json:"title"`
	Auteur           string   `json:"author_name"`
	AnneePublication int      `json:"publication_year"`
	Genres           []string `json:"genres,omitempty"`
	ISBN             string   `json:"isbn_code,omitempty"`
	EstDisponible    bool     `json:"is_available"`
}

type Editeur struct {
	Nom     string `json:"name"`
	Adresse string `json:"location"`
}

type LivreAvecEditeur struct {
	Livre
	Editeur Editeur `json:"publisher_info"`
}

type UnixTime struct {
	time.Time
}

func (u UnixTime) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(u.Unix(), 10)), nil
}

func (u *UnixTime) UnmarshalJSON(data []byte) error {
	timestamp, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return err
	}

	u.Time = time.Unix(timestamp, 0).UTC()
	return nil
}

type LivreAvecDateAjout struct {
	Titre     string   `json:"title"`
	DateAjout UnixTime `json:"added_at"`
}

func main() {
	exercice1()
	exercice2()
	exercice3()
	exercice4()
	exercice5()
}

func exercice1() {
	fmt.Println("=== Exercice 1 : Serialisation basique ===")

	p := PersonneBasique{
		Nom:   "Alice Dupont",
		Age:   30,
		Email: "alice.dupont@example.com",
		Actif: true,
	}

	data, err := json.Marshal(p)
	if err != nil {
		fmt.Printf("Erreur de serialisation : %v\n\n", err)
		return
	}

	fmt.Println(string(data))
	fmt.Println("Reponse : les cles JSON correspondent aux noms exacts des champs exportes de la struct, car aucun tag json ne les remplace.")
	fmt.Println()
}

func exercice2() {
	fmt.Println("=== Exercice 2 : Struct tags ===")

	personneComplete := Personne{
		Nom:        "Alice Dupont",
		Age:        30,
		Email:      "alice.dupont@example.com",
		Actif:      true,
		MotDePasse: "secret",
	}
	personneSansEmail := Personne{
		Nom:        "Bob Martin",
		Age:        42,
		Email:      "",
		Actif:      false,
		MotDePasse: "motdepasse",
	}

	printJSON("Personne complete", personneComplete)
	printJSON("Personne sans email", personneSansEmail)
	fmt.Println("Reponse : omitempty supprime contact_email quand Email vaut une chaine vide.")
	fmt.Println("Reponse : MotDePasse n'apparait jamais grace au tag json:\"-\".")
	fmt.Println()
}

func exercice3() {
	fmt.Println("=== Exercice 3 : Deserialisation ===")

	jsonString := `{
		"product_id": 101,
		"item_name": "Clavier Mecanique",
		"unit_price": 79.99,
		"in_stock": true
	}`

	var produit Produit
	if err := json.Unmarshal([]byte(jsonString), &produit); err != nil {
		fmt.Printf("Erreur de deserialisation : %v\n\n", err)
		return
	}

	fmt.Printf("ID: %d\nNom: %s\nPrix: %.2f\nEn stock: %t\n", produit.ID, produit.Nom, produit.Prix, produit.EnStock)
	fmt.Println("Reponse : une cle inconnue comme description est ignoree par json.Unmarshal par defaut.")
	fmt.Println("Reponse : si unit_price vaut \"79.99\", json.Unmarshal retourne une erreur de type car une chaine ne peut pas remplir un float64.")
	fmt.Println()
}

func exercice4() {
	fmt.Println("=== Exercice 4 : Gestion des erreurs ===")

	malformedJSON := `{
		"product_id": 102,
		"item_name": "Souris Gaming",
		"unit_price": 49.99,
		"in_stock": true,
	`

	wrongTypeJSON := `{
		"product_id": "103",
		"item_name": "Ecran UltraWide",
		"unit_price": 399.99,
		"in_stock": true
	}`

	tryUnmarshalProduit("JSON malforme", malformedJSON)
	tryUnmarshalProduit("Type incorrect", wrongTypeJSON)
	fmt.Println("Reponse : verifier les erreurs evite de travailler avec des donnees invalides, incompletes ou silencieusement fausses.")
	fmt.Println()
}

func exercice5() {
	fmt.Println("=== Exercice 5 : Scenario complet Livre ===")

	livreComplet := Livre{
		ID:               1,
		Titre:            "Le Go en pratique",
		Auteur:           "Claire Martin",
		AnneePublication: 2026,
		Genres:           []string{"Programmation", "Backend"},
		ISBN:             "978-2-0000-0000-1",
		EstDisponible:    true,
	}
	livreSansOptionnels := Livre{
		ID:               2,
		Titre:            "JSON simplement",
		Auteur:           "Nadia Bernard",
		AnneePublication: 2025,
		Genres:           []string{},
		ISBN:             "",
		EstDisponible:    false,
	}

	jsonComplet := mustMarshalIndent(livreComplet)
	jsonSansOptionnels := mustMarshalIndent(livreSansOptionnels)

	fmt.Println("Livre complet :")
	fmt.Println(string(jsonComplet))
	fmt.Println("Livre sans genres ni ISBN :")
	fmt.Println(string(jsonSansOptionnels))

	var livreReluComplet Livre
	if err := json.Unmarshal(jsonComplet, &livreReluComplet); err != nil {
		fmt.Printf("Erreur de deserialisation livre complet : %v\n", err)
	} else {
		fmt.Printf("Livre relu complet : %+v\n", livreReluComplet)
	}

	var livreReluSansOptionnels Livre
	if err := json.Unmarshal(jsonSansOptionnels, &livreReluSansOptionnels); err != nil {
		fmt.Printf("Erreur de deserialisation livre sans optionnels : %v\n", err)
	} else {
		fmt.Printf("Livre relu sans optionnels : %+v\n", livreReluSansOptionnels)
	}

	fmt.Println("Reponse question 1 : ajouter une struct Editeur et un champ Editeur `json:\"publisher_info\"`, comme dans LivreAvecEditeur.")
	exempleEditeur := LivreAvecEditeur{
		Livre:   livreComplet,
		Editeur: Editeur{Nom: "ESGI Press", Adresse: "Paris"},
	}
	printJSON("Livre avec publisher_info", exempleEditeur)

	fmt.Println("Reponse question 2 : implementer json.Marshaler/json.Unmarshaler sur un type dedie, par exemple UnixTime.")
	livreDate := LivreAvecDateAjout{
		Titre:     "Livre date",
		DateAjout: UnixTime{Time: time.Date(2026, 6, 18, 15, 45, 0, 0, time.UTC)},
	}
	printJSON("Livre avec timestamp Unix", livreDate)
	fmt.Println()
}

func printJSON(label string, value any) {
	data, err := json.Marshal(value)
	if err != nil {
		fmt.Printf("%s: erreur de serialisation : %v\n", label, err)
		return
	}

	fmt.Printf("%s: %s\n", label, data)
}

func mustMarshalIndent(value any) []byte {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		panic(err)
	}

	return data
}

func tryUnmarshalProduit(label string, input string) {
	var produit Produit
	if err := json.Unmarshal([]byte(input), &produit); err != nil {
		fmt.Printf("%s -> erreur : %v\n", label, err)
		return
	}

	fmt.Printf("%s -> produit : %+v\n", label, produit)
}
