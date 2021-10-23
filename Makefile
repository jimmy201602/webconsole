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