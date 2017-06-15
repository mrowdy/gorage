Gorage
===

[![Build Status](https://travis-ci.org/Slemgrim/gorage.svg?branch=master)](https://travis-ci.org/Slemgrim/gorage)

A file IO abstraction which saves disk storage by separating file content and file meta data.

Usage
---

Save a file by providing a filename and its content as byte slice
``` golang
savedFile, err := gorage.Save(filepath.Base(filename.txt), []byte("the content"), nil)
```

Load a file by providing its ID which was generated during saving

``` golang
loadedFile, err := gorage.Load(savedFile.ID)
```

Delete a file by providing its ID which was generated during saving

``` golang
err := gorage.Delete(savedFile.ID)
```

File representation
---

``` golang

type File struct {
	ID         string
	Name       string
	Hash       string
	MimeType   string
	Content    FileContent
	Size       int
	Context    interface{}
	UploadedAt time.Time
	DeletedAt  *time.Time
}

```


Context
---

You can save some custom data as context alongside a file

``` golang
context = make(map[string]string)
context["foo"] = bar

savedFile, err := gorage.Save(filepath.Base(filename.txt), []byte("the content"), context)
```

Context is of type ```interface{}``` and can be everything
