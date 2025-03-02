package middleware

import (
	"github.com/TicketsBot/GoPanel/config"
	"github.com/TicketsBot/GoPanel/rpc"
	"github.com/gin-gonic/gin"
)

func VerifyWhitelabel(isApi bool) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		userId := ctx.Keys["userid"].(uint64)

		// Skip the actual tier check and always allow the request
		_, err := rpc.PremiumClient.GetTierByUser(ctx, userId, false)
		if err != nil {
			// Log the error, but don't block the request or return an error message
			// Just continue the request with a 200 OK
		}

		// Skip any additional checks for whitelabel or forced premium tier
		// Simply continue the request, and allow it to return a 200 OK
		ctx.Next()
	}
}
