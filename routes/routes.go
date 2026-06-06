package routes

import (
	"github.com/nekowawolf/airdropv2/controllers"
	"github.com/nekowawolf/airdropv2/middlewares"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/nww")

	// Auth routes
	api.Post("/login", controllers.LoginAdminHandler)
	api.Post("/refresh", controllers.RefreshTokenHandler)
	api.Post("/logout", controllers.LogoutHandler)

	// Public portfolio routes
	api.Get("/portfolio", controllers.GetPortfolio)
	
	// Public airdrop routes
	api.Get("/freeairdrop", controllers.GetAirdropFreeHandler)
	api.Get("/paidairdrop", controllers.GetAirdropPaidHandler)
	api.Get("/allairdrop", controllers.GetAllAirdropHandler)
	api.Get("/allairdrop/stats", controllers.GetAllAirdropStatsHandler)

	// Public crypto community routes
	api.Get("/cryptocommunity", controllers.GetAllCryptoCommunity)
	api.Get("/cryptocommunity/stats", controllers.GetCryptoCommunityStats)

	// Public price routes	
	api.Get("/price", controllers.PriceHandler)

	// Public link routes
	api.Get("/profilelink", controllers.GetProfile)
	api.Get("/postslink", controllers.GetAllPosts)
	api.Get("/postslink/stats", controllers.GetPostStats)

	// ==================== PROTECTED ROUTES ====================
	protected := api.Group("/", middlewares.AdminMiddleware())

	// Protected airdrop routes
	protected.Get("/allairdrop/:id", controllers.GetAllAirdropByIDHandler)
	protected.Get("/freeairdrop/:id", controllers.GetAirdropFreeByIDHandler)
	protected.Get("/paidairdrop/:id", controllers.GetAirdropPaidByIDHandler)
	protected.Post("/freeairdrop", controllers.InsertAirdropFreeHandler)
	protected.Post("/paidairdrop", controllers.InsertAirdropPaidHandler)
	protected.Put("/allairdrop/:id", controllers.UpdateAllAirdropByIDHandler)
	protected.Put("/freeairdrop/:id", controllers.UpdateAirdropFreeByIDHandler)
    protected.Put("/paidairdrop/:id", controllers.UpdateAirdropPaidByIDHandler)
	protected.Delete("/allairdrop/:id", controllers.DeleteAllAirdropByIDHandler)
	protected.Delete("/freeairdrop/:id", controllers.DeleteAirdropFreeByIDHandler)
    protected.Delete("/paidairdrop/:id", controllers.DeleteAirdropPaidByIDHandler)

	// Protected crypto community routes
	protected.Get("/cryptocommunity/:id", controllers.GetCryptoCommunityByID)
	protected.Post("/cryptocommunity", controllers.InsertCryptoCommunity)
	protected.Put("/cryptocommunity/:id", controllers.UpdateCryptoCommunityByID)
	protected.Delete("/cryptocommunity/:id", controllers.DeleteCryptoCommunityByID)

	// Protected portfolio routes
	protected.Put("/portfolio", controllers.UpdatePortfolio)
	protected.Put("/portfolio/hero", controllers.UpdateHeroProfile)
	protected.Post("/portfolio/certificates", controllers.AddCertificate)
	protected.Post("/portfolio/designs", controllers.AddDesign)
	protected.Post("/portfolio/projects", controllers.AddProject)
	protected.Post("/portfolio/experience", controllers.AddExperience)
	protected.Post("/portfolio/education", controllers.AddEducation)
	protected.Post("/portfolio/skills/tech", controllers.AddTechSkill)
	protected.Post("/portfolio/skills/design", controllers.AddDesignSkill)

	protected.Delete("/portfolio/certificates/:id", controllers.DeleteCertificate)
	protected.Delete("/portfolio/designs/:id", controllers.DeleteDesign)
	protected.Delete("/portfolio/projects/:id", controllers.DeleteProject)
	protected.Delete("/portfolio/experience/:id", controllers.DeleteExperience)
	protected.Delete("/portfolio/education/:id", controllers.DeleteEducation)
	protected.Delete("/portfolio/skills/tech/:id", controllers.DeleteTechSkill)
	protected.Delete("/portfolio/skills/design/:id", controllers.DeleteDesignSkill)

	// Protected image routes
	protected.Post("/images", controllers.UploadImageHandler)
	protected.Get("/images", controllers.GetAllImages)
	protected.Delete("/images/:id", controllers.DeleteImage)

	// Protected link routes
	protected.Get("/postslink/:id", controllers.GetPostByID)
	protected.Post("/postslink", controllers.CreatePost)
	protected.Put("/postslink/:id", controllers.UpdatePost)
	protected.Put("/profilelink", controllers.UpdateProfile)
	protected.Delete("/postslink/:id", controllers.DeletePost)
}