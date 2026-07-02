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

	// Public AI Tool routes
	api.Get("/aitools", controllers.GetAllAITools)
	api.Get("/aitools/stats", controllers.GetAIToolStats)

	// Public Web3 Tool routes
	api.Get("/web3tools", controllers.GetAllWeb3Tools)
	api.Get("/web3tools/stats", controllers.GetWeb3ToolStats)

	// Public Github Repo routes
	api.Get("/githubrepo", controllers.GetAllGithubRepos)
	api.Get("/githubrepo/stats", controllers.GetGithubRepoStats)

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

	// Protected AI Tool routes
	protected.Get("/aitools/:id", controllers.GetAIToolsByID)
	protected.Post("/aitools", controllers.InsertAITools)
	protected.Put("/aitools/:id", controllers.UpdateAIToolsByID)
	protected.Delete("/aitools/:id", controllers.DeleteAIToolsByID)

	// Protected Web3 Tool routes
	protected.Get("/web3tools/:id", controllers.GetWeb3ToolsByID)
	protected.Post("/web3tools", controllers.InsertWeb3Tools)
	protected.Put("/web3tools/:id", controllers.UpdateWeb3ToolsByID)
	protected.Delete("/web3tools/:id", controllers.DeleteWeb3ToolsByID)

	// Protected Github Repo routes
	protected.Get("/githubrepo/:id", controllers.GetGithubRepoByID)
	protected.Post("/githubrepo", controllers.InsertGithubRepo)
	protected.Put("/githubrepo/:id", controllers.UpdateGithubRepoByID)
	protected.Delete("/githubrepo/:id", controllers.DeleteGithubRepoByID)

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

	// Protected Note routes
	protected.Get("/notes", controllers.GetAllNotes)
	protected.Get("/notes/:id", controllers.GetNoteByID)
	protected.Post("/notes", controllers.InsertNote)
	protected.Put("/notes/:id", controllers.UpdateNoteByID)
	protected.Delete("/notes/:id", controllers.DeleteNoteByID)
}