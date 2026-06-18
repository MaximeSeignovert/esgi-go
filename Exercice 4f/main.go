package main

import (
	"context"
	"fmt"
	"time"
)

func effectuerOperationLongue(ctx context.Context, id string) error {
	return effectuerOperationLongueAvecParametres(ctx, id, 5, 500*time.Millisecond)
}

func effectuerOperationLongueAvecParametres(ctx context.Context, id string, nbEtapes int, dureeEtape time.Duration) error {
	fmt.Printf("[%s] Debut de l'operation...\n", id)

	for etape := 1; etape <= nbEtapes; etape++ {
		select {
		case <-ctx.Done():
			fmt.Printf("[%s] Operation annulee : %v\n", id, ctx.Err())
			return ctx.Err()
		case <-time.After(dureeEtape):
			fmt.Printf("[%s] Traitement etape %d/%d...\n", id, etape, nbEtapes)
		}
	}

	fmt.Printf("[%s] Operation terminee avec succes.\n", id)
	return nil
}

func lancerOperationAvecTimeout(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	resultChan := make(chan error, 1)

	go func() {
		resultChan <- effectuerOperationLongue(ctx, "Tache 1")
	}()

	select {
	case err := <-resultChan:
		if err != nil {
			fmt.Printf("Main: l'operation s'est terminee avec une erreur : %v\n", err)
			return err
		}

		fmt.Println("Main: l'operation s'est terminee avec succes avant le timeout.")
		return nil
	case <-ctx.Done():
		fmt.Printf("Main: timeout atteint ou annulation : %v\n", ctx.Err())
		return ctx.Err()
	}
}

func main() {
	fmt.Println("Demarrage du programme principal.")

	err := lancerOperationAvecTimeout(2 * time.Second)
	if err != nil {
		fmt.Printf("Resultat final : %v\n", err)
	}

	fmt.Println("Fin du programme principal.")
}
