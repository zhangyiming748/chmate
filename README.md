# chmate

一个跨平台的文件元数据时间修改工具，支持修改文件的**创建时间、修改时间、访问时间**。

## 功能特性

- ✅ **Windows**: 完全支持修改创建时间、访问时间、修改时间
- ⚠️ **Linux/macOS**: 支持修改访问时间和修改时间，创建时间通常不可修改
- 🎯 简单易用的命令行界面
- 🔍 查看和修改文件时间信息

## 安装

```bash
# 克隆仓库
git clone https://github.com/yourusername/chmate.git
cd chmate

# 编译
go build -o chmate

# Windows
go build -o chmate.exe
```

## 使用方法

### 命令行参数

```
-d, --dir              文件所在目录的路径 (必填)
-c, --create-time      创建时间，格式: 2006-01-02 15:04:05 (必填)
-a, --access-time      访问时间，格式: 2006-01-02 15:04:05 (可选，默认为创建时间)
-m, --modify-time      修改时间，格式: 2006-01-02 15:04:05 (可选，默认为创建时间)
```

### 1. 基本用法 - 批量修改目录下所有文件的时间

```bash
# 只指定创建时间，访问时间和修改时间自动与创建时间相同
chmate -d <目录路径> -c "<创建时间>"

# 示例：将 testdir 目录下所有文件的三个时间都设置为 2024-05-20 10:30:00
chmate -d testdir -c "2024-05-20 10:30:00"
```

### 2. 分别设置不同的时间

```bash
# 分别指定创建、访问、修改时间
chmate -d <目录路径> -c "<创建时间>" -a "<访问时间>" -m "<修改时间>"

# 示例：
chmate -d testdir -c "2024-01-01 08:00:00" -a "2024-06-15 12:30:00" -m "2024-12-31 23:59:59"
```

### 3. 查看帮助信息

```bash
chmate --help
```

输出：
```
chmate 是一个跨平台的文件元数据时间修改工具。

支持修改文件的创建时间、访问时间和修改时间。
Windows 平台支持修改所有时间，Linux/macOS 仅支持修改访问和修改时间。

Usage:
  chmate [flags]

Flags:
  -a, --access-time string   访问时间，格式: 2006-01-02 15:04:05 (可选，默认为创建时间)
  -c, --create-time string   创建时间，格式: 2006-01-02 15:04:05 (必填)
  -d, --dir string           文件所在目录的路径 (必填)
  -h, --help                 help for chmate
  -m, --modify-time string   修改时间，格式: 2006-01-02 15:04:05 (可选，默认为创建时间)
```

### 4. 使用示例

```bash
# 修改单个目录下所有文件的时间
chmate -d ./photos -c "2024-01-01 12:00:00"

# 修改并分别设置三个时间
chmate -d ./documents -c "2024-01-01 08:00:00" -a "2024-06-15 12:00:00" -m "2024-12-31 23:59:59"

# 使用完整参数名
chmate --dir ./backup --create-time "2024-05-20 10:30:00"
```

## 技术实现

在 Go 中修改文件**创建时间、修改时间、访问时间**，核心是使用 `os.Chtimes`，但**创建时间（Birth Time）** 受系统与 Go 版本限制，标准库无法跨平台修改，需调用系统 API。

下面分两种方案：

---

## 一、标准库：修改 修改时间 / 访问时间（全平台）
`os.Chtimes` 可以直接修改：
- **修改时间（Mtime）**
- **访问时间（Atime）**

```go
package main

import (
	"log"
	"os"
	"time"
)

func main() {
	filePath := "test.txt"

	// 要设置的时间
	newTime := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)

	// 修改 访问时间 和 修改时间
	err := os.Chtimes(filePath, newTime, newTime)
	if err != nil {
		log.Fatalf("修改时间失败: %v", err)
	}
	log.Println("成功修改 访问时间/修改时间")
}
```

参数顺序：
```go
os.Chtimes(文件路径, atime, mtime)
```

---

## 二、修改 创建时间（Birth Time）
Go 标准库**不支持**修改创建时间，必须调用系统原生 API。

### 1. Windows 下修改创建时间
使用 `syscall.SetFileTime`：

```go
package main

import (
	"log"
	"os"
	"syscall"
	"time"
)

func main() {
	filePath := "test.txt"

	// 目标时间
	t := time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC)
	winTime := syscall.NsecToFiletime(t.UnixNano())

	// 打开文件
	f, err := os.OpenFile(filePath, os.O_RDWR, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// 修改 创建时间、访问时间、修改时间
	err = syscall.SetFileTime(syscall.Handle(f.Fd()), &winTime, &winTime, &winTime)
	if err != nil {
		log.Fatalf("SetFileTime 失败: %v", err)
	}
	log.Println("成功修改 创建/访问/修改 时间")
}
```

### 2. Linux / macOS 下修改创建时间
- Linux：部分文件系统（ext4）支持 `statx` 查看 birth，但**修改**需要 `utimensat` 并带标志，或使用 `debugfs`（需 root）。
- macOS：使用 `setattrlist` / `utimes` 系统调用。

通用做法：**调用 exec 执行系统命令**
```go
// Linux 示例（需要 root + debugfs）
// touch -t 可改 mtime/atime，但改不了 crtime
```

---

## 三、跨平台完整封装（推荐）
```go
package main

import (
	"log"
	"os"
	"syscall"
	"time"
)

// SetFileTime 设置文件 创建/访问/修改 时间
// Windows 有效；Linux/macOS 仅能设置 atime/mtime，crtime 大多不可改
func SetFileTime(path string, ctime, atime, mtime time.Time) error {
	// 先改 atime/mtime（通用）
	if err := os.Chtimes(path, atime, mtime); err != nil {
		return err
	}

	// Windows 额外修改创建时间
	f, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	cft := syscall.NsecToFiletime(ctime.UnixNano())
	aft := syscall.NsecToFiletime(atime.UnixNano())
	mft := syscall.NsecToFiletime(mtime.UnixNano())

	return syscall.SetFileTime(syscall.Handle(f.Fd()), &cft, &aft, &mft)
}

func main() {
	t := time.Date(2024, 5, 20, 10, 0, 0, 0, time.Local)
	err := SetFileTime("test.txt", t, t, t)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("文件元数据时间修改完成")
}
```

---

## 四、项目结构

```
chmate/
├── core/
│   ├── changeMate_windows.go   # Windows 平台实现
│   └── changeMate_unix.go      # Linux/macOS 平台实现
├── main.go                      # 命令行入口
├── go.mod                       # Go 模块文件
└── README.md                    # 说明文档
```

### 核心 API

```go
// 获取文件时间信息
times, err := core.GetFileTimes(filePath)

// 设置文件时间
err := core.SetFileTimes(filePath, &core.FileTimes{
    CreateTime: createTime,
    AccessTime: accessTime,
    ModifyTime: modifyTime,
})

// 打印文件时间信息
err := core.PrintFileTimes(filePath)
```

---

## 五、说明
1. **创建时间（crtime/birthtime）**
   - Windows：可自由修改
   - Linux/macOS：**通常不允许普通用户修改**，很多文件系统只允许内核设置。
2. 权限：需要对文件有**写权限**。
3. 软链接：`os.Chtimes` 默认修改链接**指向的文件**，不修改链接本身。

