package storage_default

import (
	"fmt"
	"io"
	"os"
)

// RangeFileReader 封装了 SectionReader，同时支持 Closer
type RangeFileReader struct {
	file   *os.File
	reader *io.SectionReader
}

// NewRangeFileReader 创建 RangeFileReader，支持 Reader、Closer、Seeker
func NewRangeFileReader(filePath string, start, length int64) (*RangeFileReader, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	sectionReader := io.NewSectionReader(file, start, length)
	return &RangeFileReader{
		file:   file,
		reader: sectionReader,
	}, nil
}

// Read 读取数据
func (r *RangeFileReader) Read(p []byte) (int, error) {
	return r.reader.Read(p)
}

// Seek 移动读取位置
func (r *RangeFileReader) Seek(offset int64, whence int) (int64, error) {
	return r.reader.Seek(offset, whence)
}

// ReadAt 直接从指定位置读取
func (r *RangeFileReader) ReadAt(p []byte, off int64) (int, error) {
	return r.reader.ReadAt(p, off)
}

// Close 关闭文件
func (r *RangeFileReader) Close() error {
	return r.file.Close()
}

func main() {
	// 创建范围读取器（从偏移 10 读取 50 字节）
	rfr, err := NewRangeFileReader("example.txt", 10, 50)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer rfr.Close() // 确保关闭文件

	// 读取数据
	data, err := io.ReadAll(rfr)
	if err != nil {
		fmt.Println("Read error:", err)
		return
	}

	fmt.Println("Read Data:", string(data))
}
