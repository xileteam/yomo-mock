#!/bin/bash

cd proxy && go build -o ../build/proxy && cd ..

cd crawler && go build -o ../build/crawler && cd ..
