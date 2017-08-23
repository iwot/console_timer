# console_timer
指定時間（分）が経過したらサウンドを鳴らすだけ。  
作業に集中しすぎて座りっぱなしならないように作ってみた。  

コマンドプロンプト上で以下のようにすると、指定分数経過後にMP3を再生します。  
デフォルトは30分。

```
go run main.go bindata.go -m 1
```

MP3ファイルを指定する場合。
```
go run main.go bindata.go -f /path/to/sound.mp3
```
