package gorage

import (
	"log"
	"testing"
)

func TestGetNameFromFile(t *testing.T) {
	f := FileContent("Test")
	if calculateHash(f) == "" {
		log.Fatal("Files should return a name")
	}
}

func TestDifferentFilesHaveDifferentNames(t *testing.T) {
	fA := FileContent("Test")
	nameA := calculateHash(fA)

	fB := FileContent("Bar")
	nameB := calculateHash(fB)

	if nameA == nameB {
		log.Fatal("Files should have different names")
	}
}

func TestSameFileSameName(t *testing.T) {
	fA := FileContent("Test")
	nameA := calculateHash(fA)

	fB := FileContent("Test")
	nameB := calculateHash(fB)

	if nameA != nameB {
		log.Fatal("Files should have same name")
	}
}

func TestNameShouldAlwaysBeTheSame(t *testing.T) {
	f := FileContent("Test")
	nameA := calculateHash(f)
	nameB := calculateHash(f)

	if nameA != nameB {
		log.Fatal("Files should have same name")
	}
}
