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
    "path/filepath"
    "log"
    "os"
)

var outputFolder = "python"

func processFile(inputPath string) {
    log.Printf("Processing file %s", inputPath)

    packageName, functions := loadFile(inputPath)

    log.Printf("Found goprotopy functions to generate in package %v: %#v", packageName, functions)

    outputPath, err := getRenderedPath(inputPath)
    if err != nil {
        log.Fatalf("Could not get output path: %s", err)
    }

    outputPathC, err := getCRenderedPath(inputPath)
    if err != nil {
        log.Fatalf("Could not get C output path: %s", err)
    }

    // Create subfolder
    err = os.MkdirAll(outputFolder, os.ModePerm)
    if err != nil {
        log.Fatalf("Could not create folder: %s", err)
    }
    
    output, err := os.OpenFile(filepath.Join(outputFolder, outputPath), os.O_WRONLY|os.O_CREATE, 0600)
    if err != nil {
        log.Fatalf("Could not open output file: %s", err)
    }
    defer output.Close()

    // TODO Construct full package name
    if err := render(output, "github.com/lleveque/goprotopiable", packageName, functions); err != nil {
        log.Fatalf("Could not generate go code: %s", err)
    }

    // Render C file
    outputC, err := os.OpenFile(filepath.Join(outputFolder, outputPathC), os.O_WRONLY|os.O_CREATE, 0600)
    if err != nil {
        log.Fatalf("Could not open output file: %s", err)
    }
    defer outputC.Close()

    // TODO Construct full package name
    if err := renderC(outputC, "github.com/lleveque/goprotopiable", packageName, functions); err != nil {
        log.Fatalf("Could not generate go code: %s", err)
    }
}

func main() {
    log.SetFlags(0)
    log.SetPrefix("goprotopy: ")

    for _, path := range os.Args[2:] {
        processFile(path)
    }
}
