package gorage

import (
	"log"
	"testing"
)

func TestGetNameFromFile(t *testing.T) {
	f := NewFile("file.png", []byte("Test"))
	if f.Hash() == "" {
		log.Fatal("Files should return a name")
	}
}

func TestDifferentFilesHAveDifferentNames(t *testing.T) {
	fA := NewFile("file.png", []byte("Test"))
	nameA := fA.Hash()

	fB := NewFile("file.png", []byte("Bar"))
	nameB := fB.Hash()

	if nameA == nameB {
		log.Fatal("Files should have different names")
	}
}

func TestSameFileSameName(t *testing.T) {
	fA := NewFile("file.png", []byte("Test"))
	nameA := fA.Hash()

	fB := NewFile("file.png", []byte("Test"))
	nameB := fB.Hash()

	if nameA != nameB {
		log.Fatal("Files should have same name")
	}
}

func TestNameShouldAlwaysBeTheSame(t *testing.T) {
	f := NewFile("file.png", []byte("Test"))
	nameA := f.Hash()
	nameB := f.Hash()

	if nameA != nameB {
		log.Fatal("Files should have same name")
	}
}
