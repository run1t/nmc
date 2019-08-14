package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	. "github.com/logrusorgru/aurora"
)
var(
	autoStop    = flag.Bool("auto-stop", false, "Auto-stop rendering?")
	numTrackers = 1
)


func main() {

	var nodeModules []string

	root, _ := os.Getwd()

	fmt.Println("Fetching all",Blue("node_module"), "in", Magenta(root))
	totalNodeModule := 0
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {

		if info.IsDir() && info.Name() == ".git" {
			return filepath.SkipDir
		}

		if info.IsDir() && info.Name() == "node_modules" {
			nodeModules = append(nodeModules, path)
			totalNodeModule++
			fmt.Printf("\rWe found %v %s",  totalNodeModule, Blue("node_module"))

			return filepath.SkipDir
		}


		return nil
	})
	fmt.Println()
	if err != nil {
		panic(err)
	}
	fmt.Println("We are now calculating the size of",Blue("node_modules"))
	totalSize := 0
	i := 0
	for _, file := range nodeModules {
		i++
		size, _ := DirSize(file, i, len(nodeModules))
		chunk := strings.Split(file,"/")
		fmt.Printf("\r%5v/%v %6v Mb in %s", i , len(nodeModules),  Red(int(size)),  Blue(chunk[len(chunk)-4] +"/" + chunk[len(chunk)-3] +"/"+ chunk[len(chunk)-2] ))
		fmt.Println()

		totalSize += int(size)
	}

	fmt.Println("whoa you can free about", totalSize, "Mb", Blue("node_modules"))
	fmt.Print("Do you want to delete all the node_modules ?(Y/n): ")
	var input string
	fmt.Scanln(&input)
	if input != "Y" && input != "n" {
		input = "Y"
	}

	if input == "Y" {
		for _, file := range nodeModules {
			fmt.Println(Red("Deleting"),"node_modules in", file)
			os.RemoveAll(file)
		}
	}

	if input == "n" {
		os.Exit(0)
	}

	//fmt.Print(totalSize)
}

func DirSize(path string, number int, total int) (float64, error) {
	var size float64
	chunk := strings.Split(path,"/")
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += float64(info.Size())

			fmt.Printf("\r%5v/%v %6v Mb in %s", number , total,  Red(int( float64(size) / 1024.0 / 1024.0)),  Blue(chunk[len(chunk)-4] +"/" + chunk[len(chunk)-3] +"/"+ chunk[len(chunk)-2] ))
		}
		return err
	})
	sizeMB := float64(size) / 1024.0 / 1024.0
	return sizeMB, err
}