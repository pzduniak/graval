// An example FTP server build on top of go-raval. graval handles the details
// of the FTP protocol, we just provide a basic in-memory persistence driver.
//
// If you're looking to create a custom graval driver, this example is a
// reasonable starting point. I suggest copying this file and changing the
// function bodies as required.
//
// USAGE:
//
//    go get github.com/yob/graval
//    go install github.com/yob/graval/graval-mem
//    ./bin/graval-mem
//
package main

import (
	"github.com/yob/graval"
	"io"
	"log"
	"os"
	"time"
)

const (
	fileOne = "This is the first file available for download.\n\nBy JĂ mes"
	fileTwo = "This is file number two.\n\n2012-12-04"
)

// A minimal driver for graval that stores everything in memory. The authentication
// details are fixed and the user is unable to upload, delete or rename any files.
//
// This really just exists as a minimal demonstration of the interface graval
// drivers are required to implement.
type MemDriver struct{}

func (driver *MemDriver) Authenticate(user string, pass string) bool {
	return user == "test" && pass == "1234"
}
func (driver *MemDriver) Bytes(path string) (bytes int) {
	switch path {
	case "/one.txt":
		bytes = len(fileOne)
		break
	case "/files/two.txt":
		bytes = len(fileTwo)
		break
	default:
		bytes = -1
	}
	return
}
func (driver *MemDriver) ModifiedTime(path string) (time.Time, error) {
	return time.Now(), nil
}
func (driver *MemDriver) ChangeDir(path string) bool {
	return path == "/" || path == "/files"
}
func (driver *MemDriver) DirContents(path string) (files []os.FileInfo) {
	files = []os.FileInfo{}
	switch path {
	case "/":
		files = append(files, graval.NewDirItem("files"))
		files = append(files, graval.NewFileItem("one.txt", len(fileOne)))
	case "/files":
		files = append(files, graval.NewFileItem("two.txt", len(fileOne)))
	}
	return files
}

func (driver *MemDriver) DeleteDir(path string) bool {
	return false
}
func (driver *MemDriver) DeleteFile(path string) bool {
	return false
}
func (driver *MemDriver) Rename(fromPath string, toPath string) bool {
	return false
}
func (driver *MemDriver) MakeDir(path string) bool {
	return false
}
func (driver *MemDriver) GetFile(path string, w io.Writer) bool {
	switch path {
	case "/one.txt":
		_, err := w.Write(fileOne)
		if err != nil {
			log.Print(err)
			return false
		}
	case "/files/two.txt":
		_, err := w.Write(fileTwo)
		if err != nil {
			log.Print(err)
			return false
		}
	}

	return true
}
func (driver *MemDriver) PutFile(destPath string, data io.Reader) bool {
	return false
}

// graval requires a factory that will create a new driver instance for each
// client connection. Generally the factory will be fairly minimal. This is
// a good place to read any required config for your driver.
type MemDriverFactory struct{}

func (factory *MemDriverFactory) NewDriver() (graval.FTPDriver, error) {
	return &MemDriver{}, nil
}

// it's alive!
func main() {
	factory := &MemDriverFactory{}
	ftpServer := graval.NewFTPServer(&graval.FTPServerOpts{Factory: factory})
	err := ftpServer.ListenAndServe()
	if err != nil {
		log.Print(err)
		log.Fatal("Error starting server!")
	}
}
