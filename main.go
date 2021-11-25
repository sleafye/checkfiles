package main

import (
	"bufio"
	"crypto/md5"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

type fileDsc struct {
	// Name string
	// Path string
	// Md5           string
	SameFilesList []string
}

var oFilesMap map[string]*fileDsc = make(map[string]*fileDsc)
var tFilesMap map[string]*fileDsc = make(map[string]*fileDsc)
var chO chan int = make(chan int)

// var oFileNumber int = 0
// var tFileNumber int = 0

func main() {

	// runtime.GOMAXPROCS(3)
	// start := time.Now() // 获取当前时间

	fmt.Println("启动`````")

	// p := flag.String("p", "", "配置文件路径，使用json配置，如果不带任何参数会默认读执行文件同级文件夹下cfg.json")
	o := flag.String("o", "", "原文件或文件夹路径")
	t := flag.String("t", "", "对比目标文件或文件夹路径")

	flag.Parse()

	if *o == "" {
		fmt.Println("-o 原文件或文件夹路径未配置")
		return
	}

	if *t == "" {
		fmt.Println("-t 对比目标文件或文件夹路径未配置")
		return
	}

	if !isFileExist(*o) {
		fmt.Println("-o ：" + *o + ",路径文件或文件夹不存在")
		return
	}

	if !isFileExist(*t) {
		fmt.Println("-t ：" + *t + ",路径文件或文件夹不存在")
		return
	}

	go GoBuildFilesMap(*o, oFilesMap)
	go GoBuildFilesMap(*t, tFilesMap)
	<-chO
	<-chO

	// BuildFilesMap(*o, oFilesMap)
	// BuildFilesMap(*t, tFilesMap)

	count := 0

	for k, v := range oFilesMap {
		if m, ok := tFilesMap[k]; ok {
			// str := "\nmd5值: [" + k + "]\n"
			str := "工程文件:\n"
			for _, pt := range v.SameFilesList {
				str = str + "[" + pt + "]\n"
			}
			str = str + "相同文件:\n"

			for _, pt := range m.SameFilesList {
				str = str + "[" + pt + "]\n"
			}
			fmt.Println(str)
			count++
		}

	}

	fmt.Printf("[%s] 文件数量: %d\n", *o, len(oFilesMap))
	fmt.Printf("[%s] 文件数量: %d\n", *t, len(tFilesMap))

	fmt.Printf("相同数量: %d \n", count)

	// fmt.Println("该函数执行完成耗时：", time.Since(start))

	input := ""
	fmt.Printf("输入回车，关闭.....")
	fmt.Scanln(&input)

}

func run() {
	// for {
	// 	select{
	// 	case:
	// 	}
	// }
}

func GetFileMd5(file string) string {

	f, err := os.Open(file)
	if err != nil {
		return ""
	}
	defer f.Close()
	r := bufio.NewReader(f)

	h := md5.New()

	_, err = io.Copy(h, r)
	if err != nil {
		return ""
	}

	return fmt.Sprintf("%x", h.Sum(nil))
}

func GoBuildFilesMap(path string, filesMap map[string]*fileDsc) {
	BuildFilesMap(path, filesMap)
	chO <- 0
}

func BuildFilesMap(path string, filesMap map[string]*fileDsc) error {
	rd, err := ioutil.ReadDir(path)

	for _, fi := range rd {
		if fi.IsDir() {

			// fmt.Printf("dir: [%s]\n", filepath.Join(path, fi.Name()))
			BuildFilesMap(filepath.Join(path, fi.Name()), filesMap)
		} else {

			filePath := filepath.Join(path, fi.Name())
			md5Value := GetFileMd5(filePath)
			if fd, ok := filesMap[md5Value]; ok {
				fd.SameFilesList = append(fd.SameFilesList, filePath)
			} else {
				filesMap[md5Value] = &fileDsc{[]string{filePath}}
			}

			// fmt.Println("file: " + filePath + ", md5: " + md5Value)
		}
	}
	return err
}

func isFileExist(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}

	return !os.IsNotExist(err)

}
