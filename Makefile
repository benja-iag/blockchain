OBJ_MAIN = cmd/blockchain/main.go
TARGET_MAIN = main.exe

OBJ_NODE = cmd/node/main.go
TARGET_NODE = node.exe

all: build

build:
	go build -o $(TARGET_MAIN) $(OBJ_MAIN)
	go build -o $(TARGET_NODE) $(OBJ_NODE)

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