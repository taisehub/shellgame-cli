package usecase

type GameUsecase struct {
}

func NewGameUsecase() *GameUsecase {
	return &GameUsecase{}
}

func (gc *GameUsecase) Hello() string {
	return "Hello!!"
}

// マッチング完了後、バトルを開始するためのメソッド
func (gc *GameUsecase) Start() {
	panic("implement me")
}

// ゲームのマッチングをするためのメソッド
// 詳細まで考えれていないので、最終的にはもう少し細かく分かれる
func (gu *GameUsecase) Matching() {
	panic("implement me")
}