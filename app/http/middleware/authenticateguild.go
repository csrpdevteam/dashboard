package middleware

import (
	"fmt"
	"github.com/TicketsBot/GoPanel/config"
	"github.com/TicketsBot/GoPanel/rpc/cache"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/gin-gonic/gin"
	"strconv"
)

// requires AuthenticateCookie middleware to be run before
func AuthenticateGuild(isApiMethod bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if guildId, ok := ctx.Params.Get("id"); ok {
			parsed, err := strconv.ParseUint(guildId, 10, 64)
			if err != nil {
				if isApiMethod {
					ctx.Redirect(302, config.Conf.Server.BaseUrl) // TODO: 404 Page
					ctx.Abort()
				} else {
					ctx.AbortWithStatusJSON(400, gin.H{
						"success": false,
						"error": "Invalid guild ID",
					})
				}
				return
			}

			ctx.Keys["guildid"] = parsed

			guild, found := cache.Instance.GetGuild(parsed, false)
			if !found {
				if isApiMethod {
					ctx.Redirect(302, config.Conf.Server.BaseUrl) // TODO: 404 Page
				} else {
					ctx.Redirect(302, fmt.Sprintf("https://invite.ticketsbot.net/?guild_id=%d&disable_guild_select=true&response_type=code&scope=bot%%20identify&redirect_uri=%s", parsed, config.Conf.Server.BaseUrl))
				}
				ctx.Abort()
				return
			}

			ctx.Keys["guild"] = guild

			// Verify the user has permissions to be here
			isAdmin := make(chan bool)
			go utils.IsAdmin(guild, ctx.Keys["userid"].(uint64), isAdmin)
			if !<-isAdmin {
				if isApiMethod {
					ctx.Redirect(302, config.Conf.Server.BaseUrl) // TODO: 403 Page
					ctx.Abort()
				} else {
					ctx.AbortWithStatusJSON(403, gin.H{
						"success": false,
						"error": "Unauthorized",
					})
				}
			}
		} else {
			if isApiMethod {
				ctx.Redirect(302, config.Conf.Server.BaseUrl) // TODO: 404 Page
				ctx.Abort()
			} else {
				ctx.AbortWithStatusJSON(400, gin.H{
					"success": false,
					"error": "Invalid guild ID",
				})
			}
		}
	}
}
