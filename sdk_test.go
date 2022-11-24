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

// Test a proper lookup.
func TestLookup200(t *testing.T) {

	SetURI("http://localhost:8080/")

	_, err := Lookup("example.com", true)
	if err != nil {
		t.Fatalf("FAILED: %s\n", err)
	}
}

// Test a lookup with invalid domain.
func TestLookup400(t *testing.T) {

	SetURI("http://localhost:8080/")

	_, err := Lookup("example", true)
	if !errors.Is(err, fault.ErrInvalidDomain) {
		t.Fatalf("FAILED: unexpected error: %s\n", err)
	}
}

// Test lookup if IP is blocked.
// The block is caused by querying an invalid API key.
func TestLookup403(t *testing.T) {

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

// Test lookup with a domain that not exist.
func TestLookup404(t *testing.T) {

	SetURI("http://localhost:8080/")

	_, err := Lookup("exampleeeeeeeeeee.commmmmmmmmmmmmmm", true)
	if !errors.Is(err, fault.ErrNotFound) {
		t.Fatalf("FAILED: unexpected error: %s, want ErrNotFound\n", err)
	}
}

// Test insert with a proper domain.
func TestInsert200(t *testing.T) {

	SetURI("http://localhost:8080/")

	err := GetDefaultUser(os.Getenv("COLUMBUS_TEST_KEY"))
	if err != nil {
		t.Fatalf("FAIL: GetDefaultUser(): %s\n", err)
	}

	err = Insert("www.example.com")
	if err != nil {
		t.Fatalf("FAIL: %s\n", err)
	}
}

// Test insert with an invalid domain.
func TestInsert400InvalidDomain(t *testing.T) {

	SetURI("http://localhost:8080/")

	err := Insert("example")
	if !errors.Is(err, fault.ErrInvalidDomain) {
		t.Fatalf("FAIL: unexpected error: %s, want ErrInvalidDomain", err)
	}
}

// Test insert with invalid API key.
func TestInsert401(t *testing.T) {

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

// Test insert with IP blocked.
// The previous test caused an IP block.
func TestInsert403(t *testing.T) {

	SetURI("http://localhost:8080/")

	err := Insert("example")
	if !errors.Is(err, fault.ErrBlocked) {
		t.Fatalf("FAIL: unexpected error: %s, want ErrBlocked", err)
	}

	time.Sleep(SLEEP_SEC * time.Second)
}

// Test get user with valid API key.
func TestGetUser200(t *testing.T) {

	SetURI("http://localhost:8080/")

	user, err := GetUser(os.Getenv("COLUMBUS_TEST_KEY"))
	if err != nil {
		t.Fatalf("FAIL: %s\n", err)
	}

	if user.Key != os.Getenv("COLUMBUS_TEST_KEY") {
		t.Fatalf("FAIL: ENV key and user key differs\n")
	}
}

// Test get user with invalid API key.
func TestGetUser401(t *testing.T) {

	SetURI("http://localhost:8080/")

	_, err := GetUser("invalid")
	if !errors.Is(err, fault.ErrInvalidAPIKey) {
		t.Fatalf("FAIL: %s\n", err)
	}
}

// Test 403 after an invalid key.
// The previous function caused a block.
func TestGetUser403(t *testing.T) {

	SetURI("http://localhost:8080/")

	_, err := GetUser("invalid")
	if !errors.Is(err, fault.ErrBlocked) {
		t.Fatalf("FAIL: %s\n", err)
	}

	// The test server is configured for a 10 sec block time
	time.Sleep(SLEEP_SEC * time.Second)
}

// Test add user with valid datas.
func TestAddUser200(t *testing.T) {

	SetURI("http://localhost:8080/")

	err := GetDefaultUser(os.Getenv("COLUMBUS_TEST_KEY"))
	if err != nil {
		t.Fatalf("FAIL: Get user: %s\n", err)
	}

	user, err := AddUser("test", false)
	if err != nil {
		t.Fatalf("FAIL: %s\n", err)
	}

	TestUser = user
}

// Test add user with invalid API key.
func TestAddUser401(t *testing.T) {

	SetURI("http://localhost:8080/")

	tmp := DefaultUser.Key
	DefaultUser.Key = "invalid"

	_, err := AddUser("test", false)
	if !errors.Is(err, fault.ErrInvalidAPIKey) {
		t.Fatalf("FAIL: unexpected error: %s, want ErrInvalidAPIKey\n", err)
	}

	DefaultUser.Key = tmp
}

// Test add user with blocked IP.
// The previous function caused a block.
func TestAddUser403Blocked(t *testing.T) {

	SetURI("http://localhost:8080/")

	_, err := AddUser("test", false)
	if !errors.Is(err, fault.ErrBlocked) {
		t.Fatalf("FAIL: unexpected error: %s, want ErrBlocked\n", err)
	}

	time.Sleep(SLEEP_SEC * time.Second)
}

// Test add user with a non admin account.
func TestAddUser403NotAdmin(t *testing.T) {

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

// Test adduser with taken name.
func TestAddUser409(t *testing.T) {

	SetURI("http://localhost:8080/")

	_, err := AddUser(DefaultUser.Name, false)
	if !errors.Is(err, fault.ErrNameTaken) {
		t.Fatalf("FAIL: unexpected error %s, want ErrNameTaken\n", err)
	}
}

// Test user name change with valid data.
func TestUserChangeName200(t *testing.T) {

	SetURI("http://localhost:8080/")

	err := ChangeName(&TestUser, "newtest")
	if err != nil {
		t.Fatalf("FAIL: %s\n", err)
	}

	if TestUser.Name != "newtest" {
		t.Fatalf("FAIL: TestUser.Name is not changed!")
	}
}

// Test change username with invalid API key.
func TestUserChangeName401(t *testing.T) {

	SetURI("http://localhost:8080/")

	tmp := TestUser.Key
	TestUser.Key = "invalid"

	err := ChangeName(&TestUser, "newtest")
	if !errors.Is(err, fault.ErrInvalidAPIKey) {
		t.Fatalf("FAIL: unexpected error: %s, want ErrInvalidAPIKey\n", err)
	}

	TestUser.Key = tmp
}

// Test change username with invalid API key.
// The previous test caused an IP block.
func TestUserChangeName403(t *testing.T) {

	SetURI("http://localhost:8080/")

	err := ChangeName(&TestUser, "newtest")
	if !errors.Is(err, fault.ErrBlocked) {
		t.Fatalf("FAIL: unexpected error: %s, want ErrBlocked\n", err)
	}

	time.Sleep(SLEEP_SEC * time.Second)
}

// Test change username with a taken username.
func TestUserChangeName409(t *testing.T) {

	SetURI("http://localhost:8080/")

	err := ChangeName(&TestUser, DefaultUser.Name)
	if !errors.Is(err, fault.ErrNameTaken) {
		t.Fatalf("FAIL: unexpected error: %s, want ErrNameTaken\n", err)
	}
}

// Test user key change.
func TestUserChangeKey200(t *testing.T) {

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

// Test user key change with invalid key.
func TestUserChangeKey401(t *testing.T) {

	SetURI("http://localhost:8080/")

	tmp := TestUser.Key
	TestUser.Key = "invalid"

	err := ChangeKey(&TestUser)
	if !errors.Is(err, fault.ErrInvalidAPIKey) {
		t.Fatalf("FAIL: unexpected error: %s, want ErrInvalidAPIKey\n", err)
	}

	TestUser.Key = tmp
}

// Test user name change with blocked IP.
// The previous test caused an IP block.
func TestUserChangeKey403(t *testing.T) {

	SetURI("http://localhost:8080/")

	err := ChangeKey(&TestUser)
	if !errors.Is(err, fault.ErrBlocked) {
		t.Fatalf("FAIL: unexpected error: %s, want ErrBlocked\n", err)
	}

	time.Sleep(SLEEP_SEC * time.Second)
}

// Test change other user name.
func TestChangeOtherUserName200(t *testing.T) {

	SetURI("http://localhost:8080/")

	err := ChangeOtherUserName(&TestUser, "test")
	if err != nil {
		t.Fatalf("FAIL: %s\n", err)
	}

	if TestUser.Name != "test" {
		t.Fatalf("FAIL: TestUser.Name not changed\n")
	}
}

// Test change other user name to the same name.
func TestChangeOtherUserName400(t *testing.T) {

	SetURI("http://localhost:8080/")

	err := ChangeOtherUserName(&TestUser, TestUser.Name)
	if !errors.Is(err, fault.ErrSameName) {
		t.Fatalf("FAIL: unexpected error: %s, want ErrSameName\n", err)
	}
}

// Test change other user name with invalid API key.
func TestChangeOtherUserName401(t *testing.T) {

	SetURI("http://localhost:8080/")

	tmp := DefaultUser.Key
	DefaultUser.Key = "invalid"

	err := ChangeOtherUserName(&TestUser, "test")
	if !errors.Is(err, fault.ErrInvalidAPIKey) {
		t.Fatalf("FAIL: unexpected error: %s, want ErrInvalidAPIKey\n", err)
	}

	DefaultUser.Key = tmp
}

// Test change other user name with blocked IP.
// The previous test caused an IP block.
func TestChangeOtherUserName403Blocked(t *testing.T) {

	SetURI("http://localhost:8080/")

	err := ChangeOtherUserName(&TestUser, "test")
	if !errors.Is(err, fault.ErrBlocked) {
		t.Fatalf("FAIL: unexpected error: %s, want ErrBlocked\n", err)
	}

	time.Sleep(SLEEP_SEC * time.Second)
}

// Test change other user name with a non admin account.
func TestChangeOtherUserName403NotAdmin(t *testing.T) {

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

// Test change other user name with an invalid username.
func TestChangeOtherUserName404(t *testing.T) {

	SetURI("http://localhost:8080/")

	TestUser.Name = "notexist"

	err := ChangeOtherUserName(&TestUser, "test")
	if !errors.Is(err, fault.ErrUserNotFound) {
		t.Fatalf("FAIL: unexpected error: %s, want ErrUserNotFound\n", err)
	}

	TestUser.Name = "test"
}

// Test change other user name with taken username.
func TestChangeOtherUserName409(t *testing.T) {

	SetURI("http://localhost:8080/")

	err := ChangeOtherUserName(&TestUser, DefaultUser.Name)
	if !errors.Is(err, fault.ErrNameTaken) {
		t.Fatalf("FAIL: unexpected error: %s, want ErrNameTaken\n", err)
	}
}

// Test change other user key.
func TestChangeOtherUserKey200(t *testing.T) {

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

// Test change other user key with invalid API key.
func TestChangeOtherUserKey401(t *testing.T) {

	SetURI("http://localhost:8080/")

	tmp := DefaultUser.Key
	DefaultUser.Key = "invalid"

	err := ChangeOtherUserKey(&TestUser)
	if !errors.Is(err, fault.ErrInvalidAPIKey) {
		t.Fatalf("FAIL: unexpected error: %s, want ErrInvalidAPIKey\n", err)
	}

	DefaultUser.Key = tmp
}

// Test change other user key with blocked IP.
// The previous test caused an IP block.
func TestChangeOtherUserKey403Blocked(t *testing.T) {

	SetURI("http://localhost:8080/")

	err := ChangeOtherUserKey(&TestUser)
	if !errors.Is(err, fault.ErrBlocked) {
		t.Fatalf("FAIL: unexpected error: %s, want ErrBlocked\n", err)
	}

	time.Sleep(SLEEP_SEC * time.Second)
}

// Test change other user name with non admin user.
func TestChangeOtherUserKey403NotAdmin(t *testing.T) {

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

// Test change other user key with invalid username.
func TestChangeOtherUserKey404(t *testing.T) {

	SetURI("http://localhost:8080/")

	TestUser.Name = "notexist"

	err := ChangeOtherUserKey(&TestUser)
	if !errors.Is(err, fault.ErrUserNotFound) {
		t.Fatalf("FAIL: unexpected error: %s, want ErrUserNotFound\n", err)
	}

	TestUser.Name = "test"
}

// Test change other user admin.
func TestChangeOtherUserAdmin200True(t *testing.T) {

	SetURI("http://localhost:8080/")

	err := ChangeOtherUserAdmin(&TestUser, true)
	if err != nil {
		t.Fatalf("FAIL: %s\n", err)
	}
}

// Test change other user admin with the same value.
func TestChangeOtherUserAdmin400(t *testing.T) {

	SetURI("http://localhost:8080/")

	err := ChangeOtherUserAdmin(&TestUser, TestUser.Admin)
	if !errors.Is(err, fault.ErrNothingToDo) {
		t.Fatalf("FAIL: unexpected error: %s, want ErrNothingToDo\n", err)
	}
}

// Test change other user admin to false.
func TestChangeOtherUserAdmin200False(t *testing.T) {

	SetURI("http://localhost:8080/")

	err := ChangeOtherUserAdmin(&TestUser, false)
	if err != nil {
		t.Fatalf("FAIL: %s\n", err)
	}
}

// Test change other user admin with invalid API key.
func TestChangeOtherUserAdmin401(t *testing.T) {

	SetURI("http://localhost:8080/")

	tmp := DefaultUser.Key
	DefaultUser.Key = "invalid"

	err := ChangeOtherUserAdmin(&TestUser, true)
	if !errors.Is(err, fault.ErrInvalidAPIKey) {
		t.Fatalf("FAIL: unexpected error: %s, want ErrInvalidAPIKey\n", err)
	}

	DefaultUser.Key = tmp
}

// Test change other user admin with blocked IP.
// The previous test caused an IP block.
func TestChangeOtherUserAdmin403Blocked(t *testing.T) {

	SetURI("http://localhost:8080/")

	err := ChangeOtherUserAdmin(&TestUser, true)
	if !errors.Is(err, fault.ErrBlocked) {
		t.Fatalf("FAIL: unexpected error: %s, want ErrBlocked\n", err)
	}

	time.Sleep(SLEEP_SEC * time.Second)
}

// Test change other user admin with a non admin user.
func TestChangeOtherUserAdmin403NotAdmin(t *testing.T) {

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

// Test change other user admin with invalid username.
func TestChangeOtherUserAdmin404(t *testing.T) {

	SetURI("http://localhost:8080/")

	TestUser.Name = "notexist"

	err := ChangeOtherUserAdmin(&TestUser, true)
	if !errors.Is(err, fault.ErrUserNotFound) {
		t.Fatalf("FAIL: unexpected error: %s, want ErrUserNotFound\n", err)
	}

	TestUser.Name = "test"
}

// Test get users.
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

// Test get users with invalid API key.
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

// Test get user with blocked IP.
func TestGetUsers403Blocked(t *testing.T) {

	SetURI("http://localhost:8080/")

	_, err := GetUsers()
	if !errors.Is(err, fault.ErrBlocked) {
		t.Fatalf("FAIL: unexpected error: %s, want ErrBlocked\n", err)
	}

	time.Sleep(SLEEP_SEC * time.Second)
}

// Test get users with a non admin account.
func TestGetUsers403NotAdmin(t *testing.T) {

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

// Test delete user.
func TestUserDelete200(t *testing.T) {

	SetURI("http://localhost:8080/")

	err := Delete(TestUser, true)
	if err != nil {
		t.Fatalf("FAIL: %s\n", err)
	}
}

// Test delete user with invalid API key.
// The user is deleted in the previous test.
func TestUserDelete401(t *testing.T) {

	SetURI("http://localhost:8080/")

	err := Delete(TestUser, true)
	if !errors.Is(err, fault.ErrInvalidAPIKey) {
		t.Fatalf("FAIL: unexpected error: %s, want ErrInvalidAPIKey\n", err)
	}
}

// Test delete user with blocked IP.
func TestUserDelete403(t *testing.T) {

	SetURI("http://localhost:8080/")

	err := Delete(TestUser, true)
	if !errors.Is(err, fault.ErrBlocked) {
		t.Fatalf("FAIL: unexpected error: %s, want ErrBlocked\n", err)
	}

	time.Sleep(SLEEP_SEC * time.Second)
}
