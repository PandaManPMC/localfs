package localfs

import (
	"crypto/rand"
	"fmt"
	"io"
	"math/big"
	"mime/multipart"
	"os"
	"strings"
	"sync"
)

const size5M = 1024 * 1024 * 5
const size20M = 1024 * 1024 * 20

var chaArr = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}

//	UploadFile 文件上传
//	fileName string 文件名；size int64 文件大小；file multipart.File 流。
//	return string 文件本地路径；string 文件短地址； bool
func UploadFile(fileName string, size int64, file multipart.File) (string, string, bool) {
	defer file.Close()
	var basePath string
	if size < size5M {
		basePath = tinyFileBasePath
	} else if size < size20M {
		basePath = smallFileBasePath
	} else {
		basePath = bigFileBasePath
	}
	dm := getDirMap(basePath)

	dir := uploadPath(dm, basePath)
	format := fileName[strings.LastIndex(fileName, ".")+1:]

	var newFilePath, newFileName string
	for {
		newFileName = fmt.Sprintf("%s.%s", RandCharacterString(32), format)
		newFilePath = fmt.Sprintf("%s%c%s", dir, os.PathSeparator, newFileName)
		isE, _ := pathExists(newFilePath)
		if !isE {
			break
		}
	}

	toFile, err := os.OpenFile(newFilePath, os.O_RDWR|os.O_CREATE, defaultPerm)
	if nil != err {
		logger.Error("OpenFile", err)
		return "", "", false
	}
	defer toFile.Close()

	buf := make([]byte, 1024)
	for {
		n, err := file.Read(buf)
		if io.EOF == err {
			toFile.Write(buf[:n])
			break
		}
		toFile.Write(buf[:n])
	}
	return newFilePath, fmt.Sprintf("%s/%s", dirMap[basePath].relDir, newFileName), true
}

var lock sync.Mutex

func uploadPath(dm *dirMapping, basePath string) string {
	lock.Lock()
	defer lock.Unlock()
	list, _ := os.ReadDir(dm.dir)
	if 999 > len(list) {
		dm.count++
		return dm.dir
	}
	return createNewDir(dm.dirName, basePath)
}

//	RandCharacterString 生成 0-9 a-z 随机字符
//	num int 指定生成字符数量
func RandCharacterString(num int) string {
	str := strings.Builder{}
	max := len(chaArr)
	for i := 0; i < num; i++ {
		result, _ := rand.Int(rand.Reader, big.NewInt(int64(max)))
		str.WriteString(chaArr[result.Int64()])
	}
	return str.String()
}
