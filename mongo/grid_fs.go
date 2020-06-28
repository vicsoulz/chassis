package mongo

import (
	"io"

	"github.com/globalsign/mgo"

	"github.com/vicsoulz/chassis/excel"
)

func GridFS(name string) *mgo.GridFS {
	return DB.GridFS(name)
}

func GridFSRead(fsName string, fileName string, writer io.Writer) (err error) {
	file, err := GridFS(fsName).Open(fileName)
	if err != nil {
		return
	}
	defer file.Close()

	_, err = io.Copy(writer, file)
	return
}

func GridFSCreate(fsName string, out io.Reader, createName string) error {
	fs := GridFS(fsName)
	file, err := fs.Create(createName)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, out)
	if err != nil {
		return err
	}
	return nil
}

func WriteExcelToGridFS(data interface{}, fsName string, createName string) error {
	e := excel.Excel{}
	err := e.FromSlice(data)
	if err != nil {
		return err
	}

	b, err := e.File.WriteToBuffer()
	if err != nil {
		panic(err)
	}

	return GridFSCreate(fsName, b, createName)
}

func WriteExcelToGridFSWithMap(data []map[string]interface{}, fsName string, createName string) error {
	e := excel.Excel{}
	err := e.FromSliceMap(data)
	if err != nil {
		return err
	}

	b, err := e.File.WriteToBuffer()
	if err != nil {
		panic(err)
	}

	return GridFSCreate(fsName, b, createName)
}
