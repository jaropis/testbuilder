# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

TestBuilder is a command-line tool written in Rust for generating randomized LaTeX exam tests. The application takes plain text test files (questions with multiple choice answers), shuffles question and answer order, generates multiple LaTeX test variants, compiles them to PDFs, and optionally merges them.

There is a legacy React frontend in the `frontend/` directory that is not actively used or maintained.

## Architecture

### Main Application (Rust)
- Single binary command-line tool: `src/main.rs`
- Uses `rand` crate for shuffling answers
- Core logic flow:
  1. Parse command-line arguments (7 required arguments)
  2. Read and parse test file (questions with A)/B)/C)/D) answers and ANSWER lines)
  3. Generate N variants with shuffled questions/answers (stored in HashMap)
  4. Create working directory if needed
  5. Write LaTeX files for each variant
  6. Generate `test.sty` style file
  7. Compile each variant to PDF using `pdflatex` command (twice per file)
  8. Optionally merge PDFs using `pdfcpu` command

### Command-Line Arguments
1. Source file path
2. Output file base name (without extension, can include directory path)
3. Number of test variants to generate (integer)
4. Exam title (LaTeX formatted)
5. Before-test text (LaTeX formatted, e.g., name/surname fields)
6. "newpage" to add page break at end, or any other value to skip
7. "merge" to merge PDFs, or any other value to skip

## Test File Format

Plain text files with this structure:
```
Question text here?
A) First answer
B) Second answer
C) Third answer
D) Fourth answer (optional)
ANSWER B

Next question here?
A) Answer option
B) Answer option
...
```

Key parsing rules in `read_and_fill()`:
- Questions are lines that don't start with A)/B)/C)/D), aren't ANSWER lines, and aren't empty
- Answers follow immediately after questions (prefixed with A)/B)/C)/D))
- ANSWER lines are skipped (not used in current implementation)
- 2-answer sets (true/false) are NOT shuffled; 3+ answers ARE shuffled
- The text after the A)/B)/C)/D) prefix (starting at index 3) becomes the answer text

## Development Commands

### Building
```bash
cargo build                      # Debug build
cargo build --release            # Release build (optimized)
cargo run -- <args>              # Run with arguments
```

The compiled binary will be at:
- Debug: `target/debug/tb`
- Release: `target/release/tb`

### Running
```bash
./target/release/tb source.txt output_dir/test 3 "Exam Title" "Name: \underline{\hspace{10cm}}" newpage merge
```

**Dependencies**:
- Rust toolchain (cargo/rustc)
- `pdflatex` command in PATH for LaTeX compilation
- `pdfcpu` command for PDF merging (optional - gracefully skips if not found)
  - Checks: `/opt/homebrew/bin/pdfcpu`, `/usr/local/bin/pdfcpu`, then PATH via `which`

## Key Implementation Details

### LaTeX Generation
- Template structure defined in `generate_test()` function (src/main.rs:207-218)
- Custom style package `test.sty` generated automatically by `save_test_style()` (src/main.rs:95-154)
- Questions wrapped in `\item`, answers formatted inline with `\hspace` separators
- Percent signs escaped to `\%` to prevent LaTeX errors (src/main.rs:87)
- Two compilation passes run per file (src/main.rs:246-251) for proper LaTeX reference resolution

### File Generation Flow
1. Parse command-line arguments and validate
2. Read source file and parse into `HashMap<String, Vec<String>>` (question -> answers)
3. Create working directory if it doesn't exist
4. Generate LaTeX files (one per variant) in the working directory
5. Create `test.sty` in the working directory
6. Change to working directory, compile all LaTeX files, restore original directory
7. Sleep for `numFiles * 5 seconds` to ensure PDF generation completes
8. Optionally merge PDFs if "merge" argument provided

### Important Notes
- The program sleeps after compilation (src/main.rs:254) to ensure PDF generation completes
- Working directory is temporarily changed for `pdflatex` execution (src/main.rs:241-252)
- Answers with 2 or fewer options are NOT shuffled (typically true/false questions)
- Answers with 3+ options ARE shuffled using `rand::seq::SliceRandom`
- User must provide LaTeX-formatted text for exam title and before-test fields (raw underscores will cause LaTeX errors)
- PDF merging requires `pdfcpu` to be installed; if not found, merging is gracefully skipped with a warning
- Generated files include: `*.tex`, `*.pdf`, `*.aux`, `*.log`, and `test.sty`
