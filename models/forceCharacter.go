package model

import "go.mongodb.org/mongo-driver/bson/primitive"

//ForceCharacterSheet is the model for the FFG Star Wars character sheet
type ForceCharacterSheet struct {
	ID                   primitive.ObjectID    `json:"_id" bson:"_id"`
	CharacterName        string                `json:"characterName" bson:"characterName"`
	PlayerName           string                `json:"playerName" bson:"playerName"`
	Species              string                `json:"species" bson:"species"`
	Career               string                `json:"career" bson:"career"`
	SpecializationTrees  []SpecializationTrees `json:"specializationTrees" bson:"specializationTrees"`
	SoakValue            int64                 `json:"soakValue" bson:"soakValue"`
	Wounds               Amount                `json:"wound" bson:"wound"`
	Strain               Amount                `json:"strain" bson:"strain"`
	Defense              DefenseStats          `json:"defense" bson:"defense"`
	Characteristics      Characteristics       `json:"characteristics" bson:"characteristics"`
	Skills               []Skills              `json:"skills" bson:"skills"`
	Weapons              []Weapons             `json:"weapons" bson:"weapons"`
	TotalXP              int64                 `json:"totalXP" bson:"totalXP"`
	AvailableXP          int64                 `json:"availableXP" bson:"availableXP"`
	Motivation           Motivation            `json:"motivation" bson:"motivation"`
	Morality             Morality              `json:"morality" bson:"morality"`
	CharacterDescription CharacterDescription  `json:"charcaterDescription" bson:"characterDescription"`
	Equipment            Equipment             `json:"equipment" bson:"equipment"`
	CriticalInjuries     []CriticalInjuries    `json:"criticalInjuries" bson:"criticalInjuries"`
	Talents              []Talents             `json:"talents" bson:"talents"`
	ForceRating          int64                 `json:"forceRating" bson:"forceRating"`
	Version              int64                 `json:"version" bson:"version"`
}

//DefenseStats is a generic that holds a characters Defensive amount for ranged and melee damage
type DefenseStats struct {
	Ranged int64 `json:"ranged" bson:"ranged"`
	Melee  int64 `json:"melee" bson:"melee"`
}

//Amount is a generic that holds threshold vs current for the FFG Star Wars character sheet
type Amount struct {
	Threshold int64 `json:"threshold" bson:"threshold"`
	Current   int64 `json:"current" bson:"current"`
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
	Brawn     int64 `json:"brawn" bson:"brawn"`
	Agility   int64 `json:"agility" bson:"agility"`
	Intellect int64 `json:"intellect" bson:"intellect"`
	Cunning   int64 `json:"cunning" bson:"cunning"`
	Willpower int64 `json:"willpower" bson:"willpower"`
	Presence  int64 `json:"presence" bson:"presence"`
}

//Skills is a subcategory of the FFG Star Wars character sheet that keeps track of different skills and their levels
type Skills struct {
	Name           string `json:"name" bson:"name"`
	Characteristic string `json:"characteristic" bson:"characteristic"`
	Career         bool   `json:"career" bson:"career"`
	Level          int64  `json:"level" bson:"level"`
	Description    string `json:"description" bson:"description"`
}

//Weapons is a subcategory of the FFG Star Wars character sheet that keeps track of a characters weapon inventory
type Weapons struct {
	Name    string `json:"name" bson:"name"`
	Skill   string `json:"skill" bson:"skill"`
	Damage  string `json:"damage" bson:"damage"`
	Crit    int64  `json:"crit" bson:"crit"`
	Range   string `json:"range" bson:"range"`
	Special string `json:"special" bson:"special"`
}

//Motivation is a subcategory of the FFG Star Wars character sheet that keeps track of a characters motivation
type Motivation struct {
	Type        string `json:"type" bson:"type"`
	Description string `json:"description" bson:"description"`
}

//Morality is a subcategory of the FFG Star Wars character sheet that keeps track of a characters morality
type Morality struct {
	EmotionalStrength string `json:"emotionalStrength" bson:"emotionalStrength"`
	EmotionalWeakness string `json:"emotionalWeakness" bson:"emotionalWeekness"`
	Conflict          int64  `json:"conflict" bson:"conflict"`
	Morality          int64  `json:"morality" bson:"morality"`
}

//Equipment is a subcatergory of the FFG Star Wars character sheet that keeps track of the equipment a character has on their person
type Equipment struct {
	Credits      int64  `json:"credits" bson:"credits"`
	Armor        []Gear `json:"armor" bson:"armor"`
	PersonalGear []Gear `json:"personalGear" bson:"personalGear"`
}

//Gear is the generic for any gear a character may carry for the FFG Star Wars character sheet
type Gear struct {
	Gear string `json:"gear" bson:"gear"`
}

//Talents is a subcategory of the FFG Star Wars character sheet that keeps track of all talents acquired through skill trees
type Talents struct {
	Name        string       `json:"name" bson:"name"`
	Page        int64        `json:"page" bson:"page"`
	Description string       `json:"description" bson:"description"`
	ForcePower  []ForcePower `json:"forcePower" bson:"forcePower"`
}

//ForcePower is a subcategory of the FFG Star Wars character sheet that keeps track of all Force abilities gained through the force tree
type ForcePower struct {
	Name        string `json:"name" bson:"name"`
	Description string `json:"description" bson:"description"`
	Completed   bool   `json:"completed" bson:"completed"`
}

//CriticalInjuries is a subcategory of the FFG Star Wars character sheet that keeps track of all critical injuries suffered
type CriticalInjuries struct {
	Severity int64 `json:"severity" bson:"severity"`
	Result   bool  `json:"result" bson:"result"`
}
