package internal

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

type Pom struct {
	XmlName                 xml.Name             `xml:"project"`
	Parent                  parent               `xml:"parent"`
	GroupId                 string               `xml:"groupId"`
	ArtifactId              string               `xml:"artifactId"`
	Version                 string               `xml:"version"`
	Properties              Properties           `xml:"properties"`
	Modules                 []string             `xml:"modules>module"`
	Dependencies            []dependency         `xml:"dependencies>dependency"`
	DependencyManagement    dependencyManagement `xml:"dependencyManagement"`
	Profiles                []profile            `xml:"profiles>profile"`
	Build                   build                `xml:"build"`
	PomFilePath             string
	propertyMap             map[string]string
	dependencyManagementMap map[string]string
}

// Parent represents the parent POM if this project is a module.
type parent struct {
	GroupId      string `xml:"groupId"`
	ArtifactId   string `xml:"artifactId"`
	Version      string `xml:"version"`
	RelativePath string `xml:"relativePath"`
}

type Properties struct {
	Entries []Property `xml:",any"` // Capture all elements inside <properties>
}

type Property struct {
	XMLName xml.Name
	Value   string `xml:",chardata"`
}

// Dependency represents a single Maven dependency.
type dependency struct {
	GroupId    string `xml:"groupId"`
	ArtifactId string `xml:"artifactId"`
	Version    string `xml:"version"`
	Scope      string `xml:"scope,omitempty"`
}

type profile struct {
	Id                      string               `xml:"id"`
	ActiveByDefault         string               `xml:"activation>activeByDefault"`
	Properties              Properties           `xml:"properties"`
	Modules                 []string             `xml:"modules>module"` // Capture the modules
	Dependencies            []dependency         `xml:"dependencies>dependency"`
	DependencyManagement    dependencyManagement `xml:"dependencyManagement"`
	Build                   build                `xml:"build"`
	propertyMap             map[string]string
	dependencyManagementMap map[string]string
}

// DependencyManagement includes a list of dependencies that are managed.
type dependencyManagement struct {
	Dependencies []dependency `xml:"dependencies>dependency"`
}

// Build represents the build configuration which can contain plugins.
type build struct {
	Plugins []plugin `xml:"plugins>plugin"`
}

// Plugin represents a build plugin.
type plugin struct {
	GroupId    string `xml:"groupId"`
	ArtifactId string `xml:"artifactId"`
	Version    string `xml:"version"`
}

func CreateEffectivePom(pomPath string) (Pom, error) {
	// todo:
	// 1. Use maven wrapper if exists.
	// 2. Download maven if "mvn" command not exist in path.
	cmd := exec.Command("mvn", "help:effective-pom", "-f", pomPath, "-pl", filepath.Base(pomPath))
	output, err := cmd.Output()
	if err != nil {
		return Pom{}, err
	}
	effectivePomString, err := getEffectivePomStringFromConsoleOutput(string(output))
	if err != nil {
		return Pom{}, err
	}
	var resultPom Pom
	err = xml.Unmarshal([]byte(effectivePomString), &resultPom)
	return resultPom, nil
}

var projectStart = regexp.MustCompile(`^\s*<project `) // the space can not be deleted.
var projectEnd = regexp.MustCompile(`^\s*</project>\s*$`)

func getEffectivePomStringFromConsoleOutput(consoleOutput string) (string, error) {
	var builder strings.Builder
	scanner := bufio.NewScanner(strings.NewReader(consoleOutput))
	projectStarted := false
	projectEnded := false
	for scanner.Scan() {
		line := scanner.Text()
		if projectStart.MatchString(line) {
			projectStarted = true
		} else if projectEnd.MatchString(line) {
			projectEnded = true
		}
		if projectStarted {
			builder.WriteString(line)
		}
		if projectEnded {
			break
		}
	}
	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("failed to scan console output: %w", err)
	}
	result := builder.String()
	if result == "" {
		return "", fmt.Errorf("failed to get effective pom from console: empty content")
	}
	return result, nil
}
