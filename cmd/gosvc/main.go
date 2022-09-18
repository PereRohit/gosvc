package main

import (
	"bytes"
	"context"
	"embed"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/PereRohit/gosvc/internal"
)

const (
	VERSION = "gosvc 0.0.8"
)

var (
	f embed.FS

	re            = regexp.MustCompile(`[^\w]`)
	svcFolderName = ""
	serviceName   = ""

	data any

	namesToReplace = map[string]string{
		"gitignore": ".gitignore",
		"service":   "service",
		"github":    ".github",
	}
)

var (
	moduleName = flag.String("init", "", "go module name")
	version    = flag.Bool("version", false, "version")
)

func WalkAndCreate(srcPath, destPath string) error {
	dirs, err := f.ReadDir(srcPath)
	if err != nil {
		return fmt.Errorf("read dir error: %s\n", err.Error())
	}
	for _, dir := range dirs {
		dirName := dir.Name()

		srcFilePath := path.Join(srcPath, dirName)

		destDirName := dirName
		if name, replace := namesToReplace[destDirName]; replace {
			destDirName = name
		}
		destFilePath := filepath.Join(destPath, destDirName)

		if !dir.IsDir() {
			tmplFile := strings.LastIndex(destFilePath, ".tmpl") != -1
			destFilePath = strings.Replace(destFilePath, ".tmpl", "", -1)

			fileData, err := f.ReadFile(srcFilePath)
			if err != nil {
				return fmt.Errorf("read file error: %s\n", err.Error())
			}

			destDir, _ := filepath.Split(destFilePath)

			if len(destDir) > 0 {
				err = os.MkdirAll(destDir, os.ModePerm)
				if err != nil {
					return fmt.Errorf("create dir error: %s\n", err.Error())
				}
			}

			if tmplFile {
				tmpl, err := template.New(srcFilePath).Parse(string(fileData))
				if err != nil {
					return fmt.Errorf("parse template error: %s\n", err.Error())
				}

				buf := bytes.NewBuffer(nil)
				err = tmpl.Execute(buf, data)
				if err != nil {
					return fmt.Errorf("execute template error: %s\n", err.Error())
				}
				fileData = buf.Bytes()
			}

			err = ioutil.WriteFile(destFilePath, fileData, os.ModePerm)
			if err != nil {
				return fmt.Errorf("file write error: %s\n", err.Error())
			}
			continue
		}
		err = WalkAndCreate(srcFilePath, destFilePath)
		if err != nil {
			return err
		}
	}
	return err
}

func ProcessServiceName(modulePath string) {
	// extract GitHub username
	ghUserName := ""
	splitModule := strings.Split(modulePath, "/")
	if len(splitModule) > 1 && splitModule[0] == "github.com" {
		ghUserName = splitModule[1]
	}

	_, serviceName = path.Split(modulePath)
	svcFolderName = serviceName
	titleCaser := cases.Title(language.English, cases.NoLower)
	serviceName = titleCaser.String(serviceName)
	serviceName = strings.TrimSpace(serviceName)
	serviceName = strings.Replace(serviceName, ".", "-", -1)
	serviceName = strings.Replace(serviceName, "_", "-", -1)

	names := strings.Split(serviceName, "-")
	newName := serviceName
	if len(names) > 0 {
		newName = ""
		for _, name := range names {
			name = titleCaser.String(name)
			newName += name
		}
	}
	serviceName = newName

	serviceName = re.ReplaceAllString(serviceName, "")

	unexportedServiceName := cases.Lower(language.English).String(serviceName[0:1]) + serviceName[1:]

	namesToReplace["service"] = svcFolderName
	data = map[string]any{
		"Service":        serviceName,
		"Module":         modulePath,
		"service":        unexportedServiceName,
		"Repository":     svcFolderName,
		"GitHubUserName": ghUserName,
	}
}

func CommandRunner(command string, args ...string) error {
	var (
		outBuf, errBuf bytes.Buffer
	)
	cmd := exec.CommandContext(context.Background(), command, args...)

	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	err := cmd.Run()
	if err != nil {
		if stdErr := errBuf.String(); stdErr != "" {
			err = fmt.Errorf("error: %s\n", stdErr)
		} else {
			err = fmt.Errorf("error: %s\n", err.Error())
		}
	}
	return err
}

func main() {
	flag.Parse()

	if *version {
		fmt.Fprintf(os.Stdout, "%s\n", VERSION)
		os.Exit(0)
	}

	initModule := strings.TrimSpace(*moduleName)
	if len(initModule) <= 0 {
		fmt.Fprintf(os.Stderr, "error: invalid usage\n")
		flag.PrintDefaults()
		os.Exit(1)
	}
	// relative module names not allowed
	if initModule[0] == '.' ||
		initModule[0] == '/' ||
		initModule[0] == '\\' {
		fmt.Fprintf(os.Stderr, "error: go module must not be relative\n")
		os.Exit(1)
	}

	f = internal.GetEmbeddedFS()

	ProcessServiceName(initModule)

	svcFolder := filepath.Join(".", svcFolderName)
	absPath, err := filepath.Abs(svcFolder)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		os.Exit(1)
	}
	defer func() {
		if err != nil {
			fmt.Fprintf(os.Stderr, "errors obtained: %s\nAttempting cleanup...\n", err.Error())
			err = os.RemoveAll(absPath)
		}
	}()

	err = WalkAndCreate("resources", svcFolder)
	if err != nil {
		return
	}

	err = os.Chdir(absPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		os.Exit(1)
	}
	err = CommandRunner("go", "mod", "tidy")
	if err != nil {
		return
	}
	err = CommandRunner("go", "generate", "./...")
	if err != nil {
		return
	}
}
