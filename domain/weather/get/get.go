package get

import (
	"context"
	"fmt"
)

type GetWeather struct{}

func NewUseCase() *GetWeather {
	return &GetWeather{}
}

func (uc *GetWeather) Get(ctx context.Context) error {
	fmt.Println("GetWeather Process")
	return nil
}
