package main

type propertyA struct {
	A1 int
	A2 map[string]string
	A3 string
}

type propertyB struct {
	B1 int
	B2 map[string]string
	B3 string
}

func BuildProperty(A_column1 int,
	A_column2 map[string]string,
	A_column3 string,
	B_column1 int,
	B_column2 map[string]string,
	B_column3 string) (struct1 propertyA, struct2 propertyB) {

	propA := propertyA{
		A1: A_column1,
		A2: A_column2,
		A3: A_column3,
	}

	propB := propertyB{
		B1: B_column1,
		B2: B_column2,
		B3: B_column3,
	}

	return propA, propB
}
