package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func pauseAleatoire() time.Duration {
	return time.Duration(rand.Intn(451)+50) * time.Millisecond
}

func effectuerTacheSimple(id int) {
	fmt.Printf("Goroutine %d: Debut de la tache...\n", id)
	time.Sleep(pauseAleatoire())
	fmt.Printf("Goroutine %d: Tache terminee.\n", id)
}

func effectuerTacheAvecWaitGroup(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	effectuerTacheSimple(id)
}

func effectuerTacheAvecResultat(id int, resultChan chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Printf("Goroutine %d: Debut de la tache...\n", id)
	time.Sleep(pauseAleatoire())
	fmt.Printf("Goroutine %d: Tache terminee.\n", id)
	resultChan <- fmt.Sprintf("Goroutine %d a termine avec succes.", id)
}

func travailleur(id int, taches <-chan int, resultats chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()

	for idTache := range taches {
		fmt.Printf("Travailleur %d traite la tache %d...\n", id, idTache)
		time.Sleep(pauseAleatoire())
		resultats <- fmt.Sprintf("Travailleur %d a termine la tache %d.", id, idTache)
	}

	fmt.Printf("Travailleur %d arrete: plus de taches.\n", id)
}

func exercice1() {
	fmt.Println("Exercice 1 : lancement simple sans synchronisation")
	for id := 1; id <= 5; id++ {
		go effectuerTacheSimple(id)
	}

	fmt.Println("Toutes les goroutines lancees.")
	fmt.Println("Sans synchronisation, main peut se terminer avant les goroutines.")
	fmt.Println("Dans cette demonstration, une courte pause laisse volontairement le temps d'observer leurs sorties.")
	time.Sleep(700 * time.Millisecond)
}

func exercice2() {
	fmt.Println("Exercice 2 : synchronisation avec sync.WaitGroup")
	var wg sync.WaitGroup

	for id := 1; id <= 5; id++ {
		wg.Add(1)
		go effectuerTacheAvecWaitGroup(id, &wg)
	}

	fmt.Println("Toutes les goroutines lancees.")
	wg.Wait()
	fmt.Println("Toutes les goroutines ont termine leur execution.")
	fmt.Println("Le comportement change: main attend explicitement que chaque goroutine appelle Done.")
}

func exercice3() {
	fmt.Println("Exercice 3 : communication avec les canaux")
	var wg sync.WaitGroup
	resultChan := make(chan string, 5)

	for id := 1; id <= 5; id++ {
		wg.Add(1)
		go effectuerTacheAvecResultat(id, resultChan, &wg)
	}

	wg.Wait()
	close(resultChan)

	fmt.Println("Resultats recus depuis le canal :")
	for resultat := range resultChan {
		fmt.Println(resultat)
	}

	fmt.Println("Les resultats arrivent selon l'ordre de fin des goroutines, pas selon l'ordre des IDs.")
	fmt.Println("La duree de travail aleatoire rend donc l'ordre non deterministe.")
}

func exercice4() {
	fmt.Println("Exercice 4 : pool de travailleurs")
	const nombreTravailleurs = 3
	const nombreTaches = 10

	var wg sync.WaitGroup
	taches := make(chan int)
	resultats := make(chan string, nombreTaches)
	debut := time.Now()

	for id := 1; id <= nombreTravailleurs; id++ {
		wg.Add(1)
		go travailleur(id, taches, resultats, &wg)
	}

	go func() {
		for idTache := 1; idTache <= nombreTaches; idTache++ {
			taches <- idTache
		}
		close(taches)
	}()

	wg.Wait()
	close(resultats)

	fmt.Println("Resultats du pool :")
	for resultat := range resultats {
		fmt.Println(resultat)
	}

	fmt.Printf("Temps total avec %d travailleurs pour %d taches : %s\n", nombreTravailleurs, nombreTaches, time.Since(debut))
	fmt.Println("Plus il y a de travailleurs disponibles, plus plusieurs taches peuvent avancer en parallele.")
}

func main() {
	rand.Seed(time.Now().UnixNano())

	exercice1()
	fmt.Println()

	exercice2()
	fmt.Println()

	exercice3()
	fmt.Println()

	exercice4()
}
