package routes

import (
	"github.com/nekowawolf/airdropv2/middlewares"
	"github.com/gofiber/fiber/v2"

	"github.com/nekowawolf/airdropv2/features/admin"
	"github.com/nekowawolf/airdropv2/features/ai_tools"
	"github.com/nekowawolf/airdropv2/features/airdrop"
	"github.com/nekowawolf/airdropv2/features/community"
	"github.com/nekowawolf/airdropv2/features/community/community_submission"
	"github.com/nekowawolf/airdropv2/features/creators"
	"github.com/nekowawolf/airdropv2/features/github"
	"github.com/nekowawolf/airdropv2/features/github/repo_submission"
	"github.com/nekowawolf/airdropv2/features/guild"
	"github.com/nekowawolf/airdropv2/features/link"
	"github.com/nekowawolf/airdropv2/features/media"
	"github.com/nekowawolf/airdropv2/features/message"
	"github.com/nekowawolf/airdropv2/features/net"
	"github.com/nekowawolf/airdropv2/features/notes"
	"github.com/nekowawolf/airdropv2/features/portfolio"
	"github.com/nekowawolf/airdropv2/features/price"
	"github.com/nekowawolf/airdropv2/features/support/support_request"
	"github.com/nekowawolf/airdropv2/features/support/supporter"
	"github.com/nekowawolf/airdropv2/features/web3_tools"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/nww")

	// Auth routes
	api.Post("/login", admin.LoginAdminHandler)
	api.Post("/refresh", admin.RefreshTokenHandler)
	api.Post("/logout", admin.LogoutHandler)

	// Public portfolio routes
	api.Get("/portfolio", portfolio.GetPortfolioHandler)
	api.Get("/portfolio/projects", portfolio.GetProjectsHandler)
	api.Get("/portfolio/projects/:id", portfolio.GetProjectByIDHandler)
	api.Get("/portfolio/designs", portfolio.GetDesignsHandler)
	api.Get("/portfolio/designs/:id", portfolio.GetDesignByIDHandler)
	api.Get("/portfolio/certificates", portfolio.GetCertificatesHandler)
	api.Get("/portfolio/certificates/:id", portfolio.GetCertificateByIDHandler)
	
	// Public airdrop routes
	api.Get("/airdrops", airdrop.GetAirdropsPublicHandler)
	api.Get("/airdrops/stats", airdrop.GetAirdropsStatsHandler)

	// Public crypto community routes
	api.Get("/cryptocommunity", community.GetAllCryptoCommunityHandler)
	api.Get("/cryptocommunity/stats", community.GetCryptoCommunityStatsHandler)

	// Public guild routes
	api.Get("/guild", guild.GetAllGuildHandler)
	api.Get("/guild/stats", guild.GetGuildStatsHandler)

	// Public AI Tool routes
	api.Get("/aitools", ai_tools.GetAllAIToolsHandler)
	api.Get("/aitools/stats", ai_tools.GetAIToolStatsHandler)

	// Public Net routes
	api.Get("/net", net.GetAllNetHandler)
	api.Get("/net/stats", net.GetNetStatsHandler)

	// Public Creators routes
	api.Get("/creators", creators.GetAllCreatorsHandler)
	api.Get("/creators/stats", creators.GetCreatorsStatsHandler)

	// Public Web3 Tool routes
	api.Get("/web3tools", web3_tools.GetAllWeb3ToolsHandler)
	api.Get("/web3tools/stats", web3_tools.GetWeb3ToolStatsHandler)

	// Public Github Repo routes
	api.Get("/githubrepo", github.GetAllGithubReposHandler)
	api.Get("/githubrepo/stats", github.GetGithubRepoStatsHandler)
	api.Get("/githubrepo/:id/history", github.GetGithubRepoHistoryHandler)
	api.Get("/githubrepo/:id/details", github.GetGithubRepoDetailsHandler)
	api.Get("/githubrepo/:id/commits", github.GetGithubRepoCommitsHandler)
	api.Get("/githubrepo/commits/:owner/:repoName", github.GetGithubRepoCommitsByOwnerRepoHandler)

	// Public price routes	
	api.Get("/price", price.PriceHandler)

	// Public link routes
	api.Get("/profilelink", link.GetProfileHandler)
	api.Get("/postslink", link.GetAllPostsHandler)
	api.Get("/postslink/stats", link.GetPostStatsHandler)

	// Public repo submission routes
	api.Post("/repo-submissions", repo_submission.SubmitRepoHandler)

	// Public community submission routes
	api.Post("/community-submissions", community_submission.SubmitCommunityHandler)

	// Public support request routes
	api.Post("/support-requests", support_request.SubmitSupportRequestHandler)

	// Public supporter routes
	api.Get("/supporters", supporter.GetAllSupportersHandler)

	// ==================== PROTECTED ROUTES ====================
	protected := api.Group("/", middlewares.AdminMiddleware())

	// Protected airdrop routes
	protected.Get("/admin/airdrops", airdrop.GetAirdropsAdminHandler)
	protected.Get("/admin/airdrops/:id", airdrop.GetAirdropByIDAdminHandler)
	protected.Post("/admin/airdrops", airdrop.InsertAirdropHandler)
	protected.Put("/admin/airdrops/:id", airdrop.UpdateAirdropHandler)
	protected.Delete("/admin/airdrops/:id", airdrop.DeleteAirdropHandler)

	// Protected crypto community routes
	protected.Get("/cryptocommunity/:id", community.GetCryptoCommunityByIDHandler)
	protected.Post("/cryptocommunity", community.InsertCryptoCommunityHandler)
	protected.Put("/cryptocommunity/:id", community.UpdateCryptoCommunityByIDHandler)
	protected.Delete("/cryptocommunity/:id", community.DeleteCryptoCommunityByIDHandler)

	// Protected guild routes
	protected.Get("/guild/:id", guild.GetGuildByIDHandler)
	protected.Post("/guild", guild.InsertGuildHandler)
	protected.Put("/guild/:id", guild.UpdateGuildByIDHandler)
	protected.Delete("/guild/:id", guild.DeleteGuildByIDHandler)

	// Protected AI Tool routes
	protected.Get("/aitools/:id", ai_tools.GetAIToolsByIDHandler)
	protected.Post("/aitools", ai_tools.InsertAIToolsHandler)
	protected.Put("/aitools/:id", ai_tools.UpdateAIToolsByIDHandler)
	protected.Delete("/aitools/:id", ai_tools.DeleteAIToolsByIDHandler)

	// Protected Net routes
	protected.Get("/net/:id", net.GetNetByIDHandler)
	protected.Post("/net", net.InsertNetHandler)
	protected.Put("/net/:id", net.UpdateNetByIDHandler)
	protected.Delete("/net/:id", net.DeleteNetByIDHandler)

	// Protected Creators routes
	protected.Get("/creators/:id", creators.GetCreatorsByIDHandler)
	protected.Post("/creators", creators.InsertCreatorsHandler)
	protected.Put("/creators/:id", creators.UpdateCreatorsByIDHandler)
	protected.Delete("/creators/:id", creators.DeleteCreatorsByIDHandler)

	// Protected Web3 Tool routes
	protected.Get("/web3tools/:id", web3_tools.GetWeb3ToolsByIDHandler)
	protected.Post("/web3tools", web3_tools.InsertWeb3ToolsHandler)
	protected.Put("/web3tools/:id", web3_tools.UpdateWeb3ToolsByIDHandler)
	protected.Delete("/web3tools/:id", web3_tools.DeleteWeb3ToolsByIDHandler)

	// Protected Github Repo routes
	protected.Get("/githubrepo/:id", github.GetGithubRepoByIDHandler)
	protected.Post("/githubrepo", github.InsertGithubRepoHandler)
	protected.Put("/githubrepo/:id", github.UpdateGithubRepoByIDHandler)
	protected.Delete("/githubrepo/:id", github.DeleteGithubRepoByIDHandler)

	// Protected portfolio routes
	protected.Put("/portfolio", portfolio.UpdatePortfolioHandler)
	protected.Put("/portfolio/hero", portfolio.UpdateHeroProfileHandler)
	
	protected.Post("/portfolio/projects", portfolio.InsertProjectHandler)
	protected.Put("/portfolio/projects/:id", portfolio.UpdateProjectHandler)
	protected.Delete("/portfolio/projects/:id", portfolio.DeleteProjectHandler)

	protected.Post("/portfolio/designs", portfolio.InsertDesignHandler)
	protected.Put("/portfolio/designs/:id", portfolio.UpdateDesignHandler)
	protected.Delete("/portfolio/designs/:id", portfolio.DeleteDesignHandler)

	protected.Post("/portfolio/certificates", portfolio.InsertCertificateHandler)
	protected.Put("/portfolio/certificates/:id", portfolio.UpdateCertificateHandler)
	protected.Delete("/portfolio/certificates/:id", portfolio.DeleteCertificateHandler)

	protected.Post("/portfolio/experience", portfolio.AddExperienceHandler)
	protected.Post("/portfolio/education", portfolio.AddEducationHandler)
	protected.Post("/portfolio/skills/tech", portfolio.AddTechSkillHandler)
	protected.Post("/portfolio/skills/design", portfolio.AddDesignSkillHandler)

	protected.Delete("/portfolio/experience/:id", portfolio.DeleteExperienceHandler)
	protected.Delete("/portfolio/education/:id", portfolio.DeleteEducationHandler)
	protected.Delete("/portfolio/skills/tech/:id", portfolio.DeleteTechSkillHandler)
	protected.Delete("/portfolio/skills/design/:id", portfolio.DeleteDesignSkillHandler)

	// Protected media routes (Cloudflare R2)
	protected.Post("/images", media.UploadMediaHandler)
	protected.Get("/images", media.GetAllMediaHandler)
	protected.Delete("/images/:id", media.DeleteMediaHandler)

	// Protected link routes
	protected.Get("/postslink/:id", link.GetPostByIDHandler)
	protected.Post("/postslink", link.CreatePostHandler)
	protected.Put("/postslink/:id", link.UpdatePostHandler)
	protected.Put("/profilelink", link.UpdateProfileHandler)
	protected.Delete("/postslink/:id", link.DeletePostHandler)

	// Protected Note routes
	protected.Get("/notes", notes.GetAllNotesHandler)
	protected.Get("/notes/:id", notes.GetNoteByIDHandler)
	protected.Post("/notes", notes.InsertNoteHandler)
	protected.Put("/notes/:id", notes.UpdateNoteByIDHandler)
	protected.Delete("/notes/:id", notes.DeleteNoteByIDHandler)

	// Protected Message routes
	protected.Get("/message", message.GetMessageHandler)
	protected.Put("/message", message.UpdateMessageHandler)

	// Protected repo submission routes
	protected.Get("/repo-submissions", repo_submission.GetAllRepoSubmissionsHandler)
	protected.Delete("/repo-submissions/:id", repo_submission.DeleteRepoSubmissionHandler)

	// Protected community submission routes
	protected.Get("/community-submissions", community_submission.GetAllCommunitySubmissionsHandler)
	protected.Delete("/community-submissions/:id", community_submission.DeleteCommunitySubmissionHandler)

	// Protected support request routes
	protected.Get("/support-requests", support_request.GetAllSupportRequestsHandler)
	protected.Delete("/support-requests/:id", support_request.DeleteSupportRequestHandler)

	// Protected supporter routes
	protected.Get("/supporters/:id", supporter.GetSupporterByIDHandler)
	protected.Post("/supporters", supporter.InsertSupporterHandler)
	protected.Put("/supporters/:id", supporter.UpdateSupporterByIDHandler)
	protected.Delete("/supporters/:id", supporter.DeleteSupporterByIDHandler)
}