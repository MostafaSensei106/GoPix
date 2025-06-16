#!/bin/bash

AppName = "GoPix"
BinName = "GoPix"
InstallationDir = "/usr/local/bin"

# cheak user platform os and architecture
platform=$(uname -s | tr '[:upper:]' '[:lower:]')
case $platform in
    linux*)
        if [ "$(uname -m)" = "x86_64" ]; then
            ARCH="amd64"
        else
            ARCH="386"
        fi
        ;;
    darwin*)
        if [ "$(uname -m)" = "x86_64" ]; then
            ARCH="amd64"
        else
            ARCH="arm64"
        fi
        ;;
    *)
        echo "Unsupported platform: $platform"
        exit 1
        ;;
esac

echo "Hi there! Iam Mr Mostafa Sensei This script will install $AppName to $InstallationDir Do you want to continue? (y/n)"
# Ask do you want to continue?
read -r answer # if empty or if it y or Y then install

if [[ $answer != [yY] ]]; then
    echo "Goodbye!"
    exit 0
fi


# check if go is installed
if ! command -v go &> /dev/null
then
    echo "Go is not installed. Please install Go and try again."
    exit 1
fi

# check if dependencies are installed : to do 

echo "Building... $AppName... Please wait..."

go build -o $BinName