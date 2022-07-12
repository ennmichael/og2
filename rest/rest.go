package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"og2/game"

	"github.com/gin-gonic/gin"
)

type userRequest struct {
	Username string `json:"username"`
}

type dashboardResponse struct {
	Resources resourcesResponse `json:"resources"`
	Factories []factoryResponse `json:"factories"`
}

type resourcesResponse struct {
	Iron   uint32 `json:"iron"`
	Copper uint32 `json:"copper"`
	Gold   uint32 `json:"gold"`
}

type factoryResponse struct {
	Resource           string  `json:"resource"`
	Level              uint32  `json:"level"`
	Upgrading          bool    `json:"upgrading"`
	UpgradeSecondsLeft *uint32 `json:"upgradeSecondsLeft"`
}

type authHeader struct {
	Username string `header:"X-Auth-Username"`
}

type upgradeRequest struct {
	Resource string `json:"resource"`
}

type successResponse struct {
	Success bool `json:"success"`
}

func RunServer(store game.Store) error {
	games := make(map[game.User]game.Game)
	router := gin.Default()

	for _, gm := range store.LoadGames() {
		games[gm.User] = game.Continue(gm, store)
	}

	router.POST("/user", func(c *gin.Context) {
		var userRequest userRequest
		if err := c.BindJSON(&userRequest); err != nil {
			return
		}
		user := game.User(userRequest.Username)
		if _, ok := games[user]; ok {
			c.Status(http.StatusBadRequest)
		} else {
			games[user] = game.Start(user, store)
			c.Status(http.StatusOK)
		}
	})

	router.GET("/dashboard", func(c *gin.Context) {
		var authHeader authHeader
		if err := c.BindHeader(&authHeader); err != nil {
			return
		}
		gm, ok := games[game.User(authHeader.Username)]
		if ok {
			resp := make(chan game.State)
			gm.StateChan <- game.StateMessage{Resp: resp}
			state := <-resp

			// TODO Testing, remove this
			j, err := json.Marshal(state)
			if err != nil {
				panic(err)
			}
			fmt.Println(string(j))

			factories := []factoryResponse{}
			for _, f := range state.Factories {
				factories = append(factories, factoryResponse{
					Resource:           resourceToString(f.Resource),
					Level:              f.Level,
					Upgrading:          f.Upgrading(),
					UpgradeSecondsLeft: f.UpgradeSecondsLeft,
				})
			}
			c.JSON(http.StatusOK, dashboardResponse{
				Resources: resourcesResponse{
					Iron:   state.Resources.Iron,
					Copper: state.Resources.Copper,
					Gold:   state.Resources.Gold,
				},
				Factories: factories,
			})
		} else {
			c.Status(http.StatusNotFound)
		}
	})

	router.POST("/upgrade", func(c *gin.Context) {
		var authHeader authHeader
		if err := c.BindHeader(&authHeader); err != nil {
			return
		}
		var upgradeRequest upgradeRequest
		if err := c.BindJSON(&upgradeRequest); err != nil {
			return
		}
		gm, ok := games[game.User(authHeader.Username)]
		if ok {
			resp := make(chan bool)
			resource, ok := stringToResource(upgradeRequest.Resource)
			if !ok {
				c.Status(http.StatusBadRequest)
				return
			}
			gm.UpgradeChan <- game.UpgradeMessage{Resp: resp, Resource: resource}
			success := <-resp
			c.JSON(http.StatusOK, successResponse{Success: success})
		} else {
			c.Status(http.StatusNotFound)
		}
	})

	return router.Run("localhost:8081")
}

func resourceToString(resource game.Resource) string {
	switch resource {
	case game.Iron:
		return "iron"
	case game.Copper:
		return "copper"
	case game.Gold:
		return "gold"
	}
	panic("Invalid resource")
}

func stringToResource(resource string) (game.Resource, bool) {
	switch resource {
	case "iron":
		return game.Iron, true
	case "copper":
		return game.Copper, true
	case "gold":
		return game.Gold, true
	default:
		return game.Resource(0), false
	}
}
