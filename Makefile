# Go parameters
GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean

build:
		$(GOBUILD) -v -o ./bin/apibox ./main.go
		$(GOBUILD) -v -o ./bin/ssh ./cmd/pty.go

clean:
		$(GOCLEAN)
		rm -f ./bin/apibox
		rm -f ./bin/ssh

package:
		rm -rf release
		mkdir release
		mkdir -p release/bin
		cp bin/apibox release/bin
		cp bin/ssh release/bin
		cp -r conf static template release
