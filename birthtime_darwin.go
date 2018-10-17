
// +build darwin

package main

import ("os";"time";"syscall")

func birthTime(fi os.FileInfo) time.Time {
	s := fi.Sys().(*syscall.Stat_t)
	return time.Unix(s.Birthtimespec.Unix())
}
