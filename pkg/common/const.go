package common

var AppTypeEnum = struct {
	Shop     int
	Open     int
	External int
}{
	Shop:     1,
	Open:     3,
	External: 2,
}

var BindStatusEnum = struct {
	No  int
	Yes int
}{
	No:  0,
	Yes: 1,
}
