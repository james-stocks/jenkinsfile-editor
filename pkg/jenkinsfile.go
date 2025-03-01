package jenkinsfile

import (
	"bytes"
	"fmt"
	"strings"
)

// Pipeline represents the structure of a Jenkinsfile
type Pipeline struct {
	Elements []Element
}

func (p *Pipeline) ToString() string {
	var buffer bytes.Buffer
	for _, element := range p.Elements {
		writeElementToBuffer(&buffer, element, "")
	}
	return buffer.String()
}

func (p *Pipeline) GetStageIndexForStep(step string) (int) {
    for _, element := range p.Elements {
        if element.Name == "pipeline" {
            for _, pipelineElement := range element.Children {
                if pipelineElement.Name == "stages" {
                    for k, stageElement := range pipelineElement.Children {
                        for _, stageChild := range stageElement.Children {
							if stageChild.Name == "steps" {
								for _, stepsElement := range stageChild.Children {
									if strings.Contains(stepsElement.Content, step) {
										return k
									}
								}
							}
                        }
                    }
                }
            }
        }
    }
    return -1
}

func (p *Pipeline) InsertStage(stageName string, steps []string, index int) {
	newStage := Element{
		Type:      "stage",
		Name:      fmt.Sprintf("stage('%s')", stageName),
		Content:   "",
		HasBraces: true,
	}
	newSteps := Element{
		Type:      "steps",
		Name:      "steps",
		Content:   "",
		HasBraces: true,
	}
	for _, step := range steps {
		newStep := Element{
			Type:    "element",
			Name:    step,
			Content: step,
		}
		newSteps.Children = append(newSteps.Children, newStep)
	}
	newStage.Children = append(newStage.Children, newSteps)
	// i and j indices are needed to ensure overwriting the original element
    for i, element := range p.Elements {
        if element.Name == "pipeline" {
            for j, pipelineElement := range element.Children {
                if pipelineElement.Name == "stages" {
                    newStages := make([]Element, len(pipelineElement.Children)+1)
                    copy(newStages, pipelineElement.Children[:index])
                    newStages[index] = newStage
                    copy(newStages[index+1:], pipelineElement.Children[index:])
                    p.Elements[i].Children[j].Children = newStages
                }
            }
        }
    }
}

// Element represents a generic element in the pipeline
type Element struct {
	Type      string
	Name      string
	Content   string
	Children  []Element
	HasBraces bool
}

func Parse(content string) (Pipeline, error) {
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

func writeElementToBuffer(buffer *bytes.Buffer, e Element, indent string) {
	buffer.WriteString(fmt.Sprintf("%s%s", indent, e.Name))
	if e.HasBraces {
		buffer.WriteString(" {\n")
	} else {
		buffer.WriteString("\n")
	}
	for _, child := range e.Children {
		if child.Type == "sh-line" {
			writeElementToBuffer(buffer, child, indent+"        ")
		} else {
			writeElementToBuffer(buffer, child, indent+"    ")
		}
	}
	if e.HasBraces {
		buffer.WriteString(fmt.Sprintf("%s}\n", indent))
	}
}
