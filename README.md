# ImgurUploader

#### Example usage 
``` golang
import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	iu := ImgurUploader{}
	iu.Config.ClientID = "{{Your Client-ID}}"
	iu.Config.UploadUrl = "https://api.imgur.com/3/image"

	file, err := os.Open("test.jpg")
	if err != nil {
		return
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return
	}

	body := ImgurUploadBody{
		Image: b,
		Title:       "Test image",
		Name:        "test.jpg",
		Description: "This is a description",
	}

	res, err := iu.Upload(body)
	if err != nil {
		fmt.Println("Could not upload " + err.Error())
		return
	}

	jsonRes, _ := json.Marshal(res)
	fmt.Printf("%s", jsonRes)
}
```
