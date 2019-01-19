# 2次元配列をバイナリデータとして読み書きする

Package binary のドキュメントには
> Clients that require high-performance serialization, especially for large data structures, should look at more advanced solutions such as the encoding/gob package or protocol buffers.
https://golang.org/pkg/encoding/binary/

って書いてあるし、gobのほうがいいのかな？　と思ってgoの練習がてら試してみたもの。
いろいろすっきりしない書き方をしている気がする。

gob.Encodeはbinary.Writeより遅いけど、gob.Decodeはbinary.Readより早い。

binary.Read が遅いのはデータの一つ一つについてエンディアンなどを確認してるから。
だとすれば、多次元配列、読み書きする環境がローカル、と限定して、バイナリ列として
扱えればそれでよいのでは？
というアイディアを採用すると劇的に早くなる。詳細は以下URL参照。
https://bleu48.hatenablog.com/entry/2019/01/19/185959

[10000][15000]Int32 で試してみた結果はこちら。
動作OSは Windows 10 Home (64bit)で、C:ドライブはSSDになっている。

```
PS C:\xxxxxx\2darraybinrw> go run .\2darraysaveload.go
2D array data read and write test
0s
set random data to 2D array A
3.5379303s
raw of A: 10000
column of A: 15000
======
write A to binaryfile with binary.Write and bufio
5.1851812s
write A to binary gob file with gob.Encode
6.3490486s
======
read from binaryfile to Bb with binary.Read and bufio
3.1756813s
A is same as Bb
 A:1179337617 1132468563 2094486022 368894420
Bb:1179337617 1132468563 2094486022 368894420
======
read from binary gob file to Bg with gob.Decode
1.922148s
A is same as Bg
 A:1179337617 1132468563 2094486022 368894420
Bg:1179337617 1132468563 2094486022 368894420
======
reset random data to 2D array bg
A is not same as Bb
 A:1179337617 1132468563 2094486022 368894420
Bb:611703184 1569009249 1390812853 720412150
binaryfilesize: 600000000
arraysize:      600000000
read from binaryfile to Bb with file.Read and unsafe pointer casting
126.753ms
A is same as Bb
 A:1179337617 1132468563 2094486022 368894420
Bb:1179337617 1132468563 2094486022 368894420

```

ファイルサイズでかいな！？
