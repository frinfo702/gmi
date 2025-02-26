use rust_search::{parse_args, search_in_file};
use std::env;
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
