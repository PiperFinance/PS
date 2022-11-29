package api

import (
	"fmt"
)

func StartApp(appPort int) {
	MainGinRouter.Run(fmt.Sprintf(":%s", PORT))
}
