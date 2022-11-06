package sdk

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/elmasy-com/columbus-sdk/fault"
	"github.com/elmasy-com/columbus-sdk/user"
)

const SLEEP_SEC = 5

// Global test user to play with
var TestUser user.User

func TestLookup200(t *testing.T) {

	// Test a proper lookup

	SetURI("http://localhost:8080/")

	subs, err := Lookup("example.com", true)
	if err != nil {
		t.Fatalf("FAILED: %s\n", err)
	}

	if len(subs) < 2 {
		t.Fatalf("FAIL: invalid number of subs for example.com: %v\n", subs)
	}
}

func TestLookup400(t *testing.T) {

	// Test a
	SetURI("http://localhost:8080/")

	_, err := Lookup("example", true)
	if !errors.Is(err, fault.ErrInvalidDomain) {
		t.Fatalf("FAILED: unexpected error: %s\n", err)
	}
}

func TestLookup403(t *testing.T) {

	// Test lookup if IP is blocked

	SetURI("http://localhost:8080/")

	_, err := GetUser("invalid")
	if !errors.Is(err, fault.ErrInvalidAPIKey) {
		t.Fatalf("FAIL: unexpected error: %s\n", err)
	}

	_, err = Lookup("exampleeeeeeeeeee.commmmmmmmmmmmmmm", true)
	if !errors.Is(err, fault.ErrBlocked) {
		t.Fatalf("FAILED: unexpected error: %s, expected ErrBlocked\n", err)
	}

	time.Sleep(SLEEP_SEC * time.Second)
}

func TestLookup404(t *testing.T) {

	// Test lookup with a domain that not exist

	SetURI("http://localhost:8080/")

	_, err := Lookup("exampleeeeeeeeeee.commmmmmmmmmmmmmm", true)
	if !errors.Is(err, fault.ErrNotFound) {
		t.Fatalf("FAILED: unexpected error: %s\n", err)
	}
}

func TestInsert200(t *testing.T) {

	// Test insert with a proper domain

	SetURI("http://localhost:8080/")

	err := GetDefaultUser(os.Getenv("COLUMBUS_API_KEY"))
	if err != nil {
		t.Fatalf("FAIL: GetDefaultUser(): %s\n", err)
	}

	err = Insert("www.example.com")
	if err != nil {
		t.Fatalf("FAIL: %s\n", err)
	}
}

func TestInsert400(t *testing.T) {

	// Test insert with an invalid domain

	SetURI("http://localhost:8080/")

	err := GetDefaultUser(os.Getenv("COLUMBUS_API_KEY"))
	if err != nil {
		t.Fatalf("FAIL: GetDefaultUser(): %s\n", err)
	}

	err = Insert("example")
	if !errors.Is(err, fault.ErrInvalidDomain) {
		t.Fatalf("FAIL: unexpected error: %s", err)
	}
}

func TestInsert401(t *testing.T) {

	// Test insert with invalid API key

	SetURI("http://localhost:8080/")

	// Save the valid API key and change the default for an invalid one
	tmp := DefaultUser.Key
	DefaultUser.Key = "invalid"

	err := Insert("example")
	if !errors.Is(err, fault.ErrInvalidAPIKey) {
		t.Fatalf("FAIL: unexpected error: %s, want ErrInvalidAPIKey", err)
	}

	// Restore the valid API key
	DefaultUser.Key = tmp
}

func TestInsert403(t *testing.T) {

	// Test insert with IP blocked
	// The previoud test caused an IP block
	SetURI("http://localhost:8080/")

	err := Insert("example")
	if !errors.Is(err, fault.ErrBlocked) {
		t.Fatalf("FAIL: unexpected error: %s, want ErrBlocked", err)
	}

	time.Sleep(SLEEP_SEC * time.Second)
}

func TestGetUser200(t *testing.T) {

	// Test get user with valid API key

	SetURI("http://localhost:8080/")

	user, err := GetUser(os.Getenv("COLUMBUS_API_KEY"))
	if err != nil {
		t.Fatalf("FAIL: %s\n", err)
	}

	if user.Key != os.Getenv("COLUMBUS_API_KEY") {
		t.Fatalf("FAIL: ENV key and user key differs\n")
	}
}

func TestGetUser401(t *testing.T) {

	// Test get user with invalid API key

	SetURI("http://localhost:8080/")

	_, err := GetUser("invalid")
	if !errors.Is(err, fault.ErrInvalidAPIKey) {
		t.Fatalf("FAIL: %s\n", err)
	}
}

// Test 403 after an invalid key
func TestGetUser403(t *testing.T) {

	// The previous function caused a block, so test get user with IP blocked and wait 10 sec

	SetURI("http://localhost:8080/")

	_, err := GetUser("invalid")
	if !errors.Is(err, fault.ErrBlocked) {
		t.Fatalf("FAIL: %s\n", err)
	}

	// The test server is configured for a 10 sec block time
	time.Sleep(SLEEP_SEC * time.Second)
}

func TestAddUser200(t *testing.T) {

	// Test add user with valid datas

	SetURI("http://localhost:8080/")

	err := GetDefaultUser(os.Getenv("COLUMBUS_API_KEY"))
	if err != nil {
		t.Fatalf("FAIL: Get user: %s\n", err)
	}

	user, err := AddUser("test", false)
	if err != nil {
		t.Fatalf("FAIL: %s\n", err)
	}

	TestUser = user
}

func TestAddUser401(t *testing.T) {

	// Test add user with invalid API key

	SetURI("http://localhost:8080/")

	tmp := DefaultUser.Key
	DefaultUser.Key = "invalid"

	_, err := AddUser("test", false)
	if !errors.Is(err, fault.ErrInvalidAPIKey) {
		t.Fatalf("FAIL: unexpected error: %s, want ErrInvalidAPIKey\n", err)
	}

	DefaultUser.Key = tmp
}

func TestAddUser403Blocked(t *testing.T) {

	// Test add user with blocked IP
	SetURI("http://localhost:8080/")

	_, err := AddUser("test", false)
	if !errors.Is(err, fault.ErrBlocked) {
		t.Fatalf("FAIL: unexpected error: %s, want ErrBlocked\n", err)
	}

	time.Sleep(SLEEP_SEC * time.Second)
}

func TestAddUser403NotAdmin(t *testing.T) {

	// Test add user with blocked IP
	SetURI("http://localhost:8080/")

	tmp := DefaultUser
	DefaultUser = &TestUser

	_, err := AddUser("test", false)
	if !errors.Is(err, fault.ErrNotAdmin) {
		t.Fatalf("FAIL: unexpected error: %s, want ErrNotAdmin\n", err)
	}

	DefaultUser = tmp

	time.Sleep(SLEEP_SEC * time.Second)
}

func TestAddUser409(t *testing.T) {

	// Test adduser with taken name

	SetURI("http://localhost:8080/")

	_, err := AddUser(DefaultUser.Name, false)
	if !errors.Is(err, fault.ErrNameTaken) {
		t.Fatalf("FAIL: unexpected error %s, want ErrNameTaken\n", err)
	}
}

func TestUserChangeName200(t *testing.T) {

	// Test user name change with valid data

	SetURI("http://localhost:8080/")

	err := ChangeName(&TestUser, "newtest")
	if err != nil {
		t.Fatalf("FAIL: %s\n", err)
	}

	if TestUser.Name != "newtest" {
		t.Fatalf("FAIL: TestUser.Name is not changed!")
	}
}

func TestUserChangeName401(t *testing.T) {

	// Test change username with invalid API key

	SetURI("http://localhost:8080/")

	tmp := TestUser.Key
	TestUser.Key = "invalid"

	err := ChangeName(&TestUser, "newtest")
	if !errors.Is(err, fault.ErrInvalidAPIKey) {
		t.Fatalf("FAIL: unexpected error: %s, want ErrInvalidAPIKey\n", err)
	}

	TestUser.Key = tmp
}

func TestUserChangeName403(t *testing.T) {

	// Test change username with invalid API key
	// The previous test caused an IP block

	SetURI("http://localhost:8080/")

	err := ChangeName(&TestUser, "newtest")
	if !errors.Is(err, fault.ErrBlocked) {
		t.Fatalf("FAIL: unexpected error: %s, want ErrBlocked\n", err)
	}

	time.Sleep(SLEEP_SEC * time.Second)
}

func TestUserChangeName409(t *testing.T) {

	// Test change username with a taken username

	SetURI("http://localhost:8080/")

	err := ChangeName(&TestUser, DefaultUser.Name)
	if !errors.Is(err, fault.ErrNameTaken) {
		t.Fatalf("FAIL: unexpected error: %s, want ErrNameTaken\n", err)
	}
}

func TestUserChangeKey200(t *testing.T) {

	// Test user key change

	SetURI("http://localhost:8080/")

	oldKey := TestUser.Key

	err := ChangeKey(&TestUser)
	if err != nil {
		t.Fatalf("FAIL: %s\n", err)
	}

	if TestUser.Key == oldKey {
		t.Fatalf("FAIL: TestUser.Key is not changed!")
	}
}

func TestUserChangeKey401(t *testing.T) {

	// Test user key change with invalid key

	SetURI("http://localhost:8080/")

	tmp := TestUser.Key
	TestUser.Key = "invalid"

	err := ChangeKey(&TestUser)
	if !errors.Is(err, fault.ErrInvalidAPIKey) {
		t.Fatalf("FAIL: unexpected error: %s, want ErrInvalidAPIKey\n", err)
	}

	TestUser.Key = tmp
}

func TestUserChangeKey403(t *testing.T) {

	// Test user name change with blocked IP

	SetURI("http://localhost:8080/")

	err := ChangeKey(&TestUser)
	if !errors.Is(err, fault.ErrBlocked) {
		t.Fatalf("FAIL: unexpected error: %s, want ErrBlocked\n", err)
	}

	time.Sleep(SLEEP_SEC * time.Second)
}

// TODO:
// Change other user admin

func TestChangeOtherUserName200(t *testing.T) {

	// Test change other user name

	SetURI("http://localhost:8080/")

	err := ChangeOtherUserName(&TestUser, "test")
	if err != nil {
		t.Fatalf("FAIL: %s\n", err)
	}

	if TestUser.Name != "test" {
		t.Fatalf("FAIL: TestUser.Name not changed\n")
	}
}

func TestChangeOtherUserName400(t *testing.T) {

	// Test change other user name

	SetURI("http://localhost:8080/")

	err := ChangeOtherUserName(&TestUser, TestUser.Name)
	if !errors.Is(err, fault.ErrSameName) {
		t.Fatalf("FAIL: unexpected error: %s, want ErrSameName\n", err)
	}
}

func TestChangeOtherUserName401(t *testing.T) {

	// Test change other user name with invalid API key

	SetURI("http://localhost:8080/")

	tmp := DefaultUser.Key
	DefaultUser.Key = "invalid"

	err := ChangeOtherUserName(&TestUser, "test")
	if !errors.Is(err, fault.ErrInvalidAPIKey) {
		t.Fatalf("FAIL: unexpected error: %s, want ErrInvalidAPIKey\n", err)
	}

	DefaultUser.Key = tmp
}

func TestChangeOtherUserName403Blocked(t *testing.T) {

	// Test change other user name with blocked IP

	SetURI("http://localhost:8080/")

	err := ChangeOtherUserName(&TestUser, "test")
	if !errors.Is(err, fault.ErrBlocked) {
		t.Fatalf("FAIL: unexpected error: %s, want ErrBlocked\n", err)
	}

	time.Sleep(SLEEP_SEC * time.Second)
}

func TestChangeOtherUserName403NotAdmin(t *testing.T) {

	// Test change other user name with blocked IP

	SetURI("http://localhost:8080/")

	tmp := DefaultUser
	DefaultUser = &TestUser

	err := ChangeOtherUserName(&TestUser, "test")
	if !errors.Is(err, fault.ErrNotAdmin) {
		t.Fatalf("FAIL: unexpected error: %s, want ErrNotAdmin\n", err)
	}

	DefaultUser = tmp

	time.Sleep(SLEEP_SEC * time.Second)
}

func TestChangeOtherUserName404(t *testing.T) {

	// Test change other user name with invalid username

	SetURI("http://localhost:8080/")

	TestUser.Name = "notexist"

	err := ChangeOtherUserName(&TestUser, "test")
	if !errors.Is(err, fault.ErrUserNotFound) {
		t.Fatalf("FAIL: unexpected error: %s, want ErrUserNotFound\n", err)
	}

	TestUser.Name = "test"
}

func TestChangeOtherUserName409(t *testing.T) {

	// Test change other user name with taken username

	SetURI("http://localhost:8080/")

	err := ChangeOtherUserName(&TestUser, DefaultUser.Name)
	if !errors.Is(err, fault.ErrNameTaken) {
		t.Fatalf("FAIL: unexpected error: %s, want ErrNameTaken\n", err)
	}
}

func TestChangeOtherUserKey200(t *testing.T) {

	// Test change other user key

	SetURI("http://localhost:8080/")

	oldKey := TestUser.Key

	err := ChangeOtherUserKey(&TestUser)
	if err != nil {
		t.Fatalf("FAIL: %s\n", err)
	}

	if TestUser.Key == oldKey {
		t.Fatalf("FAIL: TestUser.Key not changed!\n")
	}
}

func TestChangeOtherUserKey401(t *testing.T) {

	// Test change other user key with invalid API key

	SetURI("http://localhost:8080/")

	tmp := DefaultUser.Key
	DefaultUser.Key = "invalid"

	err := ChangeOtherUserKey(&TestUser)
	if !errors.Is(err, fault.ErrInvalidAPIKey) {
		t.Fatalf("FAIL: unexpected error: %s, want ErrInvalidAPIKey\n", err)
	}

	DefaultUser.Key = tmp
}

func TestChangeOtherUserKey403Blocked(t *testing.T) {

	// Test change other user key with blocked IP

	SetURI("http://localhost:8080/")

	err := ChangeOtherUserKey(&TestUser)
	if !errors.Is(err, fault.ErrBlocked) {
		t.Fatalf("FAIL: unexpected error: %s, want ErrBlocked\n", err)
	}

	time.Sleep(SLEEP_SEC * time.Second)
}

func TestChangeOtherUserKey403NotAdmin(t *testing.T) {

	// Test change other user name with blocked IP

	SetURI("http://localhost:8080/")

	tmp := DefaultUser
	DefaultUser = &TestUser

	err := ChangeOtherUserKey(&TestUser)
	if !errors.Is(err, fault.ErrNotAdmin) {
		t.Fatalf("FAIL: unexpected error: %s, want ErrNotAdmin\n", err)
	}

	DefaultUser = tmp

	time.Sleep(SLEEP_SEC * time.Second)
}

func TestChangeOtherUserKey404(t *testing.T) {

	// Test change other user key with invalid username

	SetURI("http://localhost:8080/")

	TestUser.Name = "notexist"

	err := ChangeOtherUserKey(&TestUser)
	if !errors.Is(err, fault.ErrUserNotFound) {
		t.Fatalf("FAIL: unexpected error: %s, want ErrUserNotFound\n", err)
	}

	TestUser.Name = "test"
}

func TestChangeOtherUserAdmin200True(t *testing.T) {

	// Test change other user admin

	SetURI("http://localhost:8080/")

	err := ChangeOtherUserAdmin(&TestUser, true)
	if err != nil {
		t.Fatalf("FAIL: %s\n", err)
	}
}

func TestChangeOtherUserAdmin400(t *testing.T) {

	// Test change other user admin

	SetURI("http://localhost:8080/")

	err := ChangeOtherUserAdmin(&TestUser, TestUser.Admin)
	if !errors.Is(err, fault.ErrNothingToDo) {
		t.Fatalf("FAIL: unexpected error: %s, want ErrNothingToDo\n", err)
	}
}

func TestChangeOtherUserAdmin200False(t *testing.T) {

	// Test change other user admin

	SetURI("http://localhost:8080/")

	err := ChangeOtherUserAdmin(&TestUser, false)
	if err != nil {
		t.Fatalf("FAIL: %s\n", err)
	}

}

func TestChangeOtherUserAdmin401(t *testing.T) {

	// Test change other user admin with invalid API key

	SetURI("http://localhost:8080/")

	tmp := DefaultUser.Key
	DefaultUser.Key = "invalid"

	err := ChangeOtherUserAdmin(&TestUser, true)
	if !errors.Is(err, fault.ErrInvalidAPIKey) {
		t.Fatalf("FAIL: unexpected error: %s, want ErrInvalidAPIKey\n", err)
	}

	DefaultUser.Key = tmp
}

func TestChangeOtherUserAdmin403Blocked(t *testing.T) {

	// Test change other user admin with blocked IP

	SetURI("http://localhost:8080/")

	err := ChangeOtherUserAdmin(&TestUser, true)
	if !errors.Is(err, fault.ErrBlocked) {
		t.Fatalf("FAIL: unexpected error: %s, want ErrBlocked\n", err)
	}

	time.Sleep(SLEEP_SEC * time.Second)
}

func TestChangeOtherUserAdmin403NotAdmin(t *testing.T) {

	// Test change other user name with blocked IP

	SetURI("http://localhost:8080/")

	tmp := DefaultUser
	DefaultUser = &TestUser

	err := ChangeOtherUserAdmin(&TestUser, false)
	if !errors.Is(err, fault.ErrNotAdmin) {
		t.Fatalf("FAIL: unexpected error: %s, want ErrNotAdmin\n", err)
	}

	DefaultUser = tmp

	time.Sleep(SLEEP_SEC * time.Second)
}

func TestChangeOtherUserAdmin404(t *testing.T) {

	// Test change other user admin with invalid username

	SetURI("http://localhost:8080/")

	TestUser.Name = "notexist"

	err := ChangeOtherUserAdmin(&TestUser, true)
	if !errors.Is(err, fault.ErrUserNotFound) {
		t.Fatalf("FAIL: unexpected error: %s, want ErrUserNotFound\n", err)
	}

	TestUser.Name = "test"
}

func TestGetUsers200(t *testing.T) {

	SetURI("http://localhost:8080/")

	us, err := GetUsers()
	if err != nil {
		t.Fatalf("FAIL: %s\n", err)
	}

	if len(us) != 2 {
		t.Fatalf("FAIL: invalid user number: %#v\n", us)
	}
}

func TestGetUsers401(t *testing.T) {

	SetURI("http://localhost:8080/")

	tmp := DefaultUser.Key
	DefaultUser.Key = "invalid"

	_, err := GetUsers()
	if !errors.Is(err, fault.ErrInvalidAPIKey) {
		t.Fatalf("FAIL: unexpected error: %s, want ErrInvalidAPIKey\n", err)
	}

	DefaultUser.Key = tmp
}

func TestGetUsers403Blocked(t *testing.T) {

	// Test IP block

	SetURI("http://localhost:8080/")

	_, err := GetUsers()
	if !errors.Is(err, fault.ErrBlocked) {
		t.Fatalf("FAIL: unexpected error: %s, want ErrBlocked\n", err)
	}

	time.Sleep(SLEEP_SEC * time.Second)
}

func TestGetUsers403NotAdmin(t *testing.T) {

	// Test IP block

	SetURI("http://localhost:8080/")

	tmp := DefaultUser
	DefaultUser = &TestUser

	_, err := GetUsers()
	if !errors.Is(err, fault.ErrNotAdmin) {
		t.Fatalf("FAIL: unexpected error: %s, want ErrNotAdmin\n", err)
	}

	DefaultUser = tmp

	time.Sleep(SLEEP_SEC * time.Second)
}

func TestUserDelete200(t *testing.T) {

	// Test valid delete user

	SetURI("http://localhost:8080/")

	err := Delete(TestUser, true)
	if err != nil {
		t.Fatalf("FAIL: %s\n", err)
	}
}

func TestUserDelete401(t *testing.T) {

	// Test valid delete user with invalid API key
	// The user is deleted in the previous test

	SetURI("http://localhost:8080/")

	err := Delete(TestUser, true)
	if !errors.Is(err, fault.ErrInvalidAPIKey) {
		t.Fatalf("FAIL: unexpected error: %s, want ErrInvalidAPIKey\n", err)
	}
}

func TestUserDelete403(t *testing.T) {

	// Test valid delete user with invalid API key

	SetURI("http://localhost:8080/")

	err := Delete(TestUser, true)
	if !errors.Is(err, fault.ErrBlocked) {
		t.Fatalf("FAIL: unexpected error: %s, want ErrBlocked\n", err)
	}

	time.Sleep(SLEEP_SEC * time.Second)
}
