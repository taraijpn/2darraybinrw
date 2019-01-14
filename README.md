# 2次元配列をバイナリデータとして読み書きする
goの練習がてら試してみたけど、いろいろすっきりしない書き方をしている気がする。

gob.Encodeはbinary.Writeより遅いけど、gob.Decodeはbinary.Readより早い。
[10000][15000]Int32 で試してみた結果はこちら。
動作OSは Windows 10 Home (64bit)で、C:ドライブはSSDになっている。

```
> go run .\2darraysaveload.go
2D array data read and write test
0s
set random data to 2D array A
3.5654112s
raw of A: 10000
column of A: 15000
======
save A to binaryfile with binary.Write and bufio
5.207169s
save A to binaryfile with gob.Encode
6.3129723s
======
load from binaryfile to Bb with binary.Read and bufio
3.2170841s
A is same as Bb
 A:1179337617 1132468563 2094486022 368894420
Bb:1179337617 1132468563 2094486022 368894420
load from binaryfile to Bg with gob.Decode
1.9140531s
A is same as Bg
 A:1179337617 1132468563 2094486022 368894420
Bg:1179337617 1132468563 2094486022 368894420
> dir

    ディレクトリ: C:\xxxxxxxx\2darraybinrw

Mode                LastWriteTime         Length Name
----                -------------         ------ ----
-a----       2019/01/15      0:35           5327 2darraysaveload.go
-a----       2019/01/15      0:35      600000000 binaryfile.bw
-a----       2019/01/15      0:36      749441712 binaryfile.gob
-a----       2019/01/15      0:13           1632 README.md
```

ファイルサイズでかいな！？

まあ Package binary のドキュメントにも 

> Clients that require high-performance serialization, especially for large data structures, should look at more advanced solutions such as the encoding/gob package or protocol buffers.

https://golang.org/pkg/encoding/binary/

って書いてあるし、gobのほうがいいのかも。
