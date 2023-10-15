OBJ = main.go
TARGET = main.exe

all: build createbc

build:
	go build -o $(TARGET) $(OBJ)

createbc: build
	./$(TARGET) createblockchain -address "1JvKq46JWJzjK4dLJ1rAg8YBouaPoNzgX"

clean:
	$(RM) -rf tmp $(TARGET)
