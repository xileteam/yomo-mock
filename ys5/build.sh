#!/bin/bash

go build -o build/ys5_proxy proxy/main.go

go build -o build/ys5_crawler crawler/main.go
