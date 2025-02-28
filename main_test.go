package main

import (
	"strings"
	"testing"

	jenkinsfile "github.com/james-stocks/jenkinsfile-editor/pkg"
)

func TestParseJenkinsfile(t *testing.T) {
	content := `
pipeline {
    agent any

    stages {
        stage('Build') {
            steps {
                echo 'Building..'
            }
        }
        stage('Test') {
            steps {
                echo 'Testing..'
            }
        }
        stage('Deploy') {
            steps {
                echo 'Deploying....'
            }
        }
    }
}
`
	pipeline, err := jenkinsfile.Parse(content)
	if err != nil {
		t.Fatalf("Error parsing Jenkinsfile: %v", err)
	}

	if len(pipeline.Elements) != 1 || pipeline.Elements[0].Name != "pipeline" {
		t.Errorf("Expected pipeline block, got '%v'", pipeline.Elements)
	}

	stages := pipeline.Elements[0].Children[1]
	if stages.Name != "stages" {
		t.Errorf("Expected stages block, got '%s'", stages.Name)
	}

	expectedStages := []string{"stage('Build')", "stage('Test')", "stage('Deploy')"}
	for i, stage := range stages.Children {
		if stage.Name != expectedStages[i] {
			t.Errorf("Expected stage '%s', got '%s'", expectedStages[i], stage.Name)
		}
	}
}

func TestParseShBlock(t *testing.T) {
	content := `
pipeline {
    agent any

    stages {
        stage('Build') {
            steps {
                echo 'Building..'
            }
        }
        stage('Test') {
            steps {
                echo 'Testing..'
            }
        }
        stage('Deploy') {
            steps {
                sh '''
                    echo "Deploying.."
                '''
            }
        }
    }
}
`
	pipeline, err := jenkinsfile.Parse(content)
	if err != nil {
		t.Fatalf("Error parsing Jenkinsfile: %v", err)
	}

	if len(pipeline.Elements) != 1 || pipeline.Elements[0].Name != "pipeline" {
		t.Errorf("Expected pipeline block, got '%v'", pipeline.Elements)
	}

	stages := pipeline.Elements[0].Children[1]
	if stages.Name != "stages" {
		t.Errorf("Expected stages block, got '%s'", stages.Name)
	}

	expectedStages := []string{"stage('Build')", "stage('Test')", "stage('Deploy')"}
	for i, stage := range stages.Children {
		if stage.Name != expectedStages[i] {
			t.Errorf("Expected stage '%s', got '%s'", expectedStages[i], stage.Name)
		}
	}

	deploySteps := stages.Children[2].Children[0]
	if deploySteps.Name != "steps" {
		t.Errorf("Expected steps block, got '%s'", deploySteps.Name)
	}

	shBlock := deploySteps.Children[0]
	if !strings.HasPrefix(shBlock.Name, "sh '''") {
		t.Errorf("Expected sh block, got '%s'", shBlock.Name)
	}
}

func TestWritePipelineToBuffer(t *testing.T) {
	original := `pipeline {
    agent any
    stages {
        stage('Build') {
            steps {
                echo 'Building..'
            }
        }
        stage('Test') {
            steps {
                echo 'Testing..'
            }
        }
        stage('Deploy') {
            steps {
                sh '''
                    echo "Deploying.."
                '''
            }
        }
    }
}
`
	pipeline, err := jenkinsfile.Parse(original)
	if err != nil {
		t.Fatalf("Error parsing Jenkinsfile: %v", err)
	}

	output := pipeline.ToString()
	if output != original {
		t.Errorf("Expected output to match original:\n%s\nGot:\n%s", original, output)
	}
}
