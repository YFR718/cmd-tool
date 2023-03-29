package file_control

import (
	"fmt"
	"path/filepath"
	"testing"
)

func getFile() *File {
	f, err := NewFile("./tmp")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return f

}
func Test_Newfile(t *testing.T) {
	f, err := NewFile("./tmp/tmp2/dw")
	if err != nil {
		fmt.Println(-1, err)
		return
	}
	files, err := f.GetList()
	if err != nil {
		fmt.Println(-2, err)
	}
	fmt.Println(files)
}

func Test_demo2(t *testing.T) {
	fmt.Println(filepath.Dir("./rwerew/rwer/er/ee"))
}

func Test_zip(t *testing.T) {
	f := getFile()
	if f == nil {
		return
	}
	err := f.Zip()
	if err != nil {
		fmt.Println(err)
	}
}
func Test_unzip(t *testing.T) {
	f, err := NewFile("./tmp.zip")
	if f == nil {
		t.Error(err)
	}
	err = f.UnZip()
	if err != nil {
		fmt.Println(err)
	}
}

//func Test_Read(t *testing.T) {
//	f := getFile()
//	err := f.Read()
//	if err != nil {
//		fmt.Println(err)
//	}
//	fmt.Println(len(f.data))
//}
