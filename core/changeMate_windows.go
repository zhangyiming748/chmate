//go:build windows
// +build windows

package core

import (
	"fmt"
	"os"
	"syscall"
	"time"
)

// FileTimes 文件时间结构体
type FileTimes struct {
	CreateTime time.Time // 创建时间 (Birth Time/Crtime)
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
		ModifyTime: fileInfo.ModTime(),
	}

	// Windows 特定实现
	if stat, ok := fileInfo.Sys().(*syscall.Win32FileAttributeData); ok {
		times.CreateTime = time.Unix(0, stat.CreationTime.Nanoseconds())
		times.AccessTime = time.Unix(0, stat.LastAccessTime.Nanoseconds())
	} else {
		// 如果无法获取系统特定信息，使用修改时间作为默认值
		times.CreateTime = fileInfo.ModTime()
		times.AccessTime = fileInfo.ModTime()
	}

	return times, nil
}

// SetFileTimes 设置文件的元数据时间
// Windows: 可以修改创建时间、访问时间、修改时间
func SetFileTimes(filePath string, times *FileTimes) error {
	// 打开文件以获取句柄
	f, err := os.OpenFile(filePath, os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("打开文件失败: %w", err)
	}
	defer f.Close()

	// 转换时间为 Windows FILETIME 格式
	createTime := syscall.NsecToFiletime(times.CreateTime.UnixNano())
	accessTime := syscall.NsecToFiletime(times.AccessTime.UnixNano())
	modifyTime := syscall.NsecToFiletime(times.ModifyTime.UnixNano())

	// 调用 Windows API 设置所有时间
	if err := syscall.SetFileTime(
		syscall.Handle(f.Fd()),
		&createTime,
		&accessTime,
		&modifyTime,
	); err != nil {
		return fmt.Errorf("SetFileTime 失败: %w", err)
	}

	return nil
}

// PrintFileTimes 打印文件的时间信息
func PrintFileTimes(filePath string) error {
	times, err := GetFileTimes(filePath)
	if err != nil {
		return err
	}

	fmt.Printf("文件: %s\n", filePath)
	fmt.Printf("创建时间: %s\n", times.CreateTime.Format("2006-01-02 15:04:05"))
	fmt.Printf("访问时间: %s\n", times.AccessTime.Format("2006-01-02 15:04:05"))
	fmt.Printf("修改时间: %s\n", times.ModifyTime.Format("2006-01-02 15:04:05"))
	return nil
}
