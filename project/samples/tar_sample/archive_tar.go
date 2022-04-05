package tar_sample

import (
	"archive/tar"
	"fmt"
)

func Constants() {
	constants := []interface{}{
		tar.TypeReg,  //　通常のファイル
		tar.TypeRegA, // 非推奨

		// 下記のグループはヘッダのみでボディを持たない。
		tar.TypeLink,    // ハードリンク
		tar.TypeSymlink, // シムリンク
		tar.TypeChar,    // ドライバ(マウス、キーボードなど)
		tar.TypeBlock,   // ドライバ(HDD,フロッピーディスクなど)
		tar.TypeDir,     // ディレクトリ
		tar.TypeFifo,    // 名前付きパイプ,FIFO

		// 使わん
		tar.TypeCont, // 予約してあるだけで用途なし。

		// https://qiita.com/ko1nksm/items/fbcff63639c5d141e76d
		// POSIXって規格のことなんですね...:'(
		tar.TypeXHeader,
		tar.TypeXGlobalHeader,

		tar.TypeGNUSparse,   // スパースファイル(GNUフォーマット)
		tar.TypeGNULongName, // GNU形式のメタファイル
		tar.TypeGNULongName, // GNU形式のメタファイル
	}
	for _, constant := range constants {
		fmt.Printf("%c\n", constant)
	}
}
