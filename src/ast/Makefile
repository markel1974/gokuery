.SUFFIXES: .peg .go

.peg.go:
	./peg -noast -switch -inline -strict -output $@ $<

all: kuery.go
