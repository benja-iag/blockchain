OBJ = main.go
TARGET = main.exe

all: build

build:
	go build -o $(TARGET) $(OBJ)

createbc: build
	./$(TARGET) createblockchain -address "1JvKq46JWJzjK4dLJ1rAg8YBouaPoNzgX"

clean:
	$(RM) -rf tmp $(TARGET) .completion.*

completion: build
	./$(TARGET) completion bash > .completion.bash
	./$(TARGET) completion zsh > .completion.zsh
	./$(TARGET) completion fish > .completion.fish
	./$(TARGET) completion powershell > .completion.ps1
	echo "Use 'source .completion.<YOUR SHELL>' to load completions"
