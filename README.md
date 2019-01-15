# 2次元配列をバイナリデータとして読み書きする
Package binary のドキュメントには
> Clients that require high-performance serialization, especially for large data structures, should look at more advanced solutions such as the encoding/gob package or protocol buffers.
https://golang.org/pkg/encoding/binary/

って書いてあるし、gobのほうがいいのかな？　と思ってgoの練習がてら試してみたもの。
いろいろすっきりしない書き方をしている気がする。

gob.Encodeはbinary.Writeより遅いけど、gob.Decodeはbinary.Readより早い。
[10000][15000]Int32 で試してみた結果はこちら。
動作OSは Windows 10 Home (64bit)で、C:ドライブはSSDになっている。

```
PS C:\xxxxxx\2darraybinrw> go run .\2darraysaveload.go
2D array data read and write test
0s
set random data to 2D array A
3.6031374s
raw of A: 10000
column of A: 15000
======
write A to binaryfile with binary.Write and bufio
5.2602252s
write A to binaryfile with gob.Encode
6.4188928s
======
read from binaryfile to Bb with binary.Read and bufio
3.2424086s
A is same as Bb
 A:1179337617 1132468563 2094486022 368894420
Bb:1179337617 1132468563 2094486022 368894420
read from binaryfile to Bg with gob.Decode
1.9662906s
A is same as Bg
 A:1179337617 1132468563 2094486022 368894420
Bg:1179337617 1132468563 2094486022 368894420

PS C:\xxxxxx\2darraybinrw> dir

    ディレクトリ: C:\xxxxxx\2darraybinrw

Mode                LastWriteTime         Length Name
----                -------------         ------ ----
-a----       2019/01/15     21:51           5256 2darraysaveload.go
-a----       2019/01/15     21:52      600000000 binaryfile.bw
-a----       2019/01/15     21:52      749441712 binaryfile.gob
-a----       2019/01/15      0:38           1876 README.md

```

ファイルサイズでかいな！？
