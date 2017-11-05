# goprotopy
A tool to generate Python bindings for Go APIs taking []byte and returning []byte.

## Dependencies

To build the generated bindings, you will probably need to install the `python-dev` package for your OS.

## Installation

Running `go get -u github.com/lleveque/goprotopy` should be enough to install the tool.

## Usage

Run `cd example/ && go generate` to generate Python bindings for this package.

You can then run `cd python/ && go generate` to build the Python module.

Enjoy :
```
$ python3 -c 'import greet; print(greet.Hello("Gopher"))'
b'Hello Gopher !'
```

## Note

This project can be used on its own, but it is designed to work gracefully with its sister project, [protoc-gen-go/grpcserial](https://github.com/lleveque/protoc-gen-go), that helps you build "serialized" APIs.
