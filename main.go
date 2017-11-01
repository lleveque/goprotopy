// Copyright 2014 Brett Slatkin
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
    "flag"
    "log"
    "os"
)

// var outputFolder = "python"

func processFile(packagePath string, filePath string, outputFolder string) (err error) {
    packageName, functions := loadFile(filePath)

    // Create subfolder
    if err = os.MkdirAll(outputFolder, os.ModePerm); err != nil {
        return
    }

    outputPath, err := getRenderedPath(filePath, outputFolder, "go")
    if err != nil {
        return
    }
    
    // Render python-friendly Go file
    output, err := os.OpenFile(outputPath, os.O_WRONLY|os.O_CREATE, 0600)
    if err != nil {
        return
    }
    defer output.Close()

    if err = render(output, packagePath, packageName, functions); err != nil {
        return
    }

    outputPathC, err := getRenderedPath(filePath, outputFolder, "c")
    if err != nil {
        return
    }

    // Render C file
    outputC, err := os.OpenFile(outputPathC, os.O_WRONLY|os.O_CREATE, 0600)
    if err != nil {
        return
    }
    defer outputC.Close()

    if err = renderC(outputC, packagePath, packageName, functions); err != nil {
        return
    }
    return
}

func main() {
    log.SetFlags(0)
    log.SetPrefix("goprotopy: ")

    // packagePath := "github.com/lleveque/greeting"

    packagePath := flag.String("packagePath", "", "the full path of the go package you want to generate Python bindings for")
    outputFolder := flag.String("outputFolder", "python", "the output folder for generated files")

    flag.Parse()

    for _, filePath := range flag.Args() {
        if err := processFile(*packagePath, filePath, *outputFolder); err != nil {
            log.Fatalf("Error processing file %s: %s", filePath, err)
        }
    }
}
