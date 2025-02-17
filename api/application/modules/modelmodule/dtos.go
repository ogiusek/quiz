package modelmodule

type ModelDto struct {
	Id        ModelId `json:"id"`
	CreatedAt string  `json:"created_at"`
}

func (model *Model) Dto() ModelDto {
	return ModelDto{
		Id:        model.Id,
		CreatedAt: model.CreatedAt.Format(),
	}
}

type CreatedDto struct {
	Id        ModelId `json:"id"`
	CreatedAt string  `json:"created_at"`
}

func (model *Model) CreatedDto() CreatedDto {
	return CreatedDto{
		Id:        model.Id,
		CreatedAt: model.CreatedAt.Format(),
	}
}
