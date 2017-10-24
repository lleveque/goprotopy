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
    "fmt"
    "go/ast"
    "go/parser"
    "go/token"
    "log"
    "strings"
)

func loadFile(inputPath string) (string, []GeneratedFunction) {
    fset := token.NewFileSet()
    f, err := parser.ParseFile(fset, inputPath, nil, parser.ParseComments)
    if err != nil {
        log.Fatalf("Could not parse file: %s", err)
    }

    packageName := identifyPackage(f)
    if packageName == "" {
        log.Fatalf("Could not determine package name of %s", inputPath)
    }

    protopyFuncs := map[string]bool{}
    for _, decl := range f.Decls {
        funcName, ok := identifyFunction(decl)
        if ok {
            protopyFuncs[funcName] = true
            continue
        }
    }

    functions := []GeneratedFunction{}
    for funcName, _ := range protopyFuncs {
        protopyFunction := GeneratedFunction{funcName}
        functions = append(functions, protopyFunction)
    }

    return packageName, functions
}

func identifyPackage(f *ast.File) string {
    if f.Name == nil {
        return ""
    }
    return f.Name.Name
}

func identifyFunction(decl ast.Decl) (funcName string, match bool) {
    funcDecl, ok := decl.(*ast.FuncDecl)
    if !ok {
        return
    }
    if funcDecl.Doc == nil {
        return
    }

    found := false
    for _, comment := range funcDecl.Doc.List {
        if strings.Contains(comment.Text, "@protopy") {
            found = true
            break
        }
    }
    if !found {
        return
    }

    if funcDecl.Name != nil {
        funcName = funcDecl.Name.Name
    }
    if funcName == "" {
        return
    }

    // TODO Check function is exported
    // TODO Check function signature
    fmt.Printf("Function %v : takes %#v returns %#v\n", funcName, funcDecl.Type.Params.List, funcDecl.Type.Results.List)
    
    match = true
    return
}
