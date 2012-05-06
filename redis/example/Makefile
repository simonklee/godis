ALL=string list transaction

all: $(ALL)

%: %.go
	go build $@.go

clean:
	rm -f $(ALL)
	go clean
