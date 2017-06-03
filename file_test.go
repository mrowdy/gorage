package gorage

import (
	"log"
	"testing"
)

func TestGetNameFromFile(t *testing.T) {
	f := File{Content: []byte("Test")}
	if f.CalculateHash() == "" {
		log.Fatal("Files should return a name")
	}
}

func TestDifferentFilesHAveDifferentNames(t *testing.T) {
	fA := File{Content: []byte("Test")}
	nameA := fA.CalculateHash()

	fB := File{Content: []byte("Bar")}
	nameB := fB.CalculateHash()

	if nameA == nameB {
		log.Fatal("Files should have different names")
	}
}

func TestSameFileSameName(t *testing.T) {
	fA := File{Content: []byte("Test")}
	nameA := fA.CalculateHash()

	fB := File{Content: []byte("Test")}
	nameB := fB.CalculateHash()

	if nameA != nameB {
		log.Fatal("Files should have same name")
	}
}

func TestNameShouldAlwaysBeTheSame(t *testing.T) {
	f := File{Content: []byte("Test")}
	nameA := f.CalculateHash()
	nameB := f.CalculateHash()

	if nameA != nameB {
		log.Fatal("Files should have same name")
	}
}
