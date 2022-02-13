package ziputil

import (
	"archive/zip"
	"io"
	"os"
	"online/common/log"
	"online/common/utils"
	"path/filepath"
)

//压缩文件
//files 文件数组，可以是不同dir下的文件或者文件夹
//dest 压缩文件存放地址
func Compress(files []*os.File, dest string) error {
	d, _ := os.Create(dest)
	defer d.Close()
	w := zip.NewWriter(d)
	defer w.Close()
	for _, file := range files {
		err := compress(file, "", w)
		if err != nil {
			return err
		}
	}
	return nil
}

func compress(file *os.File, prefix string, zw *zip.Writer) error {
	info, err := file.Stat()
	if err != nil {
		return err
	}
	if info.IsDir() {
		prefix = prefix + "/" + info.Name()
		fileInfos, err := file.Readdir(-1)
		if err != nil {
			return err
		}
		for _, fi := range fileInfos {
			f, err := os.Open(file.Name() + "/" + fi.Name())
			if err != nil {
				return err
			}
			err = compress(f, prefix, zw)
			if err != nil {
				return err
			}
		}
	} else {
		header, err := zip.FileInfoHeader(info)
		header.Name = prefix + "/" + header.Name
		if err != nil {
			return err
		}
		writer, err := zw.CreateHeader(header)
		if err != nil {
			return err
		}
		_, err = io.Copy(writer, file)
		file.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

//解压
func DeCompress(zipFile, dest string) error {
	_ = os.MkdirAll(dest, os.ModePerm)

	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer reader.Close()
	for _, file := range reader.File {
		filename := filepath.Join(dest, file.Name)
		if file.FileInfo().IsDir() {
			err := os.MkdirAll(filename, os.ModePerm)
			if err != nil {
				return utils.Errorf("mkdir failed: %s")
			}
			continue
		}

		dirName := filepath.Dir(filename)
		err = os.MkdirAll(dirName, os.ModePerm)
		if err != nil {
			log.Errorf("mkdir [%s] failed: %s", dirName, err)
			return err
		}

		// 打开需要解压的文件
		rc, err := file.Open()
		if err != nil {
			return err
		}

		log.Infof("start to unzip: %s", file.Name)

		w, err := os.Create(filename)
		if err != nil {
			rc.Close()
			return err
		}

		_, err = io.Copy(w, rc)
		w.Close()
		rc.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
