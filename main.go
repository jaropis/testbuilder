package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func getHeaderCondition(current_value string) bool {
	return !strings.HasPrefix(current_value, "A)") && !strings.HasPrefix(current_value, "B)") && !strings.HasPrefix(current_value, "C)") && !strings.HasPrefix(current_value, "D)")
}

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
			} else {
				text = "\\hspace{1cm}c)"
			}
			fmt.Fprintf(texFile, "%s~%s", text, val)
		}
		fmt.Fprint(texFile, "\n")
	}
	fmt.Fprint(texFile, footer)
}

func saveTestStyle(path string) {
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
	wd, _ := os.Getwd()
	texFile, _ := os.Create(filepath.Join(wd, filepath.Dir(path), "test.sty"))
	defer texFile.Close()
	texFile.WriteString(test_sty)
}

func main() {

	const header_start = ` \documentclass[12pt]{article}
	\usepackage{test}
	\begin{document}
	\fancyhead[LO]{\rightmark{\textbf{`
	const header_middle = `}\hspace{\stretch{1}}}}`
	const header_end = `\begin{enumerate}`
	const footer = `\end{enumerate}
	\end{document}`
	if len(os.Args) != 6 {
		print(`There must be 5 arguments to the call
1) the source file with the test in format described by the README
2) the result file name - it will be extended by the number of the individual file
3) the number of files to generate (and integer number, of course)
4) the title of the examination, e.g. Egzamin 2023 (this is in Polish), but your can be in a different language ;)
5) what you want to go before the test, in my case this will be: "Imię, nazwisko i typ studiów:\underline{\hspace{11.5cm} }", which stands for name, surname and studies followed by an underlined space of length 11.5 cm  )`)
		os.Exit(1)
	}
	var examTitle string
	var beforeTest string
	if os.Args[4] == "_" {
		examTitle = "Egzamin 2023"
	} else {
		examTitle = os.Args[4]
	}

	if os.Args[5] == "_" {
		beforeTest = `Imię, nazwisko i typ studiów:\underline{\hspace{11.5cm} }`
	} else {
		beforeTest = os.Args[5]
	}
	header := header_start + examTitle + header_middle + beforeTest + header_end
	all_questions := readAndFill(os.Args[1])
	loopCount, _ := strconv.Atoi(os.Args[3])
	var outputFileNames []string
	wd, _ := os.Getwd()
	for idx := 0; idx < loopCount; idx++ {
		outputfilename := filepath.Join(wd, os.Args[2]+strconv.Itoa(idx+1)+".tex")
		outputFileNames = append(outputFileNames, outputfilename)
		createTest(
			all_questions,
			outputfilename,
			header,
			footer)
	}
	saveTestStyle(os.Args[2])
	oldPath := wd
	os.Chdir(filepath.Dir(outputFileNames[0]))
	for _, filename := range outputFileNames {
		cmd := exec.Command("/Library/TeX/texbin/pdflatex", filename)
		cmd.Output()
		cmd2 := exec.Command("/Library/TeX/texbin/pdflatex", filename)
		cmd2.Output()
	}
	os.Chdir(oldPath)
}
