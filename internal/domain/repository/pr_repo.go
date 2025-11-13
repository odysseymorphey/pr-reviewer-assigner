package repository

type PRRepo interface {
	Create()
	Merge()
	Reassign()
}
