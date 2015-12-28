package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "authenticator"
	app.Version = Version
	app.Usage = "Authenticator can fetch a token from Active Directory for a given user"
	app.Author = "Erik Veld"
	app.HideHelp = true
	app.ArgsUsage = "<username> <password>"
	app.Action = authenticate
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "endpoint, e",
			Usage: "The endpoint to authenticate with",
		},
		cli.StringFlag{
			Name:  "id, i",
			Usage: "The client ID of your application",
		},
		cli.StringFlag{
			Name:  "connection, c",
			Usage: "The connection name of your application",
		},
		cli.StringFlag{
			Name:  "scope, s",
			Usage: "The scopes you wish the id_token to contain",
		},
	}
	app.Run(os.Args)
}

// Token is the jwt token that is obtained from Active Directory.
type Token struct {
	TokenType   string `json:"token_type"`
	AccessToken string `json:"access_token"`
	IDToken     string `json:"id_token"`
}

func check(err error) {
	// Handle errors in a common way.
	if err != nil {
		fmt.Printf("[ERROR] %s\n", err.Error())
	}
}

func authenticate(c *cli.Context) {
	clientID := c.String("id")
	connection := c.String("connection")
	endpoint := c.String("endpoint")
	scope := c.String("scope")

	// Make sure we have the needed information before we start.
	if len(c.Args()) < 2 || clientID == "" || connection == "" || endpoint == "" || scope == "" {
		fmt.Println("You need to provide the `id`, `connection`, `scope` and `endpoint` flags and a username/password pair in order to fetch the token")
		fmt.Println("For more information type: `authenticator --help`")
		os.Exit(0)
	}

	username := c.Args()[0]
	password := c.Args()[1]

	address := fmt.Sprintf("%s/oauth/ro", endpoint)
	client := &http.Client{}
	form := url.Values{}
	form.Set("client_id", clientID)
	form.Set("grant_type", "password")
	form.Set("username", username)
	form.Set("password", password)
	form.Set("scope", scope)
	form.Set("connection", connection)

	resp, err := client.PostForm(address, form)
	check(err)

	// Read the received token.
	content, err := ioutil.ReadAll(resp.Body)
	check(err)

	fmt.Println(string(content))

	// Parse the contents of the token.
	var token Token
	_ = json.Unmarshal(content, &token)
	fmt.Println(token.IDToken)

	os.Exit(0)
}
