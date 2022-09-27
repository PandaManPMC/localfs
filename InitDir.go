package localfs

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

var logger Log3 = nil
var tinyFileBasePath string
var smallFileBasePath string
var bigFileBasePath string

const defaultPerm fs.FileMode = 0664
const defaultDirPerm fs.FileMode = 0777

func InitDir(root string, log Log3) error {
	logger = log
	logger.Info(fmt.Sprintf("初始化目录%s", root))
	isExists, err := pathExists(root)
	if nil != err {
		return err
	}
	if !isExists {
		logger.Info("目录不存在 创建目录")
		err = os.MkdirAll(root, defaultDirPerm)
		if nil != err {
			return err
		}
	}

	tinyFileBasePath = fmt.Sprintf("%s%ca0", root, os.PathSeparator)
	smallFileBasePath = fmt.Sprintf("%s%ca1", root, os.PathSeparator)
	bigFileBasePath = fmt.Sprintf("%s%ca2", root, os.PathSeparator)
	logger.Info(fmt.Sprintf("读取目录列表 tiny=%s,small=%s,big=%s", tinyFileBasePath, smallFileBasePath, bigFileBasePath))

	createBaseDir(tinyFileBasePath)
	createBaseDir(smallFileBasePath)
	createBaseDir(bigFileBasePath)

	loadPath(tinyFileBasePath)
	loadPath(smallFileBasePath)
	loadPath(bigFileBasePath)
	return nil
}

func createBaseDir(dir string) {
	isE, err := pathExists(dir)
	if nil != err {
		logger.Error("createBaseDir", err)
		return
	}
	if !isE {
		logger.Info(fmt.Sprintf("列表没有初始目录%s，", dir))
		os.MkdirAll(dir, defaultDirPerm)
	}
}

func loadPath(path string) error {
	// 读取目录列表 检查目录下文件数量，目录下不超过1000个文件，超过1000个创建新的
	list, err := ioutil.ReadDir(path)
	if nil != err {
		return err
	}
	if 0 == len(list) {
		createDefaultDir(path)
		return nil
	}
	// 已经有目录 读取最后一个目录下的文件数量
	var info fs.FileInfo
	for i := len(list) - 1; i >= 0; i-- {
		info = list[i]
		if info.IsDir() {
			break
		}
	}
	if nil == info {
		// 没找到目录 创建
		createDefaultDir(path)
		return nil
	}
	// 找到目录 判断目录下文件数量
	dir := fmt.Sprintf("%s%c%s", path, os.PathSeparator, info.Name())
	list, err = ioutil.ReadDir(dir)
	if nil != err {
		return err
	}
	if 999 > len(list) {
		// 文件目录正常 继续使用
		putDirMap(path, &dirMapping{info.Name(), dir, getRelDir(dir), uint32(len(list))})
		return nil
	}
	createNewDir(info.Name(), path)
	return nil
}

func getRelDir(dir string) string {
	sl := strings.Split(dir, string(os.PathSeparator))
	return fmt.Sprintf("%s/%s", sl[len(sl)-2], sl[len(sl)-1])
}

func createNewDir(infoName, path string) string {
	// 文件已经超出数量 创建新目录
	num, _ := strconv.Atoi(infoName[1:])
	for {
		num++
		name := fmt.Sprintf("%s%d", infoName[:1], num)
		isE, _ := pathExists(fmt.Sprintf("%s%c%s", path, os.PathSeparator, name))
		if !isE {
			dir := fmt.Sprintf("%s%c%s", path, os.PathSeparator, name)
			logger.Info(fmt.Sprintf("新目录创建%s", dir))
			os.MkdirAll(dir, defaultDirPerm)
			putDirMap(path, &dirMapping{name, dir, getRelDir(dir), 0})
			return dir
		}
	}
}

func createDefaultDir(path string) {
	dir := fmt.Sprintf("%s%ca0", path, os.PathSeparator)
	os.MkdirAll(dir, defaultDirPerm)
	putDirMap(path, &dirMapping{"a0", dir, getRelDir(dir), 0})
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
