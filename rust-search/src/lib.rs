use fuzzy_matcher::FuzzyMatcher;
use fuzzy_matcher::skim::SkimMatcherV2;
use std::fs::File;
use std::io::{self, BufRead, BufReader};
use std::path::Path;

/// コマンドライン引数を解析する
/// 使用法: ./rust-search <path> <query> [--fuzzy]
pub fn parse_args(args: &Vec<String>) -> (String, String, bool) {
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

/// 指定したファイル内で検索を実行する
pub fn search_in_file<P: AsRef<Path>>(filepath: P, query: &str, use_fuzzy: bool) -> io::Result<()> {
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
