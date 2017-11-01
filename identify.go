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
    "go/ast"
    "go/parser"
    "go/token"
    "log"
    "strings"
    "unicode"
    "unicode/utf8"
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

    protopyFuncs := map[string]string{}
    for _, decl := range f.Decls {
        funcName, funcDoc, ok := identifyFunction(decl)
        if ok {
            protopyFuncs[funcName] = funcDoc
            continue
        }
    }

    functions := []GeneratedFunction{}
    for funcName, funcDoc := range protopyFuncs {
        protopyFunction := GeneratedFunction{funcName, funcDoc}
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

func identifyFunction(decl ast.Decl) (funcName string, funcDoc string, match bool) {
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

    // If function is not exported, don't handle it
    firstLetter, _ := utf8.DecodeRuneInString(funcName)
    if unicode.IsLower(firstLetter) {
        return
    }

    // Check func signature :
    funcParams := funcDecl.Type.Params.List
    funcResults := funcDecl.Type.Results.List
    // If function doesn't have exactly 1 parameter and 2 results, don't handle it
    if len(funcParams) != 1 || len(funcResults) != 2 {
        return
    }
    // TODO Check types
    // fmt.Printf("Function %v : takes %#v returns %#v, %#v\n", funcName, funcParams[0].Type, funcResults[0].Type, funcResults[1].Type)

    // Get documentation
    for _, comment := range funcDecl.Doc.List {
        commentText := comment.Text
        // Don't handle our protopy annotation
        if commentText == "// @protopy" {
            continue
        }
        commentText = strings.TrimPrefix(commentText, "// ")
        commentText = strings.TrimPrefix(commentText, "/* ")
        commentText = strings.TrimSuffix(commentText, " */")
        funcDoc += commentText+"\n"
    }
    funcDoc = strings.TrimSuffix(funcDoc, "\n")
    funcDoc = strings.Replace(funcDoc, "\n", "\\n", -1)
        
    match = true
    return
}
