/*
@author: sk
@date: 2024/11/16
*/
package main

import "os"

func HandleErr(err error) {
	if err != nil {
		panic(err)
	}
}

func ReadFile(path string) string {
	bs, err := os.ReadFile(path)
	HandleErr(err)
	return string(bs)
}
