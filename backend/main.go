package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/ledongthuc/pdf"
)

// This function checks if the current_value string passed as an argument does not have any of the prefixes "A)", "B)", "C)", "D)".
// It returns true if none of the prefixes are present, and false if any of them are present.
func getHeaderCondition(current_value string) bool {
	return !strings.HasPrefix(current_value, "A)") && !strings.HasPrefix(current_value, "B)") && !strings.HasPrefix(current_value, "C)") && !strings.HasPrefix(current_value, "D)")
}

// this function reads the file from the provided filepath and fills a map with the questions as the keys and their respective answers as the values.
func readAndFill(filepath string) map[string][]string {
	file, err := os.Open(filepath)
	all_questions := make(map[string][]string)
	var file_content []string
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		file_content = append(file_content, scanner.Text())
	}
	for idx := 0; idx < len(file_content); idx++ {
		current_value := file_content[idx]
		header_condition := getHeaderCondition(current_value)
		if header_condition && !strings.HasPrefix(current_value, "ANSWER") && current_value != "" {
			var new_set []string
			var counter int
			for !getHeaderCondition(file_content[idx+counter+1]) {
				modified_ans := file_content[idx+counter+1][3:]
				new_set = append(new_set, modified_ans)
				counter++
			}
			idx += counter
			all_questions[current_value] = new_set
		} else {
			idx++
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
	return all_questions
}

// function actually creating the test
func createTest(
	all_questions map[string][]string,
	exam_path string,
	header string,
	footer string) {
	texFile, err := os.Create(exam_path)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Fprintf(texFile, header)
	for key, value := range all_questions {
		oneLine := `\item`
		fmt.Fprintf(texFile, "%s %s\n\n", oneLine, key)
		// if the length of value is 3 or more, shuffle (2 is usually the true/false case)
		if len(value) > 2 {
			rand.Seed(time.Now().UnixNano())
			rand.Shuffle(len(value), func(i, j int) {
				value[i], value[j] = value[j], value[i]
			})
		}
		for i, val := range value {
			text := ""
			if i == 0 {
				text = "a)"
			} else if i == 1 {
				text = "\\hspace{1cm}b)"
			} else if i == 2 {
				text = "\\hspace{1cm}c)"
			} else if i == 3 {
				text = "\\hspace{1cm}d)"
			}
			v_percent := strings.Replace(val, "%", "\\%", -1)
			fmt.Fprintf(texFile, "%s~%s", text, v_percent)
		}
		fmt.Fprint(texFile, "\n")
	}
	fmt.Fprint(texFile, footer)
}

// function saving the test.sty file to the same location as the test LaTeX files (necessary for compilation)
func saveTestStyle(path, workingPath string) {
	test_sty := `
\ProvidesPackage{test}
\topmargin-1cm
\oddsidemargin-1cm
\evensidemargin-1cm
\textwidth17cm
\textheight23.5cm
\setlength{\parindent}{0pt}
\usepackage[T1]{fontenc}
\usepackage[utf8]{inputenc}
\usepackage{fancyhdr}
\usepackage{indentfirst}
\usepackage{amsmath}
\usepackage{amsfonts}
\usepackage{graphics}
\renewcommand{\headrulewidth}{0.5pt}
\addtolength{\headheight}{0.5pt}
\renewcommand{\thesection}{\Roman{section}}
\renewcommand{\arraystretch}{1.3}
\pagestyle{fancy}
\fancyhf{}
\cfoot{\thepage}
\pagenumbering{arabic}
\newcommand{\I}{\mathrm{i}}
\newcommand{\E}{\mathrm{e}}
\newcommand{\lin}{\mathrm{lin}}
\newcommand{\R}{\mathbb{R}}
\newcommand{\D}{$\diamond$}
\newcommand{\podpis}{\vfil\rightline{Jaros�aw Piskorski}}
\newcommand{\ramkon}[1]{\hfill\framebox{#1}}
\newcommand{\ramkitend}[1]{\framebox{#1\underline{\hspace{0.5cm}}}}
\newcommand{\ramkbeginitend}[1]{\framebox{\underline{\hspace{0.5cm}}#1\underline{\hspace{0.5cm}}}}
\newcommand{\doktor}{\vfil\rightline{dr Jaros�aw Piskorski}}
\newcommand{\OO}{\underline{\hspace{5cm}} }
\newcommand{\oo}{\underline{\hspace{1.0cm}} }
\newcommand{\od}{\underline{\hspace{0.5cm}} }
\newcommand{\odwa}[2]{\begin{tabular}{p{6cm}p{7cm}}(a) #1 & (b) #2 \end{tabular}}
\newcommand{\otrzy}[3]{\begin{tabular}{p{5cm}p{5cm}p{5cm}}(a) #1 & (b) #2 & (c)  #3\end{tabular}}
\newcommand{\ocztery}[4]{\vspace{1.8ex}\\ \begin{tabular}{p{4cm}p{4cm}p{4cm}p{4cm}}(a) #1 & (b) #2 & (c) #3 & (d) #4\end{tabular}}
\newcommand{\wybdob}[1]{\vspace*{1cm}\begin{center}\begin{tabular}{|p{14cm}|}\hline \textbf{#1} \\\hline\end{tabular}\end{center}}
\newcommand{\trans}[3]{\indent\begin{tabular}{p{3cm}l}
\textbf{#2} & \parbox[t]{16cm}{#1}\\
 & #3\\
\end{tabular}\\ \vspace{0.5cm}
}
\newcommand{\wektor}[1]{\overrightarrow{#1}}
	`
	texFile, _ := os.Create(filepath.Join(workingPath, "test.sty"))
	defer texFile.Close()
	texFile.WriteString(test_sty)
}

func getPdfFiles(fullpath string) ([]string, error) {
	var pdfFiles []string
	// Open the directory
	f, err := os.Open(fullpath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Read all files in the directory
	files, err := f.Readdir(-1) // -1 means to return all files
	if err != nil {
		return nil, err
	}

	// Filter for PDF files
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".pdf" {
			pdfFiles = append(pdfFiles, filepath.Join(fullpath, file.Name()))
		}
	}

	return pdfFiles, nil
}

func merge(fullpath string) {
	// get pdf files
	files, err := getPdfFiles(fullpath)
	if err != nil {
		fmt.Printf("error getting pdf files: %v\n", err)
		return
	}
	command := []string{"merge", path.Join(fullpath, "merged_test.pdf")}
	command = append(command, files...)
	cmd := exec.Command("/opt/homebrew/bin/pdfcpu", command...)
	cmd.Output()
}

func pageCounter(path string) (int, error) {
	f, r, err := pdf.Open(path)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	return r.NumPage(), err
}

// where it all happens
func generateTest(
    workingPath, inputFileName, outputFileName string,
    numFiles int,
    examTitle, beforeTest string, 
	newPage, merge_str bool,
) error {
	sourceFile := path.Join(workingPath, inputFileName)
	resultFile := path.Join(workingPath, outputFileName)
	const header_start = ` \documentclass[12pt]{article}
	\usepackage{test}
	\begin{document}
	\fancyhead[LO]{\rightmark{\textbf{`
	const header_middle = `}\hspace{\stretch{1}}}}`
	const header_end = `\begin{enumerate}`
	var footer1 = `\end{enumerate}
	\end{document}`
	var footer2 = `\end{enumerate}
	\newpage

	.

	\end{document}`
	var footer string
	// TO DO - handle in api
// 	if len(os.Args) != 8 {
// 		print(`There must be 5 arguments to the call
// 1) the source file with the test in format described by the README
// 2) the result file name - it will be extended by the number of the individual file
// 3) the number of files to generate (and integer number, of course)
// 4) the title of the examination, e.g. Egzamin 2023 (this is in Polish), but your can be in a different language ;)
// 5) what you want to go before the test, in my case this will be: "Imię, nazwisko i typ studiów:\underline{\hspace{11.5cm} }", which stands for name, surname and studies followed by an underlined space of length 11.5 cm  )
// 6) if you want a new page at the end of the test, write "newline", otherwise put _
// 7) if you want the resulting files to be merged, write "merge" - it will only merge file with an even number of pages so that printing is facilitated`)
// 		os.Exit(1)
// 	}
	if examTitle == "" {
		examTitle = "Egzamin 2023"
	} 

	if beforeTest == "" {
		beforeTest = `Imię, nazwisko i typ studiów:\underline{\hspace{11.5cm} }`
	} 

	if newPage {
		footer = footer2
	} else {
		footer = footer1
	}
	header := header_start + examTitle + header_middle + beforeTest + header_end
	all_questions := readAndFill(sourceFile)
	var outputFileNames []string
	for idx := 0; idx < numFiles; idx++ {
		outputfilename := outputFileName+strconv.Itoa(idx+1)+".tex"
		outputFileNames = append(outputFileNames, outputfilename)
		createTest(
			all_questions,
			path.Join(workingPath, outputfilename),
			header,
			footer)
	}
	
	saveTestStyle(resultFile, workingPath)
	oldPath, _ := os.Getwd()
	os.Chdir(workingPath)
	for _, filename := range outputFileNames {
		cmd := exec.Command("pdflatex", filename)
		cmd.Output()
		cmd2 := exec.Command("pdflatex", filename)
		cmd2.Output()
	}
	os.Chdir(oldPath)
	time.Sleep(time.Duration(numFiles) * 5 * time.Second)
	if merge_str {
		merge(workingPath)
	}
// TODO handle error
return nil
}

type RequestData struct {
    NumFiles     int    `schema:"numFiles"`
    ExamTitle    string `schema:"examTitle"`
    BeforeTest   string `schema:"beforeTest"`
    Merge        bool   `schema:"merge"`
    NewPage      bool   `schema:"newPage"`
    ResultFile   string `schema:"resultFile"`
}

func generateUniqueFolderNameWithUUID() string {
	return uuid.New().String()
}

func createFolder(folderName string) error {
	path := "./uploads/" + folderName
	return os.MkdirAll(path, 0755)
}

func generateTestHandler(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form
	err := r.ParseMultipartForm(10 << 20) // 10 MB as max size

	if err != nil {
		log.Println("Error parsing form:", err)  // This line prints the error to the console
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("sourceFile")
	if err != nil {
		http.Error(w, "Unable to get file", http.StatusBadRequest)
		return
	}

	defer file.Close()
	folderName :=generateUniqueFolderNameWithUUID()
	err = createFolder(folderName)
	if err != nil {
    	http.Error(w, "Unable to create folder on server", http.StatusInternalServerError)
    	return
	}
	workingPath := filepath.Join("uploads", folderName)
	fullPath := filepath.Join(workingPath, "filename.txt")
	dst, err := os.Create(fullPath)
	if err != nil {
    	http.Error(w, "Unable to create file on server", http.StatusInternalServerError)
    	return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
    	http.Error(w, "Unable to save file on server", http.StatusInternalServerError)
    	return
	}
	decoder := schema.NewDecoder()
	var requestData RequestData
	if err := decoder.Decode(&requestData, r.PostForm); err != nil {
		fmt.Println("Decoding error:", err)
    	http.Error(w, "Error decoding form data", http.StatusBadRequest)
    	return
	}

	// Now you have the file and requestData populated.
	// You can process the file and other form data as per your requirements.
	generateTest(
		workingPath, "filename.txt", requestData.ResultFile,
		requestData.NumFiles,
		requestData.ExamTitle, requestData.BeforeTest, requestData.NewPage, requestData.Merge,
	)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/generate-test", generateTestHandler).Methods("POST")

	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}
