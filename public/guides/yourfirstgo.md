# Running your first Go program

Congrats, if you've made it here, you should have Go installed and Hackatime set up! Now, let's run your first Go program.

## Creating a module

1. Open your terminal and make a new directory for your Go project:

   ```bash
   mkdir my-first-go
   cd my-first-go
   ```

2. Initialize a new Go module (this tracks your dependencies and project information):

   ```bash
   go mod init my-first-go
   ```

## Hello World

1. Create a new file called `main.go` in your project directory

2. Open `main.go` in your code editor and add the following code:

   ```go
   package main

   import "fmt"

   func main() {
       fmt.Println("Hello, World!")
   }
   ```

- The "package" line defines the package name, `main` for the entry point, and anything else for libraries you can import later (e.g. `package mylib`, then `import "my-first-go/mylib"`).
- The "import" line imports the `fmt` package, which contains functions for formatting and printing text.
- The `main` function is the entry point of the program, and `fmt.Println` prints "Hello, World!" to the console.

## Running the program

You have a few options to run your Go program:

1. In the terminal, run:

   ```bash
   go run .
   ```

   This compiles and runs the program in one step, but does not create a permanent executable file.

2. To create an executable file, run:

   ```bash
    go build .
   ```

   This will create an executable file named `my-first-go` (or `my-first-go.exe` on Windows) in your project directory. You can then run it from your terminal:

   ```bash
   ./my-first-go
   ```

Congratulations, you've just run your first Go program! You should see "Hello, World!" printed in your terminal.

A few things you could try doing next:

- Run an HTTP server using the `net/http` package
- Create a CLI tool using the `flag` package
- Explore Go's concurrency features with goroutines and channels
- Read the [Go documentation](https://go.dev/doc/) for more examples and tutorials!
