package database

// GENERIC ERRORS

type DatabaseUninitializedError struct {}
func (duError *DatabaseUninitializedError) DUError() string {
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

// COURSE ERRORS

type RemoveCourseError struct {
	Err error
}
func (rcErr *RemoveCourseError) RCError() string {
	return "Removing a course failed. Please try again. Error: " + rcErr.Err.Error()
}

type UpdateCourseError struct {
	Err error
}
func (ucErr *UpdateCourseError) UCError() string {
	return "Updating a course failed. Please try again. Error: " + ucErr.Err.Error()
}

type GetCourseError struct {
	Err error
}
func (gcErr *GetCourseError) GCError() string {
	return "Getting a course failed. PLease try again. Error: " + gcErr.Err.Error()
}

type AddCourseError struct {
	Err error
}
func (acErr *AddCourseError) ACError() string {
	return "Adding a course failed. Please try again. Error: " + acErr.Err.Error()
}