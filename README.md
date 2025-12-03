# testbuilder

## What it is
A command-line tool written in **Rust** for building tests in **LaTeX** from tests in the following form 

Why did the chicken cross the street?</br >
A) To talk to a friend<br />
B) To get to the other side<br />
C) The chicken did not cross the street</br>
ANSWER B

There should be many such questions, of course. Check out the example in the "example_test.txt"

The application will shuffle the order of the questions and also the order inside the answer - what was A will become e.g. C etc. 

You will get a nice looking test sheet after compiling with (pdf)LaTeX. 

## Building from source

**Requirements:**
- Rust (install from https://rustup.rs)
- LaTeX with `pdflatex` command
- `pdfcpu` (optional, only needed for PDF merging feature)

```bash
cargo build --release
```

The binary will be available at `target/release/tb`.

## How to use it

Run the `tb` binary with the following arguments:
1. Source file with the test in the above format
2. Output file name (say, "test" will get you a bunch of LaTeX files: "test1.tex", "test2.tex" etc. -- this can also include full path)
3. Number of different test variants to generate
4. Title of the test (e.g. "Medicine exam 2023")
5. What should go before the beginning of the test (LaTeX formatted, e.g. name/surname fields)
6. "newpage" to add a page break at the end, or "_" to skip
7. "merge" to merge all PDFs into one file, or "_" to skip

Example usage:
`./tb egzamin2023.txt results/res_test 3 "Biophysics, exam 2023" "Name, Surname, Group  \underline{\hspace{11.5cm}}" newpage merge`

This will output 3 LaTeX files named `res_test1.tex`, `res_test2.tex`, and `res_test3.tex`, along with the LaTeX style file `test.sty`. All files will be automatically compiled to PDF using `pdflatex`. Each test will have different question and answer orders. For example, the first resulting PDF may look like this:

![compiled_test1](res_test1.png)

and the third like this:

![compiled_test3](res_test3.png)
