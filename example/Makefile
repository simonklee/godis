include $(GOROOT)/src/Make.inc

ALL=string list transaction

all: $(ALL)

clean:
	rm -rf *.[68] $(ALL)

%: %.go
	$(GC) -o $@.6 $*.go
	$(LD) -o $@ $*.$O
