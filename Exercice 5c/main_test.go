package main

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestPersonneTagsOmitEmailAndPassword(t *testing.T) {
	personne := Personne{
		Nom:        "Bob Martin",
		Age:        42,
		Email:      "",
		Actif:      true,
		MotDePasse: "secret",
	}

	data, err := json.Marshal(personne)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	output := string(data)
	if !strings.Contains(output, `"full_name":"Bob Martin"`) {
		t.Fatalf("output = %s, want full_name key", output)
	}
	if strings.Contains(output, "contact_email") {
		t.Fatalf("output = %s, want contact_email omitted", output)
	}
	if strings.Contains(output, "MotDePasse") || strings.Contains(output, "secret") {
		t.Fatalf("output = %s, want password omitted", output)
	}
}

func TestProduitUnknownFieldIsIgnored(t *testing.T) {
	input := []byte(`{
		"product_id": 101,
		"item_name": "Clavier",
		"unit_price": 79.99,
		"in_stock": true,
		"description": "Champ inconnu"
	}`)

	var produit Produit
	if err := json.Unmarshal(input, &produit); err != nil {
		t.Fatalf("json.Unmarshal() error = %v, want nil", err)
	}

	if produit.ID != 101 || produit.Nom != "Clavier" {
		t.Fatalf("produit = %+v, want decoded known fields", produit)
	}
}

func TestProduitWrongTypeReturnsError(t *testing.T) {
	input := []byte(`{"product_id":101,"item_name":"Clavier","unit_price":"79.99","in_stock":true}`)

	var produit Produit
	if err := json.Unmarshal(input, &produit); err == nil {
		t.Fatal("json.Unmarshal() error = nil, want type error")
	}
}

func TestLivreOmitsEmptyGenresAndISBN(t *testing.T) {
	livre := Livre{
		ID:               2,
		Titre:            "JSON simplement",
		Auteur:           "Nadia Bernard",
		AnneePublication: 2025,
		Genres:           []string{},
		ISBN:             "",
		EstDisponible:    false,
	}

	data, err := json.Marshal(livre)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	output := string(data)
	if strings.Contains(output, "genres") {
		t.Fatalf("output = %s, want genres omitted", output)
	}
	if strings.Contains(output, "isbn_code") {
		t.Fatalf("output = %s, want isbn_code omitted", output)
	}
}

func TestUnixTimeMarshalsAsTimestamp(t *testing.T) {
	value := LivreAvecDateAjout{
		Titre:     "Livre date",
		DateAjout: UnixTime{Time: time.Unix(1781797500, 0).UTC()},
	}

	data, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	if !strings.Contains(string(data), `"added_at":1781797500`) {
		t.Fatalf("output = %s, want numeric Unix timestamp", data)
	}

	var decoded LivreAvecDateAjout
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	if !decoded.DateAjout.Equal(value.DateAjout.Time) {
		t.Fatalf("decoded time = %v, want %v", decoded.DateAjout, value.DateAjout)
	}
}
