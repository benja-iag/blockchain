OBJ = main.go
TARGET = main.exe

all: build

build:
	go build -o $(TARGET) $(OBJ)

clean:
	$(RM) -rf tmp $(TARGET) .completion.*

completion: build
	./$(TARGET) completion bash > .completion.bash
	./$(TARGET) completion zsh > .completion.zsh
	./$(TARGET) completion fish > .completion.fish
	./$(TARGET) completion powershell > .completion.ps1
	echo "Use 'source .completion.<YOUR SHELL>' to load completions"

delete:
	rm -rf ./tmp/blocks
print:
	go run $(OBJ) printchain