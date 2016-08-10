package imgutil

import (
  "github.com/nfnt/resize"
  "image"
  "image/jpeg"
  "net/http"
  _ "os"
  "fmt"
)

type Imager struct {
  Address string
  Width uint
  Height uint
}

func (i *Imager) Resize() (image.Image, error) {
  inputType, err := getType(i.Address)

  switch inputType {
    case "file":
      return nil, fmt.Errorf("Not yet recognized")
    case "url":
      resizedImage, err := resizeURLFile(i)
      return resizedImage, err
    default:
      return nil, err
  }
}

func resizeURLFile(i *Imager) (image.Image, error) {
  response, err := http.Get(i.Address)
  if err != nil {
    return nil, err
  }
  defer response.Body.Close()

  imgResponse, err := jpeg.Decode(response.Body)
  if err != nil {
    return nil, err
  }

  resizedImage := resize.Resize(i.Width, i.Height, imgResponse, resize.Lanczos3)
  return resizedImage, nil
}

func getType(str string) (string, error) {
  return "url", nil
}