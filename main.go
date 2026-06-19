package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"chmate/core"
)

func main() {
	// 检查命令行参数
	if len(os.Args) < 2 {
		fmt.Println("用法: chmate <文件路径> [创建时间] [访问时间] [修改时间]")
		fmt.Println("时间格式: 2006-01-02 15:04:05")
		fmt.Println("\n示例:")
		fmt.Println("  chmate test.txt                                    # 查看文件时间信息")
		fmt.Println("  chmate test.txt 2024-01-01 12:00:00               # 设置所有时间为同一值")
		fmt.Println("  chmate test.txt 2024-01-01 12:00:00 2024-02-01 13:00:00 2024-03-01 14:00:00")
		os.Exit(1)
	}

	filePath := os.Args[1]

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Fatalf("文件不存在: %s", filePath)
	}

	// 如果只有文件路径，显示当前时间信息
	if len(os.Args) == 2 {
		fmt.Println("=== 当前文件时间信息 ===")
		if err := core.PrintFileTimes(filePath); err != nil {
			log.Fatal(err)
		}
		return
	}

	// 解析时间参数
	var times core.FileTimes
	layout := "2006-01-02 15:04:05"

	// 收集所有时间参数（跳过程序名和文件路径）
	timeArgs := os.Args[2:]

	switch len(timeArgs) {
	case 1:
		// 一个时间参数，应用到所有时间
		t, err := time.ParseInLocation(layout, timeArgs[0], time.Local)
		if err != nil {
			log.Fatalf("时间格式错误: %v", err)
		}
		times = core.FileTimes{
			CreateTime: t,
			AccessTime: t,
			ModifyTime: t,
		}
	case 3:
		// 三个时间参数：创建时间 访问时间 修改时间
		ctime, err := time.ParseInLocation(layout, timeArgs[0], time.Local)
		if err != nil {
			log.Fatalf("创建时间格式错误: %v", err)
		}
		atime, err := time.ParseInLocation(layout, timeArgs[1], time.Local)
		if err != nil {
			log.Fatalf("访问时间格式错误: %v", err)
		}
		mtime, err := time.ParseInLocation(layout, timeArgs[2], time.Local)
		if err != nil {
			log.Fatalf("修改时间格式错误: %v", err)
		}
		times = core.FileTimes{
			CreateTime: ctime,
			AccessTime: atime,
			ModifyTime: mtime,
		}
	default:
		log.Fatal("参数数量错误。请使用 1 个时间（设置所有时间）或 3 个时间（分别设置）")
	}

	// 设置文件时间
	fmt.Printf("正在修改文件时间...\n")
	if err := core.SetFileTimes(filePath, &times); err != nil {
		log.Fatal(err)
	}

	fmt.Println("\n=== 修改后的文件时间信息 ===")
	if err := core.PrintFileTimes(filePath); err != nil {
		log.Fatal(err)
	}

	fmt.Println("\n操作完成！")
}
