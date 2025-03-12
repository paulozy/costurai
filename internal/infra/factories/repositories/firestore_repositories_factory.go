package factories

import (
	"github.com/paulozy/costurai/internal/infra/database"
	"github.com/paulozy/costurai/internal/infra/database/firestore/repositories"
)

type FirestoreRepositoriesFactoryOutput struct {
	DressmakerRepository       *repositories.FirestoreDressmakerRepository
	UserRepository             *repositories.FirestoreUserRepository
	DressmakerReviewRepository *repositories.FirestoreDressmakerReviewsRepository
}

func FirestoreRepositoriesFactory(
	databaseInstance database.DatabaseInterface,
) *FirestoreRepositoriesFactoryOutput {
	dressMakerRepository := repositories.NewFirestoreDressmakerRepository(databaseInstance)
	userRepository := repositories.NewFirestoreUserRepository(databaseInstance)
	dressmakerReviewsRepository := repositories.NewFirestoreReviewsRepository(databaseInstance)

	return &FirestoreRepositoriesFactoryOutput{
		DressmakerRepository:       dressMakerRepository,
		UserRepository:             userRepository,
		DressmakerReviewRepository: dressmakerReviewsRepository,
	}
}
