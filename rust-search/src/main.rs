use fuzzy_matcher::FuzzyMatcher;
use fuzzy_matcher::skim::SkimMatcherV2;
use std::env;
use std::fs::File;
use std::io::{self, BufRead, BufReader};
use std::path::Path;
use walkdir::WalkDir;

fn main() {
    let args: Vec<String> = env::args().collect();
    let (search_path, query, use_fuzzy) = parse_args(&args);

    if !Path::new(&search_path).exists() {
        eprintln!("Path not found: {}", search_path);
        std::process::exit(1);
    }

    let metadata = std::fs::metadata(&search_path).unwrap();
    if metadata.is_file() {
        if let Err(e) = search_in_file(&search_path, &query, use_fuzzy) {
            eprintln!("Error reading file {}: {}", search_path, e);
            std::process::exit(1);
        }
    } else if metadata.is_dir() {
        for entry in WalkDir::new(&search_path)
            .into_iter()
            .filter_map(|e| e.ok())
        {
            if entry.file_type().is_file() {
                let file_path = entry.path().to_string_lossy().to_string();
                if let Err(_) = search_in_file(&file_path, &query, use_fuzzy) {
                    continue;
                }
            }
        }
    } else {
        std::process::exit(0);
    }
}

/// コマンドライン引数を解析する。使用法:
/// ./rust-search <path> <query> [--fuzzy]
fn parse_args(args: &Vec<String>) -> (String, String, bool) {
    if args.len() < 3 {
        eprintln!(
            "Usage: {} <path> <query> [--fuzzy]",
            args.get(0).unwrap_or(&"rust-search".to_string())
        );
        std::process::exit(1);
    }
    let search_path = args[1].clone();
    let query = args[2].clone();
    let use_fuzzy = args.len() >= 4 && args[3] == "--fuzzy";
    (search_path, query, use_fuzzy)
}

/// 指定したファイル内で検索を実行する。use_fuzzy が true ならファジーマッチを利用する。
fn search_in_file<P: AsRef<Path>>(filepath: P, query: &str, use_fuzzy: bool) -> io::Result<()> {
    let file = File::open(&filepath)?;
    let reader = BufReader::new(file);
    let matcher = SkimMatcherV2::default();
    for (num, line_result) in reader.lines().enumerate() {
        let line_number = num + 1;
        let line = match line_result {
            Ok(text) => text,
            Err(_) => return Err(io::Error::new(io::ErrorKind::InvalidData, "Non-text data")),
        };
        let matched = if use_fuzzy {
            matcher.fuzzy_match(&line, query).is_some()
        } else {
            line.contains(query)
        };

        if matched {
            println!("{}:{}: {}", filepath.as_ref().display(), line_number, line);
        }
    }
    Ok(())
}
