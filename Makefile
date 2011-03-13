include $(GOROOT)/src/Make.inc

TARG=godis
GOFILES=\
	godis.go\
	pool.go\

include $(GOROOT)/src/Make.pkg
