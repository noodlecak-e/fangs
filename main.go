package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
)

func main() {
	if runtime.GOOS != "darwin" {
		log.Fatal("only macos is supported\n")
	}

	cmd := exec.Command("command", "-v", "brew")
	if err := cmd.Run(); err != nil {
		log.Fatal("brew is not installed\n")
	}

	utils := []string{"pylint", "flake8", "black", "mypy", "bandit", "isort"}

	var wg sync.WaitGroup
	for _, util := range utils {
		wg.Add(1)

		go func(util string) {
			defer wg.Done()
			fmt.Println("installing " + util + " if not present")
			cmd = exec.Command("brew", "install", util)
			if err := cmd.Run(); err != nil {
				log.Fatal(util + " is not installed\n")
			}
		}(util)
	}

	wg.Wait()

	var i string
	fmt.Print("Provide project path to check: ")
	fmt.Scan(&i)

	i = filepath.Clean(i)

	cmd = exec.Command("rm", "-rf", "checker_outputs")
	if err := cmd.Run(); err != nil {
		log.Fatal("rm failed\n")
	}

	cmd = exec.Command("mkdir", "-p", "checker_outputs")
	if err := cmd.Run(); err != nil {
		log.Fatal("mkdir failed\n")
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		pylintRunner(i)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		flake8Runner(i)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		blackRunner(i)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		mypyRunner(i)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		banditRunner(i)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		isortRunner(i)
	}()

	wg.Wait()
}

func flake8Runner(i string) {
	f, err := os.Create("checker_outputs/flake8_output.txt")
	if err != nil {
		log.Fatal("error creating file\n")
	}
	defer f.Close()

	cmd := exec.Command("flake8", i)
	cmd.Stdout = f
	cmd.Run()
}

func blackRunner(i string) {
	f, err := os.Create("checker_outputs/black_output.txt")
	if err != nil {
		log.Fatal("error creating file\n")
	}
	defer f.Close()

	cmd := exec.Command("black", "--check", "-l", "120", i)
	cmd.Stdout = f
	cmd.Run()
}

func mypyRunner(i string) {
	f, err := os.Create("checker_outputs/mypy_output.txt")
	if err != nil {
		log.Fatal("error creating file\n")
	}
	defer f.Close()

	cmd := exec.Command("mypy", i)
	cmd.Stdout = f
	cmd.Run()
}

func banditRunner(i string) {
	f, err := os.Create("checker_outputs/bandit_output.txt")
	if err != nil {
		log.Fatal("error creating file\n")
	}
	defer f.Close()

	cmd := exec.Command("bandit", "-r", i)
	cmd.Stdout = f
	cmd.Run()
}

func isortRunner(i string) {
	f, err := os.Create("checker_outputs/isort_output.txt")
	if err != nil {
		log.Fatal("error creating file\n")
	}
	defer f.Close()

	cmd := exec.Command("isort", "--check-only", i)
	cmd.Stdout = f
	cmd.Run()
}

func pylintRunner(i string) {
	f, err := os.Create("checker_outputs/pylint_output.txt")
	if err != nil {
		log.Fatal("error creating file\n")
	}
	defer f.Close()

	cmd := exec.Command("pylint", i)
	cmd.Stdout = f
	cmd.Run()
}
