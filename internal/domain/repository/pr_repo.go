package repository

type PRRepository interface {
	Create()
	Merge()
	Reassign()
}
