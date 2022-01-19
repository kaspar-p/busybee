package user

var Users map[string] *User

type CurrentlyBusy struct {
	IsBusy bool
	BusyWith string
}

type User struct {
	UserID string
	Name string
	CurrentlyBusy CurrentlyBusy
	BusyTimes []*BusyTime
}

func InitializeUsers() {
	Users = make(map[string] *User);
}

func CreateUser(userName string, userID string) *User {
	user := User{
		Name: userName,
		UserID: userID,
		CurrentlyBusy: CurrentlyBusy {
			IsBusy: false,
			BusyWith: "",
		},
	}

	return &user;
}