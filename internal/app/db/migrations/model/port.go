package model

type Port struct {
	Name        string `gorm:"name"`
	City        string `gorm:"city"`
	Country     string `gorm:"country"`
	Alias       string `gorm:"alias"`
	Regions     string `gorm:"regions"`
	Coordinates string `gorm:"coordinates"`
	Province    string `gorm:"province"`
	Timezone    string `gorm:"timezone"`
	Unlocs      string `gorm:"unlocs"`
	Code        string `gorm:"code"`
}
