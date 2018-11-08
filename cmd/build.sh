#!/bin/bash

go-bindata -o ./phonedat.go ../phone.dat
GOOS=windows GOARCH=amd64 go build -o phoneregion-windows
GOOS=darwin GOARCH=amd64 go build -o phoneregion-mac
# my work os is linux
go build -o phoneregion -o phoneregion-linux
