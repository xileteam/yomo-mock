#!/bin/bash

go build -o build/source source/main.go

go build -o build/sfn-1 sfn-1/main.go

go build -o build/sfn-2 sfn-2/main.go
