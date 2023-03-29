package file_control

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
)

// ************************
// 文件类的功能
// 1. 文件操作：读取、创建/复写、删除、

// File 文件类
type File struct {
	Path  string
	Name  string
	Isdir bool
	data  []byte
}

// NewFile 文件类, 传入文件相对路径
func NewFile(path string) (*File, error) {
	// 获取完整路径
	root, _ := os.Getwd()
	path = root + "/data/" + path

	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if fileInfo.IsDir() {
		return &File{path, fileInfo.Name(), true, nil}, nil
	}
	return &File{path, fileInfo.Name(), false, nil}, nil
}

// GetList 获取文件列表
func (f *File) GetList() ([]File, error) {

	if !f.Isdir {
		return nil, errors.New("a file don`t have list! please choice a dir")
	}

	files, err := os.ReadDir(f.Path)
	if err != nil {
		return nil, err
	}
	list := make([]File, 0)
	for _, file := range files {
		list = append(list, File{"", file.Name(), file.IsDir(), nil})
	}
	return list, nil
}

// MakeDir 创建文件夹
func MakeDir(path string) error {
	root, _ := os.Getwd()
	path = root + "/data/" + path
	err := os.Mkdir(path, 0755)
	if err != nil {
		return err
	}
	return nil
}

// 删除文件、文件夹
func (f *File) Remove() error {
	err := os.RemoveAll(f.Path)
	if err != nil {
		return err
	}
	return nil
}

// Zip 把路径的文件夹zip压缩
func (f *File) Zip() error {
	mypath := f.Path
	println(mypath) // t1
	// 创建压缩文件
	zipFile, err := os.Create(mypath + ".zip")
	if err != nil {
		println(-1)
		log.Fatal(err)
		return err
	}
	defer zipFile.Close()

	// 创建一个 ZipWriter
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// 遍历文件，递归压缩
	err = filepath.Walk(mypath, walkFunc(mypath, zipWriter))

	if err != nil {
		println(-2)
		log.Fatal(err)
		return err
	}

	return nil
}
func walkFunc(rootPath string, zipWriter *zip.Writer) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if rootPath == path {
			fmt.Println("rootPath==path")
			return nil
		}
		//fmt.Println(rootPath)
		//fmt.Println(path, info.IsDir())

		// If a file is a symbolic link it will be skipped.
		if info.Mode()&os.ModeSymlink != 0 {
			return nil
		}

		// Create a local file header.
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		// 选择压缩算法
		header.Method = zip.Deflate
		// 得到子路径

		// 设置文件名
		header.Name, err = filepath.Rel(filepath.Dir(rootPath), path)
		if err != nil {
			return err
		}

		fmt.Println(header.Name)
		if info.IsDir() {
			header.Name += string(os.PathSeparator)
		}

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}
		// 文件夹就直接返回
		if info.IsDir() {
			return nil
		}
		// 打开文件并写入
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(writer, file)
		return err
	}
}

// UnZip 解压缩文件
func (f *File) UnZip() (err error) {

	zipPath := f.Path
	// 获取同级目录
	dstDir := filepath.Dir(zipPath)

	// 打开压缩文件
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer func(reader *zip.ReadCloser) {
		err := reader.Close()
		if err != nil {
			fmt.Println("close zip file error")
		}
	}(reader)

	for _, file := range reader.File {
		if err = unzipFile(file, dstDir); err != nil {
			return err
		}
	}

	return nil
}

func unzipFile(file *zip.File, dstDir string) error {
	// create the directory of file
	filePath := path.Join(dstDir, file.Name)
	if file.FileInfo().IsDir() {
		if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
			return err
		}
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return err
	}

	// open the file
	rc, err := file.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	// create the file
	w, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer w.Close()

	// save the decompressed file content
	_, err = io.Copy(w, rc)
	return err
}
