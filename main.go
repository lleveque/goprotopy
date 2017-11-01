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
    "log"
    "os"
)

var outputFolder = "python"

func processFile(inputPath string) (err error) {
    packageName, functions := loadFile(inputPath)

    // TODO Construct full package name
    fullPackageName := "github.com/lleveque/greeting"

    // Create subfolder
    if err = os.MkdirAll(outputFolder, os.ModePerm); err != nil {
        return
    }

    outputPath, err := getRenderedPath(inputPath, outputFolder, "go")
    if err != nil {
        return
    }
    
    // Render python-friendly Go file
    output, err := os.OpenFile(outputPath, os.O_WRONLY|os.O_CREATE, 0600)
    if err != nil {
        return
    }
    defer output.Close()

    if err = render(output, fullPackageName, packageName, functions); err != nil {
        return
    }

    outputPathC, err := getRenderedPath(inputPath, outputFolder, "c")
    if err != nil {
        return
    }

    // Render C file
    outputC, err := os.OpenFile(outputPathC, os.O_WRONLY|os.O_CREATE, 0600)
    if err != nil {
        return
    }
    defer outputC.Close()

    if err = renderC(outputC, fullPackageName, packageName, functions); err != nil {
        return
    }
    return
}

func main() {
    log.SetFlags(0)
    log.SetPrefix("goprotopy: ")

    if len(os.Args) != 2 {
        log.Fatalf("Usage : goprotopy filename.go")
        return
    }

    for _, path := range os.Args[1:] {
        if err := processFile(path); err != nil {
            log.Fatalf("Error processing file %s: %s", path, err)
        }
    }
}
