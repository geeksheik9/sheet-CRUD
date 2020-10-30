# SHEET-CRUD

- This application is designed to create, delete, read, and update character sheets that follow the FFG Star Wars model.
- Current application support is only for Force and Destiny Character sheets.
- TODO: Add support for Edge of Empire and Age of Rebellion sheets

## Deploy

### Local

```shell
go run ./main/main.go
```

- Will run the application locally at port 3000

### Local Docker Container

```shell
docker-compose up
```

- Will run the application in a docker container at port 3000
- Exit using `control(^) + c`

### TODO: Create deployment scripts for any location

## Config

- PORT
- CHARACTER_DATABASE
- CHARACTER_COLLECTION
- CHARACTER_ARCHIVE
- LOG_LEVEL

## Routes

### Health Information

- **GET** /ping

  - checks if the pod is running

- **GET** /health

  - check that the pod and its dependencies are running

### Character Example

```shell
{
  "characterName": "Mando",
  "playerName": "Ben",
  "species": "Human",
  "career": "Mandalorian",
  "specializationTrees": [
    {
      "treeName": "Mandalorian",
      "completed": false
    }
  ],
  "soakValue": 5,
  "wounds": {
    "threshold": 15,
    "current": 0
  },
  "strain": {
    "threshold": 15,
    "current": 0
  },
  "defense": {
    "ranged": 5,
    "melee": 3
  },
  "characteristics": {
    "brawn": 3,
    "agility": 3,
    "intellect": 3,
    "cunning": 5,
    "willpower": 3,
    "presence": 3
  },
  "skills": [
    {
      "name": "athletics",
      "characteristic": "brawn",
      "career": true,
      "level": 5,
      "description": "for athletic checks"
    }
  ],
  "weapons": [
    {
      "name": "blaster pistol",
      "skill": "ranged heavy",
      "damage": "+1",
      "crit": 2,
      "range": "medium",
      "special": "none"
    }
  ],
  "totalXP": 200,
  "availableXP": 25,
  "motivation": {
    "type": "coolness",
    "description": "cool head under all situations"
  },
  "morality": {
    "emotionalStrength": "coolness",
    "emotionalWeakness": "coldness",
    "conflict": 4,
    "morality": 2
  },
  "characterDescription": {
    "gender": "Male",
    "age": 35,
    "height": 64,
    "build": "muscular",
    "hair": "brown",
    "eyes": "gray",
    "notableFeatures": "none"
  },
  "equipment": {
    "credits": 3000,
    "armor": [
      {
        "gear": "mandalorian armor"
      }
    ],
    "personalGear": [
      {
        "gear": "broken lightsaber"
      }
    ]
  },
  "criticalInjuries": [
    {
      "severity": 1,
      "result": false
    }
  ],
  "talents": [
    {
      "name": "sidestep",
      "page": 25,
      "description": "once per encounter can spend success to avoid an attack",
      "forcePower": [
        {
          "name": "move",
          "description": "spend one force point to move something",
          "completed": false
        }
      ]
    }
  ],
  "forceRating": 2,
  "version": 0
}
```

- **POST** /force-character-sheet

  - function name: InsertForceCharacterSheet
  - Inserts force character sheet into the database
  - Character model passed in through the body:
    - see character example

- **GET** /force-character-sheet

  - function name: GetForceCharacterSheets
  - Gets all force character sheets in the database
  - No parameters need to be passed

- **GET** /force-character-sheet/{ID}

  - function name: GetForceCharacterSheetsByID
  - Gets a specific force character sheet in the databasee
  - Parameters in url:
    - /force-character-sheet/`5e5d82a1802cc20001cb9b9c`

- **PUT** /force-character-sheet/{ID}

  - function name: UpdateForceCharacterSheetByID
  - Updates a specific force character sheet in the database
  - Parameters passed in url and character model passed in through the body:
    - /force-character-sheet/5e5d82a1802cc20001cb9b9c
    - see character example + `"_id": "5e5d82a1802cc20001cb9b9c",` at the start of the object

- **DELETE** /force-character-sheet/{id}

  - function name: DeleteForceCharacterSheetByID
  - Deletes a specific force character sheet in the database
  - Paramaters passed in url:
    - /force-character-sheet/`5e5d82a1802cc20001cb9b9c`

### Swagger

- **GET** /swagger/

  - serves up swagger UI
