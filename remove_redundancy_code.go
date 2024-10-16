// remove_redundancy_code.go
package main

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func get_invisible_variant(filename string) []string {
	f, err := os.Open(filename)
	if err != nil {
		panic("open " + filename + " failed")
	}
	defer f.Close()
	bfRd := bufio.NewReader(f)
	var variant_list []string
	for {
		line, err := bfRd.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			break
		}
		if strings.Index(line, ".Visible = false;") != -1 {
			//this.hsLabel5.Visible = false;
			variant1 := strings.ReplaceAll(line, "this.", "")
			variant2 := strings.ReplaceAll(variant1, ".Visible = false;", "")
			variant3 := strings.ReplaceAll(variant2, "\n", "")
			variant4 := strings.ReplaceAll(variant3, "\r", "")
			variant := strings.Trim(variant4, " ")
			variant_list = append(variant_list, variant)
		}
	}
	return variant_list
}
func get_file_ext(filepath string) string {
	id := strings.LastIndex(filepath, "\\")
	if id == -1 {
		return ""
	} else {
		file1 := filepath[id+1:]
		filename := path.Base(file1)
		idpoint := strings.Index(filename, ".")
		if idpoint == -1 {
			return ""
		} else {
			fileext := filename[idpoint+1:]
			return string(fileext)
		}
	}
}
func report_invisible_variant_status(filename_designer string, filename_cs string, variants []string) string {
	var line_list []string
	resultstr := ""
	restr := ""
	f, err := os.Open(filename_cs)
	if err != nil {
		fmt.Println("open " + filename_cs + " failed")
		return restr
	}
	defer f.Close()
	bfRd := bufio.NewReader(f)
	for {
		line, err := bfRd.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			break
		}
		linetrim := strings.Trim(line, " ")
		line_list = append(line_list, linetrim)
	}
	for _, item := range variants {
		used := false
		for _, linestr := range line_list {
			if strings.Index(linestr, item) != -1 {
				used = true
				break
			}
		}
		if !used {
			resultstr = resultstr + " " + item
		}
	}
	if len(resultstr) != 0 {
		restr = filename_designer + " doesn't show these [" + fmt.Sprintf(strings.Join(variants, ",")) + "] controls.\n"
		restr = restr + filename_cs + " doesn't used these [" + strings.Trim(resultstr, " ") + "]"
	}
	return restr
}
func browser_folder(folder string) int {
	fmt.Println("start research ", folder)
	icount := 0
	restr := ""
	err := filepath.Walk(folder, func(file_path string, info fs.FileInfo, err error) error {
		if !info.IsDir() {
			file_path_low := strings.ToLower(file_path)
			fileSuffix := get_file_ext(file_path_low)
			if strings.EqualFold(fileSuffix, "designer.cs") {
				fmt.Println("check:" + file_path_low)
				variants := get_invisible_variant(file_path_low)
				file_name_cs := strings.ReplaceAll(file_path_low, "designer.", "")
				report := report_invisible_variant_status(file_path_low, file_name_cs, variants)
				if len(report) > 0 {
					restr = restr + report + "\n"
					fmt.Println(report)
					icount++
				}
			}
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	///////////////////////////
	file_report := "report.log"
	file2, err := os.OpenFile(file_report, os.O_CREATE, 0666)
	if err != nil {
		fmt.Printf("openfile file2 err : %v\n", err)
	}
	if file2 != nil {
		defer func(file *os.File) { file.Close() }(file2)
	}
	_, err = file2.WriteString(restr + "\n")
	if err != nil {
		fmt.Printf("file2 write string err : %v\n", err)
	}
	return icount
}
func main() {
	fmt.Printf("Program start!\n")
	if len(os.Args) == 2 {
		browser_folder(os.Args[1])
	} else {
		fmt.Println("Wrong parameters.")
	}
}
