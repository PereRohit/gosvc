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
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/PereRohit/gosvc/internal"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	f embed.FS

	re            = regexp.MustCompile(`[^\w]`)
	svcFolderName = ""
	serviceName   = ""

	data any
)

var (
	moduleName = flag.String("init", "", "go module name")
)

func WalkAndCreate(srcPath, destPath string) {
	dirs, err := f.ReadDir(srcPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read dir error: %s\n", err.Error())
		os.Exit(1)
	}
	for _, dir := range dirs {
		dirName := dir.Name()

		srcFilePath := filepath.Join(srcPath, dirName)

		destDirName := dirName
		if dirName == "service" {
			destDirName = svcFolderName
		}
		destFilePath := filepath.Join(destPath, destDirName)

		if !dir.IsDir() {
			destFilePath = strings.Replace(destFilePath, ".tmpl", "", -1)

			fileData, err := f.ReadFile(srcFilePath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "read file error: %s\n", err.Error())
				return
			}

			destDir, _ := filepath.Split(destFilePath)

			if len(destDir) > 0 {
				err = os.MkdirAll(destDir, os.ModePerm)
				if err != nil {
					fmt.Fprintf(os.Stderr, "create dir error: %s\n", err.Error())
					return
				}
			}

			_, srcFile := filepath.Split(srcFilePath)
			tmpl, err := template.New(srcFile).Parse(string(fileData))
			if err != nil {
				fmt.Fprintf(os.Stderr, "parse template error: %s\n", err.Error())
				return
			}

			buf := bytes.NewBuffer(nil)
			err = tmpl.Execute(buf, data)
			if err != nil {
				fmt.Fprintf(os.Stderr, "execute template error: %s\n", err.Error())
				return
			}
			fileData = buf.Bytes()

			err = ioutil.WriteFile(destFilePath, fileData, os.ModePerm)
			if err != nil {
				fmt.Fprintf(os.Stderr, "file write error: %s\n", err.Error())
				return
			}
			continue
		}
		WalkAndCreate(srcFilePath, destFilePath)
	}
}

func ProcessServiceName(modulePath string) {
	_, serviceName = filepath.Split(modulePath)
	svcFolderName = serviceName
	serviceName = cases.Title(language.English).String(serviceName)
	serviceName = strings.TrimSpace(serviceName)
	serviceName = strings.Replace(serviceName, ".", "-", -1)

	names := strings.Split(serviceName, "-")
	newName := serviceName
	if len(names) > 0 {
		newName = ""
		for _, name := range names {
			name = cases.Title(language.English).String(name)
			newName += name
		}
	}
	serviceName = newName

	serviceName = re.ReplaceAllString(serviceName, "")

	data = map[string]any{
		"Service": serviceName,
		"Module":  modulePath,
	}
}

func CommandRunner(command string, args ...string) {
	var (
		outBuf, errBuf bytes.Buffer
	)
	cmd := exec.CommandContext(context.Background(), command, args...)

	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	err := cmd.Run()
	if err != nil {
		if stdErr := errBuf.String(); stdErr != "" {
			fmt.Fprintf(os.Stderr, "error: %s\n", stdErr)
		} else {
			fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		}
		os.Exit(1)
	}
}

func main() {
	flag.Parse()

	f = internal.GetEmbeddedFS()

	initModule := strings.TrimSpace(*moduleName)
	if len(initModule) <= 0 {
		fmt.Fprintf(os.Stderr, "error: invalid usage\n")
		flag.PrintDefaults()
		os.Exit(1)
	}
	ProcessServiceName(initModule)

	WalkAndCreate("resources", "./"+svcFolderName)

	absPath, err := filepath.Abs(svcFolderName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		os.Exit(1)
	}

	err = os.Chdir(absPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		os.Exit(1)
	}
	CommandRunner("go", "mod", "tidy")
	CommandRunner("go", "generate", "./...")
}
