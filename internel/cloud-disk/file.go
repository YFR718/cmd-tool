package cloud_disk

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func Run() {

	path := "\\"
	reader := bufio.NewReader(os.Stdin)
	cmd := ""

	for {
		fmt.Printf("%v> ", path)
		cmd, _ = reader.ReadString('\n')
		cmd = cmd[:len(cmd)-2]
		cmds := strings.Split(cmd, " ")
		if len(cmds) == 0 {
			fmt.Println("请输入！")
		}
		switch cmds[0] {
		case "exit":
			os.Exit(1)
		case "ls":
			err := getlist(path, true)
			if err != nil {
				fmt.Println(err)
			}
		case "cd":
			tmp := filepath.Join(path, cmds[1])
			err := getlist(tmp, false)
			if err != nil {
				fmt.Println(err)
			} else {
				path = tmp
			}
		case "pull":
			tmp := filepath.Join(path, cmds[1])
			err := pullFile(tmp)
			if err != nil {
				fmt.Println(err)
			}
		case "push":
			tmp, _ := os.Getwd()
			tmp = filepath.Join(tmp, cmds[1])
			err := pushFile(path, tmp)
			if err != nil {
				fmt.Println(err)
			}
		case "rm":
			tmp := filepath.Join(path, cmds[1])
			err := rmFile(tmp)
			if err != nil {
				fmt.Println(err)
			}
		case "mkdir":
			tmp := filepath.Join(path, cmds[1])
			err := mkdir(tmp)
			if err != nil {
				fmt.Println(err)
			}
		default:
			fmt.Println("命令格式错误，请重新输入")
			break
		}
	}
}
