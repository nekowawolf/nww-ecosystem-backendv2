package bot

import (
	"log"
	"time"

	"github.com/nekowawolf/airdropv2/features/github"
	"github.com/robfig/cron/v3"
)

var cronScheduler *cron.Cron

func InitScheduler() {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		log.Printf("Failed to load timezone Asia/Jakarta, falling back to UTC: %v", err)
		cronScheduler = cron.New()
	} else {
		cronScheduler = cron.New(cron.WithLocation(loc))
	}

	_, err = cronScheduler.AddFunc("0 3 * * 1", func() {
		log.Println("Running scheduled backup...")
		SendBackupArchive()
	})
	if err != nil {
		log.Fatalf("Failed to add cron job: %v", err)
	}

	_, err = cronScheduler.AddFunc("0 0,12 * * *", func() {
		log.Println("Running Github Repo Stats Sync (00:00 & 12:00)...")
		github.SyncAllGithubRepoStats()
	})
	if err != nil {
		log.Printf("Failed to add cron job for Github sync: %v", err)
	}

	cronScheduler.Start()
	log.Println("Scheduler initialized")
}

func GetNextBackupTime() string {
	if cronScheduler == nil {
		return "Unknown"
	}

	entries := cronScheduler.Entries()
	if len(entries) > 0 {
		return entries[0].Next.Format("Monday 15:04 WIB")
	}
	return "None"
}
