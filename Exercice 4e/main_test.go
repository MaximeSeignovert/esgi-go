package main

import "testing"

func TestVersionSynchroRetourneLeResultatAttendu(t *testing.T) {
	got, _ := executerVersionSynchro()
	if got != resultatAttenduCompteur {
		t.Fatalf("executerVersionSynchro() = %d, want %d", got, resultatAttenduCompteur)
	}
}

func TestVersionAtomicRetourneLeResultatAttendu(t *testing.T) {
	got, _ := executerVersionAtomic()
	if got != int64(resultatAttenduCompteur) {
		t.Fatalf("executerVersionAtomic() = %d, want %d", got, resultatAttenduCompteur)
	}
}
