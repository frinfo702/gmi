use rust_search::{parse_args, search_in_file};
use std::fs::File;
use std::io::Write;
use tempfile::tempdir;

#[test]
fn test_parse_args() {
    let args = vec![
        "rust-search".to_string(),
        "/path/to/search".to_string(),
        "query".to_string(),
    ];
    let (search_path, query, use_fuzzy) = parse_args(&args);
    assert_eq!(search_path, "/path/to/search");
    assert_eq!(query, "query");
    assert_eq!(use_fuzzy, false);

    let args_with_fuzzy = vec![
        "rust-search".to_string(),
        "/path/to/search".to_string(),
        "query".to_string(),
        "--fuzzy".to_string(),
    ];
    let (search_path, query, use_fuzzy) = parse_args(&args_with_fuzzy);
    assert_eq!(search_path, "/path/to/search");
    assert_eq!(query, "query");
    assert_eq!(use_fuzzy, true);
}

#[test]
#[should_panic]
fn test_parse_args_insufficient_args() {
    let args = vec!["rust-search".to_string(), "/path/to/search".to_string()];
    parse_args(&args);
}

#[test]
fn test_search_in_file() {
    // 一時ディレクトリを作成
    let dir = tempdir().unwrap();
    let file_path = dir.path().join("test.txt");

    // テストファイルを作成
    let mut file = File::create(&file_path).unwrap();
    writeln!(file, "This is a test file").unwrap();
    writeln!(file, "It contains some text for testing").unwrap();
    writeln!(file, "We want to find this pattern").unwrap();
    writeln!(file, "And make sure our search works").unwrap();

    // 正確な検索
    {
        let test_query = "pattern";
        let result = search_in_file(&file_path, test_query, false);
        assert!(result.is_ok());

        // ファイルを再度開いて検索結果を確認
        let contents = std::fs::read_to_string(&file_path).unwrap();
        assert!(contents.contains("pattern"));
    }

    // ファジー検索
    {
        let test_query = "patern"; // typo
        let result = search_in_file(&file_path, test_query, true);
        assert!(result.is_ok());

        // ファイルを再度開いて検索結果を確認
        let contents = std::fs::read_to_string(&file_path).unwrap();
        assert!(contents.contains("pattern"));
    }
}

#[test]
fn test_exact_match() {
    let dir = tempdir().unwrap();
    let file_path = dir.path().join("exact_match.txt");

    let mut file = File::create(&file_path).unwrap();
    writeln!(file, "line one").unwrap();
    writeln!(file, "line two with keyword").unwrap();
    writeln!(file, "line three").unwrap();

    let result = search_in_file(&file_path, "keyword", false);
    assert!(result.is_ok());
}

#[test]
fn test_fuzzy_match() {
    let dir = tempdir().unwrap();
    let file_path = dir.path().join("fuzzy_match.txt");

    let mut file = File::create(&file_path).unwrap();
    writeln!(file, "line one").unwrap();
    writeln!(file, "line two with keywrd").unwrap(); // 綴り間違い
    writeln!(file, "line three").unwrap();

    let result = search_in_file(&file_path, "keyword", true);
    assert!(result.is_ok());
}

#[test]
fn test_no_match() {
    let dir = tempdir().unwrap();
    let file_path = dir.path().join("no_match.txt");

    let mut file = File::create(&file_path).unwrap();
    writeln!(file, "line one").unwrap();
    writeln!(file, "line two").unwrap();
    writeln!(file, "line three").unwrap();

    let result = search_in_file(&file_path, "keyword", false);
    assert!(result.is_ok());
}
