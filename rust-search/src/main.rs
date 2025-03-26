use anyhow::{Context, Result};
use std::io::BufRead;
use std::{env, fs, io, path::Path};
use walkdir::WalkDir;

fn main() -> Result<()> {
    let args: Vec<String> = env::args().collect();

    if args.len() < 2 {
        println!("使用法: {} <検索クエリ> [--fuzzy]", args[0]);
        return Ok(());
    }

    let query = &args[1];
    let fuzzy_mode = args.len() > 2 && args[2] == "--fuzzy";

    // カレントディレクトリを検索対象とする
    let target_dir = ".";

    println!("検索結果:");
    search_in_directory(target_dir, query, fuzzy_mode)?;

    Ok(())
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
