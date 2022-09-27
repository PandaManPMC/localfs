package localfs

//  author: laoniqiu
//  since: 2022/9/27
//  desc: localfs

type Log3 interface {
	Info(string)
	Error(string, error)
}
