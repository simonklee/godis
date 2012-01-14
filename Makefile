include $(GOROOT)/src/Make.inc

TARG=godis
GOFILES=\
	godis.go\
	commands.go\
	conn.go\

include $(GOROOT)/src/Make.pkg

format:
	gofmt -s=true -tabs=false -tabwidth=4 -w .

.PHONY: format 
