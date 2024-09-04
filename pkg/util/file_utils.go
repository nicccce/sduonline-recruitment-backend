package util

import (
	"image"
	"image/png"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"sduonline-recruitment/pkg/conf"
)

// SaveUploadedFile 保存上传的文件到指定路径
func SaveUploadedFile(file *multipart.FileHeader, dst string) error {
	// 打开上传的文件
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()
	dst = filepath.Join(conf.Conf.UploadDir, dst)
	// 创建目标文件
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	// 将上传的文件内容复制到目标文件
	_, err = io.Copy(out, src)
	return err
}

// GetFileName 获取文件的文件名
func GetFileName(file *multipart.FileHeader) string {
	return filepath.Base(file.Filename)
}

// SaveAsPNG 保存图像为 PNG 格式
func SaveAsPNG(file *multipart.FileHeader, filePath string) error {
	// 打开上传的文件
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// 解码上传的图像
	img, _, err := image.Decode(src)
	if err != nil {
		return err
	}

	// 创建目标文件
	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	// 将图像编码为 PNG 并保存到目标文件
	err = png.Encode(out, img)
	if err != nil {
		return err
	}

	return nil
}
