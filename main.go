package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
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
			fmt.Fprintf(texFile, "%s %s", text, val)
		}
		fmt.Fprint(texFile, "\n")
	}
	fmt.Fprint(texFile, footer)
}

func main() {

	const header, footer = ` \documentclass[12pt]{article}
	\usepackage{test}
	\begin{document}
	\fancyhead[LO]{\rightmark{\textbf{Egzamin 2023}\hspace{\stretch{1}}}}
	Imię, nazwisko i typ studiów:\underline{\hspace{11.5cm} }
   \begin{enumerate}`, `\end{enumerate}
	\end{document}`

	all_questions := readAndFill("assets/egzamin2022.txt")
	for _, val := range []int{1, 2} {
		createTest(
			all_questions,
			"texspace/test"+strconv.Itoa(val)+".tex",
			header,
			footer)
	}
}
