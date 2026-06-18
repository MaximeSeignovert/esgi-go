package main

import "fmt"

func main() {
	var nomUtilisateur string = "Alice"
	var ageUtilisateur int = 28
	var estConnecte bool = true
	var soldeCompte float64 = 1520.75

	fmt.Println("Nom utilisateur :", nomUtilisateur)
	fmt.Println("Age utilisateur :", ageUtilisateur)
	fmt.Println("Est connecte :", estConnecte)
	fmt.Println("Solde compte :", soldeCompte)

	villeResidence := "Paris"
	codePostal := 75001
	tauxRemise := 12.5

	fmt.Printf("Ville de residence : %v (type : %T)\n", villeResidence, villeResidence)
	fmt.Printf("Code postal : %v (type : %T)\n", codePostal, codePostal)
	fmt.Printf("Taux de remise : %v (type : %T)\n", tauxRemise, tauxRemise)

	var nombreCommandes, nombreMessages, nombreNotifications int = 3, 7, 2
	fmt.Println("Commandes :", nombreCommandes)
	fmt.Println("Messages :", nombreMessages)
	fmt.Println("Notifications :", nombreNotifications)

	const (
		StatutEnAttente = iota
		StatutValide
		StatutRefuse
	)

	fmt.Println("Statut en attente :", StatutEnAttente)
	fmt.Println("Statut valide :", StatutValide)
	fmt.Println("Statut refuse :", StatutRefuse)

	prixArticle := 49
	tauxTVA := 0.2
	prixTTC := float64(prixArticle) + float64(prixArticle)*tauxTVA

	fmt.Println("Prix article :", prixArticle)
	fmt.Println("Taux TVA :", tauxTVA)
	fmt.Println("Prix TTC :", prixTTC)
}
