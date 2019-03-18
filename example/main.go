package main

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/playnet-public/demo/faction"
	"github.com/playnet-public/demo/fakedb"
)

func main() {
	ctx := context.Background()

	repo := &fakedb.Repository{}
	manager := &fakeManager{Repository: repo}

	incomplete := faction.NewIncomplete(ctx, "foo", "bar")

	complete, err := manager.Create(ctx, incomplete)
	if err != nil {
		panic(errors.Wrap(err, "creating faction"))
	}

	fmt.Printf("created %#v: %#v\n", complete.ID(), complete.Data())
	complete.Data().SetDescription("foobar")
	manager.Update(ctx, complete)
	fmt.Printf("updated %#v: %#v\n", complete.ID(), complete.Data())

}

type fakeManager struct {
	faction.Repository
}
