package utils

import (
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
