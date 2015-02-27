package main

import (
	"encoding/json"
	"fmt"

	"github.com/go-martini/martini"
	"github.com/ideo/dragonfruit"

	"bufio"
	"flag"
	"github.com/gedex/inflector"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ideo/dragonfruit/backends/backend_couchdb"
	"github.com/martini-contrib/gzip"
	"github.com/martini-contrib/oauth2"
	"github.com/martini-contrib/sessions"
	goauth2 "golang.org/x/oauth2"
)

type loginresult struct {
	User       gplususer `json:"user"`
	IsLoggedIn bool      `json:"isLoggedIn"`
}

type gplususer struct {
	Emails []struct {
		Value string `json:"value"`
		Type  string `json:"type"`
	} `json:"emails"`
	DisplayName string `json:"displayName"`
	Image       struct {
		Url string `json:"url"`
	} `json:"image"`
	Expires time.Time `json:"expires"`
}

func main() {
	serve, add := parseFlags()

	// dragonfruit setup
	d := backend_couchdb.Db_backend_couch{}
	d.Connect("http://localhost:5984")

	if add {
		addResouce(&d)
	}

	if serve {
		launchServer(&d)
	}

}

func launchServer(d dragonfruit.Db_backend) *martini.ClassicMartini {

	m := dragonfruit.GetMartiniInstance()

	dragonfruit.ServeDocSet(d)

	m.Use(sessions.Sessions("my_session", sessions.NewCookieStore([]byte("secret123"))))

	m.Use(gzip.All())
	m.Use(martini.Static("../static/dist"))
	//m.Use(martini.Static("../static/.tmp"))
	m.Use(martini.Static("../static"))

	conf := &goauth2.Config{
		ClientID:     "288198830216-s4klktd4qm3asq72sm7acifugcumdseq.apps.googleusercontent.com",
		ClientSecret: "zzKEh5RPFhaEeJOsQi-D34Rc",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email"},
		RedirectURL: "http://localhost:3000/oauth2callback",
	}
	m.Use(oauth2.Google(conf))
	m.Post("/checkauth", oauth2.LoginRequired, func(tokens oauth2.Tokens, res http.ResponseWriter, s sessions.Session) (int, string) {
		h := res.Header()
		h.Add("Content-Type", "application/json; charset=utf-8")
		code := 200

		result := loginresult{}

		if tokens.Expired() {
			result.IsLoggedIn = false
			code = 403
		}
		test, user := checkAuth(conf, tokens)

		if !test {
			result.IsLoggedIn = false
			code = 403
		} else {
			result.IsLoggedIn = true
			result.User = user
		}
		fmt.Println(result)
		out, err := json.Marshal(result)
		fmt.Println(out)
		if err != nil {
			code = 500
		}
		//fmt.Println(out)
		//fmt.Println(resp, string(out))

		return code, string(out)
	})

	m.Run()
	return m
}

func returnSuccess(res http.ResponseWriter) (int, string) {
	h := res.Header()

	h.Add("Content-Type", "text/plain")
	//res.Write([]byte("Content-Type: text/plain"))
	return 200, "auth"
}

func checkAuth(conf *goauth2.Config, tokens oauth2.Tokens) (bool, gplususer) {
	//fmt.Println(tokens)
	tok := &goauth2.Token{
		AccessToken:  tokens.Access(),
		Expiry:       tokens.ExpiryTime(),
		RefreshToken: tokens.Refresh(),
	}

	client := conf.Client(goauth2.NoContext, tok)
	resp, err := client.Get("https://www.googleapis.com/plus/v1/people/me")
	defer resp.Body.Close()

	fmt.Println("call error:", err)
	out, err := ioutil.ReadAll(resp.Body)

	fmt.Println("read error:", err)
	var user gplususer
	json.Unmarshal(out, &user)

	fmt.Println(user)

	for _, v := range user.Emails {
		if strings.Contains(v.Value, "ideo.com") {
			user.Expires = tokens.ExpiryTime()
			return true, user
		}
	}

	return false, user
}

/* shamelessly cut and paste from the dragonfruit api builder */
func addResouce(d dragonfruit.Db_backend) {
	rd, err := dragonfruit.LoadDescriptionFromDb(d, "resourceTemplate.json")
	res, err := dragonfruit.NewFromTemplate("resourceTemplate.json")

	scanner := bufio.NewScanner(os.Stdin)

	res.BasePath = "/api"

	fmt.Print("enter a base path for all APIs (press [enter] for \"" + res.BasePath + "\"):")
	scanner.Scan()

	tmpBasepath := scanner.Text()
	if tmpBasepath != "" {
		res.BasePath = tmpBasepath
	}

	fmt.Print("enter a resource type: ")

	// resource type
	scanner.Scan()
	resourceType := inflector.Singularize(scanner.Text())

	path := inflector.Pluralize(resourceType)
	if err != nil {
		panic(err)
	}

	defaultText := "Describes operations on " + path
	// description
	fmt.Print("describe what the ", path, " service will do (press [enter] for \""+defaultText+"\"):")
	scanner.Scan()
	resourceDescription := scanner.Text()
	if resourceDescription == "" {
		resourceDescription = defaultText
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

	ok := true
	for _, v := range rd.APIs {
		if v.Path == "/"+path {
			ok = false
			break
		}
	}

	if ok {
		tmp := dragonfruit.ResourceSummary{
			Path:        "/" + path,
			Description: resourceDescription,
		}
		rd.APIs = append(rd.APIs, &tmp)
	}

	fmt.Print("What is the base model that this resource returns? (press [enter] for \"", resourceType, "\"):")
	res.ResourcePath = "/" + path
	scanner.Scan()
	modelType := scanner.Text()
	if modelType == "" {
		modelType = resourceType
	}

	fmt.Print("Enter a path to some sample data: ")
	scanner.Scan()
	fname := scanner.Text()
	byt, err := ioutil.ReadFile(fname)
	if err != nil {
		fmt.Println(err)
	}

	modelMap, err := dragonfruit.Decompose(byt, modelType)

	res.Models = modelMap
	upstreamParams := make([]*dragonfruit.Property, 0)
	apis := dragonfruit.MakeCommonAPIs("",
		path,
		strings.Title(resourceType),
		modelMap,
		upstreamParams)

	res.Apis = append(res.Apis, apis...)

	rd.Save(d)
	res.Save(d)
	preperror := d.Prep(path, res)
	if preperror != nil {
		panic(preperror)
	}

	fmt.Println("Done!")

}

func parseFlags() (bool, bool) {
	var serve = flag.Bool("serve", true, "Start a server after running")
	var addresource = flag.Bool("add", false, "Add a new resource")

	flag.Parse()
	return *serve, *addresource
}
