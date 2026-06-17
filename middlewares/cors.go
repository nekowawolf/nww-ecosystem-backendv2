package middlewares

import (
    "strings"

    "github.com/gofiber/fiber/v2/middleware/cors"
)

var origins = []string{
    "https://nekowawolf.xyz",
    "https://cc.nekowawolf.xyz",
    "https://ai.nekowawolf.xyz",
    "https://www.nekowawolf.xyz",
    "https://link.nekowawolf.xyz",
    "https://web3.nekowawolf.xyz",
    "https://admin.nekowawolf.xyz",
    "https://github.nekowawolf.xyz",
    "https://nekowawolf.github.io",
    "https://airdrop.nekowawolf.xyz",
    "https://portfolio.nekowawolf.xyz",
}

var Cors = cors.Config{
    AllowOrigins:     strings.Join(origins[:], ","),
    AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
    ExposeHeaders:    "Content-Length",
    AllowCredentials: true,
}