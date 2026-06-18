package main

import (
	"fmt"
	"math/rand"
	"time"
)

func randomDuration(minSeconds int, maxSeconds int) time.Duration {
	return time.Duration(rand.Intn(maxSeconds-minSeconds+1)+minSeconds) * time.Second
}

func dataProducer(dataChannel chan<- string, done <-chan struct{}) {
	measurements := []string{
		"Temperature: 25 degC",
		"Humidite: 48%",
		"Pression: 1014 hPa",
		"Debit reseau: 72 Mb/s",
	}

	for {
		select {
		case <-done:
			return
		case <-time.After(randomDuration(1, 3)):
			measurement := measurements[rand.Intn(len(measurements))]
			select {
			case dataChannel <- measurement:
			case <-done:
				return
			}
		}
	}
}

func alertProducer(alertChannel chan<- string, done <-chan struct{}) {
	alerts := []string{
		"Niveau critique atteint",
		"Surchauffe detectee",
		"Perte de signal capteur",
	}

	for {
		select {
		case <-done:
			return
		case <-time.After(randomDuration(5, 10)):
			alert := alerts[rand.Intn(len(alerts))]
			select {
			case alertChannel <- alert:
			case <-done:
				return
			}
		}
	}
}

func stopAfter(delay time.Duration, quitChannel chan<- struct{}) {
	time.Sleep(delay)
	close(quitChannel)
}

func monitor(
	dataChannel <-chan string,
	alertChannel <-chan string,
	quitChannel <-chan struct{},
	statusTicker *time.Ticker,
	done chan<- struct{},
) {
	defer close(done)

	for {
		select {
		case data := <-dataChannel:
			fmt.Printf("[MESURE] %s\n", data)
		case alert := <-alertChannel:
			fmt.Printf("[ALERTE CRITIQUE] %s\n", alert)
		case <-statusTicker.C:
			fmt.Println("[STATUS] Verification systeme...")
		case <-quitChannel:
			fmt.Println("Signal d'arret recu. Arret du systeme.")
			return
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	dataChannel := make(chan string)
	alertChannel := make(chan string)
	quitChannel := make(chan struct{})
	done := make(chan struct{})
	statusTicker := time.NewTicker(2 * time.Second)
	defer statusTicker.Stop()

	go dataProducer(dataChannel, done)
	go alertProducer(alertChannel, done)
	go stopAfter(15*time.Second, quitChannel)

	fmt.Println("Systeme de surveillance demarre.")
	monitor(dataChannel, alertChannel, quitChannel, statusTicker, done)
}
