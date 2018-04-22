package entities

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"math"
	"BasicECS/components"
)

func TestNewIdGeneration(t *testing.T) {

	entityManager := CreateEntityManager()
	id := entityManager.generateNewEntityId()
	secondId := entityManager.generateNewEntityId()

	assert.Equal(t, 1, id)
	assert.Equal(t, 2, secondId)
	assert.Equal(t, entityManager.lowestUnassignedEntityId, 2)
}

func TestFullIdGeneration(t *testing.T) {

	entityManager := CreateEntityManager()
	entityManager.lowestUnassignedEntityId = math.MaxInt32 - 1

	id := entityManager.generateNewEntityId()
	secondId := entityManager.generateNewEntityId()

	assert.Equal(t, math.MaxInt32, id)
	assert.Equal(t, 0, secondId)
}

func TestCreateEntities(t *testing.T) {

	entityManager := CreateEntityManager()

	entity := entityManager.CreateEntity()
	secondEntity := entityManager.CreateEntity()

	assert.Equal(t, 1, entity.Id)
	assert.True(t, ContainsEntity(entityManager.entities, entity))
	assert.NotNil(t, 2, secondEntity.Id)
	assert.True(t, ContainsEntity(entityManager.entities, secondEntity))
}

func TestCreateEntitiesFullCapacity(t *testing.T) {

	entityManager := CreateEntityManager()
	entityManager.lowestUnassignedEntityId = math.MaxInt32 - 1

	entity := entityManager.CreateEntity()
	secondEntity := entityManager.CreateEntity()

	assert.Equal(t, math.MaxInt32, entity.Id)
	assert.True(t, ContainsEntity(entityManager.entities, entity))
	assert.NotNil(t, 0, secondEntity.Id)
	assert.True(t, ContainsEntity(entityManager.entities, secondEntity))
}

func TestAddComponentToEntity(t *testing.T) {

	entityManager := CreateEntityManager()

	entity := entityManager.CreateEntity()
	secondEntity := entityManager.CreateEntity()

	speedComponent := components.CreateSpeedComponent(32)
	healthComponent := components.CreateHealthComponent(10)

	healthComponentName := healthComponent.GetComponentName()
	speedComponentName := speedComponent.GetComponentName()

	assert.Nil(t, entityManager.componentsByClass[speedComponentName])
	assert.Nil(t, entityManager.componentsByClass[healthComponentName])

	entityManager.AddComponentToEntity(entity.Id, speedComponent)
	entityManager.AddComponentToEntity(secondEntity.Id, healthComponent)

	assert.NotNil(t, entityManager.componentsByClass[speedComponentName])
	assert.NotNil(t, entityManager.componentsByClass[healthComponentName])

	assert.Nil(t, entityManager.componentsByClass[healthComponentName][entity.Id])
	assert.NotNil(t, entityManager.componentsByClass[speedComponentName][entity.Id])

	assert.NotNil(t, entityManager.componentsByClass[healthComponentName][secondEntity.Id])

	assert.Nil(t, entityManager.componentsByClass["invalidcomponent"])
}

func TestAddComponentsToEntity(t *testing.T)  {

	entityManager := CreateEntityManager()

	entity := entityManager.CreateEntity()
	secondEntity := entityManager.CreateEntity()

	speedComponent := components.CreateSpeedComponent(32)
	healthComponent := components.CreateHealthComponent(10)

	healthComponentName := healthComponent.GetComponentName()
	speedComponentName := speedComponent.GetComponentName()

	assert.Nil(t, entityManager.componentsByClass[speedComponentName])
	assert.Nil(t, entityManager.componentsByClass[healthComponentName])

	entityManager.AddComponentsToEntity(entity.Id, []components.Component{speedComponent, healthComponent})
	entityManager.AddComponentsToEntity(secondEntity.Id, []components.Component{healthComponent})

	assert.NotNil(t, entityManager.componentsByClass[speedComponentName])
	assert.NotNil(t, entityManager.componentsByClass[healthComponentName])

	assert.NotNil(t, entityManager.componentsByClass[healthComponentName][entity.Id])
	assert.NotNil(t, entityManager.componentsByClass[speedComponentName][entity.Id])

	assert.NotNil(t, entityManager.componentsByClass[healthComponentName][secondEntity.Id])

	assert.Nil(t, entityManager.componentsByClass["invalidcomponent"])
}

func TestGetComponentOfClass(t *testing.T) {

	entityManager := CreateEntityManager()

	entity := entityManager.CreateEntity()
	secondEntity := entityManager.CreateEntity()

	speedComponent := components.CreateSpeedComponent(32)
	healthComponent := components.CreateHealthComponent(10)

	healthComponentName := healthComponent.GetComponentName()
	speedComponentName := speedComponent.GetComponentName()

	assert.Nil(t, entityManager.componentsByClass[speedComponentName])
	assert.Nil(t, entityManager.componentsByClass[healthComponentName])

	entityManager.AddComponentToEntity(entity.Id, speedComponent)
	entityManager.AddComponentToEntity(secondEntity.Id, healthComponent)

	storedSpeedComponent := entityManager.GetComponentOfClass(speedComponentName, entity.Id)
	storedHealthComponent := entityManager.GetComponentOfClass(healthComponentName, secondEntity.Id)
	noComponent := entityManager.GetComponentOfClass(healthComponentName, entity.Id)

	assert.NotNil(t, storedSpeedComponent)
	assert.NotNil(t, storedHealthComponent)
	assert.Equal(t, healthComponent, storedHealthComponent)
	assert.Equal(t, speedComponent, storedSpeedComponent)
	assert.Nil(t, noComponent)
}

func TestRemoveEntity(t *testing.T) {

	entityManager := CreateEntityManager()

	entity := entityManager.CreateEntity()
	secondEntity := entityManager.CreateEntity()

	speedComponent := components.CreateSpeedComponent(32)

	speedComponentName := speedComponent.GetComponentName()

	assert.Nil(t, entityManager.componentsByClass[speedComponentName])

	assert.Nil(t, entityManager.GetComponentOfClass(speedComponentName, entity.Id))
	assert.Nil(t, entityManager.GetComponentOfClass(speedComponentName, secondEntity.Id))

	entityManager.AddComponentToEntity(entity.Id, speedComponent)
	entityManager.AddComponentToEntity(secondEntity.Id, speedComponent)

	entityManager.RemoveEntity(entity.Id)

	returnedSpeedComponent := entityManager.GetComponentOfClass(speedComponentName, entity.Id)
	assert.Nil(t, returnedSpeedComponent)

	assert.False(t, ContainsEntity(entityManager.entities, entity))

	secondReturnedSpeedComponent := entityManager.GetComponentOfClass(speedComponentName, secondEntity.Id)
	assert.NotNil(t, secondReturnedSpeedComponent)
}

func TestGetAllEntitiesPossessingComponentsOfClass(t *testing.T) {

	entityManager := CreateEntityManager()

	entity := entityManager.CreateEntity()
	secondEntity := entityManager.CreateEntity()

	speedComponent := components.CreateSpeedComponent(32)
	healthComponent := components.CreateHealthComponent(10)

	healthComponentName := healthComponent.GetComponentName()
	speedComponentName := speedComponent.GetComponentName()

	entityManager.AddComponentToEntity(entity.Id, speedComponent)
	entityManager.AddComponentToEntity(entity.Id, healthComponent)
	entityManager.AddComponentToEntity(secondEntity.Id, healthComponent)

	healthEntities := entityManager.GetAllEntitiesPossessingComponentsOfClass(healthComponentName)
	speedEntities := entityManager.GetAllEntitiesPossessingComponentsOfClass(speedComponentName)
	positionEntities := entityManager.GetAllEntitiesPossessingComponentsOfClass("Position")

	assert.Equal(t, 2, len(healthEntities))
	assert.Equal(t, 1, len(speedEntities))
	assert.Equal(t, 0, len(positionEntities))
}

