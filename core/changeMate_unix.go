//go:build linux || darwin
// +build linux darwin

package core

import (
	"fmt"
	"os"
	"time"
)

// FileTimes 文件时间结构体
type FileTimes struct {
	CreateTime time.Time // 创建时间 (Birth Time/Crtime) - Linux/macOS 通常不可修改
	AccessTime time.Time // 访问时间 (Atime)
	ModifyTime time.Time // 修改时间 (Mtime)
}

// GetFileTimes 获取文件的元数据时间
func GetFileTimes(filePath string) (*FileTimes, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("获取文件信息失败: %w", err)
	}

	times := &FileTimes{
		CreateTime: fileInfo.ModTime(), // Linux/macOS 难以获取真实创建时间，使用修改时间代替
		AccessTime: fileInfo.ModTime(), // os.FileInfo 不提供 AccessTime，使用修改时间代替
		ModifyTime: fileInfo.ModTime(),
	}

	return times, nil
}

// SetFileTimes 设置文件的元数据时间
// Linux/macOS: 只能修改访问时间和修改时间，创建时间通常不可修改
func SetFileTimes(filePath string, times *FileTimes) error {
	// 使用标准库修改访问时间和修改时间
	if err := os.Chtimes(filePath, times.AccessTime, times.ModifyTime); err != nil {
		return fmt.Errorf("修改访问/修改时间失败: %w", err)
	}

	fmt.Println("注意: Linux/macOS 系统不支持修改文件创建时间")
	return nil
}

// PrintFileTimes 打印文件的时间信息
func PrintFileTimes(filePath string) error {
	times, err := GetFileTimes(filePath)
	if err != nil {
		return err
	}

	fmt.Printf("文件: %s\n", filePath)
	fmt.Printf("创建时间: %s (可能不准确)\n", times.CreateTime.Format("2006-01-02 15:04:05"))
	fmt.Printf("访问时间: %s\n", times.AccessTime.Format("2006-01-02 15:04:05"))
	fmt.Printf("修改时间: %s\n", times.ModifyTime.Format("2006-01-02 15:04:05"))
	return nil
}
