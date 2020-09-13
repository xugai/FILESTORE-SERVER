package utils

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

const (
	FileChunkDelCMD = `
	#!/bin/bash
	chunkDir="/Users/behe/Desktop/work_station/FILESTORE-SERVER/tmp/"
	targetDir=$1
	if [[ $targetDir =~ $chunkDir ]] && [[ $targetDir != $chunkDir ]]; then
		rm -rf $targetDir
	fi
	`
)

func RemovePathByShell(targetDir string) bool {
	cmdStr := strings.Replace(FileChunkDelCMD, "$1", targetDir, 1)
	delCmd := exec.Command("bash", "-c", cmdStr)
	if _, err := delCmd.Output(); err != nil {
		fmt.Printf("del command execute failed: %v\n", err)
		return false
	}
	return true
}

// 执行shell命令
func ExecuteShell(s string) (string, error) {
	// 函数返回一个io.Writer类型的*cmd
	cmd := exec.Command("/bin/bash", "-c", s)
	var result bytes.Buffer
	cmd.Stdout = &result
	// 执行cmd命令，期间会阻塞直至完成
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return result.String(), nil
}
