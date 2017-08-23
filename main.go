package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"

	pb "gopkg.in/cheggaaa/pb.v1"

	mp3 "github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto"
	flags "github.com/jessevdk/go-flags"
)

const defaultMP3Dir = "data"
const defaultMP3File = "knocking_a_wooden_door2.mp3"

func run(mp3FilePath string) error {
	f, err := os.Open(mp3FilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	d, err := mp3.NewDecoder(f)
	if err != nil {
		return err
	}
	defer d.Close()

	p, err := oto.NewPlayer(d.SampleRate(), 2, 2, 8192)
	if err != nil {
		return err
	}
	defer p.Close()

	if _, err := io.Copy(p, d); err != nil {
		return err
	}
	return nil
}

var opts struct {
	Minute  float32 `short:"m" long:"minute" default:"30" description:"init limit time(minute)"`
	MP3File string  `short:"f" long:"file" description:"sound file for notice"`
}

func existsFile(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// data/knocking_a_wooden_door2.mp3 がないならば作成する。
// MP3ファイルは http://taira-komori.jpn.org/index.html から拝借したものを埋め込んでいる。
func setupDefaultMP3(temporaryDir string) {
	if !existsFile(getDefaultMP3Path(temporaryDir)) {
		data, err := Asset(getDefaultMP3AssetPath())
		if err != nil {
			panic(err)
		}
		ioutil.WriteFile(getDefaultMP3Path(temporaryDir), data, os.ModePerm)
	}
}

func getDefaultMP3AssetPath() string {
	return defaultMP3Dir + "/" + defaultMP3File
}

func getDefaultMP3Path(temporaryDir string) string {
	return temporaryDir + "/" + defaultMP3File
}

func getTemporaryDir() (string, error) {
	dir := os.TempDir() + "/time_to_sound_temporary"
	if existsFile(dir) {
		return dir, nil
	}
	err := os.Mkdir(dir, 0700)
	if err != nil {
		return "", err
	}
	return dir, nil
}

func main() {
	// temporaryDir, _ := ioutil.TempDir("", "time_to_sound")
	temporaryDir, _ := getTemporaryDir()
	defer os.RemoveAll(temporaryDir)

	_, err := flags.ParseArgs(&opts, os.Args)
	if err != nil {
		panic(err)
	}
	showMP3FilePath := opts.MP3File
	if len(opts.MP3File) == 0 {
		opts.MP3File = getDefaultMP3Path(temporaryDir)
		showMP3FilePath = getDefaultMP3AssetPath()
	}

	fmt.Printf("Minute: %v\n", opts.Minute)
	fmt.Printf("MP3File: %v\n", showMP3FilePath)

	if opts.Minute < 0 {
		panic("invalid initial time")
	}

	setupDefaultMP3(temporaryDir)

	// MP3ファイルの存在チェック
	if !existsFile(opts.MP3File) {
		panic("MP3file was not found")
	}

	const timeLayout = "2006-01-02 15:04:05"
	fmt.Printf("START: %v\n", time.Now().Format(timeLayout))

	// 1秒ごとに進捗バーを伸ばす。
	secondCount := int(opts.Minute * 60)
	bar := pb.StartNew(secondCount)
	// 1秒ごとに反応するティッカー。
	ticker := time.NewTicker(time.Second)
	stop := make(chan bool)
	go func() {
	loop:
		for {
			select {
			case _ = <-ticker.C:
				bar.Increment()
			case <-stop:
				break loop
			}
		}
	}()
	// 指定された時間だけスリープする。
	time.Sleep(time.Duration(secondCount*1000) * time.Millisecond)
	ticker.Stop()
	close(stop)
	// 残りのバーをすべて埋める。
	bar.Add(secondCount - int(bar.Get()))
	bar.FinishPrint("END: " + time.Now().Format(timeLayout))

	err = run(opts.MP3File)
	if err != nil {
		log.Fatal(err)
	}
}
