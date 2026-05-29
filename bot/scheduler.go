package bot

import (
	"log"
	"time"

	"github.com/robfig/cron/v3"
)

var cronScheduler *cron.Cron

func InitScheduler() {
	// Schedule based on Jakarta Time (WIB)
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		log.Printf("Failed to load timezone Asia/Jakarta, falling back to UTC: %v", err)
		cronScheduler = cron.New()
	} else {
		cronScheduler = cron.New(cron.WithLocation(loc))
	}

	// 03:00 on Monday -> "0 3 * * 1"
	_, err = cronScheduler.AddFunc("0 3 * * 1", func() {
		log.Println("Running scheduled backup...")
		SendBackupArchive()
	})
	if err != nil {
		log.Fatalf("Failed to add cron job: %v", err)
	}

	cronScheduler.Start()
	log.Println("Scheduler initialized (Monday 03:00 WIB)")
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
