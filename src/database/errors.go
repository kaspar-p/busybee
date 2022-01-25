package database

// GENERIC ERRORS

type DatabaseUninitializedError struct {}
func (duError *DatabaseUninitializedError) Error() string {
	return "Database operation failed because database was nil. Please try again."
}

// USER ERRORS

type AddUserError struct{
	Err error
}
func (auErr *AddUserError) AUError() string {
	return "Database operation failed. Please try again." + auErr.Err.Error()
}

type RemoveUserError struct{
	Err error
}
func (ruErr *RemoveUserError) RUError() string {
	return "Removing a user failed. Please try again. Error: " + ruErr.Err.Error()
}

type GetUserError struct {
	Err error
}
func (guErr *GetUserError) GUError() string {
	return "Getting a user failed. Please try again. Error: " + guErr.Err.Error()
}

type UpdateUserError struct{
	Err error
}
func (uuErr *UpdateUserError) UUError() string {
	return "Updating a user failed. Please try again. Error: " + uuErr.Err.Error()
}

// BUSY TIMES ERRORS

type AddBusyTimeError struct{
	Err error
}
func (abtErr *AddBusyTimeError) ABTError() string {
	return "Database operation failed. Please try again." + abtErr.Err.Error()
}

type RemoveBusyTimeError struct{
	Err error
}
func (rbtErr *RemoveBusyTimeError) RBTError() string {
	return "Removing a user failed. Please try again. Error: " + rbtErr.Err.Error()
}

type GetBusyTimeError struct {
	Err error
}
func (gbtErr *GetBusyTimeError) GBTError() string {
	return "Getting a user failed. Please try again. Error: " + gbtErr.Err.Error()
}

type UpdateBusyTimeError struct{
	Err error
}
func (ubtErr *UpdateBusyTimeError) UBTError() string {
	return "Updating a user failed. Please try again. Error: " + ubtErr.Err.Error()
}

// GUILDS ERRORS
type AddGuildRolePairError struct{
	Err error
}
func (addGuildRolePairError *AddGuildRolePairError) Error() string {
	return "Database operation failed. Please try again." + addGuildRolePairError.Err.Error()
}

type GetGuildRolePairError struct {
	Err error
}
func (getGuildRolePairError *GetGuildRolePairError) Error() string {
	return "Getting guild role pair error occurred. Please try again." + getGuildRolePairError.Err.Error()
}