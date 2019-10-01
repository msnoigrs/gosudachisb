# gosudachisb

gosudachisbは、[gosudachi](https://github.com/msnoigrs/gosudachi)を使うサンプルプログラムです。サンプルとはいえ、形態素解析器として問題なく使用できます。[gosudachi](https://github.com/msnoigrs/gosudachi)を、プログラムから使用する方法を提示する目的で作成しました。

2つのプログラムがあります。

-   gosudachisb
-   gosudachiclisb

どちらも、せっかくのGo言語による実装ですので、辞書も含めて、シングルバイナリが生成されます。辞書をバイナリに含めるテクニックには、[ichiban/assets](https://github.com/ichiban/assets)を利用しています。作成過程で、辞書をzipファイルに入れてますが、あえて圧縮はしていません。gosudachiの辞書ファイルは、メモリイメージをシリアライズした形式になっているためです。


## gosudachisb

プラグイン設定を、設定ファイルから読み出すのではなく、プログラム中で行っています。


### ビルド方法

ビルドを行うと、同梱する辞書の違いにより、以下の3つのバイナリが生成されます。

-   gosudachisbsmall
-   gosudachisbcore
-   gosudachisbfull

    git clone https://github.com/msnoigrs/gosudachisb.git
    cd gosudachisb/gosudachisb
    bash ./build.sh

windows版のバイナリを生成することもできます。

    bash ./build_win.sh


### コマンド

    $ gosudachisbsmall [-m mode] [-a] [-d] [-o output] [file...]
    $ gosudachisbcore [-m mode] [-a] [-d] [-o output] [file...]
    $ gosudachisbfull [-m mode] [-a] [-d] [-o output] [file...]


#### オプション

-   -m {A|B|C}分割モード
-   -a 読み、辞書形も出力
-   -d デバッグ情報の出力
-   -o 出力ファイル（指定がない場合は標準出力）
-   -f エラーを無視して処理を続行する


#### 出力例

    $ echo 東京都へ行く | gosudachisbcore
    東京都  名詞,固有名詞,地名,一般,*,*     東京都
    へ      助詞,格助詞,*,*,*,*     へ
    行く    動詞,非自立可能,*,*,五段-カ行,終止形-一般       行く
    EOS
    
    $ echo 東京都へ行く | gosudachisbcore -a
    東京都  名詞,固有名詞,地名,一般,*,*     東京都  東京都  トウキョウト
    へ      助詞,格助詞,*,*,*,*     へ      へ      エ
    行く    動詞,非自立可能,*,*,五段-カ行,終止形-一般       行く    行く    イク
    EOS
    
    $ echo 東京都へ行く | gosudachisbcore -m A
    東京    名詞,固有名詞,地名,一般,*,*     東京
    都      名詞,普通名詞,一般,*,*,*        都
    へ      助詞,格助詞,*,*,*,*     へ
    行く    動詞,非自立可能,*,*,五段-カ行,終止形-一般       行く
    EOS


## gosudachiclisb

[gosudachi](https://github.com/msnoigrs/gosudachi)のgosudachicliコマンドのシンブルバイナリバージョンです。同梱しているシステム辞書を強制的に使用することを除いて、オプションや動作はまったく同じです。


### ビルド方法

ビルドを行うと、同梱する辞書の違いにより、以下の2つのバイナリが生成されます。

-   gosudachiclisbcore
-   gosudachiclisbfull

    git clone https://github.com/msnoigrs/gosudachisb.git
    cd gosudachisb/gosudachiclisb
    bash ./build.sh


### コマンド

    $ gosudachiclisbcore [-r conf] [-m mode] [-a] [-d] [-o output] [file...]
    $ gosudachiclisbfull [-r conf] [-m mode] [-a] [-d] [-o output] [file...]


#### オプション

-   -r conf設定ファイルを指定
-   -s デフォルト設定を上書きする設定(json文字列)
-   -p リソースディレクトリ(設定ファイル内の各種リソースのベースディレクトリ、デフォルトは実行時ディレクトリ)
-   -m {A|B|C}分割モード
-   -a 読み、辞書形も出力
-   -d デバッグ情報の出力
-   -o 出力ファイル（指定がない場合は標準出力）
-   -f エラーを無視して処理を続行する


#### 出力例

    $ echo 東京都へ行く | gosudachiclisbcore
    東京都  名詞,固有名詞,地名,一般,*,*     東京都
    へ      助詞,格助詞,*,*,*,*     へ
    行く    動詞,非自立可能,*,*,五段-カ行,終止形-一般       行く
    EOS
    
    $ echo 東京都へ行く | gosudachiclisbcore -a
    東京都  名詞,固有名詞,地名,一般,*,*     東京都  東京都  トウキョウト
    へ      助詞,格助詞,*,*,*,*     へ      へ      エ
    行く    動詞,非自立可能,*,*,五段-カ行,終止形-一般       行く    行く    イク
    EOS
    
    $ echo 東京都へ行く | gosudachiclisbcore -m A
    東京    名詞,固有名詞,地名,一般,*,*     東京
    都      名詞,普通名詞,一般,*,*,*        都
    へ      助詞,格助詞,*,*,*,*     へ
    行く    動詞,非自立可能,*,*,五段-カ行,終止形-一般       行く
    EOS


## ライセンス

[Apache License, Version2.0](http://www.apache.org/licenses/LICENSE-2.0.html)
