package rest

import (
	"encoding/base64"
	"os"
	"strconv"
	"strings"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/naoina/toml"
)

type Account struct {
	Login    string
	Password string
	Read     bool
	Write    bool
	Delete   bool
}

func (a *Account) token() string {
	base := a.Login + ":" + a.Password
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(base))
}

type Accounts struct {
	Accounts []*Account
}

func (a *Accounts) searchCredential(authValue string) (*Account, bool) {
	if len(authValue) == 0 {
		return nil, false
	}
	for _, account := range a.Accounts {
		if account.token() == authValue {
			return account, true
		}
	}
	return nil, false
}

func BasicAuth(AccountsCfgPath string) gin.HandlerFunc {

	f, err := os.Open(AccountsCfgPath)
	if err != nil {
		log.Fatalf("Error opening Accountsentification file (%v) : %v", AccountsCfgPath, err)
		return nil
	}
	defer f.Close()

	listAccounts := Accounts{}
	if err := toml.NewDecoder(f).Decode(&listAccounts); err != nil {
		log.Fatalf("Error decoding Accountsentification file : %v", err)
		return nil
	}

	realm := "Basic realm=" + strconv.Quote("aydostorx")

	return func(c *gin.Context) {

		// Search user in the slice of allowed credentials
		account, found := listAccounts.searchCredential(c.Request.Header.Get("Authorization"))
		if !found {
			// Credentials doesn't match, we return 401 and abort handlers chain.
			c.Header("WWW-Authenticate", realm)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		authorized := true
		switch c.Request.Method {
		case "POST":
			if strings.HasSuffix(c.Request.URL.Path, "exists") {
				if !account.Read {
					authorized = false
				}
			} else {
				if !account.Write {
					authorized = false
				}
			}
		case "GET":
		case "HEAD":
			if !account.Read {
				authorized = false
			}
		case "DELETE":
			if !account.Delete {
				authorized = false
			}
		}

		if !authorized {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}
