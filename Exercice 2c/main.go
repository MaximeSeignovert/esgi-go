package main

import (
	"errors"
	"fmt"
	"math/rand"
	"sort"
	"time"
)

type Produit struct {
	ID        int
	Nom       string
	Prix      float64
	Categorie string
}

var totalPrixMesure float64
var produitRechercheMesure Produit

func categorieExiste(nom string, categories []string) bool {
	for _, categorie := range categories {
		if categorie == nom {
			return true
		}
	}

	return false
}

func supprimerCategorie(nom string, categories []string) []string {
	for index, categorie := range categories {
		if categorie == nom {
			return append(categories[:index], categories[index+1:]...)
		}
	}

	return categories
}

func afficherCategories(categories []string) {
	fmt.Println("Categories :", categories)
	fmt.Printf("Longueur: %d, Capacite: %d\n", len(categories), cap(categories))
}

func obtenirProduit(id int, inventaire map[int]Produit, stock map[int]int) (Produit, int, bool) {
	produit, existe := inventaire[id]
	if !existe {
		return Produit{}, 0, false
	}

	return produit, stock[id], true
}

func afficherProduits(inventaire map[int]Produit, stock map[int]int) {
	ids := idsTries(inventaire)
	for _, id := range ids {
		produit := inventaire[id]
		fmt.Printf("ID: %d | Nom: %s | Prix: %.2f euros | Categorie: %s | Stock: %d\n",
			produit.ID, produit.Nom, produit.Prix, produit.Categorie, stock[id])
	}
}

func vendreProduit(id int, quantite int, stock map[int]int) bool {
	if quantite <= 0 {
		return false
	}

	quantiteDisponible, existe := stock[id]
	if !existe || quantiteDisponible < quantite {
		return false
	}

	stock[id] -= quantite
	return true
}

func reapprovisionnerProduit(id int, quantite int, stock map[int]int) {
	if quantite > 0 {
		stock[id] += quantite
	}
}

func creerIndexParCategorie(inventaire map[int]Produit) map[string][]int {
	produitsParCategorie := make(map[string][]int)
	for id, produit := range inventaire {
		produitsParCategorie[produit.Categorie] = append(produitsParCategorie[produit.Categorie], id)
	}

	for categorie := range produitsParCategorie {
		sort.Ints(produitsParCategorie[categorie])
	}

	return produitsParCategorie
}

func listerProduitsParCategorie(categorie string, inventaire map[int]Produit, produitsParCategorie map[string][]int) {
	ids := produitsParCategorie[categorie]
	if len(ids) == 0 {
		fmt.Printf("Aucun produit trouve pour la categorie %q\n", categorie)
		return
	}

	fmt.Printf("Produits dans la categorie %q :\n", categorie)
	for _, id := range ids {
		produit := inventaire[id]
		fmt.Printf("- %s (ID %d, %.2f euros)\n", produit.Nom, produit.ID, produit.Prix)
	}
}

func trierProduitsParPrix(inventaire map[int]Produit, croissant bool) []Produit {
	produits := make([]Produit, 0, len(inventaire))
	for _, produit := range inventaire {
		produits = append(produits, produit)
	}

	sort.Slice(produits, func(i, j int) bool {
		if croissant {
			return produits[i].Prix < produits[j].Prix
		}

		return produits[i].Prix > produits[j].Prix
	})

	return produits
}

func valeurStockCategorie(categorie string, inventaire map[int]Produit, stock map[int]int) (float64, error) {
	total := 0.0
	trouve := false

	for id, produit := range inventaire {
		if produit.Categorie == categorie {
			total += produit.Prix * float64(stock[id])
			trouve = true
		}
	}

	if !trouve {
		return 0, errors.New("categorie introuvable")
	}

	return total, nil
}

func idsTries(inventaire map[int]Produit) []int {
	ids := make([]int, 0, len(inventaire))
	for id := range inventaire {
		ids = append(ids, id)
	}

	sort.Ints(ids)
	return ids
}

func genererProduit(id int, categories []string) Produit {
	return Produit{
		ID:        id,
		Nom:       fmt.Sprintf("Produit-%06d", id),
		Prix:      1 + rand.Float64()*499,
		Categorie: categories[rand.Intn(len(categories))],
	}
}

func mesurerAjoutProduits(nombre int, avecCapacite bool, categories []string) (map[int]Produit, map[int]int, time.Duration) {
	var inventaire map[int]Produit
	var stock map[int]int

	if avecCapacite {
		inventaire = make(map[int]Produit, nombre)
		stock = make(map[int]int, nombre)
	} else {
		inventaire = make(map[int]Produit)
		stock = make(map[int]int)
	}

	debut := time.Now()
	for id := 1; id <= nombre; id++ {
		inventaire[id] = genererProduit(id, categories)
		stock[id] = rand.Intn(250)
	}

	return inventaire, stock, time.Since(debut)
}

func mesurerRecherches(inventaire map[int]Produit, nombreRecherches int) time.Duration {
	debut := time.Now()
	for i := 0; i < nombreRecherches; i++ {
		id := rand.Intn(len(inventaire)) + 1
		produitRechercheMesure = inventaire[id]
	}

	return time.Since(debut)
}

func mesurerIteration(inventaire map[int]Produit, repetitions int) time.Duration {
	debut := time.Now()
	total := 0.0
	for i := 0; i < repetitions; i++ {
		for _, produit := range inventaire {
			total += produit.Prix
		}
	}

	totalPrixMesure = total

	return time.Since(debut)
}

func main() {
	rand.Seed(time.Now().UnixNano())

	fmt.Println("Partie 1 : Gestion des categories")
	categories := []string{"Electronique", "Vetements", "Livres"}
	categories = append(categories, "Maison", "Sport")
	afficherCategories(categories)

	fmt.Println("Categorie Maison existe ?", categorieExiste("Maison", categories))
	fmt.Println("Categorie Jardin existe ?", categorieExiste("Jardin", categories))

	categories = supprimerCategorie("Livres", categories)
	fmt.Println("Apres suppression de Livres :")
	afficherCategories(categories)

	categories = supprimerCategorie("Jardin", categories)
	fmt.Println("Apres tentative de suppression de Jardin :")
	afficherCategories(categories)
	fmt.Println("len represente le nombre d'elements actuellement stockes dans le slice.")
	fmt.Println("cap represente la taille du tableau interne disponible avant une nouvelle allocation.")

	fmt.Println()
	fmt.Println("Partie 2 : Gestion des produits et du stock")
	inventaireProduits := map[int]Produit{
		1: {ID: 1, Nom: "Clavier mecanique", Prix: 89.90, Categorie: "Electronique"},
		2: {ID: 2, Nom: "T-shirt coton", Prix: 19.99, Categorie: "Vetements"},
		3: {ID: 3, Nom: "Roman Go", Prix: 14.50, Categorie: "Livres"},
		4: {ID: 4, Nom: "Lampe de bureau", Prix: 34.90, Categorie: "Maison"},
	}
	stockProduits := map[int]int{
		1: 12,
		2: 40,
		3: 18,
		4: 9,
	}

	produitModifie := inventaireProduits[1]
	produitModifie.Prix = 79.90
	inventaireProduits[1] = produitModifie
	stockProduits[2] = 35
	afficherProduits(inventaireProduits, stockProduits)

	if produit, stock, existe := obtenirProduit(1, inventaireProduits, stockProduits); existe {
		fmt.Printf("Produit trouve : %s, stock %d\n", produit.Nom, stock)
	}
	if _, _, existe := obtenirProduit(99, inventaireProduits, stockProduits); !existe {
		fmt.Println("Produit 99 introuvable")
	}

	delete(inventaireProduits, 3)
	delete(stockProduits, 3)
	if _, _, existe := obtenirProduit(3, inventaireProduits, stockProduits); !existe {
		fmt.Println("Produit 3 supprime de l'inventaire et du stock")
	}

	fmt.Printf("Stock produit 1 avant vente: %d\n", stockProduits[1])
	fmt.Println("Vente de 3 produits 1 :", vendreProduit(1, 3, stockProduits))
	fmt.Printf("Stock produit 1 apres vente: %d\n", stockProduits[1])
	fmt.Println("Vente de 100 produits 1 :", vendreProduit(1, 100, stockProduits))
	reapprovisionnerProduit(1, 20, stockProduits)
	fmt.Printf("Stock produit 1 apres reapprovisionnement: %d\n", stockProduits[1])

	fmt.Println()
	fmt.Println("Partie 3 : Slices, maps et performance")
	produitsParCategorie := creerIndexParCategorie(inventaireProduits)
	listerProduitsParCategorie("Electronique", inventaireProduits, produitsParCategorie)
	listerProduitsParCategorie("Maison", inventaireProduits, produitsParCategorie)
	listerProduitsParCategorie("Jardin", inventaireProduits, produitsParCategorie)

	nombreProduits := 100000
	_, _, dureeSansCapacite := mesurerAjoutProduits(nombreProduits, false, categories)
	grandInventaire, _, dureeAvecCapacite := mesurerAjoutProduits(nombreProduits, true, categories)
	dureeRecherches := mesurerRecherches(grandInventaire, 10000)
	repetitionsIteration := 50
	dureeIteration := mesurerIteration(grandInventaire, repetitionsIteration)

	fmt.Printf("Ajout de %d produits sans capacite initiale: %s\n", nombreProduits, dureeSansCapacite)
	fmt.Printf("Ajout de %d produits avec capacite initiale: %s\n", nombreProduits, dureeAvecCapacite)
	fmt.Printf("10000 recherches aleatoires: %s\n", dureeRecherches)
	fmt.Printf("Iteration sur %d produits repetee %d fois: %s (total des prix: %.2f)\n", nombreProduits, repetitionsIteration, dureeIteration, totalPrixMesure)
	fmt.Println("Preallouer une map avec make(..., capacite) limite les reallocations internes pendant les insertions.")
	fmt.Println("Les recherches par cle sont tres rapides en moyenne, tandis que l'iteration parcourt naturellement tous les elements.")

	fmt.Println()
	fmt.Println("Bonus")
	fmt.Println("Produits tries par prix croissant :")
	for _, produit := range trierProduitsParPrix(inventaireProduits, true) {
		fmt.Printf("- %s : %.2f euros\n", produit.Nom, produit.Prix)
	}

	totalElectronique, err := valeurStockCategorie("Electronique", inventaireProduits, stockProduits)
	if err != nil {
		fmt.Println("Erreur valeur stock :", err)
	} else {
		fmt.Printf("Valeur totale du stock Electronique: %.2f euros\n", totalElectronique)
	}
}
