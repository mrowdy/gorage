package test

import (
	"log"
	"testing"

	"github.com/slemgrim/treasury/treasury"
)

func TestGetNameFromFile(t *testing.T) {
	f := treasury.NewFile("file.png", []byte("Test"))
	if f.Hash() == "" {
		log.Fatal("Files should return a name")
	}
}

func TestDifferentFilesHAveDifferentNames(t *testing.T) {
	fA := treasury.NewFile("file.png", []byte("Test"))
	nameA := fA.Hash()

	fB := treasury.NewFile("file.png", []byte("Bar"))
	nameB := fB.Hash()

	if nameA == nameB {
		log.Fatal("Files should have different names")
	}
}

func TestSameFileSameName(t *testing.T) {
	fA := treasury.NewFile("file.png", []byte("Test"))
	nameA := fA.Hash()

	fB := treasury.NewFile("file.png", []byte("Test"))
	nameB := fB.Hash()

	if nameA != nameB {
		log.Fatal("Files should have same name")
	}
}

func TestNameShouldAlwaysBeTheSame(t *testing.T) {
	f := treasury.NewFile("file.png", []byte("Test"))
	nameA := f.Hash()
	nameB := f.Hash()

	if nameA != nameB {
		log.Fatal("Files should have same name")
	}
}
