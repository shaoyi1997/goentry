# Go Entry Task

## Documents
- [Requirements](docs/Requirements.md)
- [Report](docs/Report.md)

## Pre-requisites
- MySQL v8.0.25
- Redis v6.2.3

## Setup
0. Ensure that MySQL & Redis servers are running
1. Create source directory
```bash
mkdir $GOPATH/src/git.garena.com/shaoyi.hong/ && cd $GOPATH/src/git.garena.com/shaoyi.hong/
```
2. Clone this repository
```bash
git clone gitlab@git.garena.com:shaoyi.hong/go-entry-task.git
```
3. Install dependencies
```bash
go mod download
```
4. Run 
```bash
make run
```

