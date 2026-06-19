package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"chmate/core"

	"github.com/spf13/cobra"
)

var (
	dirPath    string
	createTime string
	accessTime string
	modifyTime string
)

// 使用cobra实现chmate的命令行工具 chmate 的参数包括 -d --dir 文件所在目录的路径 -ct --create-time 创建时间 -at --access-time 访问时间 -mt --modify-time 修改时间
//其中 -d --dir 文件所在目录的路径 -ct --create-time 为必填参数
//-at --access-time 访问时间 -mt --modify-time 修改时间 为可选参数 如任意一项不填写 则与创建时间相同
// 比如我输入 chmate -d "C:\Users\Username\Documents" -ct "20230101120000" (我输入的时间格式一定为yyyymmddhhmmss)
// 就代表将 "C:\Users\Username\Documents" 目录下的所有文件的创建时间修改为 2023 年 1 月 1 日 12:00:00

var rootCmd = &cobra.Command{
	Use:   "chmate",
	Short: "修改文件元数据时间的命令行工具",
	Long: `chmate 是一个跨平台的文件元数据时间修改工具。

支持修改文件的创建时间、访问时间和修改时间。
Windows 平台支持修改所有时间，Linux/macOS 仅支持修改访问和修改时间。`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// 验证必填参数
		if dirPath == "" {
			return fmt.Errorf("请指定文件目录路径 (-d 或 --dir)")
		}
		if createTime == "" {
			return fmt.Errorf("请指定创建时间 (-ct 或 --create-time)")
		}

		// 检查目录是否存在
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			return fmt.Errorf("目录不存在: %s", dirPath)
		}

		// 解析时间格式 (yyyymmddhhmmss)
		layout := "20060102150405"
		ctime, err := time.ParseInLocation(layout, createTime, time.Local)
		if err != nil {
			return fmt.Errorf("创建时间格式错误: %v (期望格式: yyyymmddhhmmss，例如: 20230101120000)", err)
		}

		// 如果未指定访问时间和修改时间，则使用创建时间
		atime := ctime
		mtime := ctime

		if accessTime != "" {
			atime, err = time.ParseInLocation(layout, accessTime, time.Local)
			if err != nil {
				return fmt.Errorf("访问时间格式错误: %v (期望格式: yyyymmddhhmmss，例如: 20230101120000)", err)
			}
		}

		if modifyTime != "" {
			mtime, err = time.ParseInLocation(layout, modifyTime, time.Local)
			if err != nil {
				return fmt.Errorf("修改时间格式错误: %v (期望格式: yyyymmddhhmmss，例如: 20230101120000)", err)
			}
		}

		// 获取目录中的所有文件
		files, err := os.ReadDir(dirPath)
		if err != nil {
			return fmt.Errorf("读取目录失败: %v", err)
		}

		if len(files) == 0 {
			fmt.Println("目录中没有文件")
			return nil
		}

		fmt.Printf("找到 %d 个文件/目录\n", len(files))
		fmt.Println("正在处理...")

		successCount := 0
		failCount := 0

		for _, file := range files {
			// 跳过子目录，只处理文件
			if file.IsDir() {
				continue
			}

			filePath := filepath.Join(dirPath, file.Name())
			times := &core.FileTimes{
				CreateTime: ctime,
				AccessTime: atime,
				ModifyTime: mtime,
			}

			if err := core.SetFileTimes(filePath, times); err != nil {
				fmt.Printf("❌ 失败: %s - %v\n", file.Name(), err)
				failCount++
			} else {
				fmt.Printf("✅ 成功: %s\n", file.Name())
				successCount++
			}
		}

		fmt.Printf("\n处理完成！\n")
		fmt.Printf("成功: %d 个文件\n", successCount)
		fmt.Printf("失败: %d 个文件\n", failCount)

		return nil
	},
}

func init() {
	// 定义命令行参数
	rootCmd.Flags().StringVarP(&dirPath, "dir", "d", "", "文件所在目录的路径 (必填)")
	rootCmd.Flags().StringVarP(&createTime, "create-time", "c", "", "创建时间，格式: yyyymmddhhmmss，例如: 20230101120000 (必填)")
	rootCmd.Flags().StringVarP(&accessTime, "access-time", "a", "", "访问时间，格式: yyyymmddhhmmss，例如: 20230101120000 (可选，默认为创建时间)")
	rootCmd.Flags().StringVarP(&modifyTime, "modify-time", "m", "", "修改时间，格式: yyyymmddhhmmss，例如: 20230101120000 (可选，默认为创建时间)")

	// 标记必填参数
	rootCmd.MarkFlagRequired("dir")
	rootCmd.MarkFlagRequired("create-time")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
