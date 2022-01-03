ENTRY=cmd/main.go
TARGET=vodstream

COMPILER=go build
FLAGS=CGO_ENABLED=0

all: $(TARGET)

vodstream:
	$(FLAGS) go build -o build/$(TARGET) $(ENTRY)

clean:
	rm -f build/$(TARGET)