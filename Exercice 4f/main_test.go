package main

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestEffectuerOperationLongueRetourneNilSansAnnulation(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	err := effectuerOperationLongueAvecParametres(ctx, "test succes", 2, 5*time.Millisecond)
	if err != nil {
		t.Fatalf("effectuerOperationLongueAvecParametres() error = %v, want nil", err)
	}
}

func TestEffectuerOperationLongueRetourneErreurSurTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	err := effectuerOperationLongueAvecParametres(ctx, "test timeout", 5, 20*time.Millisecond)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("effectuerOperationLongueAvecParametres() error = %v, want %v", err, context.DeadlineExceeded)
	}
}

func TestEffectuerOperationLongueRetourneErreurSurAnnulation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := effectuerOperationLongueAvecParametres(ctx, "test annulation", 1, 5*time.Millisecond)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("effectuerOperationLongueAvecParametres() error = %v, want %v", err, context.Canceled)
	}
}
