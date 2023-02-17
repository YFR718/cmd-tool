package cloud_disk

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

// Zip 把路径的文件夹zip压缩
func Zip(src string) error {
	// 打开目标文件
	f, err := os.Create(src + ".zip")
	if err != nil {
		return err
	}
	defer f.Close()

	// 创建 zip.Writer
	w := zip.NewWriter(f)
	defer w.Close()

	// 遍历源目录中的所有文件
	err = filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 跳过目录
		if info.IsDir() {
			return nil
		}

		// 打开文件
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		// 计算相对路径
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		// 创建 zip.Header
		header := &zip.FileHeader{
			Name:   relPath,
			Method: zip.Deflate,
		}

		// 设置时间戳
		header.SetModTime(info.ModTime())

		// 写入文件
		writer, err := w.CreateHeader(header)
		if err != nil {
			return err
		}

		_, err = io.Copy(writer, file)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

//func Zip(path string) error {
//	// 1. Create a ZIP file and zip.Writer
//	f, err := os.Create(path + ".zip")
//	if err != nil {
//		return err
//	}
//	defer f.Close()
//
//	writer := zip.NewWriter(f)
//	defer writer.Close()
//
//	// 2. Go through all the files of the source
//	return filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
//		if err != nil {
//			return err
//		}
//
//		// 3. Create a local file header
//		header, err := zip.FileInfoHeader(info)
//		if err != nil {
//			return err
//		}
//
//		// set compression
//		header.Method = zip.Deflate
//
//		// 4. Set relative path of a file as the header name
//		header.Name, err = filepath.Rel(filepath.Dir(path), path)
//		if err != nil {
//			return err
//		}
//		if info.IsDir() {
//			header.Name += "/"
//		}
//
//		// 5. Create writer for the file header and save content of the file
//		headerWriter, err := writer.CreateHeader(header)
//		if err != nil {
//			return err
//		}
//
//		if info.IsDir() {
//			return nil
//		}
//
//		f, err := os.Open(path)
//		if err != nil {
//			return err
//		}
//		defer f.Close()
//
//		_, err = io.Copy(headerWriter, f)
//		return err
//	})
//}

// UnZip 解压缩文件
func UnZip(src, dst string) error {
	// 打开压缩文件
	zr, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer zr.Close()

	// 创建目录
	for _, file := range zr.File {
		if file.FileInfo().IsDir() {
			path := filepath.Join(dst, file.Name)
			if err := os.MkdirAll(path, file.Mode()); err != nil {
				return err
			}
		}
	}

	// 遍历 zr ，将文件写入到磁盘
	for _, file := range zr.File {
		path := filepath.Join(dst, file.Name)

		// 如果是目录，就跳过
		if file.FileInfo().IsDir() {
			continue
		}

		// 获取到 Reader
		fr, err := file.Open()
		if err != nil {
			return err
		}
		defer fr.Close()

		// 创建要写出的文件对应的 Write
		fw, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer fw.Close()

		// 将 Reader 的内容拷贝到 Writer
		_, err = io.Copy(fw, fr)
		if err != nil {
			return err
		}
	}

	return nil
}

//
//// 原地解压缩
//func UnZip(src, dst string) (err error) {
//	// 打开压缩文件，这个 zip 包有个方便的 ReadCloser 类型
//	// 这个里面有个方便的 OpenReader 函数，可以比 tar 的时候省去一个打开文件的步骤
//	zr, err := zip.OpenReader(src)
//	defer zr.Close()
//	if err != nil {
//		return
//	}
//
//	// 如果解压后不是放在当前目录就按照保存目录去创建目录
//	//if dst != "" {
//	//	if err := os.MkdirAll(dst, 0755); err != nil {
//	//		return err
//	//	}
//	//}
//
//	// 遍历 zr ，将文件写入到磁盘
//	for _, file := range zr.File {
//		path := filepath.Join(dst, file.Name)
//
//		// 如果是目录，就创建目录
//		if file.FileInfo().IsDir() {
//			if err := os.MkdirAll(path, file.Mode()); err != nil {
//				return err
//			}
//			// 因为是目录，跳过当前循环，因为后面都是文件的处理
//			continue
//		}
//
//		// 获取到 Reader
//		fr, err := file.Open()
//		if err != nil {
//			return err
//		}
//
//		// 创建要写出的文件对应的 Write
//		fw, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_TRUNC, file.Mode())
//		if err != nil {
//			return err
//		}
//
//		//n, err := io.Copy(fw, fr)
//		if err != nil {
//			return err
//		}
//
//		// 将解压的结果输出
//		//fmt.Printf("成功解压 %s ，共写入了 %d 个字符的数据\n", path, n)
//
//		// 因为是在循环中，无法使用 defer ，直接放在最后
//		// 不过这样也有问题，当出现 err 的时候就不会执行这个了，
//		// 可以把它单独放在一个函数中，这里是个实验，就这样了
//		fw.Close()
//		fr.Close()
//	}
//	return nil
//}
