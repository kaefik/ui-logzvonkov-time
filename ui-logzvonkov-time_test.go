package main

import (
	"fmt"
	"testing"
)

func printDataConfigFile(d []DataConfigFile) {
	fmt.Println("Длина массива: ", len(d))
	for _, v := range d {
		fmt.Println(v)
	}
}

func TestRemoveItemFromDataConfigFile(t *testing.T) {
	nameConfigFile = "list-num-tel-test.cfg"
	dataConfigFile := readConfigFile(nameConfigFile)
	printDataConfigFile(dataConfigFile)

	dataConfigFile = removeItemFromDataConfigFile(dataConfigFile, 0)
	printDataConfigFile(dataConfigFile)
}

func TestSortFioRgDataConfigFile(t *testing.T) {
	nameConfigFile = "list-num-tel-test.cfg"
	dataConfigFile := readConfigFile(nameConfigFile)
	printDataConfigFile(dataConfigFile)

	dataConfigFile = sortFioRgDataConfigFile(dataConfigFile)
	printDataConfigFile(dataConfigFile)

}
