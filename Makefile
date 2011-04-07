include $(GOROOT)/src/Make.inc

TARG=godis
GOFILES=\
	godis.go\
	commands.go\
	conn.go\

include $(GOROOT)/src/Make.pkg
