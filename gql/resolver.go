package gql

import (
        "errors"
	"context"
        "fmt"
	"github.com/99designs/gqlgen/graphql"
        "go-contacts/args"
        "go-contacts/controllers"
        u "go-contacts/utils"
) // THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

type Resolver struct{}

func NewRootResolvers() Config {
    c := Config{Resolvers: &Resolver{}}

    // Complexity
    c.Complexity.User.Contacts = func(childComplexity int, limit *int, offset *int) int {
        return *limit * childComplexity
    }
    c.Complexity.Mutation.AddContacts = func(childComplexity int, input []NewContact) int {
        return len(input) * childComplexity
    }
    c.Complexity.Mutation.DeleteContacts = func(childComplexity int, input []int64) int {
        return len(input) * childComplexity
    }
    c.Complexity.Mutation.UpdateContacts = func(childComplexity int, input []ContactUpdate) int {
        return len(input) * childComplexity
    }

    // Schema Directive
    c.Directives.IsAuthenticated = func(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
        ctxUserID := ctx.Value( args.SecureKey )
        if ctxUserID != nil {
            return next(ctx)
        } else {
            return nil, errors.New("Unauthorized Action")
        }
    }
    return c
}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}
func (r *Resolver) User() UserResolver {
	return &userResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) CreateUser(ctx context.Context, input NewUser) (*User, error) {
    uid, err := DBCreateUser(&input)
    if err != nil {
        return nil, err
    }
    user := &User{
        Uid: uid, // Uid is only set by this program. So it is safe to use this
        Username: input.Username,
        Email: input.Email,
    }
    return user, nil
}

func (r *mutationResolver) AddContacts(ctx context.Context, input []NewContact) ([]*int64, error) {
    uid := controllers.UserFromContext( ctx )
    if uid < 0 {
        return nil, fmt.Errorf("Unauthorized Action. User ID missing")
    }
    res := []*int64{}
    for _, contact := range input {
        cid, err := DBAddContact(&contact, uid)
        if err != nil {
            res = append(res, nil) // if this contact is failed to insert, return nul
        } else {
            res = append(res, &cid) // if this contact inserted successfully, return its id
        }
    }
    return res, nil
}
func (r *mutationResolver) DeleteContacts(ctx context.Context, input []int64) ([]*int64, error) {
    uid := controllers.UserFromContext( ctx )
    if uid < 0 {
        return nil, fmt.Errorf("Unauthorized Action. User ID missing")
    }
    res := []*int64{}
    for _, cid := range input {
        var t = cid
        err := DBDeleteContact(uid, cid)
        if err != nil {
             res = append(res, &t) // return the list of undeleted IDs
        }
    }
    return res, nil
}
func (r *mutationResolver) UpdateContacts(ctx context.Context, input []ContactUpdate) ([]*int64, error) {
    uid := controllers.UserFromContext( ctx )
    if uid < 0 {
        return nil, fmt.Errorf("Unauthorized Action. User ID missing")
    }
    res := []*int64{}
    for _, contact := range input {
        cid, err := DBUpdateContact(&contact, uid)
        if err != nil {
            res = append(res, &cid) // return the list of contact IDs that failed to be updated
        }
    }
    return res, nil
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) User(ctx context.Context) (*User, error) {
    uid := controllers.UserFromContext( ctx )
    if uid < 0 {
        return nil, fmt.Errorf("Unauthorized Action. User ID missing")
    }
    user := &User{Uid: uid}
    err := DBGetUser(user)
    if err != nil {
        u.Logger.Println("[ERROR]:", err)
    }
    return user, err
}

type userResolver struct{ *Resolver }

func (r *userResolver) Contacts(ctx context.Context, obj *User, limit *int, offset *int) ([]Contact, error) {
    contacts, err := DBGetContacts(obj.Uid, *limit, *offset)
    if err != nil {
        u.Logger.Println("[ERROR]:", err)
    }
    return contacts, err
}
