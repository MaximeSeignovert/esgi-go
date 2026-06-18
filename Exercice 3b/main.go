package main

import (
	"errors"
	"fmt"
)

type Notifier interface {
	Send(message string) error
}

type EmailNotifier struct {
	Recipient string
	Sender    string
}

func (e EmailNotifier) Send(message string) error {
	fmt.Printf("[EMAIL] De %s a %s : %s\n", e.Sender, e.Recipient, message)
	return nil
}

type SMSNotifier struct {
	PhoneNumber string
}

func (s SMSNotifier) Send(message string) error {
	fmt.Printf("[SMS] Envoi a %s : %s\n", s.PhoneNumber, message)
	return nil
}

type ConsoleNotifier struct{}

func (c ConsoleNotifier) Send(message string) error {
	fmt.Printf("[CONSOLE] Message : %s\n", message)
	return nil
}

type User struct {
	Name  string
	Email string
	Phone string
}

func processData(data interface{}) {
	switch valeur := data.(type) {
	case int:
		fmt.Printf("Donnee de type entier : %d\n", valeur)
	case string:
		fmt.Printf("Donnee de type chaine : %s\n", valeur)
	case bool:
		fmt.Printf("Donnee de type booleen : %t\n", valeur)
	case *EmailNotifier:
		fmt.Printf("Donnee de type EmailNotifier pour %s\n", valeur.Recipient)
	default:
		fmt.Printf("Type de donnee inconnu : %T\n", valeur)
	}
}

func sendSmartNotification(data interface{}, message string) error {
	switch valeur := data.(type) {
	case User:
		if valeur.Email != "" {
			notifier := EmailNotifier{
				Recipient: valeur.Email,
				Sender:    "support@boutique.test",
			}
			return notifier.Send(message)
		}

		if valeur.Phone != "" {
			notifier := SMSNotifier{PhoneNumber: valeur.Phone}
			return notifier.Send(message)
		}

		notifier := ConsoleNotifier{}
		return notifier.Send(fmt.Sprintf("Aucune methode de contact disponible pour %s", valeur.Name))
	case string:
		notifier := ConsoleNotifier{}
		return notifier.Send("Message generique : " + message)
	default:
		return fmt.Errorf("type non supporte pour une notification intelligente : %T", valeur)
	}
}

func envoyerNotification(notifier Notifier, message string) {
	if err := notifier.Send(message); err != nil {
		fmt.Println("Erreur d'envoi :", err)
	}
}

func testerSmartNotification(data interface{}, message string) {
	if err := sendSmartNotification(data, message); err != nil {
		fmt.Println("Erreur smart notifier :", err)
	}
}

func main() {
	fmt.Println("Partie 1 : utilisation polymorphique de Notifier")
	notifiers := []Notifier{
		EmailNotifier{Recipient: "alice@example.com", Sender: "admin@example.com"},
		SMSNotifier{PhoneNumber: "+33601020304"},
		ConsoleNotifier{},
	}

	for _, notifier := range notifiers {
		envoyerNotification(notifier, "Votre commande est prete.")
	}

	fmt.Println()
	fmt.Println("Partie 2 : interface vide et switch de type")
	processData(42)
	processData("Bonjour le monde")
	processData(true)
	processData(&EmailNotifier{Recipient: "bob@example.com", Sender: "admin@example.com"})
	processData([]int{1, 2, 3})
	processData(3.14)

	fmt.Println()
	fmt.Println("Partie 3 : smart notifier")
	testerSmartNotification(User{Name: "Alice", Email: "alice@example.com", Phone: "+33601020304"}, "Bienvenue Alice !")
	testerSmartNotification(User{Name: "Bob", Phone: "+33605060708"}, "Bienvenue Bob !")
	testerSmartNotification(User{Name: "Charlie"}, "Bienvenue Charlie !")
	testerSmartNotification("message libre", "Notification sans utilisateur.")
	testerSmartNotification(99, "Ce type ne peut pas etre notifie.")

	if err := exempleErreurExplicite(); err != nil {
		fmt.Println("Exemple d'erreur explicite :", err)
	}
}

func exempleErreurExplicite() error {
	return errors.New("les erreurs peuvent etre retournees puis traitees par l'appelant")
}

/*
Questions de reflexion :

1. L'implementation implicite des interfaces rend le code plus souple : un type satisfait une interface des qu'il possede les methodes attendues, sans declaration supplementaire.

2. interface{} est utile quand une fonction doit accepter des valeurs heterogenes, par exemple pour du logging, de la deserialisation ou une API tres generique. Son inconvenient est la perte de verification statique precise : il faut tester les types a l'execution et gerer les cas non supportes.

3. Les interfaces ameliorent la modularite car le code depend d'un comportement plutot que d'une implementation concrete. Elles facilitent aussi les tests : on peut remplacer un service reel par un faux objet qui implemente la meme interface.
*/
