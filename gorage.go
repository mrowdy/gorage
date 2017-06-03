package gorage

type Gorage struct {
	Storage     Storage
	Persistance Persistance
}

func NewGorage(storage Storage, persistance Persistance) *Gorage {
	gorage := new(Gorage)
	gorage.Storage = storage
	gorage.Persistance = persistance
	return gorage
}

func (t Gorage) Save(file string) (*File, error) {
	f := new(File)
	f.Name = "foo"
	f.Content = []byte("Test")
	f.Hash = f.CalculateHash()

	if t.Persistance.HashExists(f.Hash) {
		return f, nil
	}

	err := t.Storage.Write(f)
	if err != nil {
		return f, err
	}

	return f, nil
}

func (t Gorage) Load(id string) (*File, error) {

	file, err := t.Persistance.Load(id)
	if err != nil {
		return new(File), err
	}

	content, err := t.Storage.Read(file.Hash)
	file.Content = content
	return &file, err
}
