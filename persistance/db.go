package persistance

import (
	"time"

	"github.com/Slemgrim/gorage"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type Db struct {
	Db *gorm.DB
}

func (db Db) Init() {
	db.Db.AutoMigrate(&FileContent{}, &File{})
}

func (db Db) Save(f *gorage.File) error {

	fileContent, found := db.getFileContent(f.Hash)

	if !found {
		fileContent = new(FileContent)
		fileContent.Hash = f.Hash
		fileContent.MimeType = f.MimeType
		fileContent.Size = f.Size

		result := db.Db.Create(fileContent)
		if result.Error != nil {
			return result.Error
		}
	}

	file := new(File)
	file.Name = f.Name
	file.FileContent = fileContent.ID

	result := db.Db.Create(file)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (db Db) Load(id string) (gorage.File, error) {
	return gorage.File{}, nil
}

func (db Db) HashExists(hash string) bool {
	fileContent := FileContent{}
	result := db.Db.Where(&FileContent{Hash: hash}).First(&fileContent)

	if result.Error != nil {
		return false
	}

	return true
}

func (db Db) getFileContent(hash string) (*FileContent, bool) {
	fileContent := new(FileContent)
	found := false

	result := db.Db.Where(&FileContent{Hash: hash}).First(&fileContent)
	if result.Error == nil {
		found = true
	}

	return fileContent, found
}

type FileContent struct {
	ID        string `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
	Hash      string
	MimeType  string
	Files     []File
	Size      int
}

func (fc *FileContent) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("ID", uuid.New().String())
	return nil
}

type File struct {
	ID          string `gorm:"primary_key"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
	Name        string
	FileContent string
}

func (fc *File) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("ID", uuid.New().String())
	return nil
}
