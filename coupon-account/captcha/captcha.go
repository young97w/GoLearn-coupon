package captcha

import (
	"account/internal"
	"account/log"
	"bufio"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/dchest/captcha"
	"github.com/google/uuid"
	"io"
	"os"
	"time"
)

func GenCaptcha() (string, string, error) {
	fileName := "captcha.png"
	f, err := os.Create(fileName)
	if err != nil {
		log.Logger.Info("GenerateCaptcha failed:" + err.Error())
		return "", "", err
	}
	defer f.Close()
	var w io.WriterTo
	b := captcha.RandomDigits(captcha.DefaultLen)
	w = captcha.NewImage("", b, captcha.StdWidth, captcha.StdHeight)
	_, err = w.WriteTo(f)
	if err != nil {
		log.Logger.Info("Generate Captcha failed:" + err.Error())
		return "", "", err
	}
	fmt.Println(b)
	captchaStr := ""
	for _, item := range b {
		captchaStr += fmt.Sprintf("%d", item)
	}
	fmt.Println("captchaStr: " + captchaStr)
	randUUID := uuid.New().String()
	internal.RedisClient.Set(context.Background(), randUUID, captchaStr, 120*time.Second)
	b64, err := getBase64(fileName)
	if err != nil {
		log.Logger.Info("generate base64 failed:" + err.Error())
		return "", "", err
	}
	fmt.Println(b64)
	return randUUID, b64, nil
}

func getBase64(fileName string) (string, error) {
	imgFile, err := os.Open(fileName) // image
	if err != nil {
		log.Logger.Error(err.Error())
		return "", err
	}
	defer imgFile.Close()
	// create a new buffer base on file size
	fInfo, _ := imgFile.Stat()
	var size = fInfo.Size()
	buf := make([]byte, size)
	// read file content into buffer
	fReader := bufio.NewReader(imgFile)
	fReader.Read(buf)
	// if you create a new image instead of loading from file, encode the image to buffer instead with png.Encode()
	// png.Encode(&buf, image)
	// convert the buffer bytes to base64 string - use buf.Bytes() for new image
	imgBase64Str := base64.StdEncoding.EncodeToString(buf)
	return imgBase64Str, nil
}
