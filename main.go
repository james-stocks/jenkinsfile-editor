package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
)

// Pipeline represents the structure of a Jenkinsfile
type Pipeline struct {
	Elements []Element
}

// Element represents a generic element in the pipeline
type Element struct {
	Type      string
	Name      string
	Content   string
	Children  []Element
	HasBraces bool
}

func parseJenkinsfile(content string) (Pipeline, error) {
	lines := strings.Split(strings.TrimSpace(content), "\n")
	var pipeline Pipeline
	var stack []Element
	var inShBlock bool

	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])

		// Skip empty lines
		if line == "" {
			continue
		}

		// Handle sh block start
		if strings.HasPrefix(line, "sh '''") {
			element := Element{
				Type:    "sh",
				Name:    line,
				Content: line,
			}
			if len(stack) > 0 {
				parent := &stack[len(stack)-1]
				parent.Children = append(parent.Children, element)
			} else {
				pipeline.Elements = append(pipeline.Elements, element)
			}
			inShBlock = true
			continue
		}

		// Handle sh block end
		if inShBlock && strings.HasSuffix(line, "'''") {
			element := Element{
				Type:    "sh",
				Name:    line,
				Content: line,
			}
			if len(stack) > 0 {
				parent := &stack[len(stack)-1]
				parent.Children = append(parent.Children, element)
			} else {
				pipeline.Elements = append(pipeline.Elements, element)
			}
			inShBlock = false
			continue
		}

		// Handle lines inside sh block
		if inShBlock {
			element := Element{
				Type:    "sh-line",
				Name:    line,
				Content: line,
			}
			if len(stack) > 0 {
				parent := &stack[len(stack)-1]
				parent.Children = append(parent.Children, element)
			} else {
				pipeline.Elements = append(pipeline.Elements, element)
			}
			continue
		}

		// Handle nested blocks
		if strings.HasSuffix(line, "{") {
			element := Element{
				Type:      "block",
				Name:      strings.TrimSuffix(line, " {"),
				Content:   "",
				HasBraces: true,
			}
			stack = append(stack, element)
			continue
		}

		// Handle closing braces
		if line == "}" {
			if len(stack) > 1 {
				child := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				parent := &stack[len(stack)-1]
				parent.Children = append(parent.Children, child)
			} else if len(stack) == 1 {
				pipeline.Elements = append(pipeline.Elements, stack[0])
				stack = stack[:0]
			}
			continue
		}

		// Handle function calls and free-hanging elements
		element := Element{
			Type:    "element",
			Name:    line,
			Content: line,
		}
		if len(stack) > 0 {
			parent := &stack[len(stack)-1]
			parent.Children = append(parent.Children, element)
		} else {
			pipeline.Elements = append(pipeline.Elements, element)
		}
	}

	return pipeline, nil
}

func printElement(e Element, indent string) {
	fmt.Printf("%s%s\n", indent, e.Name)
	for _, child := range e.Children {
		printElement(child, indent+"  ")
	}
}

func writePipelineToBuffer(pipeline Pipeline) string {
	var buffer bytes.Buffer
	for _, element := range pipeline.Elements {
		writeElementToBuffer(&buffer, element, "")
	}
	return buffer.String()
}

func writeElementToBuffer(buffer *bytes.Buffer, e Element, indent string) {
	buffer.WriteString(fmt.Sprintf("%s%s", indent, e.Name))
	if e.HasBraces {
		buffer.WriteString(" {\n")
	} else {
		buffer.WriteString("\n")
	}
	for _, child := range e.Children {
		if child.Type == "sh-line" {
			writeElementToBuffer(buffer, child, indent+"    ")
		} else {
			writeElementToBuffer(buffer, child, indent+"  ")
		}
	}
	if e.HasBraces {
		buffer.WriteString(fmt.Sprintf("%s}\n", indent))
	}
}

func main() {
	// Read Jenkinsfile content from stdin
	reader := bufio.NewReader(os.Stdin)
	content, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading from stdin:", err)
		return
	}

	pipeline, err := parseJenkinsfile(content)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	for _, element := range pipeline.Elements {
		printElement(element, "")
	}
}
