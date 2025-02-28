package main

import (
	"bufio"
	"fmt"
	"os"

	jenkinsfile "github.com/james-stocks/jenkinsfile-editor/pkg"
)

func main() {
	// Read Jenkinsfile content from stdin
	reader := bufio.NewReader(os.Stdin)
	content, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading from stdin:", err)
		return
	}

	pipeline, err := jenkinsfile.Parse(content)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(pipeline.ToString())
}
