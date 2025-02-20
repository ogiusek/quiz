package usersmodule

import (
	"log"
	"quizapi/common"
	"quizapi/modules/filesmodule"
	"quizapi/modules/modelmodule"
	"quizapi/modules/wsmodule"
)

// register
// login
// refresh
// profile
// change name
// change profile picture
// change password
// close connection

// register

type RegisterArgs struct {
	Name     UserName     `json:"name"`
	Password UserPassword `json:"password"`
}

func (args *RegisterArgs) Valid() []error {
	var errors []error
	for _, err := range args.Name.Valid() {
		errors = append(errors, common.ErrPath(err).Property("name"))
	}
	for _, err := range args.Password.Valid() {
		errors = append(errors, common.ErrPath(err).Property("password"))
	}
	return errors
}

func (args *RegisterArgs) Handle(c common.Ioc) error {
	var userRepo UserRepository
	var hasher common.Hasher
	c.Inject(&userRepo)
	c.Inject(&hasher)

	err := userRepo.Create(c, NewUser(
		modelmodule.NewModel(c),
		args.Name,
		filesmodule.DefaultUserImage,
		args.Password,
		hasher,
	))

	if err == common.ErrRepositoryConflict {
		return errUserConflict
	}

	return err
}

// login

type LogInArgs struct {
	Login    UserName     `json:"login"`
	Password UserPassword `json:"password"`
}

type LogInRes TokensDto

func (args *LogInArgs) Valid() []error {
	var errors []error
	for _, err := range args.Login.Valid() {
		errors = append(errors, common.ErrPath(err).Property("login"))
	}
	for _, err := range args.Password.Valid() {
		errors = append(errors, common.ErrPath(err).Property("password"))
	}
	return errors
}

func (args *LogInArgs) Handle(c common.Ioc) error {
	var userRepo UserRepository
	c.Inject(&userRepo)

	user := userRepo.GetByName(c, args.Login)

	if user == nil {
		return errUserInvalidCredentials
	}

	session := NewSessionDto(*user)
	sessionToken := session.Encode(c)
	refreshPayload := NewRefreshDto(c, sessionToken, *user)
	refreshToken := refreshPayload.Encode(c)

	var resultStorage common.ServiceStorage[common.Response]
	c.Inject(&resultStorage)
	resultStorage.Set(LogInRes{
		SessionToken: sessionToken,
		RefreshToken: refreshToken,
	})
	return nil
}

// refresh

type RefreshArgs TokensDto
type RefreshRes TokensDto

func (args *RefreshArgs) Valid() []error {
	return ((*TokensDto)(args)).Valid()
}

func (args *RefreshArgs) Handle(c common.Ioc) error {
	var previousRefreshDto RefreshDto
	if err := previousRefreshDto.Decode(c, args.RefreshToken); err != nil {
		return errNotARefreshToken
	}

	if !previousRefreshDto.Matches(c, args.SessionToken) {
		return errFabricatedSessionToken
	}

	var userRepo UserRepository
	c.Inject(&userRepo)

	user := userRepo.GetById(c, previousRefreshDto.UserId)
	if user == nil {
		return ErrUserNotFound
	}

	session := NewSessionDto(*user)
	sessionToken := session.Encode(c)

	refreshDto := NewRefreshDto(c, sessionToken, *user)
	refreshToken := refreshDto.Encode(c)

	var resStorage common.ServiceStorage[common.Response]
	c.Inject(&resStorage)
	resStorage.Set(RefreshRes{
		SessionToken: sessionToken,
		RefreshToken: refreshToken,
	})
	return nil
}

// profile

type ProfileArgs struct{}
type ProfileRes UserDto

func (args *ProfileArgs) Handle(c common.Ioc) error {
	var sessionStorage common.ServiceStorage[SessionDto]
	c.Inject(&sessionStorage)

	session := sessionStorage.Get()
	if session == nil {
		return ErrUnauthorized
	}

	var userRepo UserRepository
	c.Inject(&userRepo)

	user := userRepo.GetById(c, session.UserId)
	if user == nil {
		return ErrUserNotFound
	}

	var resStorage common.ServiceStorage[common.Response]
	c.Inject(&resStorage)
	resStorage.Set(ProfileRes(user.Dto()))
	return nil
}

// change name

type ChangeNameArgs struct {
	NewName UserName `json:"new_name"`
}

func (args *ChangeNameArgs) Valid() []error {
	var errors []error
	for _, err := range args.NewName.Valid() {
		errors = append(errors, common.ErrPath(err).Property("new_name"))
	}
	return errors
}

func (args *ChangeNameArgs) Handle(c common.Ioc) error {
	var sessionStorage common.ServiceStorage[SessionDto]
	c.Inject(&sessionStorage)

	session := sessionStorage.Get()
	if session == nil {
		return ErrUnauthorized
	}

	var userRepo UserRepository
	c.Inject(&userRepo)

	user := userRepo.GetById(c, session.UserId)
	if user == nil {
		return ErrUserNotFound
	}

	user.ChangeName(args.NewName)

	if err := userRepo.Update(c, *user); err == common.ErrRepositoryConflict {
		return errUserConflict
	} else if err == common.ErrRepositoryParallelModification {
		return common.ErrHttpParallelModification
	}

	return nil
}

// change profile picture

type ChangeProfilePictureArgs struct {
	Image filesmodule.File
}

func (args *ChangeProfilePictureArgs) Handle(c common.Ioc) error {
	var sessionStorage common.ServiceStorage[SessionDto]
	c.Inject(&sessionStorage)

	session := sessionStorage.Get()
	if session == nil {
		return ErrUnauthorized
	}

	var userRepo UserRepository
	c.Inject(&userRepo)

	user := userRepo.GetById(c, session.UserId)
	if user == nil {
		return ErrUserNotFound
	}

	var fileStorage filesmodule.FileStorage
	c.Inject(&fileStorage)

	previousImage := user.Image

	file, err := fileStorage.UploadFile(args.Image, func(fileId filesmodule.FileId) bool {
		image := filesmodule.NewImageId(fileId)
		err := image.Valid()
		return err == nil
	})

	if err != nil {
		return err
	}

	if err := fileStorage.UnTrack(previousImage.File()); err != nil {
		var logger log.Logger
		c.Inject(&logger)
		logger.Print("shutting down app because of an error. failed to remove used file. determine why used file got deleted")
		logger.Panic(err.Error())
	}

	user.ChangeProfilePicture(filesmodule.NewImageId(file))

	if err := userRepo.Update(c, *user); err == common.ErrRepositoryParallelModification {
		return modelmodule.ErrParallelModification
	} else if err == common.ErrRepositoryParallelModification {
		return common.ErrHttpParallelModification
	}

	return nil
}

// change password

type ChangePasswordArgs struct {
	NewPassword UserPassword `json:"new_password"`
}

func (args *ChangePasswordArgs) Valid() []error {
	var errors []error
	for _, err := range args.NewPassword.Valid() {
		errors = append(errors, common.ErrPath(err).Property("new_password"))
	}
	return errors
}

func (args *ChangePasswordArgs) Handle(c common.Ioc) error {
	var sessionStorage common.ServiceStorage[SessionDto]
	c.Inject(&sessionStorage)

	session := sessionStorage.Get()
	if session == nil {
		return ErrUnauthorized
	}

	var userRepo UserRepository
	c.Inject(&userRepo)

	user := userRepo.GetById(c, session.UserId)
	if user == nil {
		return ErrUserNotFound
	}

	var hasher common.Hasher
	c.Inject(&hasher)

	user.ChangePassword(args.NewPassword, hasher)

	if err := userRepo.Update(c, *user); err == common.ErrRepositoryParallelModification {
		return modelmodule.ErrParallelModification
	} else if err == common.ErrRepositoryParallelModification {
		return common.ErrHttpParallelModification
	}

	return nil
}

// close connection

type CloseConnectionArgs struct{}

func (args *CloseConnectionArgs) Handle(c common.Ioc) error {
	var storage common.ServiceStorage[wsmodule.SocketId]
	c.Inject(&storage)
	var repo UserSocketRepo
	c.Inject(&repo)
	repo.DeleteBySocketId(c, storage.MustGet())
	return nil
}
