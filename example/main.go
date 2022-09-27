package main

import "localfs"

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
}
