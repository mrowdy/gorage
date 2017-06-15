package gorage

import (
	"log"
	"testing"
)

func TestGetNameFromFile(t *testing.T) {
	f := FileContent("Test")
	if f.CalculateHash() == "" {
		log.Fatal("Files should return a name")
	}
}

func TestDifferentFilesHaveDifferentNames(t *testing.T) {
	fA := FileContent("Test")
	nameA := fA.CalculateHash()

	fB := FileContent("Bar")
	nameB := fB.CalculateHash()

	if nameA == nameB {
		log.Fatal("Files should have different names")
	}
}

func TestSameFileSameName(t *testing.T) {
	fA := FileContent("Test")
	nameA := fA.CalculateHash()

	fB := FileContent("Test")
	nameB := fB.CalculateHash()

	if nameA != nameB {
		log.Fatal("Files should have same name")
	}
}

func TestNameShouldAlwaysBeTheSame(t *testing.T) {
	f := FileContent("Test")
	nameA := f.CalculateHash()
	nameB := f.CalculateHash()

	if nameA != nameB {
		log.Fatal("Files should have same name")
	}
}
