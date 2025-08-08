package stat

import (
	"go/adv-demo/pkg/db"
	"time"

	"gorm.io/datatypes"
)

type StatRepository struct {
	Database *db.Db
}

func NewStatRepository(database *db.Db) *StatRepository {
	return &StatRepository{
		Database: database,
	}
}

func (repo *StatRepository) AddClick(linkId uint) {
	var stat Stat
	currDate := datatypes.Date(time.Now())
	repo.Database.Find(&stat, "link_id = ? and date = ?", linkId, currDate)
	if stat.ID == 0 {
		repo.Database.Create(&Stat{
			LinkId: linkId,
			Clicks: 1,
			Date:   currDate,
		})
	} else {
		stat.Clicks += 1
		repo.Database.Save(&stat)
	}
}

func (repo *StatRepository) GetStats(by string, from, to time.Time) []GetStatResponse {
	var stats []GetStatResponse
	var query string
	switch by {
	case GroupByDay:
		query = "to_char(date, 'YYYY-MM-DD') as period, sum(clicks)"
	case GroupByMonth:
		query = "to_char(date, 'YYYY-MM') as period, sum(clicks)"
	}
	repo.Database.Table("stats").Select(query).Where("date BETWEEN ? AND ?", from, to).
		Group("period").Order("period").Scan(&stats)
	return stats
}
