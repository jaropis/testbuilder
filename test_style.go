package main

import "os"

func saveTestStyle() string {
	test_sty := `
	\ProvidesPackage{test}
%komendy: 
%odwa{pierwsz}{druga}{trzecia} - daje multiple choice z dwiema mo�liwo�ciami 
%\otrzy{pierwsz}{druga}{trzecia} - daje multiple choice z trzema mo�liwo�ciami 
%\ocztery{pierwsza}{druga}{trzecia}{czwarta} - z czterema
%\OO odst�p 6cm z podkre�leniem
%\oo odst�p 1.5cm z podkre�leniem
%\wybdob{tu wypisujemy list� s��w, kt�re pojawi� si� w ramce na grubo}
%trans{zdanie}{s�owo}{zdanie do transformacji z luk�}
\topmargin-1cm
\oddsidemargin-1cm
\evensidemargin-1cm
\textwidth17cm
\textheight24.5cm
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
	texFile, _ := os.Create("test.sty")
	defer texFile.Close()
	texFile.WriteString(test_sty)

}
