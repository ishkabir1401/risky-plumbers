package settings

import "risky-plumbers/model"

var RiskStoreStruct *model.RiskStore

func InitializeRiskStore() {
	RiskStoreStruct = &model.RiskStore{
		Risks: make(map[string]*model.Risk),
	}
}
