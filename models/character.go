package model

import "go.mongodb.org/mongo-driver/bson/primitive"

//CharacterSheet is the model for the FFG Star Wars character sheet
type CharacterSheet struct {
	ID                  primitive.ObjectID    `json:"_id" bson:"_id"`
	CharacterName       string                `json:"characterName" bson:"characterName"`
	PlayerName          string                `json:"playerName" bson:"playerName"`
	Species             string                `json:"species" bson:"species"`
	Career              string                `json:"career" bson:"career"`
	SpecializationTrees []SpecializationTrees `json:"specializationTrees" bson:"specializationTrees"`
}

//CharacterDescription is a subcategory of the FFG Star Wars character sheet that describes the physical appearance of a character
type CharacterDescription struct {
	Gender          string `json:"gender" bson:"gender"`
	Age             int64  `json:"age" bson:"age"`
	Height          int64  `json:"height" bson:"height"`
	Build           string `json:"build" bson:"build"`
	Hair            string `json:"hair" bson:"hair"`
	Eyes            string `json:"eyes" bson:"eyes"`
	NotableFeatures string `json:"notableFeatures" bson:"notableFeatures"`
}

//SpecializationTrees is a subcategory of the FFG Star Wars character sheet that keeps track of different abilities
type SpecializationTrees struct {
	TreeName  string `json:"treeName" bson:"treeName"`
	Completed bool   `json:"completed" bson:"completed"`
}

//Characteristics is a subcategory of the FFG Star Wars character sheet that keeps track of characteristics and their levels
type Characteristics struct {
}

//Skills is a subcategory of the FFG Star Wars character sheet that keeps track of different skills and their levels
type Skills struct {
}

//Weapons is a subcategory of the FFG Star Wars character sheet that keeps track of a characters weapon inventory
type Weapons struct {
}

//Motivation is a subcategory of the FFG Star Wars character sheet that keeps track of a characters motivation
type Motivation struct {
}

//Obligations is a subcategory of the FFG Star Wars character sheet that keeps track of a characters obligations
type Obligations struct {
}

//Equipment is a subcatergory of the FFG Star Wars character sheet that keeps track of the equipment a character has on their person
type Equipment struct {
}

//Talents is a subcategory of the FFG Star Wars character sheet that keeps track of all talents acquired through skill trees
type Talents struct {
}

//ForceAbilities is a subcategory of the FFG Star Wars character sheet that keeps track of all Force abilities gained through the force tree
type ForceAbilities struct {
}

//CriticalInjuries is a subcategory of the FFG Star Wars character sheet that keeps track of all critical injuries suffered
type CriticalInjuries struct {
}
