use std::{env, path::Path};

use anyhow::Result;

fn main() -> Result<()> {
    // コマンドライン引数で対象ディレクトリを指定。なければデフォルトで"./data"にする
    let args: Vec<String> = env::args().collect();
    let target_dir = if args.len() > 1 { &args[1] } else { "./data" };

    println!("Indexing files from directory: {}", target_dir);

    // ディレクトリ内のファイルを再帰的にクロールして文書データに変換する
    let document = crawl_directory(target_dir)?;
    println!("Found {} documents.", document);

    // Tantivyを利用してインデックスを作成する

    // 標準入力からクエリを受け取り検索を実施する
    Ok(())
}
