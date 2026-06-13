package module

import (
	"github.com/nekowawolf/airdropv2/utils"

	"github.com/google/uuid"
	"github.com/nekowawolf/airdropv2/config"
	"github.com/nekowawolf/airdropv2/models"
	"go.mongodb.org/mongo-driver/bson"
)

const portfolioCollection = "portfolio"

func GetPortfolio() (*models.Portfolio, error) {
	ctx, cancel := utils.GetDBContext()
	defer cancel()

	var result models.Portfolio
	err := config.Database.Collection(portfolioCollection).
		FindOne(ctx, bson.M{}).
		Decode(&result)
	return &result, err
}

func UpdatePortfolio(data models.Portfolio) error {
	ctx, cancel := utils.GetDBContext()
	defer cancel()

	_, err := config.Database.Collection(portfolioCollection).
		UpdateOne(ctx, bson.M{}, bson.M{"$set": data})
	return err
}

func UpdateHeroProfile(hero models.HeroProfile) error {
	ctx, cancel := utils.GetDBContext()
	defer cancel()

	_, err := config.Database.Collection(portfolioCollection).
		UpdateOne(
			ctx, 
			bson.M{}, 
			bson.M{"$set": bson.M{"hero": hero}},
		)
	return err
}

func addItem(field string, item interface{}) error {
	ctx, cancel := utils.GetDBContext()
	defer cancel()
	_, err := config.Database.Collection(portfolioCollection).
		UpdateOne(ctx, bson.M{}, bson.M{"$push": bson.M{field: item}})
	return err
}

func deleteItem(field, id string) error {
	ctx, cancel := utils.GetDBContext()
	defer cancel()

	_, err := config.Database.Collection(portfolioCollection).
		UpdateOne(ctx, bson.M{}, bson.M{"$pull": bson.M{field: bson.M{"id": id}}})
	return err
}

func AddCertificate(c models.Certificate) error {
	c.ID = uuid.NewString()
	return addItem("certificates", c)
}

func AddDesign(d models.Design) error {
	d.ID = uuid.NewString()
	return addItem("designs", d)
}

func AddProject(p models.Project) error {
	p.ID = uuid.NewString()
	return addItem("projects", p)
}

func AddExperience(e models.Experience) error {
	e.ID = uuid.NewString()
	return addItem("experience", e)
}

func AddEducation(e models.Education) error {
	e.ID = uuid.NewString()
	return addItem("education", e)
}

func AddTechSkill(s models.SkillItem) error {
	s.ID = uuid.NewString()
	return addItem("skills.tech", s)
}

func AddDesignSkill(s models.SkillItem) error {
	s.ID = uuid.NewString()
	return addItem("skills.design", s)
}

func DeleteCertificate(id string) error {
	return deleteItem("certificates", id)
}

func DeleteDesign(id string) error {
	return deleteItem("designs", id)
}

func DeleteProject(id string) error {
	return deleteItem("projects", id)
}

func DeleteExperience(id string) error {
	return deleteItem("experience", id)
}

func DeleteEducation(id string) error {
	return deleteItem("education", id)
}

func DeleteTechSkill(id string) error {
	return deleteItem("skills.tech", id)
}

func DeleteDesignSkill(id string) error {
	return deleteItem("skills.design", id)
}