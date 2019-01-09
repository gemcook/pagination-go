package pagination_test

import pagination "github.com/gemcook/pagination-go"

type fruit struct {
	Name  string
	Price int
}

type fruitCondition struct {
	PriceLowerLimit  *int
	PriceHigherLimit *int
}

func newFruitCondition(low, high int) *fruitCondition {
	return &fruitCondition{
		PriceLowerLimit:  &low,
		PriceHigherLimit: &high,
	}
}

func (fc *fruitCondition) ApplyCondition(s interface{}) {
	fetcher, ok := s.(*fruitFetcher)
	if !ok {
		return
	}
	if fc.PriceHigherLimit != nil {
		fetcher.priceHigherLimit = *fc.PriceHigherLimit
	}
	if fc.PriceLowerLimit != nil {
		fetcher.priceLowerLimit = *fc.PriceLowerLimit
	}
}

type fruitFetcher struct {
	priceLowerLimit  int
	priceHigherLimit int
}

func newFruitFetcher() *fruitFetcher {
	return &fruitFetcher{
		priceLowerLimit:  -1 << 31,
		priceHigherLimit: 1<<31 - 1,
	}
}

func (ff *fruitFetcher) Count(cond pagination.ConditionApplier) (int, error) {
	if cond != nil {
		cond.ApplyCondition(ff)
	}
	dummyFruits := ff.GetDummy()
	return len(dummyFruits), nil
}

func (ff *fruitFetcher) FetchPage(limit, offset int, cond pagination.ConditionApplier, orders []*pagination.Order, result *pagination.PageFetchResult) error {
	if cond != nil {
		cond.ApplyCondition(ff)
	}
	dummyFruits := ff.GetDummy()
	var toIndex int
	toIndex = offset + limit
	if toIndex > len(dummyFruits) {
		toIndex = len(dummyFruits)
	}
	for _, fruit := range dummyFruits[offset:toIndex] {
		*result = append(*result, fruit)
	}
	return nil
}

func (ff *fruitFetcher) GetDummy() []fruit {
	result := make([]fruit, 0)
	for _, f := range dummyFruits {
		if ff.priceHigherLimit >= f.Price && f.Price >= ff.priceLowerLimit {
			result = append(result, f)
		}
	}
	return result
}

var dummyFruits = []fruit{
	fruit{"Apple", 112},
	fruit{"Pear", 245},
	fruit{"Banana", 60},
	fruit{"Orange", 80},
	fruit{"Kiwi", 106},
	fruit{"Strawberry", 350},
	fruit{"Grape", 400},
	fruit{"Grapefruit", 150},
	fruit{"Pineapple", 200},
	fruit{"Cherry", 140},
	fruit{"Mango", 199},
}

type LargeData struct {
	ID int
}

type LargeDataFetcher struct{}

func newLargeDataFetcher() *LargeDataFetcher {
	return &LargeDataFetcher{}
}

func (ff *LargeDataFetcher) Count(cond pagination.ConditionApplier) (int, error) {
	return len(dummyLargeList), nil
}

func (ff *LargeDataFetcher) FetchPage(limit, offset int, cond pagination.ConditionApplier, orders []*pagination.Order, result *pagination.PageFetchResult) error {
	var toIndex int
	toIndex = offset + limit
	if toIndex > len(dummyLargeList) {
		toIndex = len(dummyLargeList)
	}
	for _, LargeData := range dummyLargeList[offset:toIndex] {
		*result = append(*result, LargeData)
	}
	return nil
}

var dummyLargeList = []LargeData{
	LargeData{1},
	LargeData{2},
	LargeData{3},
	LargeData{4},
	LargeData{5},
	LargeData{6},
	LargeData{7},
	LargeData{8},
	LargeData{9},
	LargeData{10},
	LargeData{11},
	LargeData{12},
	LargeData{13},
	LargeData{14},
	LargeData{15},
	LargeData{16},
	LargeData{17},
	LargeData{18},
	LargeData{19},
	LargeData{20},
	LargeData{21},
	LargeData{22},
	LargeData{23},
	LargeData{24},
	LargeData{25},
	LargeData{26},
	LargeData{27},
	LargeData{28},
	LargeData{29},
	LargeData{30},
	LargeData{31},
	LargeData{32},
	LargeData{33},
	LargeData{34},
	LargeData{35},
	LargeData{36},
	LargeData{37},
	LargeData{38},
	LargeData{39},
	LargeData{40},
	LargeData{41},
	LargeData{42},
	LargeData{43},
	LargeData{44},
	LargeData{45},
	LargeData{46},
	LargeData{47},
	LargeData{48},
	LargeData{49},
	LargeData{50},
	LargeData{51},
	LargeData{52},
	LargeData{53},
	LargeData{54},
	LargeData{55},
	LargeData{56},
	LargeData{57},
	LargeData{58},
	LargeData{59},
	LargeData{60},
	LargeData{61},
	LargeData{62},
	LargeData{63},
	LargeData{64},
	LargeData{65},
	LargeData{66},
	LargeData{67},
	LargeData{68},
	LargeData{69},
	LargeData{70},
	LargeData{71},
	LargeData{72},
	LargeData{73},
	LargeData{74},
	LargeData{75},
	LargeData{76},
	LargeData{77},
	LargeData{78},
	LargeData{79},
	LargeData{80},
	LargeData{81},
	LargeData{82},
	LargeData{83},
	LargeData{84},
	LargeData{85},
	LargeData{86},
	LargeData{87},
	LargeData{88},
	LargeData{89},
	LargeData{90},
	LargeData{91},
	LargeData{92},
	LargeData{93},
	LargeData{94},
	LargeData{95},
	LargeData{96},
	LargeData{97},
	LargeData{98},
	LargeData{99},
	LargeData{100},
	LargeData{101},
	LargeData{102},
	LargeData{103},
}
