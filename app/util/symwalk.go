package util

import (
	"os"
	"path/filepath"
)

/* 软件包symwalk提供了符号链接感知文件路径遍历的实现。
Walk调用filepath.Walk，为它提供一个特殊的WalkFn来包装给定的WalkFn。
此函数为常规文件调用给定的WalkFn。但是，当遇到符号链接时，它将使用filepath.EvalSymlinks完全解析该链接，
并以递归方式调用Walk。这样可以确保与filepath.Walk不同，遍历不会在符号链接处停止。
*/
func SymWalk(path string, walkFn filepath.WalkFunc) error {
	return walk(path, path, walkFn)
}
func walk(filename string, linkDirname string, walkFn filepath.WalkFunc) error {
	symWalkFunc := func(path string, info os.FileInfo, err error) error {

		if fname, err := filepath.Rel(filename, path); err == nil {
			path = filepath.Join(linkDirname, fname)
		} else {
			return err
		}

		if err == nil && info.Mode()&os.ModeSymlink == os.ModeSymlink {
			finalPath, err := filepath.EvalSymlinks(path)
			if err != nil {
				return err
			}
			info, err := os.Lstat(finalPath)
			if err != nil {
				return walkFn(path, info, err)
			}
			if info.IsDir() {
				return walk(finalPath, path, walkFn)
			}
		}

		return walkFn(path, info, err)
	}
	return filepath.Walk(filename, symWalkFunc)
}
