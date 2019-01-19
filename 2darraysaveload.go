package main

import (
	"bufio"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"math/rand"
	"os"
	"time"
	"unsafe"
)

// 10000,15000 くらいにすると差ははっきりするが、
// 数百メガバイト程度のファイルサイズになるので注意
// （現時点でも 1,000*1,000*4 = 4Mbyte以上にはなる）
const rmax = 1000
const cmax = 1000

// A は元データとなる配列。これをファイルに保存する。
var A [rmax][cmax]int32

// B はファイルから取ったデータを保存する配列
var B [rmax][cmax]int32

func main() {

	start := time.Now()
	fmt.Println("2D array data read and write test")
	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Println(elapsed)

	// 0. 二次元配列を乱数で作る

	rand.Seed(11111)

	start = time.Now()

	fmt.Println("set random data to 2D array A")

	for r := 0; r < rmax; r++ {
		for c := 0; c < cmax; c++ {
			A[r][c] = rand.Int31()
		}
	}

	t = time.Now()
	elapsed = t.Sub(start)
	fmt.Println(elapsed)

	// 配列サイズを測っておく
	fmt.Println("raw of A:", len(A))
	fmt.Println("column of A:", len(A[0]))
	fmt.Println("======")

	// 1-1. 二次元配列をバイナリファイルとして保存 binary.Write and bufio

	filebW, err := os.Create("./binaryfile.bw")
	if err != nil {
		fmt.Println("file couldn't open")
		panic(err)
	}

	fmt.Println("write A to binaryfile with binary.Write and bufio")

	start = time.Now()

	// 本当はbytes で書きたかったんだけどやり方がよくわからなかったので
	// bufio で書いてしまう作戦。
	bufbW := bufio.NewWriter(filebW)

	// 書き込みバッファに書く
	err = binary.Write(bufbW, binary.LittleEndian, A)
	if err != nil {
		fmt.Println("bufio.Write failed:", err)
		panic(err)
	}
	filebW.Close()

	t = time.Now()
	elapsed = t.Sub(start)
	fmt.Println(elapsed)

	// 1-2. バイナリファイルから二次元配列 B に読み込み binary.Read and bufio
	fmt.Println("set random data to 2D array B")
	for r := 0; r < rmax; r++ {
		for c := 0; c < cmax; c++ {
			B[r][c] = rand.Int31()
		}
	}

	// 配列比較
	if A == B {
		fmt.Println("A is same as B")
	} else {
		fmt.Println("A is not same as B")
	}

	fmt.Printf("A:%d %d %d %d\n", A[0][0], A[0][cmax-1], A[rmax-1][0], A[rmax-1][cmax-1])
	fmt.Printf("B:%d %d %d %d\n", B[0][0], B[0][cmax-1], B[rmax-1][0], B[rmax-1][cmax-1])

	filebR, err := os.Open("./binaryfile.bw")
	if err != nil {
		fmt.Println("file couldn't open", err)
		panic(err)
	}
	defer filebR.Close()

	fmt.Println("read from binaryfile to B with binary.Read and bufio")

	start = time.Now()

	// 本当はbytes で読みたかったんだけどやり方がよくわからなかったので
	// bufio で読んでしまう作戦。
	// bufR := new(bytes.Buffer)

	bufbR := bufio.NewReader(filebR)

	// binary.Read で　ファイルに紐づけたバッファから 変数に読み込む。
	// 変数はポインタで与える
	err = binary.Read(bufbR, binary.LittleEndian, &B)
	if err != nil {
		fmt.Println("binary.Read failed:", err)
		panic(err)
	}

	t = time.Now()
	elapsed = t.Sub(start)
	fmt.Println(elapsed)

	// 配列比較
	if A == B {
		fmt.Println("A is same as B")
	} else {
		fmt.Println("A is not same as B")
	}

	fmt.Printf("A:%d %d %d %d\n", A[0][0], A[0][cmax-1], A[rmax-1][0], A[rmax-1][cmax-1])
	fmt.Printf("B:%d %d %d %d\n", B[0][0], B[0][cmax-1], B[rmax-1][0], B[rmax-1][cmax-1])

	fmt.Println("======")

	// 2-1. 二次元配列をバイナリファイルとして保存 gob

	filegW, err := os.Create("./binaryfile.gob")
	if err != nil {
		fmt.Println("file couldn't open", err)
		panic(err)
	}

	fmt.Println("write A to binary gob file with gob.Encode()")

	start = time.Now()

	// io.Writer に向けたエンコーダーを作る。io.Writerは書き込み可として開いたファイル
	encgob := gob.NewEncoder(filegW)

	// Encodeメソッドに送り付けたい変数 A を与える。
	// この場合はファイルに対して送り付けられる。
	if err = encgob.Encode(A); err != nil {
		fmt.Println("data couldn't Encode", err)
		panic(err)
	}
	filegW.Close()

	t = time.Now()
	elapsed = t.Sub(start)
	fmt.Println(elapsed)

	// 2-2. バイナリファイルから二次元配列に読み込み gob
	fmt.Println("set random data to 2D array B")
	for r := 0; r < rmax; r++ {
		for c := 0; c < cmax; c++ {
			B[r][c] = rand.Int31()
		}
	}

	// 配列比較
	if A == B {
		fmt.Println("A is same as B")
	} else {
		fmt.Println("A is not same as B")
	}

	fmt.Printf("A:%d %d %d %d\n", A[0][0], A[0][cmax-1], A[rmax-1][0], A[rmax-1][cmax-1])
	fmt.Printf("B:%d %d %d %d\n", B[0][0], B[0][cmax-1], B[rmax-1][0], B[rmax-1][cmax-1])

	filegR, err := os.Open("./binaryfile.gob")
	if err != nil {
		fmt.Println("file couldn't open", err)
		panic(err)
	}
	defer filegR.Close()

	fmt.Println("read from binary gob file to B with gob.Decode")

	start = time.Now()

	// io.Writer から受けるデコーダーを作る。io.Writerは開いたファイルのことになる
	decgob := gob.NewDecoder(filegR)

	// Decodeメソッドに受け取りたい変数のポインタを与える。
	// 送られてくる（ファイルの中に書き込まれている）データと同じ形、
	// サイズの変数じゃないとエラーになる。賢い。
	if err = decgob.Decode(&B); err != nil {
		fmt.Println("decgob.Decode failed:", err)
		panic(err)
	}

	t = time.Now()
	elapsed = t.Sub(start)
	fmt.Println(elapsed)

	// 配列比較
	if A == B {
		fmt.Println("A is same as B")
	} else {
		fmt.Println("A is not same as B")
	}

	fmt.Printf("A:%d %d %d %d\n", A[0][0], A[0][cmax-1], A[rmax-1][0], A[rmax-1][cmax-1])
	fmt.Printf("B:%d %d %d %d\n", B[0][0], B[0][cmax-1], B[rmax-1][0], B[rmax-1][cmax-1])

	fmt.Println("======")

	// 3-1. 二次元配列をバイナリファイルとして保存
	// file.Read/Write で多次元配列をバイト列として扱う作戦
	// https://bleu48.hatenablog.com/entry/2019/01/19/185959

	filefW, err := os.Create("./binaryfile.fw")
	if err != nil {
		fmt.Println("file couldn't open")
		panic(err)
	}

	fmt.Println("write A to binaryfile with file.Write()")

	start = time.Now()

	// 書き込みたい配列のポインタをバイト列のポインタとしてキャストした値を得る
	ptrA := (*[unsafe.Sizeof(A)]byte)(unsafe.Pointer(&A))[:][:]

	// 書き込みたい配列をバイト列としてポインタで与える
	_, err = filefW.Write(ptrA)
	if err != nil {
		fmt.Println("fileB.Write failed:", err)
		panic(err)
	}
	filefW.Close()

	t = time.Now()
	elapsed = t.Sub(start)
	fmt.Println(elapsed)

	// 3-2. file.Readで読み込み時のエンディアンチェックなどを回避する作戦。
	fmt.Println("set random data to 2D array B")
	for r := 0; r < rmax; r++ {
		for c := 0; c < cmax; c++ {
			B[r][c] = rand.Int31()
		}
	}

	// 配列比較
	if A == B {
		fmt.Println("A is same as B")
	} else {
		fmt.Println("A is not same as B")
	}

	fmt.Printf("A:%d %d %d %d\n", A[0][0], A[0][cmax-1], A[rmax-1][0], A[rmax-1][cmax-1])
	fmt.Printf("B:%d %d %d %d\n", B[0][0], B[0][cmax-1], B[rmax-1][0], B[rmax-1][cmax-1])

	// binary.Write しているファイルを開きなおす
	filefR, err := os.Open("./binaryfile.fw")
	if err != nil {
		fmt.Println("file couldn't open", err)
		panic(err)
	}
	defer filefR.Close()

	// いちおうファイルサイズくらいは見ておこうと思った
	finfo, ferr := filefR.Stat()
	if ferr != nil {
		fmt.Println("couldn't get file.Stat", ferr)
		panic(ferr)
	}

	// ファイルサイズを表示
	fmt.Println("binaryfilesize:", finfo.Size())
	fmt.Println("arraysize:     ", unsafe.Sizeof(B))

	// ファイルを読む
	fmt.Println("read from binaryfile to B with file.Read and unsafe pointer casting")

	start = time.Now()

	// 読み込み先の配列のポインタをバイト列のポインタとしてキャストした値を得る
	ptrB := (*[unsafe.Sizeof(B)]byte)(unsafe.Pointer(&B))[:][:]

	// 格納する配列をバイト列としてポインタで与える
	_, err = filefR.Read(ptrB)
	if err != nil {
		fmt.Println("filebR.Read failed:", err)
		panic(err)
	}

	t = time.Now()
	elapsed = t.Sub(start)
	fmt.Println(elapsed)

	// 配列比較
	if A == B {
		fmt.Println("A is same as B")
	} else {
		fmt.Println("A is not same as B")
	}

	fmt.Printf("A:%d %d %d %d\n", A[0][0], A[0][cmax-1], A[rmax-1][0], A[rmax-1][cmax-1])
	fmt.Printf("B:%d %d %d %d\n", B[0][0], B[0][cmax-1], B[rmax-1][0], B[rmax-1][cmax-1])

}
