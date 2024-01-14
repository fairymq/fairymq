#!/bin/bash
VERSION=v0.0.0

echo "Bundling fairyMQ $VERSION"

( GOOS=darwin GOARCH=amd64 go build -o bin/macos-darwin/amd64/fairymq && tar -czf bin/macos-darwin/amd64/fairymq-$VERSION-amd64.tar.gz -C bin/macos-darwin/amd64/ $(ls  bin/macos-darwin/amd64/))
( GOOS=darwin GOARCH=arm64 go build -o bin/macos-darwin/arm64/fairymq && tar -czf bin/macos-darwin/arm64/fairymq-$VERSION-arm64.tar.gz -C bin/macos-darwin/arm64/ $(ls  bin/macos-darwin/arm64/))
( GOOS=linux GOARCH=386 go build -o bin/linux/386/fairymq && tar -czf bin/linux/386/fairymq-$VERSION-386.tar.gz -C bin/linux/386/ $(ls  bin/linux/386/))
( GOOS=linux GOARCH=amd64 go build -o bin/linux/amd64/fairymq && tar -czf bin/linux/amd64/fairymq-$VERSION-amd64.tar.gz -C bin/linux/amd64/ $(ls  bin/linux/amd64/))
( GOOS=linux GOARCH=arm go build -o bin/linux/arm/fairymq && tar -czf bin/linux/arm/fairymq-$VERSION-arm.tar.gz -C bin/linux/arm/ $(ls  bin/linux/arm/))
( GOOS=linux GOARCH=arm64 go build -o bin/linux/arm64/fairymq && tar -czf bin/linux/arm64/fairymq-$VERSION-arm64.tar.gz -C bin/linux/arm64/ $(ls  bin/linux/arm64/))
( GOOS=freebsd GOARCH=arm go build -o bin/freebsd/arm/fairymq && tar -czf bin/freebsd/arm/fairymq-$VERSION-arm.tar.gz -C bin/freebsd/arm/ $(ls  bin/freebsd/arm/))
( GOOS=freebsd GOARCH=amd64 go build -o bin/freebsd/amd64/fairymq && tar -czf bin/freebsd/amd64/fairymq-$VERSION-amd64.tar.gz -C bin/freebsd/amd64/ $(ls  bin/freebsd/amd64/))
( GOOS=freebsd GOARCH=386 go build -o bin/freebsd/386/fairymq && tar -czf bin/freebsd/386/fairymq-$VERSION-386.tar.gz -C bin/freebsd/386/ $(ls  bin/freebsd/386/))
( GOOS=windows GOARCH=amd64 go build -o bin/windows/amd64/fairymq.exe && zip -r -j bin/windows/amd64/fairymq-$VERSION-x64.zip bin/windows/amd64/fairymq.exe)
( GOOS=windows GOARCH=arm64 go build -o bin/windows/arm64/fairymq.exe && zip -r -j bin/windows/arm64/fairymq-$VERSION-x64.zip bin/windows/arm64/fairymq.exe)
( GOOS=windows GOARCH=386 go build -o bin/windows/386/fairymq.exe && zip -r -j bin/windows/386/fairymq-$VERSION-x86.zip bin/windows/386/fairymq.exe)


echo "Fin"