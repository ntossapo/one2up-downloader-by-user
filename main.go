package main;

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"regexp"
	"os"
	"strings"

	"time"
)

const(
	urlSearchUserItem = "http://www.one2up.com/list_content.php?page=%d&txt_Search=%s&radio_type=member_displayname"
	urlItemContent = "http://www.one2up.com/view_content.php?%s"
	requestTimeOut = 60
)
func main() {
	var user string
	var dest string

	fmt.Printf("Please Enter One2Up's User you want to downlaod\n:")
	fmt.Scanln(&user);

	fmt.Printf("Plase set Destination Saved File\n:")
	fmt.Scanln(&dest);
	dest = dest+"/"
	fmt.Printf("Confirm, the program will donwload all file from %s\ninto %s\n", user, dest)
	for page := 1 ; true ; page++ {
		url := fmt.Sprintf(urlSearchUserItem, page, user);
		resp, err := http.Get(url)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)
		re := regexp.MustCompile("content_ID=[0-9]*")
		list := re.FindAllStringSubmatch(string(body), -1)

		//for i:= 0; i < len(list);i++{
		//	fmt.Println(list[i])
		//}

		if len(list) != 0 {
			for i := 0; i < len(list); i++ {
				itemUrl := fmt.Sprintf(urlItemContent, list[i][0])
				fmt.Println("try to Access : " + itemUrl)
				itemResp, err := http.Get(itemUrl)
				if err != nil {
					continue
				}
				defer itemResp.Body.Close()

				itemBody, _ := ioutil.ReadAll(itemResp.Body)

				itemRe1 := regexp.MustCompile("http://[0-9A-Za-z]*.dl-one2up.com/onetwo/content/[0-9]*/[0-9]*/[0-9]*/[0-9A-Za-z.]*")
				itemRe2 := regexp.MustCompile("http://[0-9A-Za-z-]*.one2up.com/onetwo/content/[0-9]*/[0-9]*/[0-9]*/[0-9A-Za-z.]*")
				itemUrlAccess := itemRe1.FindAllStringSubmatch(string(itemBody), -1)
				if len(itemUrlAccess) == 0 {
					itemUrlAccess = itemRe2.FindAllStringSubmatch(string(itemBody), -1)
				}
				fileUrl := itemUrlAccess[0][0]
				fmt.Println("\tget item access url: " + fileUrl)
				fmt.Println("\tstarting download...")
				sfn := strings.Split(fileUrl, "/")
				filename := sfn[len(sfn) - 1]
				fmt.Println("\t"+filename)

				out, err := os.Create(dest + filename)
				if err != nil {
					continue
				}
				defer out.Close()

				timeout :=  time.Duration(requestTimeOut * time.Second)
				client := http.Client{
					Timeout:timeout,
				}

				fileResp, err := client.Get(fileUrl)
				if err != nil {
					continue
				}
				defer fileResp.Body.Close()

				f, err := os.Create(dest + filename)
				if err != nil {
					continue
				}
				defer f.Close()

				fileData, err := ioutil.ReadAll(fileResp.Body)
				if err != nil {
					continue
				}
				f.Write(fileData)
				f.Sync()
				f.Close()
				fmt.Println("\tdownload complete")
			}
		}else{
			break
		}
	}
}