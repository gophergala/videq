package main

import "runtime"

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {

}
