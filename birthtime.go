
// +build !darwin

package main

import ("os";"time")

func birthTime(fi os.FileInfo) time.Time {
	return fi.ModTime()
}
