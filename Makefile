include $(GOROOT)/src/Make.inc

TARG=godis
GOFILES=\
	godis.go\
	commands.go\

include $(GOROOT)/src/Make.pkg
