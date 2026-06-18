package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

const (
	nbGoroutines            = 100
	incrementsParGoroutine  = 1000
	resultatAttenduCompteur = nbGoroutines * incrementsParGoroutine
)

var (
	compteur int
	mu       sync.Mutex
)

func incrementerCompteurNonSynchro(wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 0; i < incrementsParGoroutine; i++ {
		compteur++
	}
}

func incrementerUneFoisSynchro() {
	mu.Lock()
	defer mu.Unlock()

	compteur++
}

func incrementerCompteurSynchro(wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 0; i < incrementsParGoroutine; i++ {
		incrementerUneFoisSynchro()
	}
}

func incrementerCompteurAtomic(compteurAtomic *int64, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 0; i < incrementsParGoroutine; i++ {
		atomic.AddInt64(compteurAtomic, 1)
	}
}

func executerGoroutines(worker func(*sync.WaitGroup)) time.Duration {
	var wg sync.WaitGroup
	debut := time.Now()

	for i := 0; i < nbGoroutines; i++ {
		wg.Add(1)
		go worker(&wg)
	}

	wg.Wait()
	return time.Since(debut)
}

func executerVersionNonSynchro() (int, time.Duration) {
	compteur = 0
	duration := executerGoroutines(incrementerCompteurNonSynchro)
	return compteur, duration
}

func executerVersionSynchro() (int, time.Duration) {
	compteur = 0
	duration := executerGoroutines(incrementerCompteurSynchro)
	return compteur, duration
}

func executerVersionAtomic() (int64, time.Duration) {
	var compteurAtomic int64
	var wg sync.WaitGroup
	debut := time.Now()

	for i := 0; i < nbGoroutines; i++ {
		wg.Add(1)
		go incrementerCompteurAtomic(&compteurAtomic, &wg)
	}

	wg.Wait()
	return compteurAtomic, time.Since(debut)
}

func afficherResultat(titre string, obtenu int64, duree time.Duration) {
	fmt.Println(titre)
	fmt.Printf("Resultat obtenu  : %d\n", obtenu)
	fmt.Printf("Resultat attendu : %d\n", resultatAttenduCompteur)
	fmt.Printf("Temps d'execution: %s\n", duree)
	fmt.Println()
}

func main() {
	resultatNonSynchro, dureeNonSynchro := executerVersionNonSynchro()
	afficherResultat("Etape 1 - Compteur non synchronise", int64(resultatNonSynchro), dureeNonSynchro)
	fmt.Println("Observation: cette version contient une condition de concurrence.")
	fmt.Println("Plusieurs goroutines peuvent lire la meme valeur, l'incrementer, puis ecraser le resultat des autres.")
	fmt.Println()

	resultatSynchro, dureeSynchro := executerVersionSynchro()
	afficherResultat("Etape 2 - Compteur protege par sync.Mutex", int64(resultatSynchro), dureeSynchro)
	fmt.Println("Observation: le mutex rend la section critique exclusive.")
	fmt.Println("Une seule goroutine modifie le compteur a la fois, donc le resultat final est coherent.")
	fmt.Println()

	resultatAtomic, dureeAtomic := executerVersionAtomic()
	afficherResultat("Optionnel - Compteur avec sync/atomic", resultatAtomic, dureeAtomic)
	fmt.Println("Observation: atomic.AddInt64 est adapte aux increments simples.")
	fmt.Println("Un mutex reste plus flexible quand la section critique contient plusieurs operations liees.")
}
