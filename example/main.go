package main

import (
	"crypto/md5"
	"fmt"
	"github.com/PandaManPMC/localfs"
)

//  author: laoniqiu
//  since: 2022/9/27
//  desc: example

type FsLog3 struct {
}

func (*FsLog3) Info(msg string) {
	println(msg)
}

func (*FsLog3) Error(msg string, err error) {
	println(msg)
	println(err)
}

func main() {
	log := new(FsLog3)
	errInitDir := localfs.InitDir("/usr/local/localfs", log)
	if nil != errInitDir {
		panic(errInitDir)
	}

	data := []byte("巍巍的终南山高入云霄，与天帝的住所临近。")
	absPath, relPath, isOk := localfs.UploadFileByByte("a.txt", int64(len(data)), data)
	println(isOk)
	println(absPath)
	println(relPath)

	hash := md5.New()
	hash.Write(data)
	md := hash.Sum(nil)
	println(fmt.Sprintf("%x", md))
}
