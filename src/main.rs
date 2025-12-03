use rand::seq::SliceRandom;
use rand::thread_rng;
use std::collections::HashMap;
use std::env;
use std::fs::{self, File};
use std::io::{BufRead, BufReader, Write};
use std::path::{Path, PathBuf};
use std::process::Command;
use std::thread;
use std::time::Duration;

/// Checks if the current value does not have any of the prefixes A), B), C), D)
fn get_header_condition(current_value: &str) -> bool {
    !current_value.starts_with("A)")
        && !current_value.starts_with("B)")
        && !current_value.starts_with("C)")
        && !current_value.starts_with("D)")
}

/// Reads the file and fills a map with questions as keys and their respective answers as values
fn read_and_fill(filepath: &str) -> Result<HashMap<String, Vec<String>>, std::io::Error> {
    let file = File::open(filepath)?;
    let reader = BufReader::new(file);
    let file_content: Vec<String> = reader.lines().collect::<Result<Vec<_>, _>>()?;

    let mut all_questions = HashMap::new();
    let mut idx = 0;

    while idx < file_content.len() {
        let current_value = &file_content[idx];
        let header_condition = get_header_condition(current_value);

        if header_condition && !current_value.starts_with("ANSWER") && !current_value.is_empty() {
            let mut new_set = Vec::new();
            let mut counter = 0;

            while idx + counter + 1 < file_content.len()
                && !get_header_condition(&file_content[idx + counter + 1])
            {
                let modified_ans = &file_content[idx + counter + 1][3..];
                new_set.push(modified_ans.to_string());
                counter += 1;
            }

            idx += counter;
            all_questions.insert(current_value.clone(), new_set);
        }
        idx += 1;
    }

    Ok(all_questions)
}

/// Creates a test LaTeX file
fn create_test(
    all_questions: &HashMap<String, Vec<String>>,
    exam_path: &Path,
    header: &str,
    footer: &str,
) -> Result<(), std::io::Error> {
    let mut tex_file = File::create(exam_path)?;
    write!(tex_file, "{}", header)?;

    for (key, value) in all_questions.iter() {
        let one_line = r"\item";
        writeln!(tex_file, "{} {}\n", one_line, key)?;

        let mut answers = value.clone();
        // If the length is 3 or more, shuffle (2 is usually the true/false case)
        if answers.len() > 2 {
            let mut rng = thread_rng();
            answers.shuffle(&mut rng);
        }

        for (i, val) in answers.iter().enumerate() {
            let text = match i {
                0 => "a)",
                1 => r"\hspace{1cm}b)",
                2 => r"\hspace{1cm}c)",
                3 => r"\hspace{1cm}d)",
                _ => "",
            };
            let v_percent = val.replace('%', r"\%");
            write!(tex_file, "{}~{}", text, v_percent)?;
        }
        writeln!(tex_file)?;
    }

    write!(tex_file, "{}", footer)?;
    Ok(())
}

/// Saves the test.sty file to the specified path
fn save_test_style(working_path: &Path) -> Result<(), std::io::Error> {
    let test_sty = r#"
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
\newcommand{\podpis}{\vfil\rightline{Jarosław Piskorski}}
\newcommand{\ramkon}[1]{\hfill\framebox{#1}}
\newcommand{\ramkitend}[1]{\framebox{#1\underline{\hspace{0.5cm}}}}
\newcommand{\ramkbeginitend}[1]{\framebox{\underline{\hspace{0.5cm}}#1\underline{\hspace{0.5cm}}}}
\newcommand{\doktor}{\vfil\rightline{dr Jarosław Piskorski}}
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
	"#;

    let mut tex_file = File::create(working_path.join("test.sty"))?;
    tex_file.write_all(test_sty.as_bytes())?;
    Ok(())
}

/// Gets all PDF files in a directory
fn get_pdf_files(fullpath: &Path) -> Result<Vec<PathBuf>, std::io::Error> {
    let mut pdf_files = Vec::new();

    for entry in fs::read_dir(fullpath)? {
        let entry = entry?;
        let path = entry.path();
        if path.is_file() && path.extension().and_then(|s| s.to_str()) == Some("pdf") {
            pdf_files.push(path);
        }
    }

    Ok(pdf_files)
}

/// Checks if pdfcpu is available
fn find_pdfcpu() -> Option<String> {
    // Try common locations
    if Path::new("/opt/homebrew/bin/pdfcpu").exists() {
        return Some("/opt/homebrew/bin/pdfcpu".to_string());
    }
    if Path::new("/usr/local/bin/pdfcpu").exists() {
        return Some("/usr/local/bin/pdfcpu".to_string());
    }

    // Try to find in PATH
    if let Ok(output) = Command::new("which").arg("pdfcpu").output() {
        if output.status.success() {
            let path = String::from_utf8_lossy(&output.stdout).trim().to_string();
            if !path.is_empty() {
                return Some(path);
            }
        }
    }

    None
}

/// Merges PDF files using pdfcpu
fn merge(fullpath: &Path) -> Result<(), Box<dyn std::error::Error>> {
    let files = get_pdf_files(fullpath)?;
    if files.is_empty() {
        return Ok(());
    }

    let pdfcpu_path = match find_pdfcpu() {
        Some(path) => path,
        None => {
            eprintln!("Warning: pdfcpu not found. Skipping PDF merge.");
            eprintln!("Install pdfcpu to enable PDF merging: https://pdfcpu.io/");
            return Ok(());
        }
    };

    let mut command_args = vec!["merge".to_string()];
    command_args.push(fullpath.join("merged_test.pdf").to_string_lossy().to_string());

    for file in files {
        command_args.push(file.to_string_lossy().to_string());
    }

    let output = Command::new(&pdfcpu_path)
        .args(&command_args)
        .output()?;

    if !output.status.success() {
        eprintln!("Warning: PDF merge failed.");
        eprintln!("pdfcpu stderr: {}", String::from_utf8_lossy(&output.stderr));
    } else {
        println!("PDFs merged successfully into merged_test.pdf");
    }

    Ok(())
}

/// Main test generation function
fn generate_test(
    source_file: &str,
    output_file: &str,
    num_files: usize,
    exam_title: &str,
    before_test: &str,
    new_page: bool,
    merge_pdfs: bool,
) -> Result<(), Box<dyn std::error::Error>> {
    let header_start = r#" \documentclass[12pt]{article}
	\usepackage{test}
	\begin{document}
	\fancyhead[LO]{\rightmark{\textbf{"#;
    let header_middle = r#"}\hspace{\stretch{1}}}}"#;
    let header_end = r#"\begin{enumerate}"#;
    let footer1 = r#"\end{enumerate}
	\end{document}"#;
    let footer2 = r#"\end{enumerate}
	\newpage

	.

	\end{document}"#;

    let footer = if new_page { footer2 } else { footer1 };
    let header = format!("{}{}{}{}{}", header_start, exam_title, header_middle, before_test, header_end);

    let all_questions = read_and_fill(source_file)?;

    let output_path = Path::new(output_file);
    let working_dir = output_path.parent().unwrap_or_else(|| Path::new("."));

    // Create working directory if it doesn't exist
    fs::create_dir_all(working_dir)?;

    let output_base = output_path.file_stem().unwrap_or_default().to_string_lossy();
    let mut output_filenames = Vec::new();

    for idx in 0..num_files {
        let output_filename = format!("{}{}.tex", output_base, idx + 1);
        output_filenames.push(output_filename.clone());
        let output_tex_path = working_dir.join(&output_filename);
        create_test(&all_questions, &output_tex_path, &header, footer)?;
    }

    save_test_style(working_dir)?;

    // Save current directory and change to working directory
    let old_dir = env::current_dir()?;
    env::set_current_dir(working_dir)?;

    // Compile each LaTeX file twice
    for filename in &output_filenames {
        Command::new("pdflatex")
            .arg(filename)
            .output()?;
        Command::new("pdflatex")
            .arg(filename)
            .output()?;
    }

    // Restore original directory
    env::set_current_dir(old_dir)?;

    // Sleep to ensure PDF generation completes
    thread::sleep(Duration::from_secs((num_files * 5) as u64));

    if merge_pdfs {
        merge(working_dir)?;
    }

    Ok(())
}

fn main() {
    let args: Vec<String> = env::args().collect();

    if args.len() != 8 {
        eprintln!("There must be 7 arguments to the call:");
        eprintln!("1) the source file with the test in format described by the README");
        eprintln!("2) the result file name - it will be extended by the number of the individual file");
        eprintln!("3) the number of files to generate (an integer number, of course)");
        eprintln!("4) the title of the examination, e.g. Egzamin 2023");
        eprintln!("5) what you want to go before the test, e.g. \"Imię, nazwisko i typ studiów:\\underline{{\\hspace{{11.5cm}}}}\"");
        eprintln!("6) if you want a new page at the end of the test, write \"newpage\", otherwise put _");
        eprintln!("7) if you want the resulting files to be merged, write \"merge\"");
        std::process::exit(1);
    }

    let source_file = &args[1];
    let output_file = &args[2];
    let num_files: usize = args[3].parse().unwrap_or_else(|_| {
        eprintln!("Error: Number of files must be a valid integer");
        std::process::exit(1);
    });
    let exam_title = &args[4];
    let before_test = &args[5];
    let new_page = args[6] == "newpage";
    let merge_pdfs = args[7] == "merge";

    if let Err(e) = generate_test(
        source_file,
        output_file,
        num_files,
        exam_title,
        before_test,
        new_page,
        merge_pdfs,
    ) {
        eprintln!("Error generating test: {}", e);
        std::process::exit(1);
    }

    println!("Test generation completed successfully!");
}
