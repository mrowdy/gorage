package gorage

type Gorage struct {
	Storage Storage
}

func NewGorage(storage Storage) *Gorage {
	gorage := new(Gorage)
	gorage.Storage = storage
	return gorage
}

func (t Gorage) Save(file string) (*File, error) {
	f := NewFile("foo", []byte("Test"))
	err := t.Storage.Write(f)
	if err != nil {
		return f, err
	}

	return f, nil
}

func (t Gorage) Load(file string) (*File, error) {
	loadedFile, err := t.Storage.Read(file)
	return &loadedFile, err
}
