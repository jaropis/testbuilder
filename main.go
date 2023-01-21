package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
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

func main() {
	header := ` \documentclass[12pt]{article}
	\usepackage{test}
	\begin{document}
	\fancyhead[LO]{\rightmark{\textbf{Egzamin 2023}\hspace{\stretch{1}}Units 1-2.}}
   \begin{enumerate}`
	footer := `\end{enumerate}
	\end{document}`
	all_questions := readAndFill("assets/egzamin2022.txt")
	texFile, err := os.Create("texspace/exam.tex")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Fprintf(texFile, header)
	for key, value := range all_questions {
		oneLine := `\item`
		fmt.Fprintf(texFile, "%s %s\n\n", oneLine, key)
		for i, val := range value {
			text := ""
			if i == 0 {
				text = "a)"
			} else if i == 1 {
				text = "\\hspace{1cm}b)"
			} else {
				text = "\\hspace{1cm}c)"
			}
			fmt.Fprintf(texFile, "%s %s", text, val)
		}
		fmt.Fprint(texFile, "\n")
	}
	fmt.Fprint(texFile, footer)
}
