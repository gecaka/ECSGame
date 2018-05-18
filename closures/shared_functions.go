package closures

import "BasicECS/core"

type InteractionFunc func(
	entityManager *core.EntityManager,
	entityId int,
	interactingEntityId int) bool

type UsageFunc func(entityManager core.EntityManager) bool