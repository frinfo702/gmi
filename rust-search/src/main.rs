use anyhow::Result;
use std::{env, fs, path::Path};
use walkdir::WalkDir;

fn main() -> Result<()> {
    let args: Vec<String> = env::args().collect();

    let mut query = None;
    let mut target_dir = "."; // デフォルトはカレントディレクトリ
    let mut fuzzy_mode = false;
    let mut arg_index = 1;
    while arg_index < args.len() {
        match args[arg_index].as_str() {
            "--dir" | "-d" => {
                arg_index += 1;
                if arg_index < args.len() {
                    target_dir = args[arg_index].as_str();
                } else {
                    eprintln!("エラー: --dir オプションにはディレクトリパスが必要です。");
                    std::process::exit(1);
                }
            }
            "--fuzzy" | "-f" => {
                fuzzy_mode = true;
            }
            _ => {
                if query.is_none() {
                    query = Some(&args[arg_index]);
                } else {
                    eprintln!(
                        "エラー: 不明な引数または複数のクエリが指定されました: {}",
                        args[arg_index]
                    );
                    print_usage(&args[0]);
                    std::process::exit(1);
                }
            }
        }
        arg_index += 1;
    }

    if query.is_none() {
        print_usage(&args[0]);
        std::process::exit(1);
    }
    let query = query.unwrap();

    println!(
        "検索結果 (ディレクトリ: {}, クエリ: {}, Fuzzy: {}):",
        target_dir, query, fuzzy_mode
    );
    search_in_directory(target_dir, query, fuzzy_mode)?;

    Ok(())
}

fn print_usage(program_name: &str) {
    println!(
        "使用法: {} <検索クエリ> [-d <ディレクトリ>] [-f | --fuzzy]",
        program_name
    );
    println!("  <検索クエリ>: 検索する文字列");
    println!(
        "  -d, --dir <ディレクトリ>: 検索対象のディレクトリ (デフォルト: カレントディレクトリ)"
    );
    println!("  -f, --fuzzy: ファジー検索を有効にする");
}

fn search_in_directory(dir: &str, query: &str, fuzzy_mode: bool) -> Result<()> {
    for entry in WalkDir::new(dir)
        .into_iter()
        .filter_map(|e| e.ok())
        .filter(|e| e.file_type().is_file())
    {
        let path = entry.path();

        // バイナリファイルや特定の拡張子ファイルを除外
        if should_skip_file(path) {
            continue;
        }

        match fs::read_to_string(path) {
            Ok(content) => {
                search_in_content(path, &content, query, fuzzy_mode)?;
            }
            Err(_) => {
                // ファイルの読み込みに失敗した場合は無視して続行
                continue;
            }
        }
    }

    Ok(())
}

fn should_skip_file(path: &Path) -> bool {
    // バイナリファイル、巨大ファイル、特定の拡張子を持つファイルを除外
    let extension = path.extension().and_then(|e| e.to_str()).unwrap_or("");
    let skip_extensions = ["exe", "dll", "so", "dylib", "bin", "obj", "o", "a", "lib"];

    // node_modules や .git ディレクトリ内のファイルを除外
    let path_str = path.to_str().unwrap_or("");
    if path_str.contains("node_modules") || path_str.contains(".git") {
        return true;
    }

    skip_extensions.contains(&extension)
}

fn search_in_content(path: &Path, content: &str, query: &str, fuzzy_mode: bool) -> Result<()> {
    let path_str = path.to_str().unwrap_or("unknown_path");
    let query_lower = query.to_lowercase();

    if fuzzy_mode {
        // fuzzy-matcherを使用したファジー検索（今回はMVPなので実装しない）
        // fuzzy_search_in_content(path_str, content, &query_lower)?;
    } else {
        // 単純なサブストリング検索
        for (line_number, line) in content.lines().enumerate() {
            let line_lower = line.to_lowercase();
            if line_lower.contains(&query_lower) {
                println!("{}:{}:{}", path_str, line_number + 1, line);
            }
        }
    }

    Ok(())
}
